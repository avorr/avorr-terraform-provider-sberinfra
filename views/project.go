package views

import (
	"context"

	"stash.sigma.sbrf.ru/sddevops/terraform-provider-di/models"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ProjectCreate(ctx context.Context, res *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	obj := models.Project{}
	obj.ReadTF(res)

	requestBytes, err := obj.Serialize()
	if err != nil {
		return diag.FromErr(err)
	}

	responseBytes, err := obj.CreateDI(requestBytes)
	if err != nil {
		return diag.FromErr(err)
	}

	err = obj.ParseIdFromCreateResponse(responseBytes)
	if err != nil {
		return diag.FromErr(err)
	}

	responseBytes, err = obj.ReadDI()
	if err != nil {
		return diag.FromErr(err)
	}

	err = obj.Deserialize(responseBytes)
	if err != nil {
		return diag.FromErr(err)
	}
	obj.WriteTF(res)
	return diags
}

func ProjectRead(ctx context.Context, res *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	obj := models.Project{}
	obj.ReadTF(res)

	responseBytes, err := obj.ReadDI()
	if err != nil {
		return diag.FromErr(err)
	}

	err = obj.Deserialize(responseBytes)
	if err != nil {
		return diag.FromErr(err)
	}

	obj.WriteTF(res)
	return diags
}

func ProjectUpdate(ctx context.Context, res *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	obj := models.Project{}
	obj.ReadTF(res)

	requestBytes, err := obj.Serialize()
	if err != nil {
		return diag.FromErr(err)
	}

	responseBytes, err := obj.UpdateDI(requestBytes)
	if err != nil {
		return diag.FromErr(err)
	}

	err = obj.Deserialize(responseBytes)
	if err != nil {
		return diag.FromErr(err)
	}

	obj.WriteTF(res)
	return diags
}

func ProjectDelete(ctx context.Context, res *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	obj := models.Project{}
	obj.ReadTF(res)

	err := obj.DeleteDI()
	if err != nil {
		return diag.FromErr(err)
	}
	res.SetId("")
	return diags
}
