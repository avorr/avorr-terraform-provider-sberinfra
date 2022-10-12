package views

//
//import (
//	"bytes"
//	"context"
//	"fmt"
//	"io/ioutil"
//
//	"github.com/google/uuid"
//	"base.sw.sbc.space/pid/terraform-provider-si/models"
//
//	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
//)
//
//func ProjectCreate(ctx context.Context, res *schema.ResourceData, m interface{}) diag.Diagnostics {
//	var diags diag.Diagnostics
//
//	obj := models.Project{}
//	obj.ReadTF(res)
//
//	requestBytes, err := obj.Serialize()
//	if err != nil {
//		return diag.FromErr(err)
//	}
//
//	responseBytes, err := obj.CreateDI(requestBytes)
//	if err != nil {
//		return diag.FromErr(err)
//	}
//
//	err = obj.ParseIdFromCreateResponse(responseBytes)
//	if err != nil {
//		return diag.FromErr(err)
//	}
//
//	responseBytes, err = obj.ReadDI()
//	if err != nil {
//		return diag.FromErr(err)
//	}
//
//	err = obj.Deserialize(responseBytes)
//	if err != nil {
//		return diag.FromErr(err)
//	}
//	obj.WriteTF(res)
//	return diags
//}
//
//func ProjectRead(ctx context.Context, res *schema.ResourceData, m interface{}) diag.Diagnostics {
//	var diags diag.Diagnostics
//
//	obj := models.Project{}
//	obj.ReadTF(res)
//
//	responseBytes, err := obj.ReadDI()
//	if err != nil {
//		return diag.FromErr(err)
//	}
//
//	err = obj.Deserialize(responseBytes)
//	if err != nil {
//		return diag.FromErr(err)
//	}
//
//	obj.WriteTF(res)
//	return diags
//}
//
//func ProjectUpdate(ctx context.Context, res *schema.ResourceData, m interface{}) diag.Diagnostics {
//	var diags diag.Diagnostics
//
//	obj := models.Project{}
//	obj.ReadTF(res)
//
//	requestBytes, err := obj.Serialize()
//	if err != nil {
//		return diag.FromErr(err)
//	}
//
//	responseBytes, err := obj.UpdateDI(requestBytes)
//	if err != nil {
//		return diag.FromErr(err)
//	}
//
//	err = obj.Deserialize(responseBytes)
//	if err != nil {
//		return diag.FromErr(err)
//	}
//
//	obj.WriteTF(res)
//	return diags
//}
//
//func ProjectDelete(ctx context.Context, res *schema.ResourceData, m interface{}) diag.Diagnostics {
//	var diags diag.Diagnostics
//	obj := models.Project{}
//	obj.ReadTF(res)
//
//	err := obj.DeleteDI()
//	if err != nil {
//		return diag.FromErr(err)
//	}
//	res.SetId("")
//	return diags
//}
//
//func ProjectImport(ctx context.Context, res *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
//	obj := models.Project{Id: uuid.MustParse(res.Id())}
//	responseBytes, err := obj.ReadDI()
//	if err != nil {
//		return nil, err
//	}
//	err = obj.Deserialize(responseBytes)
//	if err != nil {
//		return nil, err
//	}
//	obj.WriteTF(res)
//
//	objBytes, _ := obj.ToHCL(nil)
//	// log.Println(string(objBytes))
//
//	index := bytes.IndexByte(objBytes, byte('{'))
//
//	firstString := objBytes[:index+1]
//
//	fileBytes, err := ioutil.ReadFile("project.tf")
//	if err != nil {
//		return nil, err
//	}
//
//	toReplace := []byte(fmt.Sprintf("%s}", firstString))
//
//	newBytes := bytes.Replace(fileBytes, toReplace, objBytes, -1)
//
//	err = ioutil.WriteFile("project.tf", newBytes, 0600)
//	if err != nil {
//		return nil, err
//	}
//
//	return []*schema.ResourceData{res}, nil
//}
