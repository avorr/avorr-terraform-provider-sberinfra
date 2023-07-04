package views

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/avorr/terraform-provider-sberinfra/models"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func LookupTag(TagName string) uuid.UUID {
	obj := models.Tag{}
	allTags, err := obj.ReadDI()
	if err != nil {
		panic(err)
	}

	dat := make(map[string][]*Tag)
	if err := json.Unmarshal(allTags, &dat); err != nil {
		panic(err)
	}

	for _, i := range dat["tags"] {
		if i.Name == TagName {
			fmt.Println(i.Name)
			return i.Id
		}
	}
	return uuid.Nil
}

type FullTag struct {
	Tag Tag `json:"tag"`
}

type Tag struct {
	Id      uuid.UUID `json:"id"`
	Name    string    `json:"tag_name" hcl:"name"`
	ResId   string    `json:"-" hcl:"id"`
	ResType string    `json:"-" hcl:"type,label"`
	ResName string    `json:"-" hcl:"name,label"`
}

func TagCreate(ctx context.Context, res *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	obj := models.Tag{}
	obj.ReadTF(res)

	requestBytes, err := obj.Serialize()

	var requestData Tag

	if err := json.Unmarshal(requestBytes, &requestData); err != nil {
		panic(err)
	}
	lookupTagId := LookupTag(requestData.Name)

	if err != nil {
		return diag.FromErr(err)
	}

	if lookupTagId == uuid.Nil {
		responseBytes, err := obj.CreateDI(requestBytes)
		if err != nil {
			return diag.FromErr(err)
		}
		err = obj.Deserialize(responseBytes)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		responseBytes, err := json.Marshal(FullTag{Tag: Tag{Name: requestData.Name, Id: lookupTagId}})
		err = obj.Deserialize(responseBytes)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	obj.WriteTF(res)
	return diags
}

func TagRead(ctx context.Context, res *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	obj := models.Tag{}
	obj.ReadTF(res)

	responseBytes, err := obj.ReadDI()
	if err != nil {
		return diag.FromErr(err)
	}
	err = obj.DeserializeAll(responseBytes)
	if err != nil {
		return diag.FromErr(err)
	}
	// obj.WriteTF(res)
	return diags
}

func TagDelete(ctx context.Context, res *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	obj := models.Tag{}
	obj.ReadTF(res)
	//err := obj.DeleteDI()
	//if err != nil {
	//	return diag.FromErr(err)
	//}
	res.SetId("")
	return diags
}
