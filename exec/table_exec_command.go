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

	/*
		outputDoneCh := make(chan struct{})
		go copyUIOutput2(ctx, d, outR, outputDoneCh)
		errDoneCh := make(chan struct{})
		go copyUIOutput2(ctx, d, errR, errDoneCh)
	*/

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

	/*
		wg.Add(1)
		go copyUIOutput4(ctx, d, outR, &wg)
		wg.Add(1)
		go copyUIOutput4(ctx, d, errR, &wg)
	*/

	/*
		remotePath := comm.ScriptPath()

		if err := comm.UploadScript(remotePath, script); err != nil {
			return fmt.Errorf("Failed to upload script: %v", err)
		}
	*/

	commandCtx := ctx // context.Background()

	retryCtx, cancel := context.WithTimeout(commandCtx, comm.Timeout())
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
		<-commandCtx.Done()
		plugin.Logger(ctx).Warn("listExecCommand", "ctx_done", "done!")
		/*
			<-outputDoneCh
			plugin.Logger(ctx).Warn("listExecCommand", "ctx_done", "outputDoneCh done!")
			<-errDoneCh
			plugin.Logger(ctx).Warn("listExecCommand", "ctx_done", "errDoneCh done!")
		*/
		plugin.Logger(ctx).Warn("listExecCommand", "ctx_done", "disconnecting...")
		comm.Disconnect()
		plugin.Logger(ctx).Warn("listExecCommand", "ctx_done", "disconnected")
	}()

	cmd = &remote.Cmd{
		//Command: remotePath,
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

	// TODO - Prevent crashes from timing problems. Needs a channel type approach.
	//time.Sleep(250 * time.Millisecond)

	plugin.Logger(ctx).Warn("listExecCommand", "ctx_done", "wg waiting...")
	wg.Wait()
	plugin.Logger(ctx).Warn("listExecCommand", "ctx_done", "wg done!")

	plugin.Logger(ctx).Warn("listExecCommand", "ctx_done", "finished")

	return nil, nil

}