package main

import (
	"context"
	"github.com/avorr/terraform-provider-sberinfra/client"
	"github.com/avorr/terraform-provider-sberinfra/models"
	"github.com/avorr/terraform-provider-sberinfra/views"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	runPlugin()
}

func runPlugin() {
	models.Api = client.NewApi()
	plugin.Serve(&plugin.ServeOpts{ProviderFunc: ProviderFunc})
}

func ProviderFunc() *schema.Provider {
	timeout := 300 * time.Second
	envTimeout := os.Getenv("SI_TIMEOUT")
	if envTimeout != "" {
		duration, err := strconv.ParseUint(envTimeout, 10, 64)
		if err != nil {
			log.Println(err)
		} else {
			timeout = time.Duration(duration) * time.Second
		}
	}

	return &schema.Provider{
		//ConfigureContextFunc: func(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
		//	log.Println(pp.Sprintln(data))
		//	return nil, nil
		//},
		DataSourcesMap: map[string]*schema.Resource{
			"si_domain": {
				ReadContext: ReadDataResource(&models.Domain{}),
				Schema:      models.SchemaDomain,
			},
			"si_group": {
				ReadContext: ReadDataResource(&models.Group{}),
				Schema:      models.SchemaGroup,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"si_vdc": {
				Importer: &schema.ResourceImporter{
					StateContext: views.VdcImport,
				},
				CreateContext: views.VdcCreate,
				ReadContext:   views.VdcRead,
				UpdateContext: views.VdcUpdate,
				DeleteContext: views.VdcDelete,
				Schema:        models.SchemaVdc,
				Timeouts: &schema.ResourceTimeout{
					Create: schema.DefaultTimeout(timeout),
				},
			},
			"si_vm": {
				Importer: &schema.ResourceImporter{
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
			"si_security_group": {
				Importer: &schema.ResourceImporter{
					StateContext: views.SecurityGroupImport,
				},
				CreateContext: views.SecurityGroupCreate,
				ReadContext:   views.SecurityGroupRead,
				UpdateContext: views.SecurityGroupUpdate,
				DeleteContext: views.SecurityGroupDelete,
				Schema:        models.SchemaSecurityGroup,
				Timeouts: &schema.ResourceTimeout{
					Create: schema.DefaultTimeout(timeout),
				},
			},
			"si_tag": {
				Importer: &schema.ResourceImporter{
					State: schema.ImportStatePassthrough,
				},
				CreateContext: views.TagCreate,
				ReadContext:   views.TagRead,
				//UpdateContext: views.TagUpdate,
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
