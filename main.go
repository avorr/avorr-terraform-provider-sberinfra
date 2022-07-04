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
	"stash.sigma.sbrf.ru/sddevops/terraform-provider-di/imports"
	inventory_yaml "stash.sigma.sbrf.ru/sddevops/terraform-provider-di/inventory-yaml"
	"stash.sigma.sbrf.ru/sddevops/terraform-provider-di/models"
	"stash.sigma.sbrf.ru/sddevops/terraform-provider-di/views"
)

func main() {
	if len(os.Args) == 3 && os.Args[1] == "import" {
		importer := imports.Importer{}
		err := importer.Import(os.Args[2])
		if err != nil {
			panic(err)
		}
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

	inventory := inventory_yaml.NewInventory()
	inventory.IsDisabled()
	err := inventory.FromBIN()
	if err != nil {
		log.Println(err)
	}
	views.Inventory = inventory

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
			"di_stand_type": {
				ReadContext: ReadDataResource(&models.StandType{}),
				Schema:      models.SchemaStandType,
			},
			"di_as": {
				ReadContext: ReadDataResource(&models.AS{}),
				Schema:      models.SchemaAS,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"di_project": {
				Importer: &schema.ResourceImporter{
					// State:        schema.ImportStatePassthrough,
					StateContext: views.ProjectImport,
				},
				CreateContext: views.ProjectCreate,
				ReadContext:   views.ProjectRead,
				UpdateContext: views.ProjectUpdate,
				DeleteContext: views.ProjectDelete,
				Schema:        models.SchemaProject,
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
			"di_nginx": {
				Importer: &schema.ResourceImporter{
					State: schema.ImportStatePassthrough,
				},
				CreateContext: views.CreateResource(&models.Nginx{}),
				ReadContext:   views.ReadResource(&models.Nginx{}),
				UpdateContext: views.UpdateResource(&models.Nginx{}),
				DeleteContext: views.DeleteResource(&models.Nginx{}),
				Schema:        models.SchemaNginx,
				Timeouts: &schema.ResourceTimeout{
					Create: schema.DefaultTimeout(timeout),
				},
			},
			"di_sowa": {
				Importer: &schema.ResourceImporter{
					State: schema.ImportStatePassthrough,
				},
				CreateContext: views.CreateResource(&models.Sowa{}),
				ReadContext:   views.ReadResource(&models.Sowa{}),
				UpdateContext: views.UpdateResource(&models.Sowa{}),
				DeleteContext: views.DeleteResource(&models.Sowa{}),
				Schema:        models.SchemaSowa,
				Timeouts: &schema.ResourceTimeout{
					Create: schema.DefaultTimeout(timeout),
				},
			},
			"di_openshift": {
				Importer: &schema.ResourceImporter{
					State: schema.ImportStatePassthrough,
				},
				CreateContext: views.CreateResource(&models.Openshift{}),
				ReadContext:   views.ReadResource(&models.Openshift{}),
				UpdateContext: views.UpdateResource(&models.Openshift{}),
				DeleteContext: views.DeleteResource(&models.Openshift{}),
				Schema:        models.SchemaOpenshift,
				Timeouts: &schema.ResourceTimeout{
					Create: schema.DefaultTimeout(timeout),
				},
			},
			"di_postgres": {
				Importer: &schema.ResourceImporter{
					State: schema.ImportStatePassthrough,
				},
				CreateContext: views.CreateResource(&models.Postgres{}),
				ReadContext:   views.ReadResource(&models.Postgres{}),
				UpdateContext: views.UpdateResource(&models.Postgres{}),
				DeleteContext: views.DeleteResource(&models.Postgres{}),
				Schema:        models.SchemaPostgres,
				Timeouts: &schema.ResourceTimeout{
					Create: schema.DefaultTimeout(timeout),
				},
			},
			"di_postgres_se": {
				Importer: &schema.ResourceImporter{
					State: schema.ImportStatePassthrough,
				},
				CreateContext: views.CreateResource(&models.PostgresSE{}),
				ReadContext:   views.ReadResource(&models.PostgresSE{}),
				UpdateContext: views.UpdateResource(&models.PostgresSE{}),
				DeleteContext: views.DeleteResource(&models.PostgresSE{}),
				Schema:        models.SchemaPostgresSE,
				Timeouts: &schema.ResourceTimeout{
					Create: schema.DefaultTimeout(timeout),
				},
			},
			"di_elk": {
				Importer: &schema.ResourceImporter{
					State: schema.ImportStatePassthrough,
				},
				CreateContext: views.CreateResource(&models.ELK{}),
				ReadContext:   views.ReadResource(&models.ELK{}),
				UpdateContext: views.UpdateResource(&models.ELK{}),
				DeleteContext: views.DeleteResource(&models.ELK{}),
				Schema:        models.SchemaELK,
				Timeouts: &schema.ResourceTimeout{
					Create: schema.DefaultTimeout(timeout),
				},
			},
			"di_kafka": {
				Importer: &schema.ResourceImporter{
					State: schema.ImportStatePassthrough,
				},
				CreateContext: views.CreateClusterResource(&models.Kafka{}),
				ReadContext:   views.ReadClusterResource(&models.Kafka{}),
				UpdateContext: views.UpdateClusterResource(&models.Kafka{}),
				DeleteContext: views.DeleteClusterResource(&models.Kafka{}),
				Schema:        models.SchemaKafka,
				Timeouts: &schema.ResourceTimeout{
					Create: schema.DefaultTimeout(timeout),
				},
			},
			"di_ignite": {
				Importer: &schema.ResourceImporter{
					State: schema.ImportStatePassthrough,
				},
				CreateContext: views.CreateClusterResource(&models.Ignite{}),
				ReadContext:   views.ReadClusterResource(&models.Ignite{}),
				UpdateContext: views.UpdateClusterResource(&models.Ignite{}),
				DeleteContext: views.DeleteClusterResource(&models.Ignite{}),
				Schema:        models.SchemaIgnite,
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
			"di_patroni": {
				Importer: &schema.ResourceImporter{
					State: schema.ImportStatePassthrough,
				},
				CreateContext: views.CreateClusterResource(&models.Patroni{}),
				ReadContext:   views.ReadClusterResource(&models.Patroni{}),
				UpdateContext: views.UpdateClusterResource(&models.Patroni{}),
				DeleteContext: views.DeleteClusterResource(&models.Patroni{}),
				Schema:        models.SchemaPatroni,
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
