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

type Project struct {
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

type ResProject struct {
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

		resBody, err := Api.NewRequestCreate(fmt.Sprintf("projects/%s/networks", o.ID), result)

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
	body, err := Api.NewRequestRead(fmt.Sprintf("/v2/projects/%s/quota?group_id=%s", o.ID, o.GroupID))

	if err != nil {
		return nil, err
	}
	return body, nil
}

func (o *ResProject) GetProjectQuota() ([]byte, error) {
	body, err := Api.NewRequestRead(fmt.Sprintf("/v2/projects/%s/quota?group_id=%s", o.ID, o.GroupID))

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

	_, err = Api.NewRequestUpdate(fmt.Sprintf("projects/%s/networks/set_default", o.ID), body)

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
	_, err = Api.NewRequestUpdate(fmt.Sprintf("projects/%s/networks/set_default", o.ID), body)

	if err != nil {
		return err
	}
	return nil
}

/*
	func (o *Project) GetType() string {
		return "si_vdc"
	}

	func (o *Project) NewObj() DIDataResource {
		return &Project{}
	}

	func (o *Project) GetId() string {
		return o.ID.String()
	}

	func (o *Project) GetDomainId() uuid.UUID {
		return o.DomainID
	}

	func (o *Project) GetResType() string {
		return "si_group"
	}

	func (o *Project) GetResName() string {
		return o.Name
	}

	func (o *Project) GetOutput() (string, string) {
		//return o.ResOutputName, o.ResOutputValue
		return "", ""
	}

	func (o *Project) SetResFields() {
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

}

	func (o *Project) DeserializeAll(responseBytes []byte) ([]DIDataResource, error) {
		m := make(map[string][]*Project)
		err := json.Unmarshal(responseBytes, &m)
		if err != nil {
			return nil, err
		}

		m2 := make([]DIDataResource, len(m["groups"]))
		for k, v := range m["groups"] {
			m2[k] = v
		}
		return m2, nil
	}

	func (o *Project) NewObj() DIResource {
		return &Project{}
	}
*/

func (o *Project) ReadTF(res *schema.ResourceData) diag.Diagnostics {

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

func (o *ResProject) ReadTFRes(res *schema.ResourceData) diag.Diagnostics {

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

func (o *Project) WriteTF(res *schema.ResourceData) {
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

func (o *ResProject) WriteTFRes(res *schema.ResourceData) {
	res.SetId(o.ID.String())
	res.Set("name", o.Name)
	res.Set("ir_group", o.IrGroup)
	res.Set("group_id", o.GroupID.String())
	//res.Set("domain_id", o.Project.DomainID.String())
	//res.Set("state", o.Project.State)
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

	//if o.Project.Networks != nil && len(o.Project.Networks) > 0 {sort.Sort(ByPath(o.Project.Networks))

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

func (o *Project) Serialize() ([]byte, error) {
	requestBytes, err := json.Marshal(map[string]*Project{"project": o})

	if err != nil {
		return nil, err
	}
	return requestBytes, nil
}

/*
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

		if value["name"].(string) == o.Name {
			o.GroupID = uuid.MustParse(value["group_id"].(string))
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
*/

func (o *Project) Deserialize(responseBytes []byte) error {
	var response map[string]Project
	err := json.Unmarshal(responseBytes, &response)
	*o = response["project"]
	if err != nil {
		return err
	}
	return nil
}

func (o *ResProject) DeserializeRead(responseBytes []byte) error {
	var response map[string]ResProject
	err := json.Unmarshal(responseBytes, &response)
	*o = response["project"]
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
	o.ID = uuid.MustParse(objMap["id"].(string))
	o.GroupID = uuid.MustParse(objMap["group_id"].(string))

	return nil
}

func (o *Project) CreateDI(data []byte) ([]byte, error) {
	return Api.NewRequestCreate("/v2/projects", data)
}

func (o *Project) ReadDI() ([]byte, error) {
	return Api.NewRequestRead(fmt.Sprintf("projects/%s", o.ID))
	//return Api.NewRequestRead(fmt.Sprintf("projects?group_ids=%s", o.GroupId))
}

func (o *ResProject) ReadDIRes() ([]byte, error) {
	return Api.NewRequestRead(fmt.Sprintf("projects/%s", o.ID))
	//return Api.NewRequestRead(fmt.Sprintf("projects?group_ids=%s", o.GroupId))
}

//func (o *Project) UpdateDI(data []byte) ([]byte, error) {
//	return Api.NewRequestUpdate(fmt.Sprintf("projects/%s", o.ID), data)
//}

func (o *Project) UpdateProjectName(data []byte) ([]byte, error) {
	return Api.NewRequestUpdate(fmt.Sprintf("projects/%s", o.ID), data)
}

func (o *Project) UpdateProjectDesc(data []byte) ([]byte, error) {
	return Api.NewRequestUpdate(fmt.Sprintf("projects/%s", o.ID), data)
}

func (o *Project) UpdateProjectLimits(data []byte) ([]byte, error) {
	return Api.NewRequestUpdate(fmt.Sprintf("/v2/projects/%s/quota", o.ID), data)
}

func (o *Project) DeleteDI() error {
	return Api.NewRequestDelete(fmt.Sprintf("projects/%s", o.ID), nil, 204)
}

func (o *Project) DeleteNetwork(networkId string) error {
	return Api.NewRequestDelete(fmt.Sprintf("projects/%s/networks/%s", o.ID, networkId), nil, 200)
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

/*
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
				o.IrType,
			),
			Value: fmt.Sprintf(
				"%s.%s.id",
				//o.ResType,
				o.IrType,
				//o.ResName,
				o.IrType,
			),
		},
	}
	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(dataRoot, f.Body())
	return utils.Regexp(f.Bytes())
}

func (o *Project) HostVars(server *Server) map[string]interface{} {
	return nil
}

func (o *Project) GetGroup() uuid.UUID {
	return o.GroupID
}


func (o *Project) ToHCL(server *Server) ([]byte, error) {
	//o.ResType = o.GetType()
	o.IrType = o.GetType()
	//o.ResName = utils.Reformat(o.Name)
	o.IrType = utils.Reformat(o.Name)
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
*/
