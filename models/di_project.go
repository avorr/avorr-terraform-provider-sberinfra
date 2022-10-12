package models

//
//import (
//	"encoding/json"
//	"errors"
//	"fmt"
//
//	"github.com/hashicorp/hcl/v2/gohcl"
//	"github.com/hashicorp/hcl/v2/hclwrite"
//
//	"base.sw.sbc.space/pid/terraform-provider-si/utils"
//
//	"github.com/google/uuid"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
//)
//
//type Project struct {
//	Id                 uuid.UUID `json:"id"`
//	GroupId            uuid.UUID `json:"group_id"`
//	DomainId           uuid.UUID `json:"domain_id"`
//	StandTypeId        uuid.UUID `json:"stand_type_id"`
//	Name               string    `json:"name" hcl:"name"`
//	StandType          string    `json:"stand_type"`
//	State              string    `json:"state"`
//	Type               string    `json:"type"`
//	AppSystemsCi       string    `json:"app_systems_ci" hcl:"app_systems_ci"`
//	ResId              string    `json:"-"`
//	ResType            string    `json:"-" hcl:"type,label"`
//	ResName            string    `json:"-" hcl:"name,label"`
//	ResGroupIdUUID     string    `json:"-"`
//	ResGroupId         string    `json:"-" hcl:"group_id"`
//	ResAsIdUUID        string    `json:"-"`
//	ResAsId            string    `json:"-"`
//	ResStandTypeIdUUID string    `json:"-"`
//	ResStandTypeId     string    `json:"-" hcl:"stand_type_id"`
//}
//
//func (o *Project) GetType() string {
//	return "di_project"
//}
//
//func (o *Project) NewObj() DIResource {
//	return &Project{}
//}
//
//func (o *Project) ReadTF(res *schema.ResourceData) {
//	if res.Id() != "" {
//		o.Id = uuid.MustParse(res.Id())
//	}
//	groupId := res.Get("group_id")
//	if groupId != "" {
//		o.GroupId = uuid.MustParse(groupId.(string))
//	}
//	domainId := res.Get("domain_id")
//	if domainId != "" {
//		o.DomainId = uuid.MustParse(domainId.(string))
//	}
//	o.Name = res.Get("name").(string)
//	o.StandType = res.Get("stand_type").(string)
//	standTypeId := res.Get("stand_type_id")
//	if standTypeId != "" {
//		o.StandTypeId = uuid.MustParse(standTypeId.(string))
//	}
//	o.State = res.Get("state").(string)
//	o.Type = res.Get("type").(string)
//	o.AppSystemsCi = res.Get("app_systems_ci").(string)
//}
//
//func (o *Project) WriteTF(res *schema.ResourceData) {
//	res.SetId(o.Id.String())
//	res.Set("name", o.Name)
//	res.Set("stand_type_id", o.StandTypeId.String())
//	res.Set("group_id", o.GroupId.String())
//	res.Set("domain_id", o.DomainId.String())
//	res.Set("app_systems_ci", o.AppSystemsCi)
//	res.Set("stand_type", o.StandType)
//	res.Set("state", o.State)
//	res.Set("type", o.Type)
//}
//
//func (o *Project) Serialize() ([]byte, error) {
//	requestMap := map[string]map[string]interface{}{
//		"project": {
//			"name":           o.Name,
//			"group_id":       o.GroupId,
//			"stand_type_id":  o.StandTypeId.String(),
//			"app_systems_ci": o.AppSystemsCi,
//			// "stand_type":     o.StandType,
//			// "domain_id":      o.DomainId,
//			// "state":          o.State,
//			// "type":           o.Type,
//		},
//	}
//	requestBytes, err := json.Marshal(requestMap)
//	if err != nil {
//		return nil, err
//	}
//	return requestBytes, nil
//}
//
//func (o *Project) Deserialize(responseBytes []byte) error {
//	response := make(map[string]map[string]interface{})
//	err := json.Unmarshal(responseBytes, &response)
//	if err != nil {
//		return err
//	}
//	objMap, ok := response["project"]
//	if !ok {
//		return errors.New("no project in response")
//	}
//	o.Id = uuid.MustParse(objMap["id"].(string))
//	o.ResId = objMap["id"].(string)
//	o.DomainId = uuid.MustParse(objMap["domain_id"].(string))
//	o.GroupId = uuid.MustParse(objMap["group_id"].(string))
//	o.ResGroupId = objMap["group_id"].(string)
//	o.StandTypeId = uuid.MustParse(objMap["stand_type_id"].(string))
//	o.ResStandTypeId = objMap["stand_type_id"].(string)
//	o.StandType = objMap["stand_type"].(string)
//	o.Name = objMap["name"].(string)
//	o.Type = objMap["type"].(string)
//	o.State = objMap["state"].(string)
//	o.AppSystemsCi = objMap["app_systems_ci"].(string)
//	return nil
//}
//
//func (o *Project) ParseIdFromCreateResponse(data []byte) error {
//	response := make(map[string]map[string]interface{})
//	err := json.Unmarshal(data, &response)
//	if err != nil {
//		return err
//	}
//	objMap, ok := response["project"]
//	if !ok {
//		return errors.New("no project in response")
//	}
//	o.Id = uuid.MustParse(objMap["id"].(string))
//	return nil
//}
//
//func (o *Project) CreateDI(data []byte) ([]byte, error) {
//	return Api.NewRequestCreate("projects", data)
//}
//
//func (o *Project) ReadDI() ([]byte, error) {
//	return Api.NewRequestRead(fmt.Sprintf("projects/%s", o.Id))
//}
//
//func (o *Project) UpdateDI(data []byte) ([]byte, error) {
//	return Api.NewRequestUpdate(fmt.Sprintf("projects/%s", o.Id), data)
//}
//
//func (o *Project) DeleteDI() error {
//	return Api.NewRequestDelete(fmt.Sprintf("projects/%s", o.Id), nil)
//}
//
//func (o *Project) ReadAll() ([]byte, error) {
//	return Api.NewRequestRead("projects/")
//}
//
//func (o *Project) DeserializeAll(responseBytes []byte) ([]*Project, error) {
//	response := make(map[string]interface{})
//	err := json.Unmarshal(responseBytes, &response)
//	if err != nil {
//		return nil, err
//	}
//	objList := make([]*Project, 0)
//	objResNamesList := make([]string, 0)
//	counter := make(map[string][]*Project)
//	for _, v := range response["projects"].([]interface{}) {
//		objMap := v.(map[string]interface{})
//		obj := &Project{
//			Id:           uuid.MustParse(objMap["id"].(string)),
//			GroupId:      uuid.MustParse(objMap["group_id"].(string)),
//			DomainId:     uuid.MustParse(objMap["domain_id"].(string)),
//			StandTypeId:  uuid.MustParse(objMap["stand_type_id"].(string)),
//			Name:         objMap["name"].(string),
//			StandType:    objMap["stand_type"].(string),
//			State:        objMap["state"].(string),
//			Type:         objMap["type"].(string),
//			AppSystemsCi: objMap["app_systems_ci"].(string),
//			ResId:        objMap["id"].(string),
//			ResType:      "di_project",
//			ResName:      utils.Reformat(objMap["name"].(string)),
//			// ResDomainId: objMap["domain_id"].(string),
//			ResGroupIdUUID:     objMap["group_id"].(string),
//			ResStandTypeIdUUID: objMap["stand_type_id"].(string),
//			// ResStandTypeId:     objMap["stand_type_id"].(string),
//		}
//		objList = append(objList, obj)
//		objResNamesList = append(objResNamesList, obj.ResName)
//		counter[obj.ResName] = append(counter[obj.ResName], obj)
//	}
//	for _, arr := range counter {
//		if len(arr) > 1 {
//			var c int
//			for _, v := range arr {
//				c++
//				v.ResName = fmt.Sprintf("%s-%d", v.ResName, c)
//			}
//		}
//	}
//	return objList, nil
//}
//
//func (o *Project) OnSerialize(map[string]interface{}, *Server) map[string]interface{} {
//	return nil
//}
//func (o *Project) OnDeserialize(map[string]interface{}, *Server) {}
//func (o *Project) Urls(string) string {
//	return ""
//}
//func (o *Project) OnReadTF(*schema.ResourceData, *Server)  {}
//func (o *Project) OnWriteTF(*schema.ResourceData, *Server) {}
//
//func (o *Project) ToHCLOutput() []byte {
//	dataRoot := &HCLOutputRoot{
//		Resources: &HCLOutput{
//			ResName: fmt.Sprintf(
//				"%s_id",
//				o.ResType,
//			),
//			Value: fmt.Sprintf(
//				"%s.%s.id",
//				o.ResType,
//				o.ResName,
//			),
//		},
//	}
//	f := hclwrite.NewEmptyFile()
//	gohcl.EncodeIntoBody(dataRoot, f.Body())
//	return utils.Regexp(f.Bytes())
//}
//
//func (o *Project) HostVars(server *Server) map[string]interface{} {
//	return nil
//}
//
//func (o *Project) GetGroup() string {
//	return ""
//}
//
//func (o *Project) ToHCL(server *Server) ([]byte, error) {
//	o.ResType = o.GetType()
//	o.ResName = utils.Reformat(o.Name)
//	type HCLServerRoot struct {
//		Resources *Project `hcl:"resource,block"`
//	}
//	root := &HCLServerRoot{Resources: o}
//	f := hclwrite.NewEmptyFile()
//	gohcl.EncodeIntoBody(root, f.Body())
//	// return utils.Regexp(f.Bytes())
//	return f.Bytes(), nil
//}
//
//func (o *Project) HCLAppParams() *HCLAppParams {
//	return nil
//}
//
//func (o *Project) HCLVolumes() []*HCLVolume {
//	return nil
//}
