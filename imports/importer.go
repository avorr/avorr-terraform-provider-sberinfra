package imports

import (
	"github.com/avorr/terraform-provider-sberinfra/client"
	"github.com/avorr/terraform-provider-sberinfra/models"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

type Importer struct {
	Api     *client.Api
	Project *models.Vdc
	Domain  *models.Domain
	Group   *models.Group
}

func (o *Importer) Import(projectId string) diag.Diagnostics {
	var diags diag.Diagnostics
	id, err := uuid.Parse(projectId)
	if err != nil {
		return diags
	}
	project := &models.Vdc{ID: id}
	o.Project = project
	o.Api = client.NewApi()
	models.Api = o.Api
	models.Api.Debug = true

	responseBytes, err := project.ReadDI()
	if err != nil {
		return diags
	}
	err = project.Deserialize(responseBytes)
	if err != nil {
		return diags
	}

	servers := Servers{
		Project: project,
		Api:     o.Api,
	}

	err = servers.Read()
	if err != nil {
		return diags
	}
	return diags
}
