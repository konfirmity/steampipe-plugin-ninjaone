package main

import (
	"github.com/konfirmity/steampipe-plugin-ninjaone/ninjaone"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		PluginFunc: ninjaone.Plugin,
	})
}
