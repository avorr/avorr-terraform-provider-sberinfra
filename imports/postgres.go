package imports

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"

	"stash.sigma.sbrf.ru/sddevops/terraform-provider-di/models"
	"stash.sigma.sbrf.ru/sddevops/terraform-provider-di/utils"
)

func (o *Importer) importPostgres() error {
	responseBytes, err := o.Api.NewRequestRead(fmt.Sprintf(
		"projects/%s/orders?service_type=db&f[ir_group]=postgres",
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

	if o.Postgres == nil {
		o.Postgres = make([]*models.Server, 0)
	}
	o.Files["dbImports"].Buffer.Write([]byte("\n"))
	for _, v := range resp["project_orders"].([]interface{}) {
		vmMap := v.(map[string]interface{})

		obj := &models.Postgres{}
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

		server.ResId = server.Id.String()
		server.ResType = "di_postgres"
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
			pswd := *server.ResAppParams.PostgresDbPassword
			postgresDbPasswordString := fmt.Sprintf(
				"data.vault_generic_secret.%s-%s.data.postgres_db_password",
				server.ResType,
				server.ResName,
			)
			server.ResAppParams.PostgresDbPassword = &postgresDbPasswordString

			secret := &models.VaultGenericSecret{
				Field:   "postgres_db_password",
				Data:    pswd,
				ResType: "vault_generic_secret",
				ResName: fmt.Sprintf("%s-%s", server.ResType, server.ResName),
				Path: fmt.Sprintf(
					"%s/%s/%s/%s/%s/%s/%s",
					o.VaultKV,
					o.Domain.ResName,
					o.Group.ResName,
					o.StandType.ResName,
					o.Project.ResName,
					server.ResType,
					server.ResName,
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

		data := server.ToHCL()
		_, err = o.Files["postgres"].Buffer.Write(data)
		if err != nil {
			return err
		}
		_, err = o.Files["dbImports"].Buffer.WriteString(fmt.Sprintf(
			importStr,
			server.ResType,
			server.ResName,
			server.Id.String(),
		))

		o.Postgres = append(o.Postgres, &server)
		log.Printf("%s %s %s\n", server.IrType, server.Name, server.ServiceName)
	}
	return nil
}
