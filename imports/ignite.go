package imports

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"

	"stash.sigma.sbrf.ru/sddevops/terraform-provider-di/models"
	"stash.sigma.sbrf.ru/sddevops/terraform-provider-di/utils"
)

func (o *Importer) importIgnite() error {
	responseBytes, err := o.Api.NewRequestRead(fmt.Sprintf(
		"projects/%s/orders?service_type=db&f[ir_group]=ignite_se",
		o.Project.Id,
	))
	if err != nil {
		return err
	}
	resp := make(map[string]interface{})
	err = json.Unmarshal(responseBytes, &resp)
	if err != nil {
		return err
	}

	if o.Ignite == nil {
		o.Ignite = make([]*models.Cluster, 0)
	}
	ids := []string{}
	for _, v := range resp["project_orders"].([]interface{}) {
		vmMap := v.(map[string]interface{})
		ids = append(ids, vmMap["id"].(string))
	}
	responseBytes, err = o.Api.NewRequestRead(fmt.Sprintf(
		"projects/%s/orders?service_type=db&f[ir_group]=ignite_se_persistence",
		o.Project.Id,
	))
	if err != nil {
		return err
	}
	err = json.Unmarshal(responseBytes, &resp)
	if err != nil {
		return err
	}
	for _, v := range resp["project_orders"].([]interface{}) {
		vmMap := v.(map[string]interface{})
		ids = append(ids, vmMap["id"].(string))
	}

	o.Files["dbImports"].Buffer.Write([]byte("\n"))
	for _, v := range ids {
		obj := &models.Ignite{}
		resource := obj.NewObj()
		resource = obj
		cluster := models.Cluster{
			Id:     uuid.MustParse(v),
			Object: resource,
		}
		responseBytes, err := cluster.ReadDI()
		if err != nil {
			return err
		}
		err = cluster.Deserialize(responseBytes)
		if err != nil {
			return err
		}

		cluster.ResId = cluster.Id.String()
		cluster.ResType = "di_ignite"
		cluster.ResName = utils.Reformat(cluster.Name)
		cluster.ResGroupIdUUID = cluster.GroupId.String()
		// cluster.ResGroupId = fmt.Sprintf(
		// 	"data.%s.%s.id",
		// 	o.Group.GetResType(),
		// 	o.Group.GetResName(),
		// )
		cluster.ResProjectIdUUID = o.Project.Id.String()
		// cluster.ResProjectId = fmt.Sprintf(
		// 	"%s.%s.id",
		// 	o.Project.ResType,
		// 	o.Project.ResName,
		// )
		cluster.ResGroupId = "data.terraform_remote_state.project.outputs.di_group_id"
		cluster.ResProjectId = "data.terraform_remote_state.project.outputs.di_project_id"

		if len(obj.AppParams) > 0 {
			cluster.ResAppParams = &models.HCLAppParams{}
			paramsBytes, err := json.Marshal(obj.AppParams)
			if err != nil {
				return err
			}
			err = json.Unmarshal(paramsBytes, &cluster.ResAppParams)
			if err != nil {
				return err
			}
			pswd := *cluster.ResAppParams.GGClientPassword
			passwordString := fmt.Sprintf(
				"data.vault_generic_secret.%s-%s.data.ise_client_password",
				cluster.ResType,
				cluster.ResName,
			)
			cluster.ResAppParams.GGClientPassword = &passwordString

			secret := &models.VaultGenericSecret{
				Field:   "ise_client_password",
				Data:    pswd,
				ResType: "vault_generic_secret",
				ResName: fmt.Sprintf("%s-%s", cluster.ResType, cluster.ResName),
				Path: fmt.Sprintf(
					"%s/%s/%s/%s/%s/%s/%s",
					o.VaultKV,
					o.Domain.ResName,
					o.Group.ResName,
					o.StandType.ResName,
					o.Project.ResName,
					cluster.ResType,
					cluster.ResName,
				),
			}
			data := secret.ToHCL()
			_, err = o.Files["vaultTF"].Buffer.Write(data)
			if err != nil {
				return err
			}
			data = secret.ToBash()
			_, err = o.Files["vaultSH"].Buffer.Write(data)
			if err != nil {
				return err
			}
		}

		data := cluster.ToHCL()
		_, err = o.Files["ignite"].Buffer.Write(data)
		if err != nil {
			return err
		}
		_, err = o.Files["dbImports"].Buffer.WriteString(fmt.Sprintf(
			importStr,
			cluster.ResType,
			cluster.ResName,
			cluster.Id.String(),
		))

		o.Ignite = append(o.Ignite, &cluster)
		log.Printf("%s %s %s\n", cluster.Type, cluster.IrGroup, cluster.ServiceName)
	}
	return nil
}
