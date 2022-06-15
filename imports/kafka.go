package imports

import (
	"encoding/json"
	"fmt"
	"log"

	"stash.sigma.sbrf.ru/sddevops/terraform-provider-di/utils"

	"github.com/google/uuid"
	"stash.sigma.sbrf.ru/sddevops/terraform-provider-di/models"
)

func (o *Importer) importKafka() error {
	responseBytes, err := o.Api.NewRequestRead(fmt.Sprintf(
		"projects/%s/orders?service_type=app&f[ir_group]=kafka",
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

	if o.Kafka == nil {
		o.Kafka = make([]*models.Cluster, 0)
	}
	o.Files["appImports"].Buffer.Write([]byte("\n"))
	for _, v := range resp["project_orders"].([]interface{}) {
		vmMap := v.(map[string]interface{})

		obj := &models.Kafka{}
		resource := obj.NewObj()
		resource = obj

		cluster := models.Cluster{
			Id:     uuid.MustParse(vmMap["id"].(string)),
			Object: resource,
		}

		// err := server.GetPubKey()
		// if err != nil {
		// 	return err
		// }
		responseBytes, err := cluster.ReadDI()
		if err != nil {
			return err
		}
		err = cluster.Deserialize(responseBytes)
		if err != nil {
			return err
		}
		// if len(obj.AppParams) > 0 {
		// cluster.ResAppParams = &models.HCLAppParams{}
		// paramsBytes, err := json.Marshal(obj.AppParams)
		// if err != nil {
		// 	return err
		// }
		// err = json.Unmarshal(paramsBytes, &cluster.ResAppParams)
		// if err != nil {
		// 	return err
		// }
		// }

		cluster.ResId = cluster.Id.String()
		cluster.ResType = "di_kafka"
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
		data := cluster.ToHCL()
		_, err = o.Files["kafka"].Buffer.Write(data)
		if err != nil {
			return err
		}
		_, err = o.Files["appImports"].Buffer.WriteString(fmt.Sprintf(
			importStr,
			cluster.ResType,
			cluster.ResName,
			cluster.Id.String(),
		))

		o.Kafka = append(o.Kafka, &cluster)
		log.Printf("%s %s %s %s\n", cluster.Type, cluster.IrGroup, cluster.Name, cluster.ServiceName)
	}
	return nil
}
