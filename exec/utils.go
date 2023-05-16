package exec

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"

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
