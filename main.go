package main

import (
	"bitbucket.org/24g/terraform-provider-dev24g/dev24g"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: dev24g.Provider})
}
