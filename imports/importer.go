package imports

import (
	"gitlab.gos-tech.xyz/pid/iac/terraform-provider-sberinfra/client"
	"gitlab.gos-tech.xyz/pid/iac/terraform-provider-sberinfra/models"
	"bytes"
	"fmt"
	"github.com/google/uuid"
)

type Importer struct {
	Api     *client.Api
	Project *models.Vdc
	Domain  *models.Domain
	Group   *models.Group
	//StandType *models.StandType
	//AS        *models.AS
}

func (o *Importer) Import(projectId string) error {
	id, err := uuid.Parse(projectId)
	if err != nil {
		return fmt.Errorf("can't parse uuid, %s\n", err)
	}
	//project := &models.Vdc{Id: id}
	project := &models.Vdc{}
	project.ID = id
	o.Project = project
	o.Api = client.NewApi()
	models.Api = o.Api
	models.Api.Debug = true

	responseBytes, err := project.ReadDI()
	if err != nil {
		return err
	}
	err = project.Deserialize(responseBytes)
	if err != nil {
		return err
	}

	servers := Servers{
		Project: project,
		Api:     o.Api,
	}

	//servers := Servers{}
	//servers

	err = servers.Read()
	if err != nil {
		return err
	}

	var bbTF, bbSH bytes.Buffer
	bbSH.Write([]byte("#!/usr/bin/env bash\n\n"))

	for _, v := range servers.NonCluster {
		bbTF.Write(v.HCLHeader())
		bbSH.Write(v.ImportCmd())
	}
	//for _, v := range servers.Clusters {
	//	bbTF.Write(v.HCLHeader())
	//	bbSH.Write(v.ImportCmd())
	//}

	//err = ioutil.WriteFile("imports.tf", bbTF.Bytes(), 0777)
	//if err != nil {
	//	return err
	//}
	//err = ioutil.WriteFile("imports.sh", bbSH.Bytes(), 0777)
	//if err != nil {
	//	return err
	//}
	return nil
}
