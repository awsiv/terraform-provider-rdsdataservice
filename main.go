package main

import (
	"github.com/awsiv/terraform-provider-rdsdataservice/rdsdataservice"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: rdsdataservice.Provider})
}
