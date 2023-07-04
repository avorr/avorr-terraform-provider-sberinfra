package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

type Tag struct {
	Id      uuid.UUID `json:"id"`
	Name    string    `json:"tag_name" hcl:"name"`
	ResId   string    `json:"-" hcl:"id"`
	ResType string    `json:"-" hcl:"type,label"`
	ResName string    `json:"-" hcl:"name,label"`
}

func (o *Tag) GetType() string {
	return "si_tag"
}

func (o *Tag) ReadTF(res *schema.ResourceData) {
	if res.Id() != "" {
		o.Id = uuid.MustParse(res.Id())
	}
	o.Name = res.Get("name").(string)
}

func (o *Tag) WriteTF(res *schema.ResourceData) {
	res.SetId(o.Id.String())
	err := res.Set("name", o.Name)
	if err != nil {
		log.Println(err)
	}
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
	return Api.NewRequestRead("dict/tags")
}

func (o *Tag) DeleteDI() error {
	return Api.NewRequestDelete(fmt.Sprintf("dict/tags?uuid=%s", o.Id), nil, 204)
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
