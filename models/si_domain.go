package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"

	"stash.sigma.sbrf.ru/sddevops/terraform-provider-di/utils"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SIDomain struct {
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

func (o *SIDomain) NewObj() DIDataResource {
	return &SIDomain{ResType: "di_si-domain"}
}

func (o *SIDomain) ReadTF(res *schema.ResourceData) {
	if res.Id() != "" {
		o.Id = uuid.MustParse(res.Id())
	}
	o.Name = res.Get("name").(string)
}

func (o *SIDomain) WriteTF(res *schema.ResourceData) {
	res.SetId(o.Id.String())
	// resource.Set(k, fmt.Sprintf("%v", v))
}

func (o *SIDomain) Deserialize(responseBytes []byte) error {
	data := make(map[string]interface{})
	err := json.Unmarshal(responseBytes, &data)
	if err != nil {
		return err
	}
	domains := data["domains"].([]interface{})
	if len(domains) < 1 {
		return errors.New("no domain in response")
	}
	domain := domains[0].(map[string]interface{})
	o.Id = uuid.MustParse(domain["id"].(string))
	o.Name = domain["name"].(string)
	return nil
}

func (o *SIDomain) ReadDI() ([]byte, error) {
	log.Println("READDI_FUNC")
	//return Api.NewRequestRead(fmt.Sprintf("domains?searchstring=%s", url.QueryEscape(o.Name),
	return Api.NewRequestRead(fmt.Sprintf("domains?filter[name]=%s", url.QueryEscape(o.Name))) //	https://portal.pd23.gtp/api/v1/domains?%D0%93%D0%BE%D1%81%D0%A2%D0%B5%D1%85

}

func (o *SIDomain) GetId() string {
	return o.Id.String()
}

func (o *SIDomain) ReadAll() ([]byte, error) {
	return Api.NewRequestRead("domains")
}

func (o *SIDomain) DeserializeAll(responseBytes []byte) ([]DIDataResource, error) {
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
	m := make([]*SIDomain, 0)
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

func (o *SIDomain) GetDomainId() uuid.UUID {
	return o.Id
}

func (o *SIDomain) GetResType() string {
	return "di_si-domain"
}

func (o *SIDomain) GetResName() string {
	return o.ResName
}

func (o *SIDomain) GetOutput() (string, string) {
	return o.ResOutputName, o.ResOutputValue
}

func (o *SIDomain) SetResFields() {
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
