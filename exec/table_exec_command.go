package exec

import (
	"context"
	"io"
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
		Description: "Execute a command on a local or remote machine and return as a single row.",
		List: &plugin.ListConfig{
			Hydrate: listExecCommand,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "command", Require: plugin.Required},
			},
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "output", Type: proto.ColumnType_STRING, Description: "Output from the command (both stdout and stderr)."},
			{Name: "exit_code", Type: proto.ColumnType_INT, Description: "Exit code of the command."},
			{Name: "command", Type: proto.ColumnType_STRING, Transform: transform.FromQual("command"), Description: "Command to be run."},
		},
	}
}

func listExecCommand(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	comm, _, isLocalConnection, err := GetCommunicator(d.Connection)
	if err != nil {
		plugin.Logger(ctx).Error("listExecCommand", "command_error", err)
		return nil, err
	}
	if isLocalConnection {
		return listLocalCommandResult(ctx, d, h)
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
		copyUIOutput5(ctx, d, outR, false)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		copyUIOutput5(ctx, d, errR, true)
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
		plugin.Logger(ctx).Warn("listExecCommand", "ctx_done", "wait for it...")
		<-ctx.Done()
		plugin.Logger(ctx).Warn("listExecCommand", "ctx_done", "done!")
		plugin.Logger(ctx).Warn("listExecCommand", "ctx_done", "disconnecting...")
		comm.Disconnect()
		plugin.Logger(ctx).Warn("listExecCommand", "ctx_done", "disconnected")
	}()

	cmd = &remote.Cmd{
		Command: command,
		Stdout:  outW,
		Stderr:  errW,
	}

	result := commandResult{}

	plugin.Logger(ctx).Warn("listExecCommand", "ctx_done", "cmd.Start...")
	if err := comm.Start(cmd); err != nil {
		plugin.Logger(ctx).Error("listExecCommand", "command_error", err)
		return nil, err
	}
	plugin.Logger(ctx).Warn("listExecCommand", "ctx_done", "cmd.Start done")

	plugin.Logger(ctx).Warn("listExecCommand", "ctx_done", "cmd.Wait...")
	if err := cmd.Wait(); err != nil {
		plugin.Logger(ctx).Error("listExecCommand", "command_error", err)
		if e, ok := err.(*remote.ExitError); ok {
			result.ExitCode = e.ExitStatus
		}
	}
	plugin.Logger(ctx).Warn("listExecCommand", "ctx_done", "cmd.Wait done")

	plugin.Logger(ctx).Warn("listExecCommand", "ctx_done", "comm.Disconnect...")
	outW.Close()
	errW.Close()
	comm.Disconnect()
	plugin.Logger(ctx).Warn("listExecCommand", "ctx_done", "comm.Disconnect done")

	plugin.Logger(ctx).Warn("listExecCommand", "ctx_done", "wg waiting...")
	wg.Wait()
	plugin.Logger(ctx).Warn("listExecCommand", "ctx_done", "wg done!")

	plugin.Logger(ctx).Warn("listExecCommand", "ctx_done", "finished")

	return nil, nil

}

type commandResult struct {
	Output   string `json:"output"`
	ExitCode int    `json:"exit_code"`
}

func listLocalCommandResult(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	cmd, err := prepareCommand(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("listLocalCommandResult", "command_error", err)
		return nil, err
	}

	if cmd == nil {
		// Empty command returns zero rows
		plugin.Logger(ctx).Debug("listLocalCommandResult", "cmd", cmd)
		return nil, nil
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		// Log the error, but don't fail. The command error output will be captured
		// and returned to the user.
		plugin.Logger(ctx).Error("listLocalCommandResult", "command_error", err)
	}

	outputStr := string(output)

	// NOTE - I considered stripping the final newline. The output string looks
	// weird with the final newline (appears like an extra newline) in SQL
	// results, so is confusing. But removing the newline is altering the output,
	// so on balance better to leave it accurate. I'm leaving this code here as a
	// warning to anyone who is tempted to strip the newline in the future.
	// outputStr := strings.TrimSuffix(outputStr, "\n")

	result := commandResult{
		Output:   outputStr,
		ExitCode: cmd.ProcessState.ExitCode(),
	}
	d.StreamListItem(ctx, result)

	return nil, nil
}
