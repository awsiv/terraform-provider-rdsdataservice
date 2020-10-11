package main

import (
	"github.com/awsiv/terraform-provider-dataapi/dataapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: dataapi.Provider})
}
