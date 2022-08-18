package models

//
//import (
//	"encoding/json"
//	"fmt"
//
//	"stash.sigma.sbrf.ru/sddevops/terraform-provider-di/utils"
//
//	"github.com/google/uuid"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
//)
//
//type StandType struct {
//	Id uuid.UUID `json:"id"`
//	// GroupId        uuid.UUID `json:"group_id"`
//	Name           string `json:"name" hcl:"name"`
//	NameShort      string `json:"name_short"`
//	Code           string `json:"code"`
//	HpsmStandType  string `json:"hpsm_stand_type"`
//	IsDisabled     bool   `json:"is_disabled"`
//	ResId          string `json:"-" hcl:"id"`
//	ResType        string `json:"-" hcl:"type,label"`
//	ResName        string `json:"-" hcl:"name,label"`
//	ResDomainId    string `json:"-"`
//	ResDomainName  string `json:"-"`
//	ResOutputName  string `json:"-"`
//	ResOutputValue string `json:"-"`
//	// ResGroupName   string    `json:"-" hcl:"group_id"`
//	// ResGroupIdUUID string    `json:"-" hcl:"group_id_uuid"`
//}
//
//func (o *StandType) NewObj() DIDataResource {
//	return &StandType{}
//}
//
//func (o *StandType) ReadTF(res *schema.ResourceData) {
//	if res.Id() != "" {
//		o.Id = uuid.MustParse(res.Id())
//	}
//	// if res.Get("group_id").(string) != "" {
//	// 	o.GroupId = uuid.MustParse(res.Get("group_id").(string))
//	// }
//	o.Name = res.Get("name").(string)
//	o.NameShort = res.Get("name").(string)
//	o.Code = res.Get("code").(string)
//	o.HpsmStandType = res.Get("hpsm_stand_type").(string)
//	o.IsDisabled = res.Get("is_disabled").(bool)
//}
//
//func (o *StandType) WriteTF(res *schema.ResourceData) {
//	res.SetId(o.Id.String())
//	res.Set("name", o.Name)
//	res.Set("name_short", o.NameShort)
//	res.Set("code", o.Code)
//	res.Set("hpsm_stand_type", o.HpsmStandType)
//	res.Set("is_disabled", o.IsDisabled)
//}
//
//func (o *StandType) DeserializeFromGroup(responseBytes []byte) error {
//	data := make(map[string]map[string]interface{})
//	err := json.Unmarshal(responseBytes, &data)
//	if err != nil {
//		return err
//	}
//	group := data["group"]
//	standTypes := group["stand_types"].([]interface{})
//	for _, v := range standTypes {
//		val := v.(map[string]interface{})
//		if val["name"] == o.Name {
//			o.Id = uuid.MustParse(val["uuid"].(string))
//			o.NameShort = val["name_short"].(string)
//			o.Code = val["code"].(string)
//			o.HpsmStandType = val["hpsm_stand_type"].(string)
//			o.IsDisabled = val["is_disabled"].(bool)
//			return nil
//		}
//	}
//	return fmt.Errorf("%s \"%s\" not found", "stand_type", o.Name)
//}
//
//func (o *StandType) Deserialize(responseBytes []byte) error {
//	// return o.DeserializeFromGroup(responseBytes)
//	data := make(map[string][]*StandType)
//	err := json.Unmarshal(responseBytes, &data)
//	if err != nil {
//		return err
//	}
//	for _, v := range data["stand_types"] {
//		if o.Name == v.Name || o.Id == v.Id {
//			o.Id = v.Id
//			o.NameShort = v.NameShort
//			o.Code = v.Code
//			o.HpsmStandType = v.HpsmStandType
//			o.IsDisabled = v.IsDisabled
//			o.Name = v.Name
//			return nil
//		}
//	}
//	return fmt.Errorf("%s \"%s\" not found", "stand_type", o.Name)
//}
//
//func (o *StandType) ReadDI() ([]byte, error) {
//	// return Api.NewRequestRead(fmt.Sprintf("groups/%s", o.GroupId))
//	return Api.NewRequestRead("dict/stand_types")
//}
//
//func (o *StandType) GetId() string {
//	return o.Id.String()
//}
//
//func (o *StandType) ReadAll() ([]byte, error) {
//	// return Api.NewRequestRead("groups")
//	return o.ReadDI()
//}
//
//func (o *StandType) DeserializeAll(responseBytes []byte) ([]DIDataResource, error) {
//	data := make(map[string][]*StandType)
//	err := json.Unmarshal(responseBytes, &data)
//	if err != nil {
//		return nil, err
//	}
//
//	m2 := make([]DIDataResource, len(data["stand_types"]))
//	for k, v := range data["stand_types"] {
//		v.ResId = v.GetId()
//		v.ResType = "di_stand_type"
//		v.ResName = utils.Reformat(v.Name)
//		// stand_type.ResName = stand_type.Code
//		m2[k] = v
//	}
//	return m2, nil
//}
//
//func (o *StandType) GetDomainId() uuid.UUID {
//	return uuid.UUID{}
//}
//
//func (o *StandType) GetResType() string {
//	return "di_stand_type"
//}
//
//func (o *StandType) GetResName() string {
//	return o.ResName
//}
//
//func (o *StandType) GetOutput() (string, string) {
//	return o.ResOutputName, o.ResOutputValue
//}
//
//func (o *StandType) SetResFields() {
//	o.ResType = "di_stand_type"
//	// o.ResName = utils.Reformat(o.Name)
//	o.ResOutputName = fmt.Sprintf(
//		"%s_id",
//		o.GetResType(),
//	)
//	o.ResOutputValue = fmt.Sprintf(
//		"data.%s.%s.id",
//		o.GetResType(),
//		o.GetResName(),
//	)
//}
