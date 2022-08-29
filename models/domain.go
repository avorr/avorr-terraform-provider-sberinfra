package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"stash.sigma.sbrf.ru/sddevops/terraform-provider-di/utils"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Domain struct {
	Id             uuid.UUID `json:"id"`
	Name           string    `json:"name" hcl:"name"`
	ResId          string    `json:"-" hcl:"id"`
	ResType        string    `json:"-" hcl:"type,label"`
	ResName        string    `json:"-" hcl:"name,label"`
	ResOutputName  string    `json:"-"`
	ResOutputValue string    `json:"-"`
	// PortalId      string             `json:"portal_id"`
	// SapId         string             `json:"sap_id"`
	// Type          string             `json:"type"`
	// BusinessBlock string             `json:"business_block"`
}

func (o *Domain) NewObj() DIDataResource {
	return &Domain{ResType: "di_domain"}
}

func (o *Domain) ReadTF(res *schema.ResourceData) {
	if res.Id() != "" {
		o.Id = uuid.MustParse(res.Id())
	}
	o.Name = res.Get("name").(string)
}

func (o *Domain) WriteTF(res *schema.ResourceData) {
	res.SetId(o.Id.String())
	// resource.Set(k, fmt.Sprintf("%v", v))
}

//func (o *Domain) Deserialize(responseBytes []byte) error {
//	data := make(map[string]interface{})
//	err := json.Unmarshal(responseBytes, &data)
//	if err != nil {
//		return err
//	}
//	domains := data["domains"].([]interface{})
//	if len(domains) < 1 {
//		return errors.New("no domain in response")
//	}
//	domain := domains[0].(map[string]interface{})
//	o.Id = uuid.MustParse(domain["id"].(string))
//	o.Name = domain["name"].(string)
//	return nil
//}

func (o *Domain) Deserialize(responseBytes []byte) error {
	data := make(map[string]interface{})
	err := json.Unmarshal(responseBytes, &data)
	if err != nil {
		return err
	}
	domains := data["domains"].([]interface{})
	if len(domains) < 1 {
		return errors.New("no domain in response")
	}

	domainsBytes, err := json.Marshal(domains)
	if err != nil {
		return err
	}

	m := make([]*Domain, 0)
	err = json.Unmarshal(domainsBytes, &m)
	if err != nil {
		return err
	}

	for _, v := range m {
		if v.Name == o.Name {
			o.Id = v.Id
		}
	}
	// domain := domains[0].(map[string]interface{})
	// o.Id = uuid.MustParse(domain["id"].(string))
	// o.Name = domain["name"].(string)
	return nil
}

func (o *Domain) ReadDI() ([]byte, error) {
	//return Api.NewRequestRead(fmt.Sprintf("domains?searchstring=%s", url.QueryEscape(o.Name)))
	return Api.NewRequestRead(fmt.Sprintf("domains?filter[name]=%s", url.QueryEscape(o.Name)))
}

func (o *Domain) GetId() string {
	return o.Id.String()
}

func (o *Domain) ReadAll() ([]byte, error) {
	return Api.NewRequestRead("domains")
}

func (o *Domain) DeserializeAll(responseBytes []byte) ([]DIDataResource, error) {
	data := make(map[string]interface{})
	err := json.Unmarshal(responseBytes, &data)
	if err != nil {
		return nil, err
	}
	domains := data["domains"].([]interface{})
	if len(domains) < 1 {
		return nil, errors.New("no domain in response")
	}
	objBytes, err := json.Marshal(domains)
	if err != nil {
		return nil, err
	}
	m := make([]*Domain, 0)
	err = json.Unmarshal(objBytes, &m)
	if err != nil {
		return nil, err
	}

	m2 := make([]DIDataResource, len(m))
	for k, v := range m {
		m2[k] = v
	}
	return m2, nil
}

func (o *Domain) GetDomainId() uuid.UUID {
	return o.Id
}

func (o *Domain) GetResType() string {
	return "di_domain"
}

func (o *Domain) GetResName() string {
	return o.ResName
}

func (o *Domain) GetOutput() (string, string) {
	return o.ResOutputName, o.ResOutputValue
}

func (o *Domain) SetResFields() {
	o.ResId = o.GetId()
	o.ResType = o.GetResType()
	o.ResName = utils.Reformat(o.Name)
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
