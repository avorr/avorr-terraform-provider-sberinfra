package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"stash.sigma.sbrf.ru/sddevops/terraform-provider-di/client"
	"stash.sigma.sbrf.ru/sddevops/terraform-provider-di/models"
	"stash.sigma.sbrf.ru/sddevops/terraform-provider-di/views"
)

func main() {
	if len(os.Args) == 3 && os.Args[1] == "import" {
		log.Println("Imports doesnt support")
		//importer := imports.Importer{}
		//err := importer.Import(os.Args[2])
		//if err != nil {
		//	panic(err)
		//}
	} else {
		runPlugin()
	}
}

func runPlugin() {
	models.Api = client.NewApi()
	plugin.Serve(&plugin.ServeOpts{ProviderFunc: ProviderFunc})
}

func ProviderFunc() *schema.Provider {
	timeout := 300 * time.Second
	envTimeout := os.Getenv("DI_TIMEOUT")
	if envTimeout != "" {
		duration, err := strconv.ParseUint(envTimeout, 10, 64)
		if err != nil {
			log.Println(err)
		} else {
			timeout = time.Duration(duration) * time.Second
		}
	}

	//inventory := inventory_yaml.NewInventory()
	//inventory.IsDisabled()
	//err := inventory.FromBIN()
	//if err != nil {
	//	log.Println(err)
	//}
	//views.Inventory = inventory

	return &schema.Provider{
		// ConfigureContextFunc: func(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
		// 	log.Println(pp.Println(data))
		// 	return nil, nil
		// },
		DataSourcesMap: map[string]*schema.Resource{
			"di_domain": {
				ReadContext: ReadDataResource(&models.Domain{}),
				Schema:      models.SchemaDomain,
			},
			"di_group": {
				ReadContext: ReadDataResource(&models.Group{}),
				Schema:      models.SchemaGroup,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			//"di_project": {
			//	Importer: &schema.ResourceImporter{
			//// State:        schema.ImportStatePassthrough,
			//StateContext: views.ProjectImport,
			//},
			//CreateContext: views.ProjectCreate,
			//ReadContext:   views.ProjectRead,
			//UpdateContext: views.ProjectUpdate,
			//DeleteContext: views.ProjectDelete,
			//Schema:        models.SchemaProject,
			//Timeouts: &schema.ResourceTimeout{
			//	Create: schema.DefaultTimeout(timeout),
			//},
			//},

			"di_project": {
				//Importer: &schema.ResourceImporter{
				//	State:        schema.ImportStatePassthrough,
				//	StateContext: views.SiProjectImport,
				//},
				CreateContext: views.SiProjectCreate,
				ReadContext:   views.SiProjectRead,
				UpdateContext: views.SiProjectUpdate,
				DeleteContext: views.SiProjectDelete,
				Schema:        models.SchemaSiProject,
				Timeouts: &schema.ResourceTimeout{
					Create: schema.DefaultTimeout(timeout),
				},
			},

			"di_vm": {
				Importer: &schema.ResourceImporter{
					// State: schema.ImportStatePassthrough,
					StateContext: views.ImportResource(&models.VM{}),
				},
				CreateContext: views.CreateResource(&models.VM{}),
				ReadContext:   views.ReadResource(&models.VM{}),
				UpdateContext: views.UpdateResource(&models.VM{}),
				DeleteContext: views.DeleteResource(&models.VM{}),
				Schema:        models.SchemaVM,
				Timeouts: &schema.ResourceTimeout{
					Create: schema.DefaultTimeout(timeout),
				},
			},
			"di_tag": {
				Importer: &schema.ResourceImporter{
					State: schema.ImportStatePassthrough,
				},
				CreateContext: views.TagCreate,
				ReadContext:   views.TagRead,
				// UpdateContext: views.ProjectUpdate,
				DeleteContext: views.TagDelete,
				Schema:        models.SchemaTag,
				Timeouts: &schema.ResourceTimeout{
					Create: schema.DefaultTimeout(timeout),
				},
			},
		},
	}
}

func ReadDataResource(o models.DIDataResource) schema.ReadContextFunc {
	f := func(ctx context.Context, res *schema.ResourceData, m interface{}) diag.Diagnostics {
		var diags diag.Diagnostics
		obj := o.NewObj()
		obj.ReadTF(res)
		responseBytes, err := obj.ReadDI()
		if err != nil {
			return diag.FromErr(err)
		}
		err = obj.Deserialize(responseBytes)
		if err != nil {
			return diag.FromErr(err)
		}
		obj.WriteTF(res)
		return diags
	}
	return f
}
