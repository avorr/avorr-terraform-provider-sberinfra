package imports

import (
	"encoding/json"
	"fmt"
	"log"

	"stash.sigma.sbrf.ru/sddevops/terraform-provider-di/utils"

	"github.com/google/uuid"
	"stash.sigma.sbrf.ru/sddevops/terraform-provider-di/models"
)

func (o *Importer) importOpenshift() error {
	responseBytes, err := o.Api.NewRequestRead(fmt.Sprintf(
		"projects/%s/orders?service_type=compute&f[ir_group]=project",
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

	if o.Openshift == nil {
		o.Openshift = make([]*models.Server, 0)
	}
	o.Files["computeImports"].Buffer.Write([]byte("\n"))
	for _, v := range resp["project_orders"].([]interface{}) {
		vmMap := v.(map[string]interface{})

		obj := &models.Openshift{}
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
			paramsBytes, err := json.Marshal(obj.AppParams)
			if err != nil {
				return err
			}
			err = json.Unmarshal(paramsBytes, &server.ResAppParams)
			if err != nil {
				return err
			}
			// if obj.AppParams["admin_user"] != nil {
			// 	adminUser := obj.AppParams["admin_user"].(string)
			// 	server.ResAppParams.AdminUser = &adminUser
			// }
			// if obj.AppParams["name_project"] != nil {
			// 	nameProject := obj.AppParams["name_project"].(string)
			// 	server.ResAppParams.NameProject = &nameProject
			// }
			// if obj.AppParams["endpoint"] != nil {
			// 	endpoint := obj.AppParams["endpoint"].(string)
			// 	server.ResAppParams.Endpoint = &endpoint
			// }
		}

		server.ResId = server.Id.String()
		server.ResType = "di_openshift"
		server.ResName = utils.Reformat(server.Name)
		server.ResGroupIdUUID = server.GroupId.String()
		// server.ResGroupId = fmt.Sprintf(
		// 	"data.%s.%s.id",
		// 	o.Group.GetResType(),
		// 	o.Group.GetResName(),
		// )
		server.ResProjectIdUUID = o.Project.Id.String()
		// server.ResProjectId = fmt.Sprintf(
		// 	"%s.%s.id",
		// 	o.Project.ResType,
		// 	o.Project.ResName,
		// )
		server.ResGroupId = "data.terraform_remote_state.project.outputs.di_group_id"
		server.ResProjectId = "data.terraform_remote_state.project.outputs.di_project_id"

		data := server.ToHCL()
		_, err = o.Files["openshift"].Buffer.Write(data)
		if err != nil {
			return err
		}
		_, err = o.Files["computeImports"].Buffer.WriteString(fmt.Sprintf(
			importStr,
			server.ResType,
			server.ResName,
			server.Id.String(),
		))

		o.Openshift = append(o.Openshift, &server)
		log.Printf("%s %s %s\n", server.IrType, server.Name, server.ServiceName)
	}
	return nil
}
