package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"log"
	"stash.sigma.sbrf.ru/sddevops/terraform-provider-di/utils"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SIProject struct {
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
		OpenstackProjectID interface{} `json:"openstack_project_id"`
		DefaultNetwork     interface{} `json:"default_network"`
		Limits             struct {
			CoresVcpuCount  int `json:"cores_vcpu_count"`
			RAMGbAmount     int `json:"ram_gb_amount"`
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

func (o *SIProject) GetType() string {
	return "di_siproject"
}

//func (o *SIProject) NewObj() DIDataResource {
//	return &SIProject{}
//}

func (o *SIProject) GetId() string {
	return o.Project.ID.String()
}

func (o *SIProject) GetDomainId() uuid.UUID {
	return o.Project.DomainID
}

func (o *SIProject) GetResType() string {
	return "di_group"
}

func (o *SIProject) GetResName() string {
	return o.Project.Name
}

func (o *SIProject) GetOutput() (string, string) {
	//return o.ResOutputName, o.ResOutputValue
	return "", ""
}

func (o *SIProject) SetResFields() {
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

//func (o *SIProject) DeserializeAll(responseBytes []byte) ([]DIDataResource, error) {
//	m := make(map[string][]*SIProject)
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

//func (o *SIProject) NewObj() DIResource {
//	return &SIProject{}
//}

func (o *SIProject) ReadTF(res *schema.ResourceData) diag.Diagnostics {

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
	o.Project.Desc = res.Get("desc").(string)
	//o.JumpHost = res.Get("jump_host")

	if res.Get("jump_host") == "true" {
		o.Project.JumpHost = true
	} else {
		o.Project.JumpHost = false
	}

	net, ok := res.GetOk("network")
	log.Println("NETWORK", net)
	log.Println("NETWORK", net.(*schema.Set).List())
	log.Println("NETWORK", net.(*schema.Set).Len())

	limits := res.Get("limits")

	log.Println("LLOK", ok)
	log.Println("LL", limits)
	log.Println("LL", len(limits.(*schema.Set).List()))
	log.Println("LL", limits.(*schema.Set).Len())

	//if ok {
	//	if limits.(*schema.Set).Len() > 1 {
	//		res.Get("limits").(*schema.Set).Len()
	//		return diag.Errorf("Limits set should not be more than one")
	//	}
	//}

	if ok {
		limitsSet := limits.(*schema.Set)

		for _, v := range limitsSet.List() {
			values := v.(map[string]interface{})

			CoresVcpuCount, err := strconv.Atoi(values["cores_vcpu_count"].(string))
			if err != nil {
				panic(err)
			}
			RamGbAmount, err := strconv.Atoi(values["ram_gb_amount"].(string))
			if err != nil {
				panic(err)
			}
			StorageGbAmount, err := strconv.Atoi(values["storage_gb_amount"].(string))
			if err != nil {
				panic(err)
			}

			o.Project.Limits.CoresVcpuCount = CoresVcpuCount
			o.Project.Limits.RAMGbAmount = RamGbAmount
			o.Project.Limits.StorageGbAmount = StorageGbAmount
		}
	}

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

	network, ok := res.GetOk("network")

	if ok {
		networkSet := network.(*schema.Set).List()

		for _, v := range networkSet {
			if v.(map[string]interface{})["is_default"].(bool) {
				o.Project.Networks.NetworkName = v.(map[string]interface{})["network_name"].(string)
				o.Project.Networks.Cidr = v.(map[string]interface{})["cidr"].(string)
				o.Project.Networks.EnableDhcp = v.(map[string]interface{})["enable_dhcp"].(bool)
				//o.Project.Networks.IsDefault = v.(map[string]interface{})["is_default"].(bool)
				o.Project.Networks.IsDefault = true
				//o.Project.Networks.NetworkUuid = v.(map[string]interface{})["network_uuid"].(uuid.UUID)
				//log.Printf("#@ %v, %T\n", v.(map[string]interface{})["network_uuid"], v.(map[string]interface{})["network_uuid"])

				var dnsNameServers = []string{}
				for _, dnsIp := range v.(map[string]interface{})["dns_nameservers"].(*schema.Set).List() {
					dnsNameServers = append(dnsNameServers, dnsIp.(string))
				}
				o.Project.Networks.DNSNameservers = dnsNameServers
			}
		}
	}

	return diag.Diagnostics{}
}

func (o *SIProject) WriteTF(res *schema.ResourceData) {
	log.Println("@@@", o.Project.Networks)
	res.SetId(o.Project.ID.String())
	res.Set("ir_group", o.Project.IrGroup)
	//res.Set("stand_type_id", o.StandTypeId.String())
	res.Set("group_id", o.Project.GroupID.String())
	res.Set("domain_id", o.Project.GroupID.String())
	//res.Set("app_systems_ci", o.AppSystemsCi)
	//res.Set("stand_type", o.StandType)
	//res.Set("state", o.State)
	res.Set("type", o.Project.Type)
	//res.Set("network", o.Project.Networks)

	//if o.Project.Networks != nil && len(o.Project.Networks) > 0 {
	//sort.Sort(ByPath(o.Project.Networks))

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

	//}

	//res.SetConnInfo("network")
	//res.ConnInfo()
	//res.
	//log.Println("##NS", res.Get("network"))

	//res.Set("network_uuid")
}

//{
//    "project": {
//        "ir_group":"vdc",
//        "type":"vdc",
//        "ir_type":"vdc_openstack",
//        "virtualization":"openstack",
//        "name":"test-di-project1", // requared false
//        "group_id":"52ffd9f6-fbc0-4ddc-bf99-b092c6d0351a",
//        "datacenter":"PD23R3PROM",
//        "jump_host":false,
//        "limits": { // requared false
//            "cores_vcpu_count":100,
//            "ram_gb_amount":10000,
//            "storage_gb_amount":1000
//        },
//        "network": {
//            "network_name":"internal-network",
//            "cidr":"172.31.0.0/20",
//            "dns_nameservers":["8.8.8.8","8.8.4.4"],
//            "enable_dhcp":true
//        }
//    }
//}

func (o *SIProject) Serialize() ([]byte, error) {
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
	//	Project SIProject `json:"project"`
	//}

	//requestMap := SIProject{}
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

func (o *SIProject) DeserializeOld(responseBytes []byte) error {
	//log.Println("!!!bytes", responseBytes)
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
			//log.Println("@@@", value["name"].(string))
			//log.Println("@@@", reflect.TypeOf(value["name"].(string)))
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

func (o *SIProject) Deserialize(responseBytes []byte) error {

	//response := make(map[string]map[string]interface{})
	//response := make(map[string]interface{})
	response := SIProject{}
	err := json.Unmarshal(responseBytes, &response)
	if err != nil {
		return err
	}

	//log.Println("!!!!!!!DES", response)
	//log.Println("!!!!!!!IDDD", response.Project.ID)

	o.Project.ID = response.Project.ID
	o.Project.DomainID = response.Project.DomainID
	o.Project.GroupID = response.Project.GroupID
	//o.Project. = value["group_id"].(string)
	//o.Project.StandTypeId = uuid.MustParse(value["stand_type_id"].(string))
	//o.Project.ResStandTypeId = value["stand_type_id"].(string)
	//o.Project.StandType = value["stand_type"].(string)
	o.Project.Name = response.Project.Name
	o.Project.Type = response.Project.Type
	o.Project.State = response.Project.State
	//o.Project.AppSystemsCi = value["app_systems_ci"].(string)

	//objMap, ok := response["projects"].([]interface{})
	//if !ok {
	//	return errors.New("no project in response")
	//}

	//for _, v := range objMap {
	//	value := v.(map[string]interface{})

	//if value["name"].(string) == o.Project.Name {
	//	log.Println("@@@", value["name"].(string))
	//	log.Println("@@@", reflect.TypeOf(value["name"].(string)))
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

func (o *SIProject) ParseIdFromCreateResponse(data []byte) error {
	response := make(map[string]map[string]interface{})
	//log.Println("DATA", data)
	err := json.Unmarshal(data, &response)
	if err != nil {
		return err
	}
	objMap, ok := response["project"]
	if !ok {
		return errors.New("no project in response")
	}

	//o2 := &SIProject{}
	o.Project.ID = uuid.MustParse(objMap["id"].(string))
	o.Project.GroupID = uuid.MustParse(objMap["group_id"].(string))
	o.Project.Networks.NetworkUuid = uuid.MustParse(objMap["networks"].([]interface{})[0].(map[string]interface{})["network_uuid"].(string))
	log.Println("NUUID", o.Project.Networks.NetworkUuid)

	return nil
}

func (o *SIProject) CreateDI(data []byte) ([]byte, error) {
	//log.Println("###data", pp.Sprintln(string(data)))
	return Api.NewRequestCreate("/v2/projects", data)
}

func (o *SIProject) ReadDI() ([]byte, error) {
	//return Api.NewRequestRead(fmt.Sprintf("projects/%s", o.Id))
	//log.Println("###ID", o.Project.ID)
	//log.Println("###GROUPID", o.Project.GroupID)

	//log.Println("###ProjectID", o.Project.ID)
	return Api.NewRequestRead(fmt.Sprintf("projects/%s", o.Project.ID))
	//return Api.NewRequestRead(fmt.Sprintf("projects?group_ids=%s", o.GroupId))
}

func (o *SIProject) UpdateDI(data []byte) ([]byte, error) {
	return Api.NewRequestUpdate(fmt.Sprintf("projects/%s", o.Project.ID), data)
}

func (o *SIProject) UpdateSIProjectName(data []byte) ([]byte, error) {
	return Api.NewRequestUpdate(fmt.Sprintf("projects/%s", o.Project.ID), data)
}

func (o *SIProject) UpdateSIProjectDesc(data []byte) ([]byte, error) {
	return Api.NewRequestUpdate(fmt.Sprintf("projects/%s", o.Project.ID), data)
}

func (o *SIProject) UpdateSIProjectLimits(data []byte) ([]byte, error) {
	return Api.NewRequestUpdate(fmt.Sprintf("/v2/projects/%s/quota", o.Project.ID), data)
}

func (o *SIProject) DeleteDI() error {
	return Api.NewRequestDelete(fmt.Sprintf("projects/%s", o.Project.ID), nil)
}

func (o *SIProject) ReadAll() ([]byte, error) {
	return Api.NewRequestRead("projects/")
}

func (o *SIProject) StateChange(res *schema.ResourceData) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Timeout:      res.Timeout(schema.TimeoutCreate),
		PollInterval: 15 * time.Second,
		Pending:      []string{"Creating"},
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
			return o, "Creating", nil
		},
	}
}

/*
	func (o *SIProject) DeserializeAll(responseBytes []byte) ([]*SIProject, error) {
		response := make(map[string]interface{})
		err := json.Unmarshal(responseBytes, &response)
		if err != nil {
			return nil, err
		}
		objList := make([]*SIProject, 0)
		objResNamesList := make([]string, 0)
		counter := make(map[string][]*SIProject)
		for _, v := range response["projects"].([]interface{}) {
			objMap := v.(map[string]interface{})
			obj := &SIProject{
				Id:           uuid.MustParse(objMap["id"].(string)),
				GroupId:      uuid.MustParse(objMap["group_id"].(string)),
				DomainId:     uuid.MustParse(objMap["domain_id"].(string)),
				StandTypeId:  uuid.MustParse(objMap["stand_type_id"].(string)),
				Name:         objMap["name"].(string),
				StandType:    objMap["stand_type"].(string),
				State:        objMap["state"].(string),
				Type:         objMap["type"].(string),
				AppSystemsCi: objMap["app_systems_ci"].(string),
				ResId:        objMap["id"].(string),
				ResType:      "di_siproject",
				ResName:      utils.Reformat(objMap["name"].(string)),
				// ResDomainId: objMap["domain_id"].(string),
				ResGroupIdUUID:     objMap["group_id"].(string),
				ResStandTypeIdUUID: objMap["stand_type_id"].(string),
				// ResStandTypeId:     objMap["stand_type_id"].(string),
			}
			objList = append(objList, obj)
			objResNamesList = append(objResNamesList, obj.ResName)
			counter[obj.ResName] = append(counter[obj.ResName], obj)
		}
		for _, arr := range counter {
			if len(arr) > 1 {
				var c int
				for _, v := range arr {
					c++
					v.ResName = fmt.Sprintf("%s-%d", v.ResName, c)
				}
			}
		}
		return objList, nil
	}
*/
func (o *SIProject) OnSerialize(map[string]interface{}, *Server) map[string]interface{} {
	return nil
}
func (o *SIProject) OnDeserialize(map[string]interface{}, *Server) {}
func (o *SIProject) Urls(string) string {
	return ""
}
func (o *SIProject) OnReadTF(*schema.ResourceData, *Server)  {}
func (o *SIProject) OnWriteTF(*schema.ResourceData, *Server) {}

func (o *SIProject) ToHCLOutput() []byte {
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

func (o *SIProject) HostVars(server *Server) map[string]interface{} {
	return nil
}

func (o *SIProject) GetGroup() string {
	return ""
}

func (o *SIProject) ToHCL(server *Server) ([]byte, error) {
	//o.ResType = o.GetType()
	o.Project.IrType = o.GetType()
	//o.ResName = utils.Reformat(o.Name)
	o.Project.IrType = utils.Reformat(o.Project.Name)
	type HCLServerRoot struct {
		Resources *SIProject `hcl:"resource,block"`
	}
	root := &HCLServerRoot{Resources: o}
	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(root, f.Body())
	// return utils.Regexp(f.Bytes())
	return f.Bytes(), nil
}

func (o *SIProject) HCLAppParams() *HCLAppParams {
	return nil
}

func (o *SIProject) HCLVolumes() []*HCLVolume {
	return nil
}
