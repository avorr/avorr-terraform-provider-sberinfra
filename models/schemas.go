package models

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
		"group_id":       {Type: schema.TypeString, Required: true, ValidateFunc: validation.IsUUID},
		//"domain_id":       {Type: schema.TypeString, Optional: true},
		"default_network": {Type: schema.TypeString, Computed: true},
		"datacenter":      {Type: schema.TypeString, Required: true},
		"jump_host":       {Type: schema.TypeString, Required: true},
		"desc":            {Type: schema.TypeString, Optional: true},
		"limits": {
			Type:     schema.TypeSet,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"cores_vcpu_count":  {Type: schema.TypeInt, Required: true, ValidateFunc: validation.IntBetween(1, 1000)},
					"ram_gb_amount":     {Type: schema.TypeInt, Required: true, ValidateFunc: validation.IntBetween(500, 1000000)},
					"storage_gb_amount": {Type: schema.TypeInt, Required: true, ValidateFunc: validation.IntBetween(50, 1000000000)},
				},
			},
		},
		"network": {
			Type:     schema.TypeSet,
			Required: true,
			MinItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"network_name": {Type: schema.TypeString, Required: true},
					"network_uuid": {Type: schema.TypeString, Computed: true},
					"cidr":         {Type: schema.TypeString, Required: true, ValidateFunc: validation.IsCIDR},
					"dns_nameservers": {
						Type:     schema.TypeSet,
						Required: true,
						Elem:     &schema.Schema{Type: schema.TypeString, ValidateFunc: validation.IsIPv4Address},
					},
					"enable_dhcp": {Type: schema.TypeBool, Required: true},
					"is_default":  {Type: schema.TypeBool, Optional: true, Default: false},
				},
			},
		},
	}

	SchemaVM = map[string]*schema.Schema{
		"name": {Type: schema.TypeString, Computed: true},
		//"service_name": {Type: schema.TypeString, Required: true},
		"service_name": {Type: schema.TypeString, Optional: true},
		"group_id":     {Type: schema.TypeString, Required: true},
		"project_id":   {Type: schema.TypeString, Required: true},
		//"ir_group":        {Type: schema.TypeString, Required: true},
		"ir_group":     {Type: schema.TypeString, Optional: true, Default: "vm"}, //Required
		"ir_type":      {Type: schema.TypeString, Computed: true},
		"cpu":          {Type: schema.TypeInt, Computed: true},
		"ram":          {Type: schema.TypeInt, Computed: true},
		"disk":         {Type: schema.TypeInt, Required: true},
		"flavor":       {Type: schema.TypeString, Required: true},
		"network_uuid": {Type: schema.TypeString, Optional: true},
		//"virtualization":  {Type: schema.TypeString, Required: true},
		"virtualization": {Type: schema.TypeString, Optional: true, Default: "openstack"},
		"os_name":        {Type: schema.TypeString, Required: true},
		"os_version":     {Type: schema.TypeString, Required: true},
		//"fault_tolerance": {Type: schema.TypeString, Required: true},
		"fault_tolerance": {Type: schema.TypeString, Optional: true, Default: "Stand-alone"},
		"state":           {Type: schema.TypeString, Computed: true},
		"state_resize":    {Type: schema.TypeString, Computed: true},
		//"zone":            {Type: schema.TypeString, Required: true},
		"zone":            {Type: schema.TypeString, Optional: true, Default: "internal"},
		"ip":              {Type: schema.TypeString, Computed: true},
		"step":            {Type: schema.TypeString, Computed: true},
		"public_ssh_name": {Type: schema.TypeString, Optional: true},
		"user":            {Type: schema.TypeString, Computed: true},
		"password":        {Type: schema.TypeString, Computed: true},
		"volume": {
			Type:     schema.TypeSet,
			Optional: true,
			// Computed: true,
			ForceNew: false,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"volume_id": {Type: schema.TypeString, Computed: true},
					//"path":         {Type: schema.TypeString, Optional: true},
					"size":         {Type: schema.TypeInt, Required: true, ForceNew: false},
					"storage_type": {Type: schema.TypeString, Optional: true},
				},
			},
		},
		"tag_ids": {Type: schema.TypeSet, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
	}

	//{
	//  "security_group": {
	//    "group_name": "string",
	//    "server_uuid": "string",
	//    "security_rules": [
	//      {
	//        "direction": "ingress",
	//        "ethertype": "IPv4",
	//        "protocol": "tcp",
	//        "port_range_min": 0,
	//        "port_range_max": 0,
	//        "remote_ip_prefix": "string",
	//        "remote_group_id": "string"
	//      }
	//    ]
	//  }
	//}

	//{
	// "security_group": {
	//   "group_name": "dsdsds",
	//   "security_rules": [
	//     {
	//       "ethertype": "IPv4",
	//       "id": "3",
	//       "direction": "ingress",
	//       "protocol": "tcp",
	//       "remote_ip_prefix": "172.21.21.0/0"
	//     }
	//   ]
	// }
	//}

	SchemaSecurityGroup = map[string]*schema.Schema{
		"project_id": {Type: schema.TypeString, Required: true},
		"group_name": {Type: schema.TypeString, Required: true},
		"security_rule": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id":               {Type: schema.TypeString, Computed: true},
					"ethertype":        {Type: schema.TypeString, Required: true, ValidateFunc: validation.StringInSlice([]string{"IPv4", "IPv6"}, false)},
					"direction":        {Type: schema.TypeString, Required: true, ValidateFunc: validation.StringInSlice([]string{"ingress", "egress"}, false)},
					"protocol":         {Type: schema.TypeString, Required: true, ValidateFunc: validation.StringInSlice([]string{"tcp", "udp", "icmp"}, false)},
					"remote_ip_prefix": {Type: schema.TypeString, Optional: true, ValidateFunc: validation.IsCIDR},
					"port_range_min":   {Type: schema.TypeInt, Optional: true, ValidateFunc: validation.IsPortNumber},
					"port_range_max":   {Type: schema.TypeInt, Optional: true, ValidateFunc: validation.IsPortNumber},
				},
			},
		},
		//"network": {
		//	Type:     schema.TypeSet,
		//	Required: true,
		//	MinItems: 1,
		//	Elem: &schema.Resource{
		//
		//		Schema: map[string]*schema.Schema{
		//			"network_name": {Type: schema.TypeString, Required: true},
		//			"network_uuid": {Type: schema.TypeString, Computed: true},
		//			"cidr":         {Type: schema.TypeString, Required: true},
		//			"dns_nameservers": {
		//				Type:     schema.TypeSet,
		//				Required: true,
		//				Elem:     &schema.Schema{Type: schema.TypeString},
		//			},
		//			"enable_dhcp": {Type: schema.TypeBool, Required: true},
		//			"is_default":  {Type: schema.TypeBool, Optional: true, Default: false},
		//		},
		//	},
		//},
	}

	SchemaTag = map[string]*schema.Schema{
		"name": {Type: schema.TypeString, Required: true, ForceNew: true},
	}

}
