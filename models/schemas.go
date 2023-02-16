package models

import (
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"net"
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
	SchemaVdc           map[string]*schema.Schema
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
		"domain_id": {Type: schema.TypeString, Required: true, ValidateFunc: validation.IsUUID},
		// "domain_name": {Type: schema.TypeString, Computed: true},
		// "limit":   {Type: schema.TypeFloat, Computed: true},
		"is_prom": {Type: schema.TypeBool, Computed: true},
		// "is_deleted":  {Type: schema.TypeBool, Computed: true},
	}

	SchemaVdc = map[string]*schema.Schema{
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
				validation.MapKeyMatch(regexp.MustCompile("(^cores$)|(^ram$)|(^storage$)"), "An argument is not expected here"),
				validateLimitsMapValue(),
			),
		},
		"network": {
			Type:     schema.TypeSet,
			Required: true,
			MinItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					//"name": {Type: schema.TypeString, Required: true, ValidateDiagFunc: uniqueNetworkName()},
					"name": {Type: schema.TypeString, Required: true, ValidateFunc: validation.StringIsNotEmpty},
					"id":   {Type: schema.TypeString, Computed: true},
					//"cidr": {Type: schema.TypeString, Required: true, ValidateDiagFunc: uniqueCidr()},
					"cidr": {Type: schema.TypeString, Required: true, ValidateFunc: validation.IsCIDR},
					"dns": {
						Type:     schema.TypeSet,
						Required: true,
						Elem:     &schema.Schema{Type: schema.TypeString, ValidateFunc: validation.IsIPv4Address},
					},
					"dhcp": {Type: schema.TypeBool, Required: true},
					//"default": {Type: schema.TypeBool, Optional: true, Default: false, ValidateDiagFunc: defaultNetworkCount()},
					"default": {Type: schema.TypeBool, Optional: true, Default: false},
				},
			},
		},
	}

	SchemaVM = map[string]*schema.Schema{
		"name": {Type: schema.TypeString, Computed: true},
		//"service_name": {Type: schema.TypeString, Required: true},
		"service_name": {Type: schema.TypeString, Optional: true},
		"description":  {Type: schema.TypeString, Optional: true},
		"group_id":     {Type: schema.TypeString, Required: true, ValidateFunc: validation.IsUUID},
		"vdc_id":       {Type: schema.TypeString, Required: true, ValidateFunc: validation.IsUUID},
		//"ir_group":        {Type: schema.TypeString, Required: true},
		"ir_group":   {Type: schema.TypeString, Optional: true, Default: "vm"}, //Required
		"ir_type":    {Type: schema.TypeString, Computed: true},
		"cpu":        {Type: schema.TypeInt, Computed: true},
		"ram":        {Type: schema.TypeInt, Computed: true},
		"flavor":     {Type: schema.TypeString, Required: true, ValidateFunc: validation.StringInSlice([]string{"m1.tiny", "m1.small", "m1.medium", "m1.large", "m1.xlarge", "m2.tiny", "m2.small", "m2.medium", "m2.large", "m2.xlarge", "m2.xxlarge", "m3.medium", "m3.large", "m4.tiny", "m4.small", "m4.medium", "m4.large", "m4.xlarge", "m6.tiny", "m6.small", "m6.medium", "m6.large", "m6.xlarge", "m8.tiny", "m8.small", "m8.medium", "m8.large", "m8.xlarge", "m12.large", "m16.tiny", "m16.small", "m16.large", "m16.xxlarge", "kasper_n2.tiny", "kasper_n2.small", "kasper_n2.medium", "kasper_n2.large", "kasper_n2.xlarge", "kasper_n1.small", "kasper_n1.medium", "kasper_n1.large", "kasper_n3.small", "kasper_n3.medium", "kasper_n3.large", "kasper_n3.xlarge"}, false)},
		"network_id": {Type: schema.TypeString, Optional: true, ValidateFunc: validation.IsUUID},
		//"virtualization":  {Type: schema.TypeString, Required: true},
		"virtualization": {Type: schema.TypeString, Optional: true, Default: "openstack"},
		"os_name":        {Type: schema.TypeString, Required: true},
		"os_version":     {Type: schema.TypeString, Required: true},
		//"fault_tolerance": {Type: schema.TypeString, Required: true},
		"fault_tolerance": {Type: schema.TypeString, Optional: true, Default: "Stand-alone"},
		"state":           {Type: schema.TypeString, Computed: true},
		"state_resize":    {Type: schema.TypeString, Computed: true},
		"zone":            {Type: schema.TypeString, Optional: true, Default: "internal"},
		"ip":              {Type: schema.TypeString, Computed: true},
		"step":            {Type: schema.TypeString, Computed: true},
		"public_ssh_name": {Type: schema.TypeString, Optional: true},
		"user":            {Type: schema.TypeString, Computed: true},
		"password":        {Type: schema.TypeString, Computed: true},
		"disk": {
			Type:     schema.TypeMap,
			Required: true,
			Elem:     &schema.Schema{Type: schema.TypeString, Required: true},
			ValidateDiagFunc: allDiagFunc(
				validation.MapKeyMatch(regexp.MustCompile("(^size$)|(^storage_type$)"), "An argument is not expected here"),
				validateMapValue(),
			),
		},
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
		"tag_ids":         {Type: schema.TypeSet, Optional: true, Elem: &schema.Schema{Type: schema.TypeString, ValidateFunc: validation.IsUUID}},
		"security_groups": {Type: schema.TypeSet, Optional: true, Elem: &schema.Schema{Type: schema.TypeString, ValidateFunc: validation.IsUUID}},
	}

	SchemaSecurityGroup = map[string]*schema.Schema{
		"id":     {Type: schema.TypeString, Computed: true},
		"vdc_id": {Type: schema.TypeString, Required: true, ForceNew: false, ValidateFunc: validation.IsUUID},
		//"name":       {Type: schema.TypeString, Required: true, ForceNew: true},
		"name":               {Type: schema.TypeString, Required: true},
		"attached_to_server": {Type: schema.TypeSet, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
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
					"remote_group_id":  {Type: schema.TypeString, Optional: true, ForceNew: false, ValidateFunc: validation.IsUUID},
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
			if key == "cores" {
				if !(value.(int) >= 1 && value.(int) <= 200000) {
					diags = append(diags, diag.Diagnostic{
						Severity:      diag.Error,
						Summary:       "cores value is not in range",
						Detail:        fmt.Sprintf("expected limits.cores to be in the range (1 - 1000), got %d", value),
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

func uniqueCidr() schema.SchemaValidateDiagFunc {
	var cidr []string
	var count int
	return func(v interface{}, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		count++
		v, ok := v.(string)
		if !ok {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("expected type of network.%d.cidr to be string", count),
				//Detail:        fmt.Sprintf("There is more than one network with the same cidr [%s]", value),
				//AttributePath: append(path, cty.IndexStep{Key: cty.StringVal("default")}),
			})
		}

		if _, _, err := net.ParseCIDR(v.(string)); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("expected network.%d.cidr to be a valid IPv4 Value, got %v", count, v.(string)),
				//Detail:        fmt.Sprintf("There is more than one network with the same cidr [%s]", value),
				//AttributePath: append(path, cty.IndexStep{Key: cty.StringVal("default")}),
			})
		}

		for _, value := range cidr {
			if v.(string) == value {
				diags = append(diags, diag.Diagnostic{
					Severity:      diag.Error,
					Summary:       "There mustn't be networks with the same cidr",
					Detail:        fmt.Sprintf("There is more than one network with the same cird [%s]", value),
					AttributePath: append(path, cty.IndexStep{Key: cty.StringVal("default")}),
				})
			}
		}

		cidr = append(cidr, v.(string))
		return diags
	}
}

func uniqueNetworkName() schema.SchemaValidateDiagFunc {
	var networkNames []string
	return func(v interface{}, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		for _, value := range networkNames {
			if v.(string) == value {
				diags = append(diags, diag.Diagnostic{
					Severity:      diag.Error,
					Summary:       "There mustn't be networks with the same name",
					Detail:        fmt.Sprintf("There is more than one network with the same name [%s]", value),
					AttributePath: append(path, cty.IndexStep{Key: cty.StringVal("default")}),
				})
			}
		}

		networkNames = append(networkNames, v.(string))
		return diags
	}
}

func defaultNetworkCount() schema.SchemaValidateDiagFunc {
	var count int
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
//	i := 0
//	for key := range m {
//		keys[i] = key
//		i++
//	}
//	sort.Strings(keys)
//
//	return keys
//}
