package exec

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	"github.com/mitchellh/go-linereader"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func prepareCommand(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (*exec.Cmd, error) {

	conf := GetConfig(d.Connection)

	plugin.Logger(ctx).Warn("listLocalCommand", "conf", conf)

	command := d.EqualsQualString("command")
	if command == "" {
		// Empty command returns zero rows
		return nil, nil
	}

	plugin.Logger(ctx).Warn("listLocalCommand", "command", command)

	//envVal := req.Config.GetAttr("environment")
	envVal := map[string]string{"TODO": "support_map_config"}
	var env []string
	//if !envVal.IsNull() {
	if len(envVal) > 0 {
		//for k, v := range envVal.AsValueMap() {
		for k, v := range envVal {
			//if !v.IsNull() {
			if v != "" {
				//entry := fmt.Sprintf("%s=%s", k, v.AsString())
				entry := fmt.Sprintf("%s=%s", k, v)
				env = append(env, entry)
			}
		}
	}

	plugin.Logger(ctx).Warn("listLocalCommand", "env", env)

	// Choose the shell interpreter and add it to the start of the command
	var cmdargs []string
	//if !intrVal.IsNull() && intrVal.LengthInt() > 0 {
	if len(conf.Interpreter) > 0 {
		//for _, v := range intrVal.AsValueSlice() {
		for _, v := range conf.Interpreter {
			//if !v.IsNull() {
			if v != "" {
				//cmdargs = append(cmdargs, v.AsString())
				cmdargs = append(cmdargs, v)
			}
		}
	} else {
		if runtime.GOOS == "windows" {
			cmdargs = []string{"cmd", "/C"}
		} else {
			cmdargs = []string{"/bin/sh", "-c"}
		}
	}

	// Command comes last
	cmdargs = append(cmdargs, command)

	plugin.Logger(ctx).Warn("listLocalCommand", "cmdargs", cmdargs)

	cmd := exec.CommandContext(ctx, cmdargs[0], cmdargs[1:]...)

	// Dir specifies the working directory of the command.
	// If Dir is the empty string (this is default), runs the command
	// in the calling process's current directory.
	if conf.WorkingDir != nil {
		cmd.Dir = *conf.WorkingDir
	}

	// Env specifies the environment of the command.
	// By default will use the calling process's environment
	var cmdEnv []string
	cmdEnv = os.Environ()
	cmdEnv = append(cmdEnv, env...)
	cmd.Env = cmdEnv

	return cmd, nil

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

func copyUIOutput3(ctx context.Context, d *plugin.QueryData, r io.Reader, isError bool) error {
	plugin.Logger(ctx).Warn("listRemoteCommandResult", "ctx_done", "copyUIOutput3 starting...")

	stream := "stdout"
	if isError {
		stream = "stderr"
	}

	lr := linereader.New(r)
	i := 1
	for line := range lr.Ch {
		d.StreamListItem(ctx, outputRow{Line: line, LineNumber: i, Stream: stream})
		i = i + 1
	}

	plugin.Logger(ctx).Warn("listRemoteCommandResult", "ctx_done", "copyUIOutput3 done")
	return nil
}

func copyUIOutput4(ctx context.Context, d *plugin.QueryData, r io.Reader, wg *sync.WaitGroup) error {
	plugin.Logger(ctx).Warn("listRemoteCommandResult", "ctx_done", "copyUIOutput4 starting...")
	defer wg.Done()
	lr := linereader.New(r)
	i := 1
	for line := range lr.Ch {
		d.StreamListItem(ctx, outputRow{Line: line, LineNumber: i, Stream: "stdout"})
		i = i + 1
		//o.Output(line)
	}
	plugin.Logger(ctx).Warn("listRemoteCommandResult", "ctx_done", "copyUIOutput4 done")
	return nil
}

func copyUIOutput5(ctx context.Context, d *plugin.QueryData, r io.Reader, isError bool) error {
	plugin.Logger(ctx).Warn("listRemoteCommandResult", "ctx_done", "copyUIOutput5 starting...")

	exitCode := 0
	if isError {
		exitCode = 1
	}

	buf := new(strings.Builder)
	n, _ := io.Copy(buf, r)
	if n == 0 {
		return nil
	}
	d.StreamListItem(ctx, commandResult{Output: buf.String(), ExitCode: exitCode})

	plugin.Logger(ctx).Warn("listRemoteCommandResult", "ctx_done", "copyUIOutput5 done")
	return nil
}

func copyUIOutput6(ctx context.Context, d *plugin.QueryData, r io.Reader, isError bool) error {
	plugin.Logger(ctx).Warn("listRemoteCommandResult", "ctx_done", "copyUIOutput6 starting...")
	exitCode := 0
	if isError {
		exitCode = 1
	}

	outputLines := []string{}
	lr := linereader.New(r)
	for line := range lr.Ch {
		outputLines = append(outputLines, line)
	}
	if len(outputLines) == 0 {
		return nil
	}
	d.StreamListItem(ctx, commandLineResult{
		Output:   outputLines,
		ExitCode: exitCode,
	})
	return nil
}
