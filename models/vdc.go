package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"log"
	"time"
	//	"github.com/hashicorp/hcl/v2/gohcl"
	//	"github.com/hashicorp/hcl/v2/hclwrite"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Vdc struct {
	ID             uuid.UUID `json:"id,omitempty"`
	Name           string    `json:"name"`
	State          string    `json:"state"`
	Type           string    `json:"type"`
	IrGroup        string    `json:"ir_group"`
	IrType         string    `json:"ir_type"`
	Virtualization string    `json:"virtualization"`
	Datacenter     string    `json:"datacenter"`
	DefaultNetwork uuid.UUID `json:"default_network"`
	Limits         struct {
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
	GroupName      string          `json:"group_name,omitempty"`
	DomainID       uuid.UUID       `json:"domain_id"`
	GroupID        uuid.UUID       `json:"group_id"`
	JumpHost       bool            `json:"jump_host"`
	Desc           string          `json:"desc"`
	SecurityGroups []SecurityGroup `json:"security_groups"`
}

type ResVdc struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	State          string    `json:"state"`
	Type           string    `json:"type"`
	IrGroup        string    `json:"ir_group"`
	IrType         string    `json:"ir_type"`
	Virtualization string    `json:"virtualization"`
	Datacenter     string    `json:"datacenter"`
	DefaultNetwork uuid.UUID `json:"default_network"`
	Limits         struct {
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
	DomainName     string    `json:"domain_name"`
	GroupName      string    `json:"group_name"`
	DomainID       uuid.UUID `json:"domain_id"`
	GroupID        uuid.UUID `json:"group_id"`
	IsProm         bool      `json:"is_prom"`
	JumpHost       bool      `json:"jump_host"`
	Desc           string    `json:"desc"`
	SecurityGroups []struct {
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

func (o *Vdc) AddNetwork(ctx context.Context, res *schema.ResourceData, additionalNets []interface{}) diag.Diagnostics {
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

		resBody, err := Api.NewRequestCreate(fmt.Sprintf("projects/%s/networks", o.ID), result)

		if err != nil {
			return diag.FromErr(err)
		}

		deserializeResBody := ResVdc{}
		err = json.Unmarshal(resBody, &deserializeResBody)
		if err != nil {
			return diag.FromErr(err)
		}

		_, err = deserializeResBody.StateChangeNetwork(res, body.Network.NetworkName, body.Network.IsDefault).WaitForStateContext(ctx)
	}

	return diag.Diagnostics{}
}

func (o *Vdc) GetVdcQuota() ([]byte, error) {
	body, err := Api.NewRequestRead(fmt.Sprintf("/v2/projects/%s/quota?group_id=%s", o.ID, o.GroupID))

	if err != nil {
		return nil, err
	}
	return body, nil
}

func (o *ResVdc) GetVdcQuota() ([]byte, error) {
	body, err := Api.NewRequestRead(fmt.Sprintf("/v2/projects/%s/quota?group_id=%s", o.ID, o.GroupID))

	if err != nil {
		return nil, err
	}
	return body, nil
}

func (o *Vdc) SetDefaultNetwork(networkUuid string) error {
	body, err := json.Marshal(map[string]string{"network_uuid": networkUuid})
	if err != nil {
		return err
	}

	_, err = Api.NewRequestUpdate(fmt.Sprintf("projects/%s/networks/set_default", o.ID), body)

	if err != nil {
		return err
	}
	return nil
}

func (o *ResVdc) SetDefaultNetwork(networkUuid string) error {
	body, err := json.Marshal(map[string]string{"network_uuid": networkUuid})
	if err != nil {
		return err
	}
	_, err = Api.NewRequestUpdate(fmt.Sprintf("projects/%s/networks/set_default", o.ID), body)

	if err != nil {
		return err
	}
	return nil
}

func (o *Vdc) ReadTF(res *schema.ResourceData) diag.Diagnostics {

	if res.Id() != "" {
		o.ID = uuid.MustParse(res.Id())
	}

	o.IrGroup = res.Get("ir_group").(string)
	o.Type = res.Get("type").(string)
	o.IrType = res.Get("ir_type").(string)
	o.Virtualization = res.Get("virtualization").(string)
	o.Name = res.Get("name").(string)
	o.GroupID = uuid.MustParse(res.Get("group_id").(string))
	//o.ID = uuid.MustParse(res.Id())
	o.Datacenter = res.Get("datacenter").(string)
	o.Desc = res.Get("description").(string)
	o.JumpHost = res.Get("jump_host").(bool)

	limits, ok := res.GetOk("limits")
	if ok {
		limits := limits.(map[string]interface{})
		o.Limits.CoresVcpuCount = limits["cores"].(int)
		o.Limits.RamGbAmount = limits["ram"].(int)
		o.Limits.StorageGbAmount = limits["storage"].(int)
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
			o.Networks.NetworkName = v["name"].(string)
			o.Networks.Cidr = v["cidr"].(string)
			o.Networks.EnableDhcp = v["dhcp"].(bool)
			o.Networks.IsDefault = true
			var dnsNameServers []string
			for _, dnsIp := range v["dns"].(*schema.Set).List() {
				dnsNameServers = append(dnsNameServers, dnsIp.(string))
			}
			o.Networks.DNSNameservers = dnsNameServers
		}
	}
	return diag.Diagnostics{}
}

func (o *ResVdc) ReadTFRes(res *schema.ResourceData) diag.Diagnostics {

	if res.Id() != "" {
		o.ID = uuid.MustParse(res.Id())
	}

	o.IrGroup = res.Get("ir_group").(string)
	o.Type = res.Get("type").(string)
	o.IrType = res.Get("ir_type").(string)
	o.Virtualization = res.Get("virtualization").(string)
	o.Name = res.Get("name").(string)
	o.GroupID = uuid.MustParse(res.Get("group_id").(string))
	//o.ID = uuid.MustParse(res.Id())
	o.Datacenter = res.Get("datacenter").(string)
	o.Desc = res.Get("description").(string)
	o.JumpHost = res.Get("jump_host").(bool)

	limits, ok := res.GetOk("limits")
	if ok {
		limits := limits.(map[string]interface{})
		o.Limits.CoresVcpuCount = limits["cores"].(int)
		o.Limits.RamGbAmount = limits["ram"].(int)
		o.Limits.StorageGbAmount = limits["storage"].(int)
	}

	return diag.Diagnostics{}
}

func (o *Vdc) WriteTF(res *schema.ResourceData) {
	res.SetId(o.ID.String())

	res.Set("datacenter", o.Datacenter)
	res.Set("ir_type", o.IrType)
	res.Set("description", o.Desc)
	res.Set("group_id", o.GroupID.String())
	res.Set("jump_host", o.JumpHost)
	res.Set("name", o.Name)
	res.Set("virtualization", o.Virtualization)

	limits := map[string]int{
		"cores":   o.Limits.CoresVcpuCount,
		"ram":     o.Limits.RamGbAmount,
		"storage": o.Limits.StorageGbAmount,
	}
	err := res.Set("limits", limits)
	if err != nil {
		log.Println(err)
	}

	//res.SetConnInfo("network")
	//res.ConnInfo()
	//res.Set("network_uuid")
}

func (o *ResVdc) WriteTFRes(res *schema.ResourceData) {
	res.SetId(o.ID.String())
	res.Set("name", o.Name)
	res.Set("ir_group", o.IrGroup)
	res.Set("group_id", o.GroupID.String())
	//res.Set("domain_id", o.Vdc.DomainID.String())
	//res.Set("state", o.Vdc.State)
	res.Set("type", o.Type)

	res.Set("description", o.Desc)
	res.Set("default_network", o.DefaultNetwork.String())
	limits := map[string]int{
		"cores":   o.Limits.CoresVcpuCount,
		"ram":     o.Limits.RamGbAmount,
		"storage": o.Limits.StorageGbAmount,
	}

	err := res.Set("limits", limits)
	if err != nil {
		log.Println(err)
	}

	//if o.Vdc.Networks != nil && len(o.Vdc.Networks) > 0 {sort.Sort(ByPath(o.Vdc.Networks))

	networks := make([]map[string]interface{}, 0)
	for _, v := range o.Networks {
		if v.NetworkUUID == o.DefaultNetwork {
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

func (o *Vdc) Serialize() ([]byte, error) {
	requestBytes, err := json.Marshal(map[string]*Vdc{"project": o})

	if err != nil {
		return nil, err
	}
	return requestBytes, nil
}

func (o *Vdc) Deserialize(responseBytes []byte) error {
	var response map[string]Vdc
	err := json.Unmarshal(responseBytes, &response)
	*o = response["project"]
	if err != nil {
		return err
	}
	return nil
}

func (o *ResVdc) DeserializeRead(responseBytes []byte) error {
	var response map[string]ResVdc
	err := json.Unmarshal(responseBytes, &response)
	*o = response["project"]
	if err != nil {
		return err
	}

	return nil
}

func (o *Vdc) ParseIdFromCreateResponse(data []byte) error {
	response := make(map[string]map[string]interface{})
	err := json.Unmarshal(data, &response)
	if err != nil {
		return err
	}
	objMap, ok := response["project"]
	if !ok {
		return errors.New("no project in response")
	}

	//o2 := &Vdc{}
	o.ID = uuid.MustParse(objMap["id"].(string))
	o.GroupID = uuid.MustParse(objMap["group_id"].(string))

	return nil
}

func (o *Vdc) CreateDI(data []byte) ([]byte, error) {
	return Api.NewRequestCreate("/v2/projects", data)
}

func (o *Vdc) ReadDI() ([]byte, error) {
	return Api.NewRequestRead(fmt.Sprintf("projects/%s", o.ID))
	//return Api.NewRequestRead(fmt.Sprintf("projects?group_ids=%s", o.GroupId))
}

func (o *ResVdc) ReadDIRes() ([]byte, error) {
	return Api.NewRequestRead(fmt.Sprintf("projects/%s", o.ID))
	//return Api.NewRequestRead(fmt.Sprintf("projects?group_ids=%s", o.GroupId))
}

func (o *Vdc) UpdateVdcName(data []byte) ([]byte, error) {
	return Api.NewRequestUpdate(fmt.Sprintf("projects/%s", o.ID), data)
}

func (o *Vdc) UpdateVdcDesc(data []byte) ([]byte, error) {
	return Api.NewRequestUpdate(fmt.Sprintf("projects/%s", o.ID), data)
}

func (o *Vdc) UpdateVdcLimits(data []byte) ([]byte, error) {
	return Api.NewRequestUpdate(fmt.Sprintf("/v2/projects/%s/quota", o.ID), data)
}

func (o *Vdc) DeleteDI() error {
	return Api.NewRequestDelete(fmt.Sprintf("projects/%s", o.ID), nil, 204)
}

func (o *Vdc) DeleteNetwork(networkId string) error {
	return Api.NewRequestDelete(fmt.Sprintf("projects/%s/networks/%s", o.ID, networkId), nil, 200)
}

func (o *Vdc) ReadAll() ([]byte, error) {
	return Api.NewRequestRead("projects/")
}

func (o *Vdc) StateChange(res *schema.ResourceData) *resource.StateChangeConf {
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

			log.Printf("[DEBUG] Refresh state for [%s]: state: %s", o.ID.String(), o.State)
			// write to TF state
			o.WriteTF(res)

			if o.State == "ready" {
				return o, "Running", nil
			}
			if o.State == "damaged" {
				return o, "Damaged", nil
			}
			if o.State == "pending" {
				return o, "Pending", nil
			}
			return o, "Creating", nil
		},
	}
}

func (o *ResVdc) StateChangeNetwork(res *schema.ResourceData, networkName string, isDefault bool) *resource.StateChangeConf {
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

			log.Printf("[DEBUG] Refresh state for [%s]: state: %s", o.ID, o.State)
			// write to TF state
			//o.WriteTFRes(res)

			for _, net := range o.Networks {
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
