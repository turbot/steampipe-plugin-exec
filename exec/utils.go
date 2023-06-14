package exec

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/mitchellh/go-linereader"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type commandResult struct {
	Output   string `json:"output"`
	ExitCode int    `json:"exit_code"`
}

type outputRow struct {
	LineNumber int
	Line       string
	Stream     string
}

func prepareCommand(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (*exec.Cmd, error) {

	conf := GetConfig(d.Connection)

	plugin.Logger(ctx).Trace("listLocalCommand", "conf", conf)

	command := d.EqualsQualString("command")
	if command == "" {
		// Empty command returns zero rows
		return nil, nil
	}

	plugin.Logger(ctx).Trace("listLocalCommand", "command", command)

	envVal := map[string]string{"TODO": "support_map_config"}
	var env []string
	if len(envVal) > 0 {
		for k, v := range envVal {
			if v != "" {
				entry := fmt.Sprintf("%s=%s", k, v)
				env = append(env, entry)
			}
		}
	}

	plugin.Logger(ctx).Trace("listLocalCommand", "env", env)

	// Choose the shell interpreter and add it to the start of the command
	var cmdargs []string
	if len(conf.Interpreter) > 0 {
		for _, v := range conf.Interpreter {
			if v != "" {
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

	plugin.Logger(ctx).Trace("listLocalCommand", "cmdargs", cmdargs)

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

func outputLinesIntoRows(ctx context.Context, d *plugin.QueryData, r io.Reader, isError bool) error {
	plugin.Logger(ctx).Trace("listRemoteCommandResult", "ctx_done", "outputLinesIntoRows starting...")

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

	plugin.Logger(ctx).Trace("listRemoteCommandResult", "ctx_done", "outputLinesIntoRows done")
	return nil
}

func outputIntoRow(ctx context.Context, d *plugin.QueryData, r io.Reader, isError bool) error {
	plugin.Logger(ctx).Trace("listRemoteCommandResult", "ctx_done", "outputIntoRow starting...")

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

	plugin.Logger(ctx).Trace("listRemoteCommandResult", "ctx_done", "outputIntoRow done")
	return nil
}
