package exec

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"

	"github.com/armon/circbuf"
	"github.com/mitchellh/go-linereader"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func listLocalCommandTerraformStyle(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	conf := GetConfig(d.Connection)

	plugin.Logger(ctx).Warn("listLocalCommand", "conf", conf)

	command := d.EqualsQualString("command")
	if command == "" {
		return nil, errors.New("command cannot be empty")
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

	// Set up the reader that will read the output from the command.
	// We use an os.Pipe so that the *os.File can be passed directly to the
	// process, and not rely on goroutines copying the data which may block.
	// See golang.org/issue/18874
	pr, pw, err := os.Pipe()
	if err != nil {
		/*
			resp.Diagnostics = resp.Diagnostics.Append(tfdiags.WholeContainingBody(
				tfdiags.Error,
				"local-exec provisioner error",
				fmt.Sprintf("Failed to initialize pipe for output: %s", err),
			))
			return resp
		*/
		plugin.Logger(ctx).Error("listLocalCommand", "pipe_error", err)
		return nil, err
	}

	var cmdEnv []string
	cmdEnv = os.Environ()
	cmdEnv = append(cmdEnv, env...)

	// Set up the command
	cmd := exec.CommandContext(ctx, cmdargs[0], cmdargs[1:]...)
	cmd.Stderr = pw
	cmd.Stdout = pw
	// Dir specifies the working directory of the command.
	// If Dir is the empty string (this is default), runs the command
	// in the calling process's current directory.
	if conf.WorkingDir != nil {
		cmd.Dir = *conf.WorkingDir
	}
	// Env specifies the environment of the command.
	// By default will use the calling process's environment
	cmd.Env = cmdEnv

	output, _ := circbuf.NewBuffer(maxBufSize)

	// Write everything we read from the pipe to the output buffer too
	tee := io.TeeReader(pr, output)

	// copy the teed output to the UI output
	copyDoneCh := make(chan struct{})
	//go copyUIOutput(req.UIOutput, tee, copyDoneCh)
	go copyUIOutput(ctx, d, tee, copyDoneCh)

	// Output what we're about to run
	//req.UIOutput.Output(fmt.Sprintf("Executing: %q", cmdargs))
	//fmt.Printf("Executing: %q\n", cmdargs)

	plugin.Logger(ctx).Warn("listLocalCommand", "cmd", cmd)

	// Start the command
	err = cmd.Start()
	if err == nil {
		err = cmd.Wait()
	}

	// Close the write-end of the pipe so that the goroutine mirroring output
	// ends properly.
	pw.Close()

	// Cancelling the command may block the pipe reader if the file descriptor
	// was passed to a child process which hasn't closed it. In this case the
	// copyOutput goroutine will just hang out until exit.
	select {
	case <-copyDoneCh:
	case <-ctx.Done():
	}

	if err != nil {
		/*
			resp.Diagnostics = resp.Diagnostics.Append(tfdiags.WholeContainingBody(
				tfdiags.Error,
				"local-exec provisioner error",
				fmt.Sprintf("Error running command '%s': %v. Output: %s", command, err, output.Bytes()),
			))
			return resp
		*/
		plugin.Logger(ctx).Error("listLocalCommand", "command_error", err, "output", string(output.Bytes()))
		//return nil, err
	}

	//return resp

	return nil, nil

}

// func copyUIOutput(o provisioners.UIOutput, r io.Reader, doneCh chan<- struct{}) {
func copyUIOutput(ctx context.Context, d *plugin.QueryData, r io.Reader, doneCh chan<- struct{}) error {
	defer close(doneCh)
	lr := linereader.New(r)
	i := 1
	for line := range lr.Ch {
		d.StreamListItem(ctx, outputRow{Line: line, LineNumber: i, Stream: "stdout"})
		i = i + 1
		//o.Output(line)
	}
	return nil
}
