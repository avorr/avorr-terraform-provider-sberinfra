package views

import (
	"base.sw.sbc.space/pid/terraform-provider-si/models"
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func SecurityGroupCreate(ctx context.Context, res *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	obj := models.SecurityGroup{}

	//diags = obj.ReadTF(res)
	obj.ReadTF(res)

	requestBytes, err := obj.Serialize()

	if err != nil {
		return diag.FromErr(err)
	}

	_, err = obj.CreateResource(requestBytes)

	if err != nil {
		return diag.FromErr(err)
	}

	_, err = obj.StateChangeSecurityGroup(res).WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	obj.WriteTF(res)
	return diags
}

func SecurityGroupRead(ctx context.Context, res *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	obj := models.SecurityGroup{}
	obj.ReadTF(res)

	responseBytes, err := obj.ReadResource()
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

func SecurityGroupUpdate(ctx context.Context, res *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	obj := models.SecurityGroup{}

	diags = obj.ReadTF(res)
	if diags.HasError() {
		return diags
	}
	if res.HasChange("security_rule") {
		v1, v2 := res.GetChange("security_rule")
		securityRulesSetState := v1.(*schema.Set)
		securityRulesSetTf := v2.(*schema.Set)

		removeSrSet := securityRulesSetState.Difference(securityRulesSetTf)
		addSrSet := securityRulesSetTf.Difference(securityRulesSetState)

		//Add rule
		for _, rule := range addSrSet.List() {
			rule := rule.(map[string]interface{})
			securityRuleMap := map[string]map[string]interface{}{
				"security_rule": {
					"ethertype":         rule["ethertype"],
					"direction":         rule["direction"],
					"port_range_min":    rule["port_range_min"],
					"port_range_max":    rule["port_range_max"],
					"protocol":          rule["protocol"],
					"remote_ip_prefix":  rule["remote_ip_prefix"],
					"security_group_id": obj.SecurityGroupID,
				},
			}

			if rule["remote_ip_prefix"] == "" {
				delete(securityRuleMap["security_rule"], "remote_ip_prefix")
			}

			if rule["port_range_min"] == 0 {
				delete(securityRuleMap["security_rule"], "port_range_min")
			}

			if rule["port_range_max"] == 0 {
				delete(securityRuleMap["security_rule"], "port_range_max")
			}

			requestBytes, err := json.Marshal(securityRuleMap)
			if err != nil {
				return diag.FromErr(err)
			}

			_, err = obj.CreateSecurityRule(requestBytes)
			if err != nil {
				return diag.FromErr(err)
			}

			_, err = obj.StateChangeSecurityGroup(res).WaitForStateContext(ctx)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		//Remove rules
		for _, rule := range removeSrSet.List() {
			foo := map[string]map[string]string{
				"security_rule": {
					"rule_uuid":         rule.(map[string]interface{})["id"].(string),
					"security_group_id": obj.SecurityGroupID,
				},
			}

			requestBytes, err := json.Marshal(foo)
			if err != nil {
				return diag.FromErr(err)
			}

			err = obj.RemoveSecurityRule(requestBytes)
			if err != nil {
				return diag.FromErr(err)
			}
			_, err = obj.StateChangeSecurityGroup(res).WaitForStateContext(ctx)
		}
	}

	obj.WriteTF(res)
	return diags
}

func SecurityGroupDelete(ctx context.Context, res *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	obj := models.SecurityGroup{}
	obj.ReadTF(res)

	responseBytes, err := obj.ReadResource()
	if err != nil {
		return diag.FromErr(err)
	}

	err = obj.Deserialize(responseBytes)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(obj.AttachedToServer) > 0 {
		for _, v := range obj.AttachedToServer {
			vm := models.Server{}
			vm.Id = v.ServerUUID

			_, err = vm.SecurityGroupVM(obj.SecurityGroupID, "detach")
			if err != nil {
				return diag.FromErr(err)
			}

			_, err = vm.StateSecurityGroupChange(res).WaitForStateContext(ctx)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	err = obj.DeleteResource()
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = obj.StateChangeSecurityGroup(res).WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	res.SetId("")
	return diags
}

func SecurityGroupImport(ctx context.Context, res *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	obj := models.SecurityGroup{}
	obj.SecurityGroupID = res.Id()

	responseBytes, err := obj.ReadAllVdc()
	if err != nil {
		return nil, err
	}
	err = obj.DeserializeImport(responseBytes)
	if err != nil {
		return nil, err
	}
	obj.IsImport = true

	obj.WriteTF(res)

	return []*schema.ResourceData{res}, nil
}
