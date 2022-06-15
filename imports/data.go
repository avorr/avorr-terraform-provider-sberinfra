package imports

import (
	"fmt"
	"log"

	"stash.sigma.sbrf.ru/sddevops/terraform-provider-di/utils"
)

func (o *Importer) importDomain() error {
	obj := o.Domain
	responseBytes, err := obj.ReadDI()
	if err != nil {
		return err
	}
	err = obj.Deserialize(responseBytes)
	if err != nil {
		return err
	}
	obj.ResId = obj.Id.String()
	obj.ResType = "di_domain"
	obj.ResName = utils.Reformat(obj.Name)

	log.Printf("domain %s -> %s\n", obj.Name, obj.ResName)
	return nil
}

func (o *Importer) importGroup() error {
	obj := o.Group
	responseBytes, err := o.Api.NewRequestRead(fmt.Sprintf("groups/%s", obj.Id))
	if err != nil {
		return err
	}
	err = obj.DeserializeOne(responseBytes)
	if err != nil {
		return err
	}
	obj.ResId = obj.Id.String()
	obj.ResType = "di_group"
	obj.ResName = utils.Reformat(obj.Name)
	obj.ResDomainIdUUID = obj.DomainId.String()
	obj.SetDomainName(fmt.Sprintf(
		"data.%s.%s.id",
		o.Domain.GetResType(),
		o.Domain.GetResName(),
	))
	log.Printf("group %s -> %s\n", obj.Name, obj.ResName)
	return nil
}

func (o *Importer) importStandType() error {
	obj := o.StandType

	responseBytes, err := obj.ReadDI()
	if err != nil {
		log.Println(err)
	}
	err = obj.Deserialize(responseBytes)
	if err != nil {
		return err
	}
	obj.ResId = obj.Id.String()
	obj.ResName = utils.Reformat(obj.Name)
	obj.ResType = "di_stand_type"
	log.Printf("stand type %s -> %s\n", obj.Name, obj.ResName)
	return nil
}

func (o *Importer) importAS() error {
	obj := o.AS

	responseBytes, err := obj.ReadDI()
	if err != nil {
		log.Println(err)
	}
	err = obj.Deserialize(responseBytes)
	if err != nil {
		return err
	}
	obj.ResId = obj.Id.String()
	obj.ResName = utils.Reformat(obj.Name)
	obj.ResType = "di_as"
	log.Printf("app system %s -> %s\n", obj.Name, obj.ResName)
	return nil
}
