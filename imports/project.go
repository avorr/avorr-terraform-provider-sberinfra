package imports

import (
	"fmt"
	"log"

	"stash.sigma.sbrf.ru/sddevops/terraform-provider-di/models"
	"stash.sigma.sbrf.ru/sddevops/terraform-provider-di/utils"
)

func (o *Importer) importProject() error {
	obj := o.Project
	responseBytes, err := obj.ReadDI()
	if err != nil {
		return err
	}
	err = obj.Deserialize(responseBytes)
	if err != nil {
		return err
	}
	obj.ResId = obj.Id.String()
	obj.ResType = "di_project"
	obj.ResName = utils.Reformat(obj.Name)
	// ResDomainId: objMap["domain_id"].(string),
	obj.ResGroupIdUUID = obj.GroupId.String()
	obj.ResStandTypeIdUUID = obj.StandTypeId.String()
	log.Printf("project %s -> %s\n", obj.Name, obj.ResName)
	return nil
}

func (o *Importer) writeProject() error {
	project := o.Project
	project.ResGroupId = fmt.Sprintf(
		"data.%s.%s.id",
		o.Group.GetResType(),
		o.Group.GetResName(),
	)
	project.ResStandTypeId = fmt.Sprintf(
		"data.%s.%s.id",
		o.StandType.GetResType(),
		o.StandType.GetResName(),
	)
	project.ResAsIdUUID = o.AS.Id.String()
	project.ResAsId = fmt.Sprintf(
		"data.%s.%s.code",
		o.AS.GetResType(),
		o.AS.GetResName(),
	)
	resourceProject := project.NewObj()
	resourceProject = project
	data := models.ToHCLResource(resourceProject)
	_, err := o.Files["project"].Buffer.Write(data)
	if err != nil {
		return err
	}
	data = project.ToHCLOutput()
	_, err = o.Files["projectOutputs"].Buffer.Write(data)
	if err != nil {
		return err
	}
	_, err = o.Files["projectImports"].Buffer.WriteString(fmt.Sprintf(
		importStr+"\n",
		project.ResType,
		project.ResName,
		project.Id.String(),
	))
	return nil
}
