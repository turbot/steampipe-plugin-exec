package main

import (
	"github.com/turbot/steampipe-plugin-exec/exec"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: exec.Plugin})
}
