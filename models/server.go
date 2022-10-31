package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"base.sw.sbc.space/pid/terraform-provider-si/utils"
	"github.com/google/uuid"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Server struct {
	Id               uuid.UUID     `json:"id"`
	Object           DIResource    `json:"-"`
	GroupId          uuid.UUID     `json:"group_id"`
	ProjectId        uuid.UUID     `json:"project_id"`
	ClusterUuid      uuid.UUID     `json:"cluster_uuid"`
	State            string        `json:"state"`
	StateResize      string        `json:"state_resize"`
	Step             string        `json:"step"`
	Name             string        `json:"name"`
	ServiceName      string        `json:"service_name" hcl:"service_name"`
	IrGroup          string        `json:"ir_group" hcl:"ir_group"`
	IrType           string        `json:"ir_type"`
	OsName           string        `json:"os_name" hcl:"os_name"`
	OsVersion        string        `json:"os_version" hcl:"os_version"`
	Virtualization   string        `json:"virtualization" hcl:"virtualization"`
	FaultTolerance   string        `json:"fault_tolerance" hcl:"fault_tolerance"`
	Region           string        `json:"region" hcl:"region"`
	NetworkUuid      uuid.UUID     `json:"network_uuid,omitempty" hcl:"network_uuid"`
	User             string        `json:"user"`
	Password         string        `json:"password,omitempty"`
	Cpu              int           `json:"cpu" hcl:"cpu"`
	Ram              int           `json:"ram" hcl:"ram"`
	Disk             int           `json:"disk" hcl:"disk"`
	Flavor           string        `json:"flavor"`
	Zone             string        `json:"zone" hcl:"zone"`
	Ip               string        `json:"ip"`
	DNS              string        `json:"dns"`
	DNSName          string        `json:"dns_name"`
	PublicSshName    string        `json:"public_ssh_name,omitempty" hcl:"public_ssh_name"`
	PublicSsh        string        `json:"public_ssh,omitempty"`
	Group            string        `json:"group,omitempty"`
	ResId            string        `json:"-" hcl:"id"`
	ResType          string        `json:"-" hcl:"type,label"`
	ResName          string        `json:"-" hcl:"name,label"`
	ResGroupIdUUID   string        `json:"-" hcl:"group_id_uuid"`
	ResGroupId       string        `json:"-" hcl:"group_id"`
	ResProjectIdUUID string        `json:"-" hcl:"project_id_uuid"`
	ResProjectId     string        `json:"-" hcl:"project_id"`
	ResAppParams     *HCLAppParams `json:"-" hcl:"app_params,block"`
	ResVolumes       []*HCLVolume  `json:"-" hcl:"volume,block"`
	TagIds           []uuid.UUID   `json:"tag_ids" hcl:"-"`
	ErrMsg           string        `json:"err_msg,omitempty" hcl:"-"`
}

func (o *Server) ReadTF(res *schema.ResourceData) {
	if !res.IsNewResource() {
		o.Id = uuid.MustParse(res.Id())
	}
	groupId := res.Get("group_id")
	if groupId != "" {
		o.GroupId = uuid.MustParse(groupId.(string))
	}
	projectId := res.Get("project_id")
	if projectId != "" {
		o.ProjectId = uuid.MustParse(projectId.(string))
	}

	networkId, ok := res.GetOk("network_uuid")
	if ok {
		o.NetworkUuid = uuid.MustParse(networkId.(string))
	}

	o.ServiceName = res.Get("service_name").(string)
	o.IrGroup = res.Get("ir_group").(string)
	o.OsName = res.Get("os_name").(string)

	osVersion, ok := res.GetOk("os_version")
	if ok {
		o.OsVersion = osVersion.(string)

	}
	o.FaultTolerance = res.Get("fault_tolerance").(string)
	o.Virtualization = res.Get("virtualization").(string)
	//o.Region = res.Get("region").(string)
	o.Zone = res.Get("zone").(string)
	disk, ok := res.GetOk("disk")
	if ok {
		o.Disk = disk.(int)
	}
	_, ok = res.GetOk("cluster_uuid")
	if ok {
		o.ClusterUuid = uuid.MustParse(res.Get("cluster_uuid").(string))
	}
	irType, ok := res.GetOk("ir_type")
	if ok {
		o.IrType = irType.(string)
	}
	user, ok := res.GetOk("user")
	if ok {
		o.User = user.(string)
	}
	flavor, ok := res.GetOk("flavor")
	if ok {
		o.Flavor = flavor.(string)
	}
	// o.Flavor = res.Get("flavor").(string)
	cpu, ok := res.GetOk("cpu")
	if ok {
		o.Cpu = cpu.(int)
	}
	ram, ok := res.GetOk("ram")
	if ok {
		o.Ram = ram.(int)
	}
	// o.Cpu = res.Get("cpu").(int)
	// o.Ram = res.Get("ram").(int)
	o.Name = res.Get("name").(string)
	o.DNS = res.Get("dns").(string)
	o.DNSName = res.Get("dns_name").(string)
	o.Ip = res.Get("ip").(string)
	o.State = res.Get("state").(string)
	o.StateResize = res.Get("state_resize").(string)
	o.Step = res.Get("step").(string)
	publicSshName, ok := res.GetOk("public_ssh_name")
	if ok {
		o.PublicSshName = publicSshName.(string)
	}
	group, ok := res.GetOk("group")
	if ok {
		o.Group = group.(string)
	}
	tags, ok := res.GetOk("tag_ids")
	if ok {
		tagSet := tags.(*schema.Set)
		for _, v := range tagSet.List() {
			// for _, v := range tags.([]interface{}) {
			id, err := uuid.Parse(v.(string))
			if err != nil {
				log.Println(err)
			}
			o.TagIds = append(o.TagIds, id)
		}
	}
	o.Object.OnReadTF(res, o)
}

func (o *Server) GetPubKey() error {
	if o.PublicSshName == "" {
		return nil
	}

	respBytes, err := Api.NewRequestRead("ssh_keys")
	if err != nil {
		return err
	}

	keys := make(map[string][]*SSHKey)
	err = json.Unmarshal(respBytes, &keys)
	if err != nil {
		return err
	}

	for _, v := range keys["ssh_keys"] {
		if v.PublicSshName == o.PublicSshName {
			o.PublicSsh = v.PublicSsh
			return nil
		}
	}
	return fmt.Errorf("can't find public ssh key \"%s\"", o.PublicSshName)
}

func (o *Server) WriteTF(res *schema.ResourceData) {
	res.SetId(o.Id.String())

	res.Set("state", o.State)
	res.Set("state_resize", o.StateResize)
	res.Set("step", o.Step)
	res.Set("name", o.Name)
	res.Set("service_name", o.ServiceName)
	res.Set("ir_group", o.IrGroup)
	res.Set("ir_type", o.IrType)
	res.Set("os_name", o.OsName)
	res.Set("os_version", o.OsVersion)
	res.Set("fault_tolerance", o.FaultTolerance)
	res.Set("virtualization", o.Virtualization)
	res.Set("flavor", o.Flavor)
	res.Set("cpu", o.Cpu)
	res.Set("ram", o.Ram)
	res.Set("disk", o.Disk)
	res.Set("dns", o.DNS)
	res.Set("dns_name", o.DNSName)
	res.Set("ip", o.Ip)
	res.Set("zone", o.Zone)
	//res.Set("region", o.Region)

	_, ok := res.GetOk("network_uuid")
	if ok && o.NetworkUuid != uuid.Nil {
		res.Set("network_uuid", o.NetworkUuid.String())
	}

	res.Set("project_id", o.ProjectId.String())
	res.Set("group_id", o.GroupId.String())
	res.Set("user", o.User)

	isDisabled, ok := os.LookupEnv("VM_PASSWORD_OUTPUT")
	if ok {
		isDisabledBool, err := strconv.ParseBool(isDisabled)
		if err != nil {
			log.Println(err)
		}
		if isDisabledBool == true {
			res.Set("password", o.Password)
		}
	}
	if o.ClusterUuid.ID() != uint32(0) {
		res.Set("cluster_uuid", o.ClusterUuid.String())
	}
	if o.PublicSshName != "" {
		res.Set("public_ssh_name", o.PublicSshName)
	}
	if o.Group != "" {
		res.Set("group", o.Group)
	}
	o.Object.OnWriteTF(res, o)
}

func (o *Server) ToMap() map[string]interface{} {
	serverMap := map[string]interface{}{
		"group_id":     o.GroupId.String(),
		"project_id":   o.ProjectId.String(),
		"service_name": o.ServiceName,
		"ir_group":     o.IrGroup,
		//"region":          o.Region,
		"network_uuid":    o.NetworkUuid.String(),
		"zone":            o.Zone,
		"cpu":             o.Cpu,
		"ram":             o.Ram,
		"disk":            o.Disk,
		"virtualization":  o.Virtualization,
		"os_name":         o.OsName,
		"os_version":      o.OsVersion,
		"fault_tolerance": o.FaultTolerance,
		"id":              o.Id.String(),
		"name":            o.Name,
		"cluster_uuid":    o.ClusterUuid.String(),
		"ir_type":         o.IrType,
		"flavor":          o.Flavor,
		"state":           o.State,
		"state_resize":    o.StateResize,
		"ip":              o.Ip,
		"dns":             o.DNS,
		"dns_name":        o.DNSName,
		"step":            o.Step,
		"user":            o.User,
	}
	if o.PublicSshName != "" {
		serverMap["public_ssh_name"] = o.PublicSshName
		serverMap["public_ssh"] = o.PublicSsh
	}
	if o.Group != "" {
		serverMap["group"] = o.Group
	}
	return serverMap
}

func (o *Server) Serialize() ([]byte, error) {
	serverMap := o.ToMap()
	delete(serverMap, "id")
	delete(serverMap, "name")
	delete(serverMap, "cluster_uuid")
	delete(serverMap, "ir_type")
	// delete(serverMap, "flavor")
	// delete(serverMap, "cpu")
	// delete(serverMap, "ram")
	delete(serverMap, "state")
	delete(serverMap, "state_resize")
	delete(serverMap, "ip")
	delete(serverMap, "dns")
	delete(serverMap, "dns_name")
	delete(serverMap, "step")
	delete(serverMap, "user")
	delete(serverMap, "group")

	// switch o.Object.GetType() {
	// case "di_openshift":
	// 	delete(serverMap, "flavor")
	// default:
	// 	delete(serverMap, "cpu")
	// 	delete(serverMap, "ram")
	// }
	serverMap = o.Object.OnSerialize(serverMap, o)

	requestMap := map[string]interface{}{
		"server": serverMap,
		"count":  1,
	}
	requestBytes, err := json.Marshal(requestMap)
	if err != nil {
		return nil, err
	}
	return requestBytes, nil
}

func (o *Server) Deserialize(data []byte) error {
	response := make(map[string]map[string]interface{})
	err := json.Unmarshal(data, &response)
	if err != nil {
		return err
	}

	serverMap := response["server"]

	o.Name = serverMap["name"].(string)
	o.ServiceName = serverMap["service_name"].(string)
	o.Zone = serverMap["zone"].(string)
	o.OsName = serverMap["os_name"].(string)
	o.OsVersion = serverMap["os_version"].(string)
	o.Virtualization = serverMap["virtualization"].(string)
	o.IrGroup = serverMap["ir_group"].(string)
	o.FaultTolerance = serverMap["fault_tolerance"].(string)
	//o.Region = serverMap["region"].(string)
	//o.NetworkUuid = uuid.MustParse(serverMap["network_uuid"].(string))

	o.State = serverMap["state"].(string)
	// o.IrType = serverMap["ir_type"].(string)
	irType, ok := serverMap["ir_type"]
	if ok && irType != nil {
		o.IrType = irType.(string)
	}
	o.Cpu = int(serverMap["cpu"].(float64))
	o.Ram = int(serverMap["ram"].(float64))
	o.Disk = int(serverMap["disk"].(float64))
	ip, ok := serverMap["ip"]
	if ok && ip != nil {
		o.Ip = ip.(string)
	}
	dns, ok := serverMap["dns"]
	if ok && dns != nil {
		o.DNS = dns.(string)
	}
	dnsName, ok := serverMap["dns-name"]
	if ok && dnsName != nil {
		o.DNSName = dnsName.(string)
	}
	flavor, ok := serverMap["flavor"]
	if ok && flavor != nil {
		o.Flavor = serverMap["flavor"].(string)
	}
	step, ok := serverMap["step"]
	if ok && step != nil {
		o.Step = step.(string)
	}
	outputs, ok := serverMap["outputs"].(map[string]interface{})
	if ok {
		user, ok := outputs["user"]
		if ok {
			o.User = user.(string)
		}
	}
	tempFlavor, ok := serverMap["temp_flavor"]
	if ok {
		if tempFlavor == nil {
			o.StateResize = "stable"
		} else {
			o.StateResize = "resizing"
		}
	}
	clusterUuid, ok := serverMap["cluster_uuid"]
	if ok && clusterUuid != nil {
		o.ClusterUuid = uuid.MustParse(clusterUuid.(string))
	}
	groupId, ok := serverMap["group_id"]
	if ok && groupId != nil {
		o.GroupId = uuid.MustParse(groupId.(string))
	}
	projectId, ok := serverMap["project_id"]
	if ok && projectId != nil {
		o.ProjectId = uuid.MustParse(projectId.(string))
	}

	networkUuid, ok := serverMap["network_uuid"]
	if ok && networkUuid != "" {
		o.NetworkUuid = uuid.MustParse(serverMap["network_uuid"].(string))
	}

	publicSshName, ok := serverMap["public_ssh_name"]
	if ok && publicSshName != "" {
		o.PublicSshName = publicSshName.(string)
	}
	password, ok := serverMap["password"]
	if ok && password != nil {
		o.Password = password.(string)
	}
	errMsg, ok := serverMap["err_msg"]
	if ok && errMsg != nil {
		o.ErrMsg = errMsg.(string)
	}

	o.TagIds = make([]uuid.UUID, 0)
	tagIds, ok := serverMap["tag_ids"]
	if ok && tagIds != nil {
		for _, v := range tagIds.([]interface{}) {
			id, err := uuid.Parse(v.(string))
			if err != nil {
				log.Println(err)
			} else {
				o.TagIds = append(o.TagIds, id)
			}
		}
	}
	o.Object.OnDeserialize(serverMap, o)

	return nil
}

func (o *Server) CreateDI(data []byte) ([]byte, error) {
	return Api.NewRequestCreate(o.Object.Urls("create"), data)
}

func (o *Server) ReadDI() ([]byte, error) {
	return Api.NewRequestRead(fmt.Sprintf(o.Object.Urls("read"), o.Id))
}

func (o *Server) ReadDIStatusCode() ([]byte, int, error) {
	return Api.NewRequestReadStatusCode(fmt.Sprintf(o.Object.Urls("read"), o.Id))
}

func (o *Server) UpdateDI(data []byte) ([]byte, error) {
	return Api.NewRequestUpdate(fmt.Sprintf(o.Object.Urls("update"), o.Id), data)
}

func (o *Server) DeleteDI() error {
	return Api.NewRequestDelete(fmt.Sprintf(o.Object.Urls("delete"), o.Id), nil, 204)
}

func (o *Server) DeleteVM() error {
	return Api.NewRequestDelete(fmt.Sprintf(o.Object.Urls("delete"), o.Id), nil, 204)
}

func (o *Server) ResizeDI(data []byte) ([]byte, error) {
	return Api.NewRequestResize(fmt.Sprintf(o.Object.Urls("resize"), o.Id), data)
}

func (o *Server) VolumeCreateDI(data []byte) ([]byte, error) {
	return Api.NewRequestCreate(fmt.Sprintf(o.Object.Urls("volume_create"), o.Id), data)
}

func (o *Server) VolumeRemoveDI(data []byte) ([]byte, error) {
	return nil, Api.NewRequestDelete(fmt.Sprintf(o.Object.Urls("volume_remove"), o.Id), data, 200)
}

func (o *Server) TagAttachDI(tagId string) ([]byte, error) {
	request := map[string]map[string]string{
		"tag": {
			"tag_uuid": tagId,
		},
	}
	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	return Api.NewRequestCreate(fmt.Sprintf(o.Object.Urls("tag_attach"), o.Id), data)
}

func (o *Server) TagDetachDI(tagId string) error {
	return Api.NewRequestDelete(fmt.Sprintf(o.Object.Urls("tag_detach"), o.Id, tagId), nil, 204)
}

func (o *Server) MoveDI(data []byte) ([]byte, error) {
	return Api.NewRequestMove(o.Object.Urls("move"), data)
}

func (o *Server) ParseIdFromCreateResponse(data []byte) error {
	response := make(map[string][]map[string]interface{})
	err := json.Unmarshal(data, &response)
	if err != nil {
		return err
	}
	if len(response["servers"]) == 0 {
		return errors.New("no server in response")
	}
	server := response["servers"][0]
	o.Id = uuid.MustParse(server["id"].(string))
	return nil
}

func (o *Server) StateChange(res *schema.ResourceData) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Timeout:      res.Timeout(schema.TimeoutCreate),
		PollInterval: 15 * time.Second,
		Pending:      []string{"Creating", "Removing"},
		Target:       []string{"Running", "Damaged", "Removed"},
		Refresh: func() (interface{}, string, error) {

			//responseBytes, err := o.ReadDIStatusCode()
			responseBytes, responseStatusCode, err := o.ReadDIStatusCode()

			if responseStatusCode == 404 {
				return o, "Removed", nil
			}

			if err != nil {
				return nil, "error", err
			}

			err = o.Deserialize(responseBytes)
			if err != nil {
				return nil, "error", err
			}

			log.Printf("[DEBUG] Refresh state for [%s]: state: %s, step: %s", o.Id.String(), o.State, o.Step)
			// write to TF state
			o.WriteTF(res)
			if o.State == "running" {
				return o, "Running", nil
			}
			if o.State == "damaged" {
				return o, "Damaged", nil
			}
			if o.State == "removing" {
				return o, "Removing", nil
			}
			return o, "Creating", nil
		},
	}
}

func (o *Server) StateResizeChange(res *schema.ResourceData) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Timeout:      res.Timeout(schema.TimeoutCreate),
		PollInterval: 15 * time.Second,
		Pending:      []string{"Resizing"},
		Target:       []string{"Stable"},
		Refresh: func() (interface{}, string, error) {
			responseBytes, err := o.ReadDI()
			if err != nil {
				return nil, "error", err
			}
			err = o.Deserialize(responseBytes)
			if err != nil {
				return nil, "error", err
			}
			log.Printf("[DEBUG] Refresh state for [%s]: %s", o.Id.String(), o.StateResize)
			if o.StateResize == "stable" {
				o.WriteTF(res)
				return o, "Stable", nil
			}
			return o, "Resizing", nil
		},
	}
}

//func (o *Server) StateClusterChange(res *schema.ResourceData) *resource.StateChangeConf {
//	return &resource.StateChangeConf{
//		Timeout:      res.Timeout(schema.TimeoutCreate),
//		PollInterval: 15 * time.Second,
//		Pending:      []string{"WaitForCluster"},
//		Target:       []string{"InCluster"},
//		Refresh: func() (interface{}, string, error) {
//			responseBytes, err := o.ReadDI()
//			if err != nil {
//				return nil, "error", err
//			}
//			err = o.Deserialize(responseBytes)
//			if err != nil {
//				return nil, "error", err
//			}
//			log.Printf("[DEBUG] Refresh state for [%s]", o.Id.String())
//			if o.ClusterUuid.ID() != uint32(0) {
//				log.Printf("[DEBUG] Got cluster uuid: %s", o.ClusterUuid)
//				return o, "InCluster", nil
//			}
//			return o, "WaitForCluster", nil
//		},
//	}
//}

func (o *Server) ToHCL() []byte {

	type HCLServerRoot struct {
		Resources *Server `hcl:"resource,block"`
	}
	root := &HCLServerRoot{Resources: o}
	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(root, f.Body())
	return utils.Regexp(f.Bytes())
}

func (o *Server) GetGroup() string {

	group := o.Object.GetGroup()
	if group == "" {
		return utils.Reformat(o.ServiceName)
	}
	return group
}

//func (o *Server) GetAnsibleVaultPassword() (string, error) {
//	if o.Password == "" {
//		return "", fmt.Errorf("no/blank password")
//	}
//	isEnabled, ok := os.LookupEnv("DI_ANSIBLE_PASSWORD")
//	if ok {
//		isEnabledBool, err := strconv.ParseBool(isEnabled)
//		if err != nil {
//			return "", err
//		}
//		if !isEnabledBool {
//			return "", fmt.Errorf("ansible_password disabled")
//		}
//	}
//	vaultPasswordFileLocation := os.Getenv("DI_VAULT_PASSWORD_FILE")
//	TODO: check vaultPasswordFileLocation != ""
//	vaultPasswordFileBytes, err := ioutil.ReadFile(vaultPasswordFileLocation)
// if last byte is '\n'- remove it
//if vaultPasswordFileBytes[len(vaultPasswordFileBytes)-1] == 0x0a {
//	vaultPasswordFileBytes = vaultPasswordFileBytes[:len(vaultPasswordFileBytes)-1]
//}
//passwordEncrypted, err := vault.Encrypt(o.Password, string(vaultPasswordFileBytes))
//if err != nil {
//	return "", err
//}
//return passwordEncrypted, nil
//}

func (o *Server) GetHCLRoot() *HCLRoot {

	root := &HCLRoot{Resources: &HCL{
		// Id:             o.Id.String(),
		// Name:           o.Name,
		// Name:           utils.Reformat(o.ServiceName),
		Name:           utils.Reformat(fmt.Sprintf("%s_%s", o.ServiceName, o.Name)),
		GroupId:        o.GroupId.String(),
		ProjectId:      o.ProjectId.String(),
		ServiceName:    o.ServiceName,
		IrGroup:        o.IrGroup,
		OsName:         o.OsName,
		OsVersion:      o.OsVersion,
		Virtualization: o.Virtualization,
		FaultTolerance: o.FaultTolerance,
		//Region:         o.Region,
		Disk:          o.Disk,
		Flavor:        o.Flavor,
		Zone:          o.Zone,
		PublicSshName: o.PublicSshName,
		AppParams:     nil,
		Volumes:       nil,
		TagIds:        nil,
	}}
	if len(o.TagIds) > 0 {
		tags := HCLTags{}
		for _, v := range o.TagIds {
			tags = append(tags, v.String())
		}
		root.Resources.TagIds = &tags
	}
	return root
}

func (o *Server) GetHCLRootBytes(root *HCLRoot) []byte {

	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(root, f.Body())
	// return utils.Regexp(f.Bytes())
	return f.Bytes()
}

func (o *Server) SetObject() bool {

	if o.IrGroup == "" {
		return false
	}
	if o.Object != nil {
		return false
	}
	var obj DIResource
	switch o.IrGroup {
	case "vm":
		obj = &VM{}
		//case "nginx":
		//	obj = &Nginx{}
		//case "sowa":
		//	obj = &Sowa{}
		//case "project":
		//	obj = &Openshift{}
		//case "postgres":
		//	obj = &Postgres{}
		//case "postgres_se":
		//	obj = &PostgresSE{}
		//case "elk":
		//	obj = &ELK{}
	}
	o.Object = obj
	return true
}

// "kafka":       "di_kafka",
// "ignite":      "di_ignite",
// "patroni":     "di_patroni",

func (o *Server) HCLHeader() []byte {

	return []byte(fmt.Sprintf(
		"resource \"%s\" \"%s\" {}\n",
		o.Object.GetType(),
		utils.Reformat(fmt.Sprintf("%s_%s", o.ServiceName, o.Name)),
	))
}

func (o *Server) ImportCmd() []byte {

	return []byte(fmt.Sprintf(
		"terraform import %s.%s %s\n",
		o.Object.GetType(),
		utils.Reformat(fmt.Sprintf("%s_%s", o.ServiceName, o.Name)),
		o.Id.String(),
	))
}
