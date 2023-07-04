package models

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gitlab.gos-tech.xyz/pid/iac/terraform-provider-sberinfra/utils"
)

type Group struct {
	Id              uuid.UUID `json:"id"`
	Name            string    `json:"name" hcl:"name"`
	DomainId        uuid.UUID `json:"domain_id"`
	IsProm          bool      `json:"is_prom"`
	ResId           string    `json:"-" hcl:"id"`
	ResType         string    `json:"-" hcl:"type,label"`
	ResName         string    `json:"-" hcl:"name,label"`
	ResDomainIdUUID string    `json:"-" hcl:"domain_id_uuid"`
	ResDomainId     string    `json:"-" hcl:"domain_id"`
	ResDomainName   string    `json:"-"`
	ResOutputName   string    `json:"-"`
	ResOutputValue  string    `json:"-"`
}

func (o *Group) NewObj() DIDataResource {
	return &Group{}
}

func (o *Group) ReadTF(res *schema.ResourceData) {
	domainId := res.Get("domain_id")
	if domainId != "" {
		o.DomainId = uuid.MustParse(domainId.(string))
	}
	o.Name = res.Get("name").(string)
}

func (o *Group) WriteTF(res *schema.ResourceData) {
	res.SetId(o.Id.String())
	err := res.Set("name", o.Name)
	err = res.Set("domain_id", o.DomainId.String())
	err = res.Set("is_prom", o.IsProm)
	if err != nil {
		log.Println(err)
	}
}

func (o *Group) deserializeList(responseBytes []byte) error {
	data := make(map[string]interface{})
	err := json.Unmarshal(responseBytes, &data)
	if err != nil {
		return err
	}
	for _, val := range data["groups"].([]interface{}) {
		v := val.(map[string]interface{})
		if v["name"] == o.Name {
			o.Id = uuid.MustParse(v["id"].(string))
			if v["is_prom"] != nil {
				o.IsProm = v["is_prom"].(bool)
			}
			return nil
		}
	}
	return nil
}

func (o *Group) DeserializeOne(responseBytes []byte) error {
	data := make(map[string]map[string]interface{})
	err := json.Unmarshal(responseBytes, &data)
	if err != nil {
		return err
	}
	resource := data["group"]
	o.Id = uuid.MustParse(resource["id"].(string))
	o.Name = resource["name"].(string)
	return nil
}

func (o *Group) Deserialize(responseBytes []byte) error {
	return o.deserializeList(responseBytes)
}

func (o *Group) ReadDI() ([]byte, error) {
	return Api.NewRequestRead(fmt.Sprintf("groups?domain_id=%s", o.DomainId.String()))
}

func (o *Group) GetId() string {
	return o.Id.String()
}

func (o *Group) ReadAll() ([]byte, error) {
	return Api.NewRequestRead("groups")
}

func (o *Group) DeserializeAll(responseBytes []byte) ([]DIDataResource, error) {
	m := make(map[string][]*Group)
	err := json.Unmarshal(responseBytes, &m)
	if err != nil {
		return nil, err
	}

	m2 := make([]DIDataResource, len(m["groups"]))
	for k, v := range m["groups"] {
		m2[k] = v
	}
	return m2, nil
}

func (o *Group) GetDomainId() uuid.UUID {
	return o.DomainId
}

func (o *Group) GetResType() string {
	return "si_group"
}

func (o *Group) GetResName() string {
	return o.ResName
}

func (o *Group) GetOutput() (string, string) {
	return o.ResOutputName, o.ResOutputValue
}

func (o *Group) SetDomainName(domain_name string) {
	o.ResDomainId = domain_name
}

func (o *Group) SetResFields() {
	o.ResId = o.GetId()
	o.ResType = o.GetResType()
	o.ResName = utils.Reformat(o.Name)
	o.ResDomainIdUUID = o.DomainId.String()
	o.ResOutputName = fmt.Sprintf(
		"%s_id",
		o.GetResType(),
	)
	o.ResOutputValue = fmt.Sprintf(
		"data.%s.%s.id",
		o.GetResType(),
		o.GetResName(),
	)
}
