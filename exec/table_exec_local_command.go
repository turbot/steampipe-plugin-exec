package exec

import (
	"bufio"
	"context"
	"log"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableExecLocalCommand(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "exec_local_command",
		Description: "Execute a command on the local machine.",
		List: &plugin.ListConfig{
			Hydrate: listLocalCommand,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "command", Require: plugin.Required},
			},
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "line", Type: proto.ColumnType_STRING, Description: "Line data."},
			{Name: "stream", Type: proto.ColumnType_STRING, Description: "Stream the line was sent to, e.g. stdout or stderr."},
			{Name: "line_number", Type: proto.ColumnType_INT, Description: "Line number within the stream."},
			{Name: "command", Type: proto.ColumnType_STRING, Transform: transform.FromQual("command"), Description: "Command to be run."},
		},
	}
}

const (
	// maxBufSize limits how much output we collect from a local
	// invocation. This is to prevent TF memory usage from growing
	// to an enormous amount due to a faulty process.
	maxBufSize = 8 * 1024
)

type outputRow struct {
	LineNumber int
	Line       string
	Stream     string
}

func listLocalCommand(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	cmd, err := prepareCommand(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("listLocalCommand", "command_error", err)
		return nil, err
	}

	if cmd == nil {
		// Empty command returns zero rows
		return nil, nil
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		plugin.Logger(ctx).Error("listLocalCommand", "pipe_error", err)
		log.Fatal(err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		plugin.Logger(ctx).Error("listLocalCommand", "pipe_error", err)
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		plugin.Logger(ctx).Error("listLocalCommand", "command_error", err)
		return nil, err
	}

	stdoutScanner := bufio.NewScanner(stdout)
	lineNumber := 0
	for stdoutScanner.Scan() {
		lineNumber++
		d.StreamListItem(ctx, outputRow{LineNumber: lineNumber, Line: stdoutScanner.Text(), Stream: "stdout"})
	}
	if err := stdoutScanner.Err(); err != nil {
		plugin.Logger(ctx).Error("listLocalCommand", "stdout_error", err)
		return nil, err
	}

	stderrScanner := bufio.NewScanner(stderr)
	for stderrScanner.Scan() {
		lineNumber++
		d.StreamListItem(ctx, outputRow{LineNumber: lineNumber, Line: stderrScanner.Text(), Stream: "stderr"})
	}
	if err := stderrScanner.Err(); err != nil {
		plugin.Logger(ctx).Error("listLocalCommand", "stderr_error", err)
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		// Log the error, but don't fail. The command error output will be captured
		// and returned to the user.
		plugin.Logger(ctx).Error("listLocalCommand", "command_error", err)
	}

	return nil, nil

}
