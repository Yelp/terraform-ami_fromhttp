package main

import (
	"github.com/hashicorp/terraform/plugin"
	"terraform-provider-yelpaws/yelpaws"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: yelpaws.Provider,
	})
}
