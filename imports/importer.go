package imports

import (
	"fmt"

	"github.com/google/uuid"

	"stash.sigma.sbrf.ru/sddevops/terraform-provider-di/client"
	"stash.sigma.sbrf.ru/sddevops/terraform-provider-di/models"
)

const importStr = "terraform import %s.%s %s\n"

type Importer struct {
	Api       *client.Api
	Project   *models.Project
	Domain    *models.Domain
	Group     *models.Group
	StandType *models.StandType
	AS        *models.AS
	VMs       []*models.Server
	Nginx     []*models.Server
	Sowa      []*models.Server
	Openshift []*models.Server
	Postgres  []*models.Server
	Kafka     []*models.Cluster
	Ignite    []*models.Cluster
	Files     map[string]*FileWriter
	VaultKV   string
	Bucket    string
}

func (o *Importer) Import(projectId string) error {

	id, err := uuid.Parse(projectId)
	if err != nil {
		return fmt.Errorf("can't parse uuid, %s\n", err)
	}
	project := models.Project{Id: id}
	o.Project = &project

	// network client
	o.Api = client.NewApi()
	models.Api = o.Api
	models.Api.Debug = false

	// Project
	err = o.importProject()
	if err != nil {
		return err
	}

	// Domain
	o.Domain = &models.Domain{
		Id: o.Project.DomainId,
	}
	err = o.importDomain()
	if err != nil {
		return err
	}

	// Group
	o.Group = &models.Group{Id: o.Project.GroupId, DomainId: o.Domain.Id}
	err = o.importGroup()
	if err != nil {
		return err
	}

	// StandType
	o.StandType = &models.StandType{Id: o.Project.StandTypeId}
	err = o.importStandType()
	if err != nil {
		return err
	}

	// AS
	o.AS = &models.AS{Code: o.Project.AppSystemsCi, DomainId: o.Domain.Id}
	err = o.importAS()
	if err != nil {
		return err
	}
	o.AS.ResDomainIdUUID = o.Domain.Id.String()
	o.AS.ResDomainId = fmt.Sprintf(
		"data.%s.%s.id",
		o.Domain.GetResType(),
		o.Domain.GetResName(),
	)

	// files
	projectDir := fmt.Sprintf(
		"%s/%s/%s/%s",
		o.Domain.ResName, o.Group.ResName, o.StandType.ResName, o.Project.ResName,
	)
	o.Files = make(map[string]*FileWriter)
	o.Files["data"] = NewFileWriter(fmt.Sprintf("%s/project/data.tf", projectDir))
	o.Files["project"] = NewFileWriter(fmt.Sprintf("%s/project/project.tf", projectDir))
	o.Files["projectMain"] = NewFileWriter(fmt.Sprintf("%s/project/main.tf", projectDir))
	o.Files["projectOutputs"] = NewFileWriter(fmt.Sprintf("%s/project/outputs.tf", projectDir))
	o.Files["projectImports"] = NewFileWriter(fmt.Sprintf("%s/project/imports.sh", projectDir))
	o.Files["projectImports"].Buffer.WriteString("#!/usr/bin/env bash\n\n")

	o.Files["appMain"] = NewFileWriter(fmt.Sprintf("%s/app/main.tf", projectDir))
	o.Files["appImports"] = NewFileWriter(fmt.Sprintf("%s/app/imports.sh", projectDir))
	o.Files["appImports"].Buffer.WriteString("#!/usr/bin/env bash\n\n")
	o.Files["nginx"] = NewFileWriter(fmt.Sprintf("%s/app/nginx.tf", projectDir))
	o.Files["sowa"] = NewFileWriter(fmt.Sprintf("%s/app/sowa.tf", projectDir))
	o.Files["kafka"] = NewFileWriter(fmt.Sprintf("%s/app/kafka.tf", projectDir))

	o.Files["computeMain"] = NewFileWriter(fmt.Sprintf("%s/compute/main.tf", projectDir))
	o.Files["computeImports"] = NewFileWriter(fmt.Sprintf("%s/compute/imports.sh", projectDir))
	o.Files["computeImports"].Buffer.WriteString("#!/usr/bin/env bash\n\n")
	o.Files["vm"] = NewFileWriter(fmt.Sprintf("%s/compute/vm.tf", projectDir))
	o.Files["openshift"] = NewFileWriter(fmt.Sprintf("%s/compute/openshift.tf", projectDir))

	o.Files["dbMain"] = NewFileWriter(fmt.Sprintf("%s/db/main.tf", projectDir))
	o.Files["dbImports"] = NewFileWriter(fmt.Sprintf("%s/db/imports.sh", projectDir))
	o.Files["dbImports"].Buffer.WriteString("#!/usr/bin/env bash\n\n")
	o.Files["postgres"] = NewFileWriter(fmt.Sprintf("%s/db/postgres.tf", projectDir))
	o.Files["ignite"] = NewFileWriter(fmt.Sprintf("%s/db/ignite.tf", projectDir))
	o.Files["vaultSH"] = NewFileWriter(fmt.Sprintf("%s/db/vault.sh", projectDir))
	o.Files["vaultSH"].Buffer.WriteString("#!/usr/bin/env bash\n\n")
	o.Files["vaultTF"] = NewFileWriter(fmt.Sprintf("%s/db/vault.tf", projectDir))

	domainResource := o.Domain.NewObj()
	domainResource = o.Domain
	data := models.ToHCLData(domainResource)
	_, err = o.Files["data"].Buffer.Write(data)
	if err != nil {
		return err
	}
	data = models.ToHCLOutput(domainResource)
	_, err = o.Files["projectOutputs"].Buffer.Write(data)
	if err != nil {
		return err
	}

	groupResource := o.Group.NewObj()
	groupResource = o.Group
	data = models.ToHCLData(groupResource)
	_, err = o.Files["data"].Buffer.Write(data)
	if err != nil {
		return err
	}
	data = models.ToHCLOutput(groupResource)
	_, err = o.Files["projectOutputs"].Buffer.Write(data)
	if err != nil {
		return err
	}

	standTypeResource := o.StandType.NewObj()
	standTypeResource = o.StandType
	data = models.ToHCLData(standTypeResource)
	_, err = o.Files["data"].Buffer.Write(data)
	if err != nil {
		return err
	}
	data = models.ToHCLOutput(standTypeResource)
	_, err = o.Files["projectOutputs"].Buffer.Write(data)
	if err != nil {
		return err
	}

	asResource := o.AS.NewObj()
	asResource = o.AS
	data = models.ToHCLData(asResource)
	_, err = o.Files["data"].Buffer.Write(data)
	if err != nil {
		return err
	}
	data = models.ToHCLOutput(asResource)
	_, err = o.Files["projectOutputs"].Buffer.Write(data)
	if err != nil {
		return err
	}

	// write project data
	err = o.writeProject()
	if err != nil {
		return err
	}

	// vm
	// models.Api.Debug = true
	err = o.importVMs()
	if err != nil {
		return err
	}

	// openshift
	err = o.importOpenshift()
	if err != nil {
		return err
	}

	// nginx
	err = o.importNginx()
	if err != nil {
		return err
	}

	// sowa
	err = o.importSowa()
	if err != nil {
		return err
	}

	// postgres
	err = o.importPostgres()
	if err != nil {
		return err
	}

	// kafka
	err = o.importKafka()
	if err != nil {
		return err
	}

	// ignite
	err = o.importIgnite()
	if err != nil {
		return err
	}

	err = o.writeMain()
	if err != nil {
		return err
	}
	for _, v := range o.Files {
		v.Close()
	}

	return nil
}

func (o *Importer) writeMain() error {
	bucket := "di-terraformtest"
	s3KeyProject := fmt.Sprintf(
		"%s/%s/%s/%s",
		o.Domain.ResName,
		o.Group.ResName,
		o.StandType.ResName,
		o.Project.ResName,
	)
	mainTF := fmt.Sprintf(
		mainTpl,
		diProviderTpl,
		"",
		fmt.Sprintf(s3backendTpl, bucket, s3KeyProject, "project"),
		"",
	)
	appMainTF := fmt.Sprintf(
		mainTpl,
		diProviderTpl,
		"",
		fmt.Sprintf(s3backendTpl, bucket, s3KeyProject, "app"),
		fmt.Sprintf(remoteStateTpl, bucket, s3KeyProject),
	)
	computeMainTF := fmt.Sprintf(
		mainTpl,
		diProviderTpl,
		"",
		fmt.Sprintf(s3backendTpl, bucket, s3KeyProject, "compute"),
		fmt.Sprintf(remoteStateTpl, bucket, s3KeyProject),
	)
	dbMainTF := fmt.Sprintf(
		mainTpl,
		diProviderTpl,
		vaultProviderTpl,
		fmt.Sprintf(s3backendTpl, bucket, s3KeyProject, "db"),
		fmt.Sprintf(remoteStateTpl, bucket, s3KeyProject),
	)

	_, err := o.Files["projectMain"].Buffer.WriteString(mainTF)
	if err != nil {
		return err
	}
	_, err = o.Files["appMain"].Buffer.WriteString(appMainTF)
	if err != nil {
		return err
	}
	_, err = o.Files["computeMain"].Buffer.WriteString(computeMainTF)
	if err != nil {
		return err
	}
	_, err = o.Files["dbMain"].Buffer.WriteString(dbMainTF)
	if err != nil {
		return err
	}
	return nil
}
