package exec

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableExecLocalCommandResult(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "exec_local_command_result",
		Description: "Execute a command on the local machine and return as a single row.",
		List: &plugin.ListConfig{
			Hydrate: listLocalCommandResult,
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

/*
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		plugin.Logger(ctx).Error("listLocalCommandResult", "pipe_error", err)
		log.Fatal(err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		plugin.Logger(ctx).Error("listLocalCommandResult", "pipe_error", err)
		return nil, err
	}
*/

/*
	if err := cmd.Start(); err != nil {
		plugin.Logger(ctx).Error("listLocalCommandResult", "command_error", err)
		return nil, err
	}

	result := commandResult{}

	stdoutString, err := io.ReadAll(stdout)
	if err != nil {
		plugin.Logger(ctx).Error("listLocalCommandResult", "stdout_error", err)
		return nil, err
	}
	result.Stdout = string(stdoutString)

	stderrString, err := io.ReadAll(stderr)
	if err != nil {
		plugin.Logger(ctx).Error("listLocalCommandResult", "stderr_error", err)
		return nil, err
	}
	result.Stderr = string(stderrString)

	if err := cmd.Wait(); err != nil {
		// Log the error, but don't fail. The command error output will be captured
		// and returned to the user.
		plugin.Logger(ctx).Error("listLocalCommandResult", "command_error", err)
	}

	result.ExitCode = cmd.ProcessState.ExitCode()

*/
