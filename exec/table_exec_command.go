package exec

import (
	"bytes"
	"context"
	"io"
	"strings"
	"sync"

	communicator "github.com/turbot/go-exec-communicator"
	"github.com/turbot/go-exec-communicator/remote"
	"github.com/turbot/go-exec-communicator/shared"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableExecCommand(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "exec_command",
		Description: "Execute a command locally or on a remote machine and return the output as a single row.",
		List: &plugin.ListConfig{
			Hydrate: listExecCommand,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "command", Require: plugin.Required},
			},
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "stdout_output", Type: proto.ColumnType_STRING, Description: "Standard output from the command."},
			{Name: "stderr_output", Type: proto.ColumnType_STRING, Description: "Standard error output from the command."},
			{Name: "exit_code", Type: proto.ColumnType_INT, Description: "Exit code of the command."},
			{Name: "command", Type: proto.ColumnType_STRING, Transform: transform.FromQual("command"), Description: "Command to be run."},
		},
	}
}

func listExecCommand(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	comm, isLocalConnection, err := GetCommunicator(d.Connection)
	if err != nil {
		plugin.Logger(ctx).Error("listExecCommand", "configuration_error", err)
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
	defer outW.Close()

	var wg sync.WaitGroup

	stdout := ""

	wg.Add(1)
	go func() {
		defer wg.Done()
		buf := new(strings.Builder)
		_, err := io.Copy(buf, outR)
		if err != nil {
			plugin.Logger(ctx).Error("listExecCommand", "output_read_error", err)
		}
		stdout = buf.String()
	}()

	retryCtx, cancel := context.WithTimeout(ctx, comm.Timeout())
	defer cancel()

	// Wait and retry until we establish the connection
	o := shared.Outputter{}
	err = communicator.Retry(retryCtx, func() error {
		return comm.Connect(&o)
	})
	if err != nil {
		plugin.Logger(ctx).Error("listExecCommand", "connection_error", err)
		return nil, err
	}

	// Wait for the context to end and then disconnect
	go func() {
		plugin.Logger(ctx).Debug("listExecCommand", "wait for it...")
		<-ctx.Done()
		plugin.Logger(ctx).Debug("listExecCommand", "done!")
		plugin.Logger(ctx).Debug("listExecCommand", "disconnecting...")
		err = comm.Disconnect()
		if err != nil {
			plugin.Logger(ctx).Error("listExecCommand", "disconnection failure")
		}
		plugin.Logger(ctx).Debug("listExecCommand", "disconnected")
	}()

	cmd = &remote.Cmd{
		Command: command,
		Stdout:  outW,
	}

	result := commandResult{}

	plugin.Logger(ctx).Debug("listExecCommand", "cmd.Start...")
	if err := comm.Start(cmd); err != nil {
		plugin.Logger(ctx).Error("listExecCommand.comm.Start", "command_error", err)
		return nil, err
	}
	plugin.Logger(ctx).Debug("listExecCommand", "cmd.Start done")

	plugin.Logger(ctx).Debug("listExecCommand", "cmd.Wait...")
	if err := cmd.Wait(); err != nil {
		plugin.Logger(ctx).Error("listExecCommand.cmd.Wait", "command_error", err)
		if e, ok := err.(*remote.ExitError); ok {
			result.ExitCode = e.ExitStatus
		}
	}
	plugin.Logger(ctx).Debug("listExecCommand", "cmd.Wait done")

	plugin.Logger(ctx).Debug("listExecCommand", "comm.Disconnect...")
	outW.Close()
	err = comm.Disconnect()
	if err != nil {
		plugin.Logger(ctx).Error("listExecCommand", "disconnection_failure")
	}
	plugin.Logger(ctx).Debug("listExecCommand", "comm.Disconnect done")

	plugin.Logger(ctx).Debug("listExecCommand", "wg waiting...")
	wg.Wait()
	plugin.Logger(ctx).Debug("listExecCommand", "wg done!")

	plugin.Logger(ctx).Debug("listExecCommand", "adding row...")

	// If the command failed, return the stderr output
	if result.ExitCode != 0 {
		d.StreamListItem(ctx, commandResult{StderrOutput: stdout, ExitCode: result.ExitCode})
		return nil, nil
	}

	d.StreamListItem(ctx, commandResult{StdoutOutput: stdout, ExitCode: result.ExitCode})
	return nil, nil
}

func listLocalCommand(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	cmd, err := prepareCommand(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("listLocalCommand", "prepareCommand", "command_error", err)
		return nil, err
	}

	if cmd == nil {
		// Empty command returns zero rows
		plugin.Logger(ctx).Debug("listLocalCommand", "cmd", cmd)
		return nil, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		plugin.Logger(ctx).Error("listLocalCommand", "error reading stdout", err)
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		plugin.Logger(ctx).Error("listLocalCommand", "error reading stderr", err)
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		plugin.Logger(ctx).Error("listLocalCommand", "error starting command", err)
		return nil, err
	}

	// Create slices to store standard output and standard error lines
	var stdoutBytes bytes.Buffer
	var stderrBytes bytes.Buffer

	// Read standard output and standard error concurrently
	_, err = io.Copy(&stdoutBytes, stdout)
	if err != nil {
		plugin.Logger(ctx).Error("listLocalCommand", "error reading stdout", err)
	}

	_, err = io.Copy(&stderrBytes, stderr)
	if err != nil {
		plugin.Logger(ctx).Error("listLocalCommand", "error reading stderr", err)
	}

	if err := cmd.Wait(); err != nil {
		// Log the error, but don't fail. The command error output will be captured and returned to the user.
		plugin.Logger(ctx).Error("listLocalCommand", "cmd.Wait", "command_error", err)
	}

	result := commandResult{
		StdoutOutput: stdoutBytes.String(),
		StderrOutput: stderrBytes.String(),
		ExitCode:     cmd.ProcessState.ExitCode(),
	}
	d.StreamListItem(ctx, result)

	return nil, nil
}
