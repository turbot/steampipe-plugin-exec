package exec

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-sdk/v5/rate_limiter"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name: "steampipe-plugin-exec",
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		RateLimiters: []*rate_limiter.Definition{
			{
				MaxConcurrency: 15,
				Name:           "exec_max_concurrency_limiter",
			},
		},
		DefaultTransform: transform.FromGo(),
		TableMap: map[string]*plugin.Table{
			"exec_command":      tableExecCommand(ctx),
			"exec_command_line": tableExecCommandLine(ctx),
		},
	}
	return p
}
