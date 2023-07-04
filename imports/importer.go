package imports

import (
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"gitlab.gos-tech.xyz/pid/iac/terraform-provider-sberinfra/client"
	"gitlab.gos-tech.xyz/pid/iac/terraform-provider-sberinfra/models"
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
	//project.ID = id
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
