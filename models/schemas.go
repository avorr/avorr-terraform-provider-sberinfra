package models

import (
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"strconv"
	//"github.com/hashicorp/go-cty/cty"
	//"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"regexp"
	//"sort"
	//"strconv"
)

var (
	//// data
	SchemaDomain map[string]*schema.Schema
	SchemaGroup  map[string]*schema.Schema

	//// resource
	SchemaProject       map[string]*schema.Schema
	SchemaVM            map[string]*schema.Schema
	SchemaSecurityGroup map[string]*schema.Schema
	SchemaTag           map[string]*schema.Schema
)

func init() {
	SchemaDomain = map[string]*schema.Schema{
		"id":   {Type: schema.TypeString, Computed: true},
		"name": {Type: schema.TypeString, Required: true},
		// "portal_id":       {Type: schema.TypeString, Computed: true},
		// "sap_id":          {Type: schema.TypeString, Computed: true},
		// "buisiness_block": {Type: schema.TypeString, Computed: true},
		// "type":            {Type: schema.TypeString, Computed: true},
	}

	SchemaGroup = map[string]*schema.Schema{
		// "id":   {Type: schema.TypeString, Required: true},
		"name": {Type: schema.TypeString, Required: true},
		// "portal_id":   {Type: schema.TypeInt, Computed: true},
		"domain_id": {Type: schema.TypeString, Required: true},
		// "domain_name": {Type: schema.TypeString, Computed: true},
		// "limit":   {Type: schema.TypeFloat, Computed: true},
		"is_prom": {Type: schema.TypeBool, Computed: true},
		// "is_deleted":  {Type: schema.TypeBool, Computed: true},
	}

	SchemaProject = map[string]*schema.Schema{
		//"ir_group":       {Type: schema.TypeString, Required: true},
		"ir_group": {Type: schema.TypeString, Optional: true, Default: "vdc"},
		//"type":           {Type: schema.TypeString, Required: true},
		"type": {Type: schema.TypeString, Optional: true, Default: "vdc"},
		//"ir_type":        {Type: schema.TypeString, Required: true},
		"ir_type": {Type: schema.TypeString, Optional: true, Default: "vdc_openstack"},
		//"virtualization": {Type: schema.TypeString, Required: true},
		"virtualization": {Type: schema.TypeString, Optional: true, Default: "openstack"},
		"name":           {Type: schema.TypeString, Optional: true},
		"group_id":       {Type: schema.TypeString, Required: true},
		//"domain_id":       {Type: schema.TypeString, Optional: true},
		"default_network": {Type: schema.TypeString, Computed: true},
		"datacenter":      {Type: schema.TypeString, Required: true},
		"jump_host":       {Type: schema.TypeBool, Optional: true, Default: false},
		"description":     {Type: schema.TypeString, Optional: true},
		"limits": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeInt, Required: true},
			ValidateDiagFunc: allDiagFunc(
				validation.MapKeyMatch(regexp.MustCompile("(^vcpu$)|(^ram$)|(^storage$)"), "An argument is not expected here"),
				validateLimitsMapValue(),
			),
		},
		"network": {
			Type:     schema.TypeSet,
			Required: true,
			MinItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {Type: schema.TypeString, Required: true},
					"id":   {Type: schema.TypeString, Computed: true},
					"cidr": {Type: schema.TypeString, Required: true},
					"dns": {
						Type:     schema.TypeSet,
						Required: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
					"dhcp":    {Type: schema.TypeBool, Required: true},
					"default": {Type: schema.TypeBool, Optional: true, Default: false, ValidateDiagFunc: defaultNetworkCount()},
				},
			},
		},
	}

	SchemaVM = map[string]*schema.Schema{
		"name": {Type: schema.TypeString, Computed: true},
		//"service_name": {Type: schema.TypeString, Required: true},
		"service_name": {Type: schema.TypeString, Optional: true},
		"group_id":     {Type: schema.TypeString, Required: true},
		"vdc_id":       {Type: schema.TypeString, Required: true},
		//"ir_group":        {Type: schema.TypeString, Required: true},
		"ir_group": {Type: schema.TypeString, Optional: true, Default: "vm"}, //Required
		"ir_type":  {Type: schema.TypeString, Computed: true},
		"cpu":      {Type: schema.TypeInt, Computed: true},
		"ram":      {Type: schema.TypeInt, Computed: true},
		//"disk":         {Type: schema.TypeInt, Optional: true},
		"flavor":     {Type: schema.TypeString, Required: true},
		"network_id": {Type: schema.TypeString, Optional: true},
		//"virtualization":  {Type: schema.TypeString, Required: true},
		"virtualization": {Type: schema.TypeString, Optional: true, Default: "openstack"},
		"os_name":        {Type: schema.TypeString, Required: true},
		"os_version":     {Type: schema.TypeString, Required: true},
		//"fault_tolerance": {Type: schema.TypeString, Required: true},
		"fault_tolerance": {Type: schema.TypeString, Optional: true, Default: "Stand-alone"},
		"state":           {Type: schema.TypeString, Computed: true},
		"state_resize":    {Type: schema.TypeString, Computed: true},
		//"zone":            {Type: schema.TypeString, Required: true},
		"zone": {Type: schema.TypeString, Optional: true, Default: "internal"},
		"ip":   {Type: schema.TypeString, Computed: true},
		//"dns":             {Type: schema.TypeString, Computed: true},
		//"dns_name":        {Type: schema.TypeString, Computed: true},
		"step":            {Type: schema.TypeString, Computed: true},
		"public_ssh_name": {Type: schema.TypeString, Optional: true},
		//"group":           {Type: schema.TypeString, Optional: true},
		"user":     {Type: schema.TypeString, Computed: true},
		"password": {Type: schema.TypeString, Computed: true},
		"disk": {
			Type:     schema.TypeMap,
			Required: true,
			Elem:     &schema.Schema{Type: schema.TypeString, Required: true},
			ValidateDiagFunc: allDiagFunc(
				validation.MapKeyMatch(regexp.MustCompile("(^size$)|(^storage_type$)"), "An argument is not expected here"),
				validateMapValue(),
			),
		},

		//"hdd": {
		//	Type:     schema.TypeSet,
		//	Required: true,
		//	ForceNew: false,
		//	MaxItems: 1,
		//	Elem: &schema.Resource{
		//		Schema: map[string]*schema.Schema{
		//			"size":         {Type: schema.TypeInt, Required: true, ForceNew: false},
		//			"storage_type": {Type: schema.TypeString, Optional: true, ForceNew: false},
		//		},
		//	},
		//},

		"volume": {
			Type:     schema.TypeSet,
			Optional: true,
			ForceNew: false,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					//"path":         {Type: schema.TypeString, Optional: true},
					"volume_id":    {Type: schema.TypeString, Computed: true},
					"size":         {Type: schema.TypeInt, Required: true, ForceNew: false},
					"storage_type": {Type: schema.TypeString, Optional: true},
				},
			},
		},
		"tag_ids":         {Type: schema.TypeSet, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
		"security_groups": {Type: schema.TypeSet, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
	}
	SchemaSecurityGroup = map[string]*schema.Schema{
		"id":     {Type: schema.TypeString, Computed: true},
		"vdc_id": {Type: schema.TypeString, Required: true, ForceNew: false},
		//"name":       {Type: schema.TypeString, Required: true, ForceNew: true},
		"name": {Type: schema.TypeString, Required: true},
		"security_rule": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id":               {Type: schema.TypeString, Computed: true},
					"ethertype":        {Type: schema.TypeString, Required: true, ForceNew: false, ValidateFunc: validation.StringInSlice([]string{"IPv4", "IPv6"}, false)},
					"direction":        {Type: schema.TypeString, Required: true, ForceNew: false, ValidateFunc: validation.StringInSlice([]string{"ingress", "egress"}, false)},
					"protocol":         {Type: schema.TypeString, Required: true, ForceNew: false, ValidateFunc: validation.StringInSlice([]string{"tcp", "udp", "icmp"}, false)},
					"remote_ip_prefix": {Type: schema.TypeString, Optional: true, ForceNew: false, ValidateFunc: validation.IsCIDR},
					"port_range_min":   {Type: schema.TypeInt, Optional: true, ForceNew: false, ValidateFunc: validation.IsPortNumber},
					"port_range_max":   {Type: schema.TypeInt, Optional: true, ForceNew: false, ValidateFunc: validation.IsPortNumber},
					"remote_group_id":  {Type: schema.TypeString, Optional: true, ForceNew: false, ValidateFunc: validation.IsCIDR},
				},
			},
		},
	}
	SchemaTag = map[string]*schema.Schema{
		"name": {Type: schema.TypeString, Required: true, ForceNew: true},
	}
}

func validateLimitsMapValue() schema.SchemaValidateDiagFunc {
	return func(v interface{}, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		for key, value := range v.(map[string]interface{}) {
			if key == "vcpu" {
				if !(value.(int) >= 1 && value.(int) <= 200000) {
					diags = append(diags, diag.Diagnostic{
						Severity:      diag.Error,
						Summary:       "vcpu value is not in range",
						Detail:        fmt.Sprintf("expected limits.vcpu to be in the range (1 - 1000), got %d", value),
						AttributePath: append(path, cty.IndexStep{Key: cty.StringVal(key)}),
					})
				}
			} else if key == "ram" {
				if !(value.(int) >= 500 && value.(int) <= 1000000) {
					diags = append(diags, diag.Diagnostic{
						Severity:      diag.Error,
						Summary:       "ram value is not in range",
						Detail:        fmt.Sprintf("expected limits.ram to be in the range (500 - 1000000), got %d", value),
						AttributePath: append(path, cty.IndexStep{Key: cty.StringVal(key)}),
					})
				}
			} else if key == "storage" {
				if !(value.(int) >= 50 && value.(int) <= 1000000000) {
					diags = append(diags, diag.Diagnostic{
						Severity:      diag.Error,
						Summary:       "storage value is not in range",
						Detail:        fmt.Sprintf("expected limits.storage to be in the range (50 - 1000000000), got %d", value),
						AttributePath: append(path, cty.IndexStep{Key: cty.StringVal(key)}),
					})
				}

			}
		} //else if key == "storage_type" {
		//if value != "iscsi-fast-01" {
		//	diags = append(diags, diag.Diagnostic{
		//		Severity:      diag.Error,
		//		Summary:       "Invalid map key",
		//		Detail:        detail,
		//		AttributePath: append(path, cty.IndexStep{Key: cty.StringVal(key)}),
		//	})
		//}
		//}
		return diags
	}
}

var count int

func defaultNetworkCount() schema.SchemaValidateDiagFunc {
	return func(v interface{}, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		if v.(bool) {
			count++
		}
		if count == 2 {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Default networks should not be more than one",
				Detail:        fmt.Sprintf("\"default = true\", this parameter must be on one network only"),
				AttributePath: append(path, cty.IndexStep{Key: cty.StringVal("default")}),
			})
		}
		return diags
	}
}

func allDiagFunc(validators ...schema.SchemaValidateDiagFunc) schema.SchemaValidateDiagFunc {
	return func(i interface{}, k cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		for _, validator := range validators {
			diags = append(diags, validator(i, k)...)
		}
		return diags
	}
}

func validateMapValue() schema.SchemaValidateDiagFunc {
	return func(v interface{}, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		for key, value := range v.(map[string]interface{}) {
			var detail string
			if key == "size" {
				_, err := strconv.Atoi(value.(string))
				if err != nil {
					diags = append(diags, diag.Diagnostic{
						Severity:      diag.Error,
						Summary:       "Inappropriate value for attribute \"size\": a number is required.",
						Detail:        detail,
						AttributePath: append(path, cty.IndexStep{Key: cty.StringVal(key)}),
					})
				}
			} //else if key == "storage_type" {
			//if value != "iscsi-fast-01" {
			//	diags = append(diags, diag.Diagnostic{
			//		Severity:      diag.Error,
			//		Summary:       "Invalid map key",
			//		Detail:        detail,
			//		AttributePath: append(path, cty.IndexStep{Key: cty.StringVal(key)}),
			//	})
			//}
			//}
		}
		return diags
	}
}

//func sortedKeys(m map[string]interface{}) []string {
//	keys := make([]string, len(m))
//
//	i := 0
//	for key := range m {
//		keys[i] = key
//		i++
//	}

//sort.Strings(keys)
//
//return keys
//}
