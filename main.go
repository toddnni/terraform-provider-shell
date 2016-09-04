package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/toddnni/terraform-provider-shell/shell"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: shell.Provider,
	})
}
