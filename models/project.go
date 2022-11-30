package models

import (
	"base.sw.sbc.space/pid/terraform-provider-si/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Project struct {
	Project struct {
		ID                 uuid.UUID   `json:"id"`
		Name               string      `json:"name"`
		State              string      `json:"state"`
		Type               string      `json:"type"`
		Storages           interface{} `json:"storages"`
		IrGroup            string      `json:"ir_group"`
		IrType             string      `json:"ir_type"`
		Virtualization     string      `json:"virtualization"`
		ChecksumMatch      bool        `json:"checksum_match"`
		Datacenter         string      `json:"datacenter"`
		DatacenterName     string      `json:"datacenter_name"`
		HpsmCi             interface{} `json:"hpsm_ci"`
		OrderCreatedAt     time.Time   `json:"order_created_at"`
		SerialNumber       string      `json:"serial_number"`
		OpenstackProjectID uuid.UUID   `json:"openstack_project_id"`
		DefaultNetwork     uuid.UUID   `json:"default_network"`
		Limits             struct {
			CoresVcpuCount  int `json:"cores_vcpu_count"`
			RamGbAmount     int `json:"ram_gb_amount"`
			StorageGbAmount int `json:"storage_gb_amount"`
		} `json:"limits"`
		Networks struct {
			NetworkName    string    `json:"network_name"`
			NetworkUuid    uuid.UUID `json:"network_uuid"`
			Cidr           string    `json:"cidr"`
			DNSNameservers []string  `json:"dns_nameservers"`
			EnableDhcp     bool      `json:"enable_dhcp"`
			IsDefault      bool      `json:"is_default"`
		} `json:"network"`
		NetworksRes []struct {
			Cidr           string    `json:"cidr"`
			Status         string    `json:"status"`
			EnableDhcp     bool      `json:"enable_dhcp"`
			SubnetName     string    `json:"subnet_name"`
			SubnetUUID     uuid.UUID `json:"subnet_uuid"`
			NetworkName    string    `json:"network_name"`
			NetworkUUID    uuid.UUID `json:"network_uuid"`
			DNSNameservers []string  `json:"dns_nameservers"`
			IsDefault      bool      `json:"is_default"`
		} `json:"networks"`
		RealState            string        `json:"real_state"`
		GroupName            string        `json:"group_name"`
		DomainID             uuid.UUID     `json:"domain_id"`
		GroupID              uuid.UUID     `json:"group_id"`
		JumpHost             bool          `json:"jump_host"`
		Desc                 string        `json:"desc"`
		JumpHostState        interface{}   `json:"jump_host_state"`
		JumpHostServiceName  interface{}   `json:"jump_host_service_name"`
		JumpHostCreatorLogin interface{}   `json:"jump_host_creator_login"`
		JumpHostCreatedAt    interface{}   `json:"jump_host_created_at"`
		PublicIPCount        int           `json:"public_ip_count"`
		PublicIps            []interface{} `json:"public_ips"`
		Edge                 interface{}   `json:"edge"`
		HighAvailability     interface{}   `json:"high_availability"`
		SecurityGroups       []interface{} `json:"security_groups"`
		Routers              interface{}   `json:"routers"`
		RouterInterfaces     interface{}   `json:"router_interfaces"`
	} `json:"project"`
}

type ResProject struct {
	Project struct {
		ID                 uuid.UUID     `json:"id"`
		Name               string        `json:"name"`
		State              string        `json:"state"`
		Type               string        `json:"type"`
		Storages           []interface{} `json:"storages"`
		IrGroup            string        `json:"ir_group"`
		IrType             string        `json:"ir_type"`
		Virtualization     string        `json:"virtualization"`
		ChecksumMatch      bool          `json:"checksum_match"`
		Datacenter         string        `json:"datacenter"`
		DatacenterName     string        `json:"datacenter_name"`
		HpsmCi             interface{}   `json:"hpsm_ci"`
		OrderCreatedAt     time.Time     `json:"order_created_at"`
		SerialNumber       string        `json:"serial_number"`
		OpenstackProjectID uuid.UUID     `json:"openstack_project_id"`
		DefaultNetwork     uuid.UUID     `json:"default_network"`
		Limits             struct {
			CoresVcpuCount  int `json:"cores_vcpu_count"`
			RamGbAmount     int `json:"ram_gb_amount"`
			StorageGbAmount int `json:"storage_gb_amount"`
		} `json:"limits"`
		Networks []struct {
			Cidr           string    `json:"cidr"`
			Status         string    `json:"status"`
			EnableDhcp     bool      `json:"enable_dhcp"`
			SubnetName     string    `json:"subnet_name"`
			SubnetUUID     uuid.UUID `json:"subnet_uuid"`
			NetworkName    string    `json:"network_name"`
			NetworkUUID    uuid.UUID `json:"network_uuid"`
			DNSNameservers []string  `json:"dns_nameservers"`
			IsDefault      bool      `json:"is_default"`
		} `json:"networks"`
		RealState            string        `json:"real_state"`
		DomainName           string        `json:"domain_name"`
		GroupName            string        `json:"group_name"`
		DomainID             uuid.UUID     `json:"domain_id"`
		GroupID              uuid.UUID     `json:"group_id"`
		IsProm               bool          `json:"is_prom"`
		JumpHost             bool          `json:"jump_host"`
		Desc                 string        `json:"desc"`
		JumpHostState        interface{}   `json:"jump_host_state"`
		JumpHostServiceName  interface{}   `json:"jump_host_service_name"`
		JumpHostCreatorLogin interface{}   `json:"jump_host_creator_login"`
		JumpHostCreatedAt    interface{}   `json:"jump_host_created_at"`
		PublicIPCount        int           `json:"public_ip_count"`
		PublicIps            []interface{} `json:"public_ips"`
		Edge                 interface{}   `json:"edge"`
		HighAvailability     interface{}   `json:"high_availability"`
		SecurityGroups       []struct {
			Rules []struct {
				ID              string      `json:"id"`
				Protocol        interface{} `json:"protocol"`
				Direction       string      `json:"direction"`
				Ethertype       string      `json:"ethertype"`
				PortRangeMax    interface{} `json:"port_range_max"`
				PortRangeMin    interface{} `json:"port_range_min"`
				RemoteGroupID   interface{} `json:"remote_group_id"`
				RemoteIpPrefix  interface{} `json:"remote_ip_prefix,omitempty"`
				SecurityGroupID string      `json:"security_group_id"`
			} `json:"rules"`
			Status           string        `json:"status"`
			GroupName        string        `json:"group_name"`
			SecurityGroupID  string        `json:"security_group_id"`
			AttachedToServer []interface{} `json:"attached_to_server"`
		} `json:"security_groups"`
		Routers          interface{} `json:"routers"`
		RouterInterfaces interface{} `json:"router_interfaces"`
	} `json:"project"`
}

type ProjectNew struct {
	ID                 uuid.UUID     `json:"id"`
	Name               string        `json:"name"`
	State              string        `json:"state"`
	Type               string        `json:"type"`
	Storages           []interface{} `json:"storages"`
	IrGroup            string        `json:"ir_group"`
	IrType             string        `json:"ir_type"`
	Virtualization     string        `json:"virtualization"`
	ChecksumMatch      bool          `json:"checksum_match"`
	Datacenter         string        `json:"datacenter"`
	DatacenterName     string        `json:"datacenter_name"`
	HpsmCi             interface{}   `json:"hpsm_ci"`
	OrderCreatedAt     time.Time     `json:"order_created_at"`
	SerialNumber       string        `json:"serial_number"`
	OpenstackProjectID uuid.UUID     `json:"openstack_project_id"`
	DefaultNetwork     uuid.UUID     `json:"default_network"`
	Limits             struct {
		CoresVcpuCount  int `json:"cores_vcpu_count"`
		RamGbAmount     int `json:"ram_gb_amount"`
		StorageGbAmount int `json:"storage_gb_amount"`
	} `json:"limits"`
	Networks []struct {
		Cidr           string    `json:"cidr"`
		Status         string    `json:"status"`
		EnableDhcp     bool      `json:"enable_dhcp"`
		SubnetName     string    `json:"subnet_name"`
		SubnetUUID     uuid.UUID `json:"subnet_uuid"`
		NetworkName    string    `json:"network_name"`
		NetworkUUID    uuid.UUID `json:"network_uuid"`
		DNSNameservers []string  `json:"dns_nameservers"`
		IsDefault      bool      `json:"is_default"`
	} `json:"networks"`
	RealState            string        `json:"real_state"`
	DomainName           string        `json:"domain_name"`
	GroupName            string        `json:"group_name"`
	DomainID             uuid.UUID     `json:"domain_id"`
	GroupID              uuid.UUID     `json:"group_id"`
	IsProm               bool          `json:"is_prom"`
	JumpHost             bool          `json:"jump_host"`
	Desc                 string        `json:"desc"`
	JumpHostState        interface{}   `json:"jump_host_state"`
	JumpHostServiceName  interface{}   `json:"jump_host_service_name"`
	JumpHostCreatorLogin interface{}   `json:"jump_host_creator_login"`
	JumpHostCreatedAt    interface{}   `json:"jump_host_created_at"`
	PublicIPCount        int           `json:"public_ip_count"`
	PublicIps            []interface{} `json:"public_ips"`
	Edge                 interface{}   `json:"edge"`
	HighAvailability     interface{}   `json:"high_availability"`
	SecurityGroups       []struct {
		Rules            []Rule        `json:"rules"`
		Status           string        `json:"status"`
		GroupName        string        `json:"group_name"`
		SecurityGroupID  string        `json:"security_group_id"`
		AttachedToServer []interface{} `json:"attached_to_server"`
	} `json:"security_groups"`
	Routers          interface{} `json:"routers"`
	RouterInterfaces interface{} `json:"router_interfaces"`
}

type Networks struct {
	Network struct {
		NetworkName    string   `json:"network_name"`
		Cidr           string   `json:"cidr"`
		DNSNameservers []string `json:"dns_nameservers"`
		EnableDhcp     bool     `json:"enable_dhcp"`
		IsDefault      bool     `json:"is_default"`
	} `json:"network"`
}

func (o *Project) AddNetwork(ctx context.Context, res *schema.ResourceData, additionalNets []interface{}) diag.Diagnostics {

	body := Networks{}

	for _, v := range additionalNets {
		v := v.(map[string]interface{})
		body.Network.Cidr = v["cidr"].(string)
		body.Network.EnableDhcp = v["dhcp"].(bool)
		body.Network.IsDefault = v["default"].(bool)
		body.Network.NetworkName = v["name"].(string)

		dnsNameServers := make([]string, 0)

		for _, dnsIp := range v["dns"].(*schema.Set).List() {
			dnsNameServers = append(dnsNameServers, dnsIp.(string))
		}
		body.Network.DNSNameservers = dnsNameServers
		result, err := json.Marshal(body)
		if err != nil {
			return diag.FromErr(err)
		}

		resBody, err := Api.NewRequestCreate(fmt.Sprintf("projects/%s/networks", o.Project.ID), result)

		if err != nil {
			return diag.FromErr(err)
		}

		deserializeResBody := ResProject{}
		err = json.Unmarshal(resBody, &deserializeResBody)
		if err != nil {
			return diag.FromErr(err)
		}

		_, err = deserializeResBody.StateChangeNetwork(res, body.Network.NetworkName, body.Network.IsDefault).WaitForStateContext(ctx)
	}

	return diag.Diagnostics{}
}

func (o *Project) GetProjectQuota() ([]byte, error) {
	body, err := Api.NewRequestRead(fmt.Sprintf("/v2/projects/%s/quota?group_id=%s", o.Project.ID, o.Project.GroupID))

	if err != nil {
		return nil, err
	}
	return body, nil
}

func (o *ResProject) GetProjectQuota() ([]byte, error) {
	body, err := Api.NewRequestRead(fmt.Sprintf("/v2/projects/%s/quota?group_id=%s", o.Project.ID, o.Project.GroupID))

	if err != nil {
		return nil, err
	}
	return body, nil
}

func (o *Project) SetDefaultNetwork(networkUuid string) error {
	body, err := json.Marshal(map[string]string{"network_uuid": networkUuid})
	if err != nil {
		return err
	}

	_, err = Api.NewRequestUpdate(fmt.Sprintf("projects/%s/networks/set_default", o.Project.ID), body)

	if err != nil {
		return err
	}
	return nil
}

func (o *ResProject) SetDefaultNetwork(networkUuid string) error {
	body, err := json.Marshal(map[string]string{"network_uuid": networkUuid})
	if err != nil {
		return err
	}
	_, err = Api.NewRequestUpdate(fmt.Sprintf("projects/%s/networks/set_default", o.Project.ID), body)

	if err != nil {
		return err
	}
	return nil
}

func (o *Project) GetType() string {
	return "si_vdc"
}

//func (o *Project) NewObj() DIDataResource {
//	return &Project{}
//}

func (o *Project) GetId() string {
	return o.Project.ID.String()
}

func (o *Project) GetDomainId() uuid.UUID {
	return o.Project.DomainID
}

func (o *Project) GetResType() string {
	return "si_group"
}

func (o *Project) GetResName() string {
	return o.Project.Name
}

func (o *Project) GetOutput() (string, string) {
	//return o.ResOutputName, o.ResOutputValue
	return "", ""
}

func (o *Project) SetResFields() {
	/*
		o.ResId = o.GetId()
		o.ResType = o.GetResType()
		o.ResName = utils.Reformat(o.Name)
		// o.ResDomainId = o.DomainId.String()
		o.ResDomainIdUUID = o.DomainId.String()
		o.ResOutputName = fmt.Sprintf(
			"%s_id",
			o.GetResType(),
		)
		o.ResOutputValue = fmt.Sprintf(
			"data.%s.%s.id",
			o.GetResType(),
			o.GetResName(),
		)
	*/
}

//func (o *Project) DeserializeAll(responseBytes []byte) ([]DIDataResource, error) {
//	m := make(map[string][]*Project)
//	err := json.Unmarshal(responseBytes, &m)
//	if err != nil {
//		return nil, err
//	}
//
//	m2 := make([]DIDataResource, len(m["groups"]))
//	for k, v := range m["groups"] {
//		m2[k] = v
//	}
//	return m2, nil
//}

//func (o *Project) NewObj() DIResource {
//	return &Project{}
//}

func (o *Project) ReadTF(res *schema.ResourceData) diag.Diagnostics {

	if res.Id() != "" {
		o.Project.ID = uuid.MustParse(res.Id())
	}

	o.Project.IrGroup = res.Get("ir_group").(string)
	o.Project.Type = res.Get("type").(string)
	o.Project.IrType = res.Get("ir_type").(string)
	o.Project.Virtualization = res.Get("virtualization").(string)
	o.Project.Name = res.Get("name").(string)
	o.Project.GroupID = uuid.MustParse(res.Get("group_id").(string))
	//o.Project.ID = uuid.MustParse(res.Id())
	o.Project.Datacenter = res.Get("datacenter").(string)
	o.Project.Desc = res.Get("description").(string)
	o.Project.JumpHost = res.Get("jump_host").(bool)

	limits, ok := res.GetOk("limits")
	if ok {
		limits := limits.(map[string]interface{})
		o.Project.Limits.CoresVcpuCount = limits["cores"].(int)
		o.Project.Limits.RamGbAmount = limits["ram"].(int)
		o.Project.Limits.StorageGbAmount = limits["storage"].(int)
	}

	network := res.Get("network")

	var existDefaultNetwork bool

	networkSet := network.(*schema.Set)
	for i, v := range networkSet.List() {
		v := v.(map[string]interface{})

		if v["default"].(bool) {
			existDefaultNetwork = v["default"].(bool)
		}
		if existDefaultNetwork == false && networkSet.Len() == i+1 {
			return diag.Errorf("There must be one default network. [default = true]")
		}

		if networkSet.Len() == 1 || v["default"].(bool) {
			o.Project.Networks.NetworkName = v["name"].(string)
			o.Project.Networks.Cidr = v["cidr"].(string)
			o.Project.Networks.EnableDhcp = v["dhcp"].(bool)
			o.Project.Networks.IsDefault = true
			var dnsNameServers []string
			for _, dnsIp := range v["dns"].(*schema.Set).List() {
				dnsNameServers = append(dnsNameServers, dnsIp.(string))
			}
			o.Project.Networks.DNSNameservers = dnsNameServers
		}
	}
	return diag.Diagnostics{}
}

func (o *ResProject) ReadTFRes(res *schema.ResourceData) diag.Diagnostics {

	if res.Id() != "" {
		o.Project.ID = uuid.MustParse(res.Id())
	}

	o.Project.IrGroup = res.Get("ir_group").(string)
	o.Project.Type = res.Get("type").(string)
	o.Project.IrType = res.Get("ir_type").(string)
	o.Project.Virtualization = res.Get("virtualization").(string)
	o.Project.Name = res.Get("name").(string)
	o.Project.GroupID = uuid.MustParse(res.Get("group_id").(string))
	//o.Project.ID = uuid.MustParse(res.Id())
	o.Project.Datacenter = res.Get("datacenter").(string)
	o.Project.Desc = res.Get("description").(string)
	o.Project.JumpHost = res.Get("jump_host").(bool)

	limits, ok := res.GetOk("limits")
	if ok {
		limits := limits.(map[string]interface{})
		o.Project.Limits.CoresVcpuCount = limits["cores"].(int)
		o.Project.Limits.RamGbAmount = limits["ram"].(int)
		o.Project.Limits.StorageGbAmount = limits["storage"].(int)
	}

	//net, ok := res.GetOk("network")
	//networks := make([]map[string]interface{}, 0)
	//for _, v := range o.Project.Networks {
	//	volume := map[string]interface{}{
	//		"size":         v.Size,
	//		"path":         v.Path,
	//		"storage_type": v.StorageType,
	//	}
	//	networks = append(networks, volume)
	//}
	//err := res.Set("network", networks)
	//if err != nil {
	//	log.Println(err)
	//}

	//network, ok := res.GetOk("network")

	//if ok {
	//	networkSet := network.(*schema.Set).List()
	//	for _, v := range networkSet {
	//		if v.(map[string]interface{})["default"].(bool) {
	//			o.Project.Networks.NetworkName = v.(map[string]interface{})["name"].(string)
	//			o.Project.Networks.Cidr = v.(map[string]interface{})["cidr"].(string)
	//			o.Project.Networks.EnableDhcp = v.(map[string]interface{})["dhcp"].(bool)
	//			o.Project.Networks.IsDefault = true
	//			var dnsNameServers = []string{}
	//			for _, dnsIp := range v.(map[string]interface{})["dns"].(*schema.Set).List() {
	//				dnsNameServers = append(dnsNameServers, dnsIp.(string))
	//			}
	//			o.Project.Networks.DNSNameservers = dnsNameServers
	//		}
	//	}
	//}

	return diag.Diagnostics{}
}

func (o *Project) WriteTF(res *schema.ResourceData) {
	res.SetId(o.Project.ID.String())

	res.Set("datacenter", o.Project.Datacenter)
	res.Set("ir_type", o.Project.IrType)
	res.Set("description", o.Project.Desc)
	res.Set("group_id", o.Project.GroupID.String())
	res.Set("jump_host", o.Project.JumpHost)
	res.Set("name", o.Project.Name)
	res.Set("virtualization", o.Project.Virtualization)

	limits := map[string]int{
		"cores":   o.Project.Limits.CoresVcpuCount,
		"ram":     o.Project.Limits.RamGbAmount,
		"storage": o.Project.Limits.StorageGbAmount,
	}
	err := res.Set("limits", limits)
	if err != nil {
		log.Println(err)
	}

	//res.SetConnInfo("network")
	//res.ConnInfo()
	//res.Set("network_uuid")
}

func (o *ResProject) WriteTFRes(res *schema.ResourceData) {
	res.SetId(o.Project.ID.String())
	res.Set("name", o.Project.Name)
	res.Set("ir_group", o.Project.IrGroup)
	res.Set("group_id", o.Project.GroupID.String())
	//res.Set("domain_id", o.Project.DomainID.String())
	//res.Set("state", o.Project.State)
	res.Set("type", o.Project.Type)

	res.Set("description", o.Project.Desc)
	res.Set("default_network", o.Project.DefaultNetwork.String())

	limits := map[string]int{
		"cores":   o.Project.Limits.CoresVcpuCount,
		"ram":     o.Project.Limits.RamGbAmount,
		"storage": o.Project.Limits.StorageGbAmount,
	}

	err := res.Set("limits", limits)
	if err != nil {
		log.Println(err)
	}

	//if o.Project.Networks != nil && len(o.Project.Networks) > 0 {sort.Sort(ByPath(o.Project.Networks))

	networks := make([]map[string]interface{}, 0)
	for _, v := range o.Project.Networks {
		if v.NetworkUUID == o.Project.DefaultNetwork {
			volume := map[string]interface{}{
				"cidr":    v.Cidr,
				"dns":     v.DNSNameservers,
				"dhcp":    v.EnableDhcp,
				"default": true,
				"name":    v.NetworkName,
				"id":      v.NetworkUUID.String(),
			}
			networks = append(networks, volume)
		} else {
			volume := map[string]interface{}{
				"cidr":    v.Cidr,
				"dns":     v.DNSNameservers,
				"dhcp":    v.EnableDhcp,
				"default": false,
				"name":    v.NetworkName,
				"id":      v.NetworkUUID.String(),
			}
			networks = append(networks, volume)
		}
	}

	err = res.Set("network", networks)
	if err != nil {
		log.Println(err)
	}
}

func (o *Project) Serialize() ([]byte, error) {
	//requestMap := map[string]map[string]interface{}{
	//	"project": {
	//		"ir_group":       o.IrGroup,
	//		"type":           o.Type,
	//		"ir_type":        o.IrType,
	//		"virtualization": o.Virtualization,
	//		"name":           o.Name,
	//		"group_id":       o.GroupId,
	//		"datacenter":     o.Datacenter,
	//		"jump_host":      o.JumpHost,
	//		"limits":         o.Limits,
	//		"network":        o.Network,
	//	},
	//}

	//type FullSiProject struct {
	//	Project Project `json:"project"`
	//}

	//requestMap := Project{}
	//
	//	IrGroup:        o.Project.IrGroup,
	//	Type:           o.Project.Type,
	//	IrType:         o.Project.IrType,
	//	Virtualization: o.Project.Virtualization,
	//	Name:           o.Project.Name,
	//	GroupId:        o.Project.GroupID,
	//	Datacenter:     o.Project.Datacenter,
	//	JumpHost:       o.Project.JumpHost,
	//	Limits:         o.Project.Limits,
	//	Networks:       o.Project.Networks,
	//},
	//}

	//requestBytes, err := json.Marshal(FullSiProject{Project: requestMap})
	requestBytes, err := json.Marshal(o)

	if err != nil {
		return nil, err
	}
	return requestBytes, nil
}

func (o *Project) DeserializeOld(responseBytes []byte) error {
	//response := make(map[string]map[string]interface{})
	response := make(map[string]interface{})
	err := json.Unmarshal(responseBytes, &response)
	if err != nil {
		return err
	}

	objMap, ok := response["projects"].([]interface{})
	if !ok {
		return errors.New("no project in response")
	}

	for _, v := range objMap {
		value := v.(map[string]interface{})

		if value["name"].(string) == o.Project.Name {
			o.Project.GroupID = uuid.MustParse(value["group_id"].(string))
			//o.ResId = value["id"].(string)
			//o.DomainId = uuid.MustParse(value["domain_id"].(string))
			//o.GroupId = uuid.MustParse(value["group_id"].(string))
			//o.ResGroupId = value["group_id"].(string)
			//o.StandTypeId = uuid.MustParse(value["stand_type_id"].(string))
			//o.ResStandTypeId = value["stand_type_id"].(string)
			//o.StandType = value["stand_type"].(string)
			//o.Name = value["name"].(string)
			//o.Type = value["type"].(string)
			//o.State = value["state"].(string)
			//o.AppSystemsCi = value["app_systems_ci"].(string)
		}
	}

	//o.Id = uuid.MustParse(objMap["id"].(string))
	//o.ResId = objMap["id"].(string)
	//o.DomainId = uuid.MustParse(objMap["domain_id"].(string))
	//o.GroupId = uuid.MustParse(objMap["group_id"].(string))
	//o.ResGroupId = objMap["group_id"].(string)
	//o.StandTypeId = uuid.MustParse(objMap["stand_type_id"].(string))
	//o.ResStandTypeId = objMap["stand_type_id"].(string)
	//o.StandType = objMap["stand_type"].(string)
	//o.Name = objMap["name"].(string)
	//o.Type = objMap["type"].(string)
	//o.State = objMap["state"].(string)
	//o.AppSystemsCi = objMap["app_systems_ci"].(string)
	return nil
}

func (o *Project) Deserialize(responseBytes []byte) error {
	//response := make(map[string]map[string]interface{})
	//response := make(map[string]interface{})
	//response := Project{}
	//err := json.Unmarshal(responseBytes, &response)
	err := json.Unmarshal(responseBytes, &o)
	if err != nil {
		return err
	}

	//o.Project.ID = response.Project.ID
	//o.Project.Datacenter = response.Project.Datacenter
	//o.Project.DomainID = response.Project.DomainID
	//o.Project.GroupID = response.Project.GroupID

	//o.Project. = value["group_id"].(string)
	//o.Project.StandTypeId = uuid.MustParse(value["stand_type_id"].(string))
	//o.Project.ResStandTypeId = value["stand_type_id"].(string)
	//o.Project.StandType = value["stand_type"].(string)

	//o.Project.Name = response.Project.Name
	//o.Project.Type = response.Project.Type
	//o.Project.State = response.Project.State

	//o.Project.AppSystemsCi = value["app_systems_ci"].(string)

	//objMap, ok := response["projects"].([]interface{})
	//if !ok {
	//	return errors.New("no project in response")
	//}

	//for _, v := range objMap {
	//	value := v.(map[string]interface{})

	//if value["name"].(string) == o.Project.Name {
	//	o.Project.GroupID = uuid.MustParse(value["group_id"].(string))
	//o.ResId = value["id"].(string)
	//o.DomainId = uuid.MustParse(value["domain_id"].(string))
	//o.GroupId = uuid.MustParse(value["group_id"].(string))
	//o.ResGroupId = value["group_id"].(string)
	//o.StandTypeId = uuid.MustParse(value["stand_type_id"].(string))
	//o.ResStandTypeId = value["stand_type_id"].(string)
	//o.StandType = value["stand_type"].(string)
	//o.Name = value["name"].(string)
	//o.Type = value["type"].(string)
	//o.State = value["state"].(string)
	//o.AppSystemsCi = value["app_systems_ci"].(string)
	//}
	//}

	//o.Id = uuid.MustParse(objMap["id"].(string))
	//o.ResId = objMap["id"].(string)
	//o.DomainId = uuid.MustParse(objMap["domain_id"].(string))
	//o.GroupId = uuid.MustParse(objMap["group_id"].(string))
	//o.ResGroupId = objMap["group_id"].(string)
	//o.StandTypeId = uuid.MustParse(objMap["stand_type_id"].(string))
	//o.ResStandTypeId = objMap["stand_type_id"].(string)
	//o.StandType = objMap["stand_type"].(string)
	//o.Name = objMap["name"].(string)
	//o.Type = objMap["type"].(string)
	//o.State = objMap["state"].(string)
	//o.AppSystemsCi = objMap["app_systems_ci"].(string)
	return nil
}

func (o *ResProject) DeserializeRead(responseBytes []byte) error {

	err := json.Unmarshal(responseBytes, &o)
	if err != nil {
		return err
	}

	return nil
}

func (o *Project) ParseIdFromCreateResponse(data []byte) error {
	response := make(map[string]map[string]interface{})
	err := json.Unmarshal(data, &response)
	if err != nil {
		return err
	}
	objMap, ok := response["project"]
	if !ok {
		return errors.New("no project in response")
	}

	//o2 := &Project{}
	o.Project.ID = uuid.MustParse(objMap["id"].(string))
	o.Project.GroupID = uuid.MustParse(objMap["group_id"].(string))

	return nil
}

func (o *Project) CreateDI(data []byte) ([]byte, error) {
	return Api.NewRequestCreate("/v2/projects", data)
}

func (o *Project) ReadDI() ([]byte, error) {
	return Api.NewRequestRead(fmt.Sprintf("projects/%s", o.Project.ID))
	//return Api.NewRequestRead(fmt.Sprintf("projects?group_ids=%s", o.GroupId))
}

func (o *ResProject) ReadDIRes() ([]byte, error) {
	return Api.NewRequestRead(fmt.Sprintf("projects/%s", o.Project.ID))
	//return Api.NewRequestRead(fmt.Sprintf("projects?group_ids=%s", o.GroupId))
}

func (o *Project) UpdateDI(data []byte) ([]byte, error) {
	return Api.NewRequestUpdate(fmt.Sprintf("projects/%s", o.Project.ID), data)
}

func (o *Project) UpdateProjectName(data []byte) ([]byte, error) {
	return Api.NewRequestUpdate(fmt.Sprintf("projects/%s", o.Project.ID), data)
}

func (o *Project) UpdateProjectDesc(data []byte) ([]byte, error) {
	return Api.NewRequestUpdate(fmt.Sprintf("projects/%s", o.Project.ID), data)
}

func (o *Project) UpdateProjectLimits(data []byte) ([]byte, error) {
	return Api.NewRequestUpdate(fmt.Sprintf("/v2/projects/%s/quota", o.Project.ID), data)
}

func (o *Project) DeleteDI() error {
	return Api.NewRequestDelete(fmt.Sprintf("projects/%s", o.Project.ID), nil, 204)
}

func (o *Project) DeleteNetwork(NetworkUuid string) error {
	return Api.NewRequestDelete(fmt.Sprintf("projects/%s/networks/%s", o.Project.ID, NetworkUuid), nil, 200)
}

func (o *Project) ReadAll() ([]byte, error) {
	return Api.NewRequestRead("projects/")
}

func (o *Project) StateChange(res *schema.ResourceData) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Timeout:      res.Timeout(schema.TimeoutCreate),
		PollInterval: 15 * time.Second,
		Pending:      []string{"Creating", "Pending"},
		Target:       []string{"Running", "Damaged"},
		Refresh: func() (interface{}, string, error) {

			responseBytes, err := o.ReadDI()
			if err != nil {
				return nil, "error", err
			}

			err = o.Deserialize(responseBytes)
			if err != nil {
				return nil, "error", err
			}

			log.Printf("[DEBUG] Refresh state for [%s]: state: %s", o.Project.ID.String(), o.Project.State)
			// write to TF state
			o.WriteTF(res)

			if o.Project.State == "ready" {
				return o, "Running", nil
			}
			if o.Project.State == "damaged" {
				return o, "Damaged", nil
			}
			if o.Project.State == "pending" {
				return o, "Pending", nil
			}
			return o, "Creating", nil
		},
	}
}

func (o *ResProject) StateChangeNetwork(res *schema.ResourceData, networkName string, isDefault bool) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Timeout:      res.Timeout(schema.TimeoutCreate),
		PollInterval: 15 * time.Second,
		Pending:      []string{"Creating", "Pending"},
		Target:       []string{"Running", "Damaged"},
		Refresh: func() (interface{}, string, error) {

			responseBytes, err := o.ReadDIRes()
			if err != nil {
				return nil, "error", err
			}

			err = o.DeserializeRead(responseBytes)
			if err != nil {
				return nil, "error", err
			}

			log.Printf("[DEBUG] Refresh state for [%s]: state: %s", o.Project.ID, o.Project.State)
			// write to TF state
			//o.WriteTFRes(res)

			for _, net := range o.Project.Networks {
				if net.NetworkName == networkName {
					if net.Status == "ready" {
						if isDefault {
							err := o.SetDefaultNetwork(net.NetworkUUID.String())
							if err != nil {
								return o, "Running", err
							}
						}
						return o, "Running", nil

					} else if net.Status == "pending" {
						return o, "Pending", nil
					}
				}
			}

			return o, "Creating", nil
		},
	}
}

func (o *Project) OnSerialize(map[string]interface{}, *Server) map[string]interface{} {
	return nil
}
func (o *Project) OnDeserialize(map[string]interface{}, *Server) {}
func (o *Project) Urls(string) string {
	return ""
}
func (o *Project) OnReadTF(*schema.ResourceData, *Server)  {}
func (o *Project) OnWriteTF(*schema.ResourceData, *Server) {}

func (o *Project) ToHCLOutput() []byte {
	dataRoot := &HCLOutputRoot{
		Resources: &HCLOutput{
			ResName: fmt.Sprintf(
				"%s_id",
				//o.ResType,
				o.Project.IrType,
			),
			Value: fmt.Sprintf(
				"%s.%s.id",
				//o.ResType,
				o.Project.IrType,
				//o.ResName,
				o.Project.IrType,
			),
		},
	}
	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(dataRoot, f.Body())
	return utils.Regexp(f.Bytes())
}

//func (o *Project) HostVars(server *Server) map[string]interface{} {
//	return nil
//}

func (o *Project) GetGroup() uuid.UUID {
	return o.Project.GroupID
}

func (o *Project) ToHCL(server *Server) ([]byte, error) {
	//o.ResType = o.GetType()
	o.Project.IrType = o.GetType()
	//o.ResName = utils.Reformat(o.Name)
	o.Project.IrType = utils.Reformat(o.Project.Name)
	type HCLServerRoot struct {
		Resources *Project `hcl:"resource,block"`
	}
	root := &HCLServerRoot{Resources: o}
	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(root, f.Body())
	// return utils.Regexp(f.Bytes())
	return f.Bytes(), nil
}

func (o *Project) HCLAppParams() *HCLAppParams {
	return nil
}

func (o *Project) HCLVolumes() []*HCLVolume {
	return nil
}
