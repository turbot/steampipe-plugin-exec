package exec

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name: "steampipe-plugin-exec",
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		DefaultTransform: transform.FromGo(),
		TableMap: map[string]*plugin.Table{
			"exec_local_command":         tableExecLocalCommand(ctx),
			"exec_local_command_result":  tableExecLocalCommandResult(ctx),
			"exec_remote_command_result": tableExecRemoteCommandResult(ctx),
			"exec_remote_command":        tableExecRemoteCommand(ctx),
			"exec_command":               tableExecCommand(ctx),
			"exec_command_line":          tableExecCommandLine(ctx),
			"exec_run":                   tableExecRun(ctx),
		},
	}
	return p
}
