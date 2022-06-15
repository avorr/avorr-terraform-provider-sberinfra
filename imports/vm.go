package imports

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"

	"stash.sigma.sbrf.ru/sddevops/terraform-provider-di/models"
	"stash.sigma.sbrf.ru/sddevops/terraform-provider-di/utils"
)

func (o *Importer) importVMs() error {
	responseBytes, err := o.Api.NewRequestRead(fmt.Sprintf(
		"projects/%s/orders?service_type=compute&f[ir_group]=vm",
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

	if o.VMs == nil {
		o.VMs = make([]*models.Server, 0)
	}
	for _, v := range resp["project_orders"].([]interface{}) {
		vmMap := v.(map[string]interface{})

		obj := &models.VM{}
		resource := obj.NewObj()
		resource = obj

		server := models.Server{
			Id:     uuid.MustParse(vmMap["id"].(string)),
			Object: resource,
		}

		responseBytes, err := server.ReadDI()
		if err != nil {
			return err
		}
		err = server.Deserialize(responseBytes)
		if err != nil {
			return err
		}

		err = server.GetPubKey()
		if err != nil {
			// return err
			log.Printf("%s: %s", server.Name, err)
		}

		if len(obj.AppParams) > 0 {
			server.ResAppParams = &models.HCLAppParams{}
			var tmp string
			if obj.AppParams["version_jdk"] != nil {
				tmp = obj.AppParams["version_jdk"].(string)
				server.ResAppParams.VersionJDK = &tmp
			}
		}
		if len(obj.Volumes) > 0 {
			server.ResVolumes = make([]*models.HCLVolume, len(obj.Volumes))
			for k, vol := range obj.Volumes {
				hclVol := &models.HCLVolume{Size: vol.Size}
				if vol.Path != "" {
					hclVol.Path = &vol.Path
				}
				server.ResVolumes[k] = hclVol

			}
		}

		server.ResId = server.Id.String()
		server.ResType = "di_vm"
		server.ResName = utils.Reformat(server.Name)
		server.ResGroupIdUUID = server.GroupId.String()
		server.ResProjectIdUUID = o.Project.Id.String()
		// server.ResGroupId = fmt.Sprintf(
		// 	"data.%s.%s.id",
		// 	o.Group.GetResType(),
		// 	o.Group.GetResName(),
		// )
		// server.ResProjectId = fmt.Sprintf(
		// 	"%s.%s.id",
		// 	o.Project.ResType,
		// 	o.Project.ResName,
		// )
		server.ResGroupId = "data.terraform_remote_state.project.outputs.di_group_id"
		server.ResProjectId = "data.terraform_remote_state.project.outputs.di_project_id"

		data := server.ToHCL()
		// fmt.Println(string(data))
		_, err = o.Files["vm"].Buffer.Write(data)
		if err != nil {
			return err
		}
		_, err = o.Files["computeImports"].Buffer.WriteString(fmt.Sprintf(
			importStr,
			server.ResType,
			server.ResName,
			server.Id.String(),
		))

		o.VMs = append(o.VMs, &server)
		log.Printf("%s %s %s\n", server.IrType, server.Name, server.ServiceName)
	}
	return nil
}
