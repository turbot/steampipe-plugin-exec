package exec

import (
	"bufio"
	"context"
	"io"
	"log"
	"sync"

	communicator "github.com/turbot/go-exec-communicator"
	"github.com/turbot/go-exec-communicator/remote"
	"github.com/turbot/go-exec-communicator/shared"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableExecCommandLine(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "exec_command_line",
		Description: "Execute a command locally or on a remote machine and return one row per output line.",
		List: &plugin.ListConfig{
			Hydrate: listExecCommandLine,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "command", Require: plugin.Required},
			},
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "line", Type: proto.ColumnType_STRING, Description: "Line data."},
			{Name: "stream", Type: proto.ColumnType_STRING, Description: "Stream the line was sent to, e.g. stdout or stderr."},
			{Name: "line_number", Type: proto.ColumnType_INT, Description: "Line number within the stream."},
			{Name: "command", Type: proto.ColumnType_STRING, Transform: transform.FromQual("command"), Description: "Command to be run."},
		},
	}
}

func listExecCommandLine(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	comm, _, isLocalConnection, err := GetCommunicator(d.Connection)
	if err != nil {
		plugin.Logger(ctx).Error("listRemoteCommandResult", "init", "command_error", err)
		return nil, err
	}
	if isLocalConnection {
		return listLocalCommand(ctx, d, h)
	}

	command := d.EqualsQualString("command")
	if command == "" {
		// Empty command returns zero rows
		return nil, nil
	}

	var cmd *remote.Cmd

	outR, outW := io.Pipe()
	errR, errW := io.Pipe()
	defer outW.Close()
	defer errW.Close()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		err = outputLinesIntoRows(ctx, d, outR, false)
		if err != nil {
			plugin.Logger(ctx).Error("listRemoteCommandResult", "outputLinesIntoRows", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err = outputLinesIntoRows(ctx, d, errR, true)
		if err != nil {
			plugin.Logger(ctx).Error("listRemoteCommandResult", "outputLinesIntoRows", err)
		}
	}()

	retryCtx, cancel := context.WithTimeout(ctx, comm.Timeout())
	defer cancel()

	// Wait and retry until we establish the connection
	o := shared.Outputter{}
	err = communicator.Retry(retryCtx, func() error {
		return comm.Connect(&o)
	})
	if err != nil {
		plugin.Logger(ctx).Error("listRemoteCommandResult", "connection_error", err)
		return nil, err
	}

	// Wait for the context to end and then disconnect
	go func() {
		plugin.Logger(ctx).Debug("listRemoteCommandResult", "ctx_done", "wait for it...")
		<-ctx.Done()
		plugin.Logger(ctx).Debug("listRemoteCommandResult", "ctx_done", "done!")
		plugin.Logger(ctx).Debug("listRemoteCommandResult", "ctx_done", "disconnecting...")
		err = comm.Disconnect()
		if err != nil {
			plugin.Logger(ctx).Error("listRemoteCommandResult", "ctx_done", "disconnection failed")
		}
		plugin.Logger(ctx).Debug("listRemoteCommandResult", "ctx_done", "disconnected")
	}()

	cmd = &remote.Cmd{
		Command: command,
		Stdout:  outW,
		Stderr:  errW,
	}

	result := commandResult{}

	plugin.Logger(ctx).Debug("listRemoteCommandResult", "ctx_done", "cmd.Start...")
	if err := comm.Start(cmd); err != nil {
		plugin.Logger(ctx).Error("listRemoteCommandResult", "comm.Start", "command_error", err)
		return nil, err
	}
	plugin.Logger(ctx).Debug("listRemoteCommandResult", "ctx_done", "cmd.Start done")

	plugin.Logger(ctx).Debug("listRemoteCommandResult", "ctx_done", "cmd.Wait...")
	if err := cmd.Wait(); err != nil {
		plugin.Logger(ctx).Error("listRemoteCommandResult", "comm.Wait", "command_error", err)
		if e, ok := err.(*remote.ExitError); ok {
			result.ExitCode = e.ExitStatus
		}
	}
	plugin.Logger(ctx).Debug("listRemoteCommandResult", "ctx_done", "cmd.Wait done")

	plugin.Logger(ctx).Debug("listRemoteCommandResult", "ctx_done", "comm.Disconnect...")
	outW.Close()
	errW.Close()
	err = comm.Disconnect()
	if err != nil {
		plugin.Logger(ctx).Error("listRemoteCommandResult", "ctx_done", "disconnection failed")
	}
	plugin.Logger(ctx).Debug("listRemoteCommandResult", "ctx_done", "comm.Disconnect done")

	plugin.Logger(ctx).Debug("listRemoteCommandResult", "ctx_done", "wg waiting...")
	wg.Wait()
	plugin.Logger(ctx).Debug("listRemoteCommandResult", "ctx_done", "wg done!")

	plugin.Logger(ctx).Debug("listRemoteCommandResult", "ctx_done", "finished")

	return nil, nil
}

func listLocalCommand(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	cmd, err := prepareCommand(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("listLocalCommand", "command_error", err)
		return nil, err
	}

	if cmd == nil {
		// Empty command returns zero rows
		return nil, nil
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		plugin.Logger(ctx).Error("listLocalCommand", "pipe_error", err)
		log.Fatal(err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		plugin.Logger(ctx).Error("listLocalCommand", "pipe_error", err)
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		plugin.Logger(ctx).Error("listLocalCommand", "command_error", err)
		return nil, err
	}

	stdoutScanner := bufio.NewScanner(stdout)
	lineNumber := 0
	for stdoutScanner.Scan() {
		lineNumber++
		d.StreamListItem(ctx, outputRow{LineNumber: lineNumber, Line: stdoutScanner.Text(), Stream: "stdout"})
	}
	if err := stdoutScanner.Err(); err != nil {
		plugin.Logger(ctx).Error("listLocalCommand", "stdout_error", err)
		return nil, err
	}

	stderrScanner := bufio.NewScanner(stderr)
	for stderrScanner.Scan() {
		lineNumber++
		d.StreamListItem(ctx, outputRow{LineNumber: lineNumber, Line: stderrScanner.Text(), Stream: "stderr"})
	}
	if err := stderrScanner.Err(); err != nil {
		plugin.Logger(ctx).Error("listLocalCommand", "stderr_error", err)
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		// Log the error, but don't fail. The command error output will be captured
		// and returned to the user.
		plugin.Logger(ctx).Error("listLocalCommand", "command_error", err)
	}

	return nil, nil
}
