package models

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"time"
)

type SecurityGroup struct {
	ProjectID       string `json:"-"`
	GroupName       string `json:"group_name"`
	SecurityRules   []Rule `json:"security_rules"`
	State           string `json:"-"`
	SecurityGroupID string `json:"-"` //`json:"security_group_id"`
}

type Rule struct {
	Ethertype      string `json:"ethertype"`
	ID             string `json:"id,omitempty"`
	Status         string `json:"status,omitempty"`
	Direction      string `json:"direction"`
	Protocol       string `json:"protocol"`
	PortRangeMin   int    `json:"port_range_min"`
	PortRangeMax   int    `json:"port_range_max"`
	RemoteIpPrefix string `json:"remote_ip_prefix"`
	RemoteGroupID  string `json:"remote_group_id,omitempty"`
}

func (o *SecurityGroup) ReadTF(res *schema.ResourceData) diag.Diagnostics {
	o.GroupName = res.Get("group_name").(string)
	o.ProjectID = res.Get("project_id").(string)
	//o.SecurityGroupID = res.Get("security_group_id").(string)
	o.SecurityGroupID = res.Id()
	rules, ok := res.GetOk("security_rule")
	if ok {
		for _, v := range rules.(*schema.Set).List() {
			v := v.(map[string]interface{})
			o.SecurityRules = append(o.SecurityRules, Rule{
				//ID:             uuid.Nil,
				Ethertype:      v["ethertype"].(string),
				Direction:      v["direction"].(string),
				PortRangeMin:   v["port_range_min"].(int),
				PortRangeMax:   v["port_range_max"].(int),
				Protocol:       v["protocol"].(string),
				RemoteIpPrefix: v["remote_ip_prefix"].(string),
				RemoteGroupID:  v["remote_group_id"].(string),
			})
		}
	} else {
		o.SecurityRules = []Rule{}
	}
	return diag.Diagnostics{}
}

func (o *SecurityGroup) Serialize() ([]byte, error) {
	requestMap := map[string]*SecurityGroup{"security_group": o}
	requestBytes, err := json.Marshal(requestMap)
	if err != nil {
		return nil, err
	}
	return requestBytes, nil
}

func (o *SecurityGroup) CreateResource(data []byte) ([]byte, error) {
	return Api.NewRequestCreate(fmt.Sprintf("projects/%s/security_groups", o.ProjectID), data)
}

func (o *SecurityGroup) CreateSecurityRule(data []byte) ([]byte, error) {
	return Api.NewRequestCreate(fmt.Sprintf("projects/%s/security_rules", o.ProjectID), data)
}

func (o *SecurityGroup) Deserialize(responseBytes []byte) error {
	var responseStruct map[string]ProjectNew
	err := json.Unmarshal(responseBytes, &responseStruct)
	if err != nil {
		return err
	}
	var existGroup bool
out:
	for _, group := range responseStruct["project"].SecurityGroups {
		if group.GroupName == o.GroupName {
			existGroup = true
			o.SecurityGroupID = group.SecurityGroupID
			if len(o.SecurityRules) > 0 {
				for _, rule := range group.Rules {
					for i2, stateRule := range o.SecurityRules {
						if rule.ID == "" || rule.Status != "" {
							o.State = "creating"
							break out
						}
						if cmp.Equal(rule, stateRule, cmpopts.IgnoreFields(Rule{}, "Status", "ID")) && rule.ID != "" {
							o.SecurityRules[i2].ID = rule.ID
							o.State = group.Status
						}
					}
				}
			} else {
				o.State = group.Status
			}
		}
	}

	if !existGroup {
		o.State = "deleted"
	}
	return nil
}

func (o *SecurityGroup) StateChangeSecurityGroup(res *schema.ResourceData) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Timeout:      res.Timeout(schema.TimeoutCreate),
		PollInterval: 15 * time.Second,
		Pending:      []string{"Creating", "Pending", "Deleting"},
		Target:       []string{"Running", "Damaged", "Deleted"},
		Refresh: func() (interface{}, string, error) {

			responseBytes, err := o.ReadResource()
			if err != nil {
				return nil, "error", err
			}

			err = o.Deserialize(responseBytes)
			if err != nil {
				return nil, "error", err
			}

			log.Printf("[DEBUG] Refresh state for [%s]: state: %s", o.GroupName, o.State)

			// write to TF state
			//o.WriteTF(res)

			if o.State == "ready" || o.State == "created" {
				return o, "Running", nil
			}

			if o.State == "deleted" {
				return o, "Deleted", nil
			}
			//if o.State == "created" {
			//	return o, "Running", nil
			//}
			//if o.State == "damaged" {
			//	return o, "Damaged", nil
			//}
			//if o.State == "pending" {
			//	return o, "Pending", nil
			//}
			return o, "Creating", nil
		},
	}
}

func (o *SecurityGroup) WriteTF(res *schema.ResourceData) {
	res.SetId(o.SecurityGroupID)

	rules := make([]map[string]interface{}, 0)
	for _, v := range o.SecurityRules {
		rule := map[string]interface{}{
			"id":               v.ID,
			"direction":        v.Direction,
			"ethertype":        v.Ethertype,
			"protocol":         v.Protocol,
			"port_range_min":   v.PortRangeMin,
			"port_range_max":   v.PortRangeMax,
			"remote_ip_prefix": v.RemoteIpPrefix,
			"remote_group_id":  v.RemoteGroupID,
		}
		rules = append(rules, rule)
	}

	err := res.Set("security_rule", rules)
	if err != nil {
		log.Println(err)
	}
	//res.SetConnInfo("network")
	//res.ConnInfo()
	//res.Set("network_uuid")
}

func (o *SecurityGroup) ReadResource() ([]byte, error) {
	return Api.NewRequestRead(fmt.Sprintf("projects/%s", o.ProjectID))
}

func (o *SecurityGroup) DeleteResource() error {
	return Api.NewRequestDelete(fmt.Sprintf("projects/%s/security_groups?security_group[security_group_id]=%s", o.ProjectID, o.SecurityGroupID), nil, 200)
}

func (o *SecurityGroup) RemoveSecurityRule(data []byte) error {
	return Api.NewRequestDelete(fmt.Sprintf("projects/%s/security_rules", o.ProjectID), data, 200)
}
