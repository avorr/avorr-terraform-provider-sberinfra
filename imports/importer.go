package imports

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/google/uuid"
	"stash.sigma.sbrf.ru/sddevops/terraform-provider-di/client"
	"stash.sigma.sbrf.ru/sddevops/terraform-provider-di/models"
)

type Importer struct {
	Api     *client.Api
	Project *models.Project
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
	project := &models.Project{Id: id}
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

	err = ioutil.WriteFile("imports.tf", bbTF.Bytes(), 0777)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("imports.sh", bbSH.Bytes(), 0777)
	if err != nil {
		return err
	}
	return nil
}
