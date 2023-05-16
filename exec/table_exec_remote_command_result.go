package exec

import (
	"context"
	"errors"
	"io"
	"sync"

	"github.com/mitchellh/go-linereader"
	communicator "github.com/turbot/go-exec-communicator"
	"github.com/turbot/go-exec-communicator/remote"
	"github.com/turbot/go-exec-communicator/shared"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableExecRemoteCommandResult(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "exec_remote_command_result",
		Description: "Execute a command on the remote machine and return as a single row.",
		List: &plugin.ListConfig{
			Hydrate: listRemoteCommandResult,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "command", Require: plugin.Required},
			},
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "line", Type: proto.ColumnType_STRING, Description: "Line data."},
			{Name: "stream", Type: proto.ColumnType_STRING, Description: "Stream the line was sent to, e.g. stdout or stderr."},
			{Name: "line_number", Type: proto.ColumnType_INT, Description: "Line number within the stream."},
			//{Name: "output", Type: proto.ColumnType_STRING, Description: "Output from the command (both stdout and stderr)."},
			//{Name: "exit_code", Type: proto.ColumnType_INT, Description: "Exit code of the command."},
			{Name: "command", Type: proto.ColumnType_STRING, Transform: transform.FromQual("command"), Description: "Command to be run."},
		},
	}
}

func listRemoteCommandResult(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	conf := GetConfig(d.Connection)

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
		copyUIOutput3(ctx, d, outR)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		copyUIOutput3(ctx, d, errR)
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

	config := shared.ConnectionInfo{
		Timeout: "10s",
		Port:    22,
	}
	if conf.Protocol != nil {
		config.Type = *conf.Protocol
	} else {
		return nil, errors.New("protocol is required in config")
	}
	if conf.Host != nil {
		config.Host = *conf.Host
	} else {
		return nil, errors.New("host is required in config")
	}
	if conf.Port != nil {
		config.Port = uint16(*conf.Port)
	}
	if conf.Https != nil {
		config.HTTPS = *conf.Https
	}
	if conf.Insecure != nil {
		config.Insecure = *conf.Insecure
	}
	if config.Type == "ssh" {
		if conf.User != nil {
			config.User = *conf.User
		} else {
			return nil, errors.New("user is required for SSH connections")
		}
		if conf.Password != nil {
			config.Password = *conf.Password
		} else if conf.PrivateKey != nil {
			config.PrivateKey = *conf.PrivateKey
		} else {
			return nil, errors.New("password or private_key is required for SSH connections")
		}
	}
	if config.Type == "winrm" {
		if conf.User != nil {
			config.User = *conf.User
		} else {
			return nil, errors.New("user is required for WinRM connections")
		}
		if conf.Password != nil {
			config.Password = *conf.Password
		} else {
			return nil, errors.New("password is required for WinRM connections")
		}
	}

	plugin.Logger(ctx).Warn("listRemoteCommandResult", "config", config)

	comm, err := communicator.New(config)
	if err != nil {
		plugin.Logger(ctx).Error("listRemoteCommandResult", "command_error", err)
		return nil, err
	}

	commandCtx := ctx // context.Background()

	retryCtx, cancel := context.WithTimeout(commandCtx, comm.Timeout())
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
		plugin.Logger(ctx).Warn("listRemoteCommandResult", "ctx_done", "wait for it...")
		<-commandCtx.Done()
		plugin.Logger(ctx).Warn("listRemoteCommandResult", "ctx_done", "done!")
		/*
			<-outputDoneCh
			plugin.Logger(ctx).Warn("listRemoteCommandResult", "ctx_done", "outputDoneCh done!")
			<-errDoneCh
			plugin.Logger(ctx).Warn("listRemoteCommandResult", "ctx_done", "errDoneCh done!")
		*/
		plugin.Logger(ctx).Warn("listRemoteCommandResult", "ctx_done", "disconnecting...")
		comm.Disconnect()
		plugin.Logger(ctx).Warn("listRemoteCommandResult", "ctx_done", "disconnected")
	}()

	cmd = &remote.Cmd{
		//Command: remotePath,
		Command: command,
		Stdout:  outW,
		Stderr:  errW,
	}

	result := commandResult{}

	plugin.Logger(ctx).Warn("listRemoteCommandResult", "ctx_done", "cmd.Start...")
	if err := comm.Start(cmd); err != nil {
		plugin.Logger(ctx).Error("listRemoteCommandResult", "command_error", err)
		return nil, err
	}
	plugin.Logger(ctx).Warn("listRemoteCommandResult", "ctx_done", "cmd.Start done")

	plugin.Logger(ctx).Warn("listRemoteCommandResult", "ctx_done", "cmd.Wait...")
	if err := cmd.Wait(); err != nil {
		plugin.Logger(ctx).Error("listRemoteCommandResult", "command_error", err)
		if e, ok := err.(*remote.ExitError); ok {
			result.ExitCode = e.ExitStatus
		}
	}
	plugin.Logger(ctx).Warn("listRemoteCommandResult", "ctx_done", "cmd.Wait done")

	plugin.Logger(ctx).Warn("listRemoteCommandResult", "ctx_done", "comm.Disconnect...")
	outW.Close()
	errW.Close()
	comm.Disconnect()
	plugin.Logger(ctx).Warn("listRemoteCommandResult", "ctx_done", "comm.Disconnect done")

	// TODO - Prevent crashes from timing problems. Needs a channel type approach.
	//time.Sleep(250 * time.Millisecond)

	plugin.Logger(ctx).Warn("listRemoteCommandResult", "ctx_done", "wg waiting...")
	wg.Wait()
	plugin.Logger(ctx).Warn("listRemoteCommandResult", "ctx_done", "wg done!")

	plugin.Logger(ctx).Warn("listRemoteCommandResult", "ctx_done", "finished")

	return nil, nil

}

func copyUIOutput2(ctx context.Context, d *plugin.QueryData, r io.Reader, doneCh chan<- struct{}) error {
	plugin.Logger(ctx).Warn("listRemoteCommandResult", "ctx_done", "copyUIOutput2 starting...")
	defer close(doneCh)
	lr := linereader.New(r)
	i := 1
	for line := range lr.Ch {
		d.StreamListItem(ctx, outputRow{Line: line, LineNumber: i, Stream: "stdout"})
		i = i + 1
		//o.Output(line)
	}
	plugin.Logger(ctx).Warn("listRemoteCommandResult", "ctx_done", "copyUIOutput2 done")
	return nil
}

func copyUIOutput3(ctx context.Context, d *plugin.QueryData, r io.Reader) error {
	plugin.Logger(ctx).Warn("listRemoteCommandResult", "ctx_done", "copyUIOutput3 starting...")
	lr := linereader.New(r)
	i := 1
	for line := range lr.Ch {
		plugin.Logger(ctx).Warn("listRemoteCommandResult", "ctx_done", "copyUIOutput3: "+line)
		d.StreamListItem(ctx, outputRow{Line: line, LineNumber: i, Stream: "stdout"})
		i = i + 1
	}
	plugin.Logger(ctx).Warn("listRemoteCommandResult", "ctx_done", "copyUIOutput3 done")
	return nil
}

func copyUIOutput4(ctx context.Context, d *plugin.QueryData, r io.Reader, wg *sync.WaitGroup) error {
	plugin.Logger(ctx).Warn("listRemoteCommandResult", "ctx_done", "copyUIOutput2 starting...")
	defer wg.Done()
	lr := linereader.New(r)
	i := 1
	for line := range lr.Ch {
		d.StreamListItem(ctx, outputRow{Line: line, LineNumber: i, Stream: "stdout"})
		i = i + 1
		//o.Output(line)
	}
	plugin.Logger(ctx).Warn("listRemoteCommandResult", "ctx_done", "copyUIOutput2 done")
	return nil
}
