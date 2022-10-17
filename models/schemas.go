package models

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	//// data
	SchemaDomain map[string]*schema.Schema
	SchemaGroup  map[string]*schema.Schema

	//// resource
	SchemaProject map[string]*schema.Schema
	SchemaVM      map[string]*schema.Schema
	SchemaTag     map[string]*schema.Schema
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
		"ir_group":       {Type: schema.TypeString, Required: true},
		"type":           {Type: schema.TypeString, Required: true},
		"ir_type":        {Type: schema.TypeString, Required: true},
		"virtualization": {Type: schema.TypeString, Required: true},
		"name":           {Type: schema.TypeString, Optional: true},
		"group_id":       {Type: schema.TypeString, Required: true},
		//"domain_id":       {Type: schema.TypeString, Optional: true},
		"default_network": {Type: schema.TypeString, Computed: true},
		"datacenter":      {Type: schema.TypeString, Required: true},
		"jump_host":       {Type: schema.TypeString, Required: true},
		"desc":            {Type: schema.TypeString, Optional: true},
		//"limits": {
		//	Type:     schema.TypeMap,
		//	Optional: true,
		//	Elem:     &schema.Schema{Type: schema.TypeString, Optional: true},
		//},
		"limits": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"cores_vcpu_count":  {Type: schema.TypeString, Required: true},
					"ram_gb_amount":     {Type: schema.TypeString, Required: true},
					"storage_gb_amount": {Type: schema.TypeString, Required: true},
				},
			},
		},
		"network": {
			Type:     schema.TypeSet,
			Required: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"network_name": {Type: schema.TypeString, Required: true},
					"network_uuid": {Type: schema.TypeString, Computed: true},
					"cidr":         {Type: schema.TypeString, Required: true},
					"dns_nameservers": {
						Type:     schema.TypeSet,
						Required: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
					"enable_dhcp": {Type: schema.TypeBool, Required: true},
					"is_default":  {Type: schema.TypeBool, Optional: true},
				},
			},
		},
	}

	SchemaVM = map[string]*schema.Schema{
		"name":            {Type: schema.TypeString, Computed: true},
		"service_name":    {Type: schema.TypeString, Required: true},
		"group_id":        {Type: schema.TypeString, Required: true},
		"project_id":      {Type: schema.TypeString, Required: true},
		"cluster_uuid":    {Type: schema.TypeString, Computed: true},
		"ir_group":        {Type: schema.TypeString, Required: true},
		"ir_type":         {Type: schema.TypeString, Computed: true},
		"cpu":             {Type: schema.TypeInt, Computed: true},
		"ram":             {Type: schema.TypeInt, Computed: true},
		"disk":            {Type: schema.TypeInt, Required: true},
		"flavor":          {Type: schema.TypeString, Required: true},
		"region":          {Type: schema.TypeString, Optional: true},
		"network_uuid":    {Type: schema.TypeString, Optional: true},
		"virtualization":  {Type: schema.TypeString, Required: true},
		"os_name":         {Type: schema.TypeString, Required: true},
		"os_version":      {Type: schema.TypeString, Required: true},
		"fault_tolerance": {Type: schema.TypeString, Required: true},
		"state":           {Type: schema.TypeString, Computed: true},
		"state_resize":    {Type: schema.TypeString, Computed: true},
		"zone":            {Type: schema.TypeString, Required: true},
		"ip":              {Type: schema.TypeString, Computed: true},
		"dns":             {Type: schema.TypeString, Computed: true},
		"dns_name":        {Type: schema.TypeString, Computed: true},
		"step":            {Type: schema.TypeString, Computed: true},
		"public_ssh_name": {Type: schema.TypeString, Optional: true},
		"group":           {Type: schema.TypeString, Optional: true},
		"user":            {Type: schema.TypeString, Computed: true},
		"password":        {Type: schema.TypeString, Computed: true},
		"app_params": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString, Required: true},
		},
		"volume": {
			Type:     schema.TypeSet,
			Optional: true,
			// Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"path":         {Type: schema.TypeString, Optional: true},
					"size":         {Type: schema.TypeInt, Required: true},
					"storage_type": {Type: schema.TypeString, Optional: true},
				},
			},
		},
		"tag_ids": {Type: schema.TypeSet, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
	}

	SchemaTag = map[string]*schema.Schema{
		"name": {Type: schema.TypeString, Required: true, ForceNew: true},
	}
}
