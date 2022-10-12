package models

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclwrite"

	"base.sw.sbc.space/pid/terraform-provider-si/utils"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Tag struct {
	Id      uuid.UUID `json:"id"`
	Name    string    `json:"tag_name" hcl:"name"`
	ResId   string    `json:"-" hcl:"id"`
	ResType string    `json:"-" hcl:"type,label"`
	ResName string    `json:"-" hcl:"name,label"`
}

func (o *Tag) GetType() string {
	return "di_tag"
}

func (o *Tag) NewObj() DIResource {
	return &Tag{}
}

func (o *Tag) ReadTF(res *schema.ResourceData) {
	if res.Id() != "" {
		o.Id = uuid.MustParse(res.Id())
	}
	o.Name = res.Get("name").(string)
}

func (o *Tag) WriteTF(res *schema.ResourceData) {
	res.SetId(o.Id.String())
	res.Set("name", o.Name)
}

func (o *Tag) Serialize() ([]byte, error) {
	requestMap := map[string]interface{}{
		"tag_name": o.Name,
	}
	requestBytes, err := json.Marshal(requestMap)
	if err != nil {
		return nil, err
	}
	return requestBytes, nil
}

func (o *Tag) Deserialize(responseBytes []byte) error {
	response := make(map[string]*Tag)
	err := json.Unmarshal(responseBytes, &response)
	if err != nil {
		return err
	}
	o.Id = response["tag"].Id
	o.Name = response["tag"].Name
	return nil
}

func (o *Tag) ParseIdFromCreateResponse(data []byte) error {
	response := make(map[string]map[string]interface{})
	err := json.Unmarshal(data, &response)
	if err != nil {
		return err
	}
	objMap, ok := response["tag"]
	if !ok {
		return errors.New("no tag in response")
	}
	o.Id = uuid.MustParse(objMap["id"].(string))
	return nil
}

func (o *Tag) CreateDI(data []byte) ([]byte, error) {
	return Api.NewRequestCreate("dict/tags", data)
}

func (o *Tag) ReadDI() ([]byte, error) {
	// Api.Debug = false
	// response, err := Api.NewRequestRead("dict/tags")
	// Api.Debug = true
	// return response, err
	return Api.NewRequestRead("dict/tags")
}

func (o *Tag) DeleteDI() error {
	return Api.NewRequestDelete(fmt.Sprintf("dict/tags?uuid=%s", o.Id), nil)
}

func (o *Tag) ReadAll() ([]byte, error) {
	return Api.NewRequestRead("dict/tags")
}

func (o *Tag) DeserializeAll(responseBytes []byte) error {
	response := make(map[string][]*Tag)
	err := json.Unmarshal(responseBytes, &response)
	if err != nil {
		return err
	}
	tags, ok := response["tags"]
	if !ok {
		return errors.New("no tag in response")
	}

	found := false
	for _, v := range tags {
		if v.Id == o.Id && v.Name == o.Name {
			found = true
		}
	}
	if !found {
		return fmt.Errorf("no tag [%s]:%s in response list", o.Id, o.Name)
	}
	return nil
}

func (o *Tag) OnSerialize(map[string]interface{}, *Server) map[string]interface{} {
	return nil
}
func (o *Tag) OnDeserialize(map[string]interface{}, *Server) {}
func (o *Tag) Urls(string) string {
	return ""
}
func (o *Tag) OnReadTF(*schema.ResourceData, *Server)  {}
func (o *Tag) OnWriteTF(*schema.ResourceData, *Server) {}

func (o *Tag) ToHCLOutput() []byte {
	dataRoot := &HCLOutputRoot{
		Resources: &HCLOutput{
			ResName: fmt.Sprintf(
				"%s_id",
				o.ResType,
			),
			Value: fmt.Sprintf(
				"%s.%s.id",
				o.ResType,
				o.ResName,
			),
		},
	}
	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(dataRoot, f.Body())
	return utils.Regexp(f.Bytes())
}

func (o *Tag) HostVars(server *Server) map[string]interface{} {
	return nil
}

func (o *Tag) GetGroup() string {
	return ""
}

func (o *Tag) HCLAppParams() *HCLAppParams {
	return nil
}

func (o *Tag) HCLVolumes() []*HCLVolume {
	return nil
}
