package models

//
//import (
//	"encoding/json"
//	"fmt"
//
//	"github.com/google/uuid"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
//
//	"stash.sigma.sbrf.ru/sddevops/terraform-provider-di/utils"
//)
//
//type AS struct {
//	Id              uuid.UUID `json:"id"`
//	Name            string    `json:"name"`
//	Code            string    `json:"code" hcl:"code"`
//	Status          string    `json:"status"`
//	ServiceType     string    `json:"service_type"`
//	DomainId        uuid.UUID `json:"domain_id"`
//	ResId           string    `json:"-" hcl:"id"`
//	ResType         string    `json:"-" hcl:"type,label"`
//	ResName         string    `json:"-" hcl:"name,label"`
//	ResDomainIdUUID string    `json:"-" hcl:"domain_id_uuid"`
//	ResDomainId     string    `json:"-" hcl:"domain_id"`
//	ResOutputName   string    `json:"-"`
//	ResOutputValue  string    `json:"-"`
//}
//
//func (o *AS) NewObj() DIDataResource {
//	return &AS{}
//}
//
//func (o *AS) ReadTF(res *schema.ResourceData) {
//	if res.Id() != "" {
//		o.Id = uuid.MustParse(res.Id())
//	}
//	o.Name = res.Get("name").(string)
//	o.Code = res.Get("code").(string)
//	o.Status = res.Get("status").(string)
//	o.ServiceType = res.Get("service_type").(string)
//	o.DomainId = uuid.MustParse(res.Get("domain_id").(string))
//}
//
//func (o *AS) WriteTF(res *schema.ResourceData) {
//	res.SetId(o.Id.String())
//	res.Set("name", o.Name)
//	res.Set("code", o.Code)
//	res.Set("status", o.Status)
//	res.Set("service_type", o.ServiceType)
//	res.Set("domain_id", o.DomainId.String())
//}
//
//func (o *AS) Deserialize(responseBytes []byte) error {
//	data := make(map[string][]map[string]interface{})
//	err := json.Unmarshal(responseBytes, &data)
//	if err != nil {
//		return err
//	}
//	for _, v := range data["app_systems"] {
//		if v["code"] == o.Code && v["domain_id"] == o.DomainId.String() {
//			// o.Limit = v["limit"].(float64)
//			o.Id = uuid.MustParse(v["id"].(string))
//			o.Name = v["name"].(string)
//			o.Code = v["code"].(string)
//			o.Status = v["status"].(string)
//			serviceType := v["service_type"]
//			if serviceType == nil {
//				o.ServiceType = ""
//			} else {
//				o.ServiceType = serviceType.(string)
//			}
//			// o.ServiceType = v["service_type"].(string)
//			// o.DomainId = uuid.MustParse(v["domain_id"].(string))
//		}
//	}
//	return nil
//}
//
//func (o *AS) ReadDI() ([]byte, error) {
//	return Api.NewRequestRead(
//		fmt.Sprintf(
//			"dict/app_systems?domain_id=%s",
//			o.DomainId,
//		),
//	)
//}
//
//func (o *AS) GetId() string {
//	return o.Id.String()
//}
//
//func (o *AS) ReadAll() ([]byte, error) {
//	return nil, nil
//}
//
//func (o *AS) DeserializeAll(responseBytes []byte) ([]DIDataResource, error) {
//	data := make(map[string][]*AS)
//	err := json.Unmarshal(responseBytes, &data)
//	if err != nil {
//		return nil, err
//	}
//	m2 := make([]DIDataResource, len(data["app_systems"]))
//	for k, v := range data["app_systems"] {
//		v.ResId = v.GetId()
//		v.ResDomainIdUUID = v.DomainId.String()
//		m2[k] = v
//	}
//	return m2, nil
//}
//
//func (o *AS) GetDomainId() uuid.UUID {
//	return o.DomainId
//}
//
//func (o *AS) GetResType() string {
//	return "di_as"
//}
//
//func (o *AS) GetResName() string {
//	return o.ResName
//}
//
//func (o *AS) GetOutput() (string, string) {
//	return o.ResOutputName, o.ResOutputValue
//}
//
//func (o *AS) SetResFields() {
//	o.ResType = "di_as"
//	o.ResName = utils.Reformat(o.Name)
//	o.ResOutputName = fmt.Sprintf(
//		"%s_code",
//		o.GetResType(),
//	)
//	o.ResOutputValue = fmt.Sprintf(
//		"data.%s.%s.code",
//		o.GetResType(),
//		o.GetResName(),
//	)
//}
//
//func (o *AS) SetDomainName(domain_name string) {
//	o.ResDomainId = domain_name
//}
