package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Server struct {
	Id             uuid.UUID  `json:"id"`
	Object         DIResource `json:"-"`
	GroupId        uuid.UUID  `json:"group_id"`
	ProjectId      uuid.UUID  `json:"project_id"`
	State          string     `json:"state"`
	StateResize    string     `json:"state_resize"`
	Step           string     `json:"step"`
	Name           string     `json:"name"`
	ServiceName    string     `json:"service_name" hcl:"service_name"`
	Comment        string     `json:"comment,omitempty"`
	IrGroup        string     `json:"ir_group" hcl:"ir_group"`
	IrType         string     `json:"ir_type"`
	OsName         string     `json:"os_name" hcl:"os_name"`
	OsVersion      string     `json:"os_version" hcl:"os_version"`
	Virtualization string     `json:"virtualization" hcl:"virtualization"`
	FaultTolerance string     `json:"fault_tolerance" hcl:"fault_tolerance"`
	//Region         string    `json:"region" hcl:"region"`
	NetworkUuid      uuid.UUID `json:"network_uuid,omitempty" hcl:"network_uuid"`
	User             string    `json:"user"`
	Password         string    `json:"password,omitempty"`
	Cpu              int       `json:"cpu" hcl:"cpu"`
	Ram              int       `json:"ram" hcl:"ram"`
	Disk             int       `json:"disk" hcl:"disk"`
	Flavor           string    `json:"flavor"`
	Zone             string    `json:"zone" hcl:"zone"`
	Ip               string    `json:"ip"`
	PublicSshName    string    `json:"public_ssh_name,omitempty" hcl:"public_ssh_name"`
	PublicSsh        string    `json:"public_ssh,omitempty"`
	ResId            string    `json:"-" hcl:"id"`
	ResType          string    `json:"-" hcl:"type,label"`
	ResName          string    `json:"-" hcl:"name,label"`
	ResGroupIdUUID   string    `json:"-" hcl:"group_id_uuid"`
	ResGroupId       string    `json:"-" hcl:"group_id"`
	ResProjectIdUUID string    `json:"-" hcl:"project_id_uuid"`
	ResProjectId     string    `json:"-" hcl:"project_id"`
	//ResAppParams      *HCLAppParams   `json:"-" hcl:"app_params,block"`
	//ResVolumes        []*HCLVolume    `json:"-" hcl:"volume,block"`
	TagIds            []uuid.UUID     `json:"tag_ids" hcl:"-"`
	SecurityGroups    []uuid.UUID     `json:"-" hcl:"-"`
	ResSecurityGroups []SecurityGroup `json:"security_groups" hcl:"-"`
	ErrMsg            string          `json:"err_msg,omitempty" hcl:"-"`
	Hdd               struct {
		Size        int    `json:"size"`
		StorageType string `json:"storage_type,omitempty"`
	} `json:"hdd,omitempty"`
	IsImport bool `json:"-"`
}

func (o *Server) ReadTF(res *schema.ResourceData) {
	if !res.IsNewResource() {
		o.Id = uuid.MustParse(res.Id())
	}
	groupId := res.Get("group_id")
	if groupId != "" {
		o.GroupId = uuid.MustParse(groupId.(string))
	}
	projectId := res.Get("vdc_id")
	if projectId != "" {
		o.ProjectId = uuid.MustParse(projectId.(string))
	}
	networkId, ok := res.GetOk("network_id")
	if ok {
		o.NetworkUuid = uuid.MustParse(networkId.(string))
	}

	o.ServiceName = res.Get("service_name").(string)
	description, ok := res.GetOk("description")
	if ok {
		o.Comment = description.(string)
	}
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
		o.Disk, _ = strconv.Atoi(disk.(map[string]interface{})["size"].(string))
		o.Hdd.Size = o.Disk
		storageType := disk.(map[string]interface{})["storage_type"]
		if storageType != nil {
			o.Hdd.StorageType = disk.(map[string]interface{})["storage_type"].(string)
		}
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
	o.Ip = res.Get("ip").(string)
	o.State = res.Get("state").(string)
	o.StateResize = res.Get("state_resize").(string)
	o.Step = res.Get("step").(string)
	publicSshName, ok := res.GetOk("public_ssh_name")
	if ok {
		o.PublicSshName = publicSshName.(string)
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
	securityGroups, ok := res.GetOk("security_groups")
	if ok {
		for _, v := range securityGroups.(*schema.Set).List() {
			id, err := uuid.Parse(v.(string))
			if err != nil {
				log.Println(err)
			}
			o.SecurityGroups = append(o.SecurityGroups, id)
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
	//return fmt.Errorf("can't find public ssh key \"%s\"", o.PublicSshName)
	return nil
}

func (o *Server) WriteTF(res *schema.ResourceData) {
	res.SetId(o.Id.String())
	err := res.Set("state", o.State)
	err = res.Set("state_resize", o.StateResize)
	err = res.Set("step", o.Step)
	err = res.Set("name", o.Name)
	err = res.Set("service_name", o.ServiceName)
	err = res.Set("description", o.Comment)
	err = res.Set("ir_group", o.IrGroup)
	err = res.Set("ir_type", o.IrType)
	err = res.Set("os_name", o.OsName)
	err = res.Set("os_version", o.OsVersion)
	err = res.Set("fault_tolerance", o.FaultTolerance)
	err = res.Set("virtualization", o.Virtualization)
	err = res.Set("flavor", o.Flavor)
	err = res.Set("cpu", o.Cpu)
	err = res.Set("ram", o.Ram)
	err = res.Set("ip", o.Ip)
	err = res.Set("zone", o.Zone)
	//res.Set("region", o.Region)

	//_, ok := res.GetOk("disk")
	//if ok && o.Disk != 0 {
	//	res.Set("disk", o.Disk)
	//}
	if o.IsImport {
		if o.Hdd.StorageType != "" {
			err := res.Set("disk", map[string]string{"size": strconv.Itoa(o.Disk), "storage_type": o.Hdd.StorageType})
			if err != nil {
				log.Println(err)
			}
		} else {
			err := res.Set("disk", map[string]string{"size": strconv.Itoa(o.Disk)})
			if err != nil {
				log.Println(err)
			}
		}
	}

	_, ok := res.GetOk("network_id")
	if ok && o.NetworkUuid != uuid.Nil || o.IsImport && o.NetworkUuid != uuid.Nil {
		err = res.Set("network_id", o.NetworkUuid.String())
	}

	err = res.Set("vdc_id", o.ProjectId.String())
	err = res.Set("group_id", o.GroupId.String())
	err = res.Set("user", o.User)
	if err != nil {
		log.Println(err)
	}

	isDisabled, ok := os.LookupEnv("VM_PASSWORD_OUTPUT")
	if ok {
		isDisabledBool, err := strconv.ParseBool(isDisabled)
		if err != nil {
			log.Println(err)
		}
		if isDisabledBool == true {
			err = res.Set("password", o.Password)
		}
	}
	if o.PublicSshName != "" {
		err = res.Set("public_ssh_name", o.PublicSshName)
	}
	if o.Hdd.Size != 0 && o.Disk == 0 {
		if o.Hdd.StorageType != "" {
			err = res.Set("disk", map[string]string{
				"size":         strconv.Itoa(o.Disk),
				"storage_type": o.Hdd.StorageType,
			})
		} else {
			err = res.Set("disk", map[string]string{"size": strconv.Itoa(o.Disk)})
		}
	}
	if err != nil {
		log.Println(err)
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
		"network_uuid": o.NetworkUuid.String(),
		"zone":         o.Zone,
		"cpu":          o.Cpu,
		"ram":          o.Ram,
		//"disk":            o.Disk,
		"virtualization":  o.Virtualization,
		"os_name":         o.OsName,
		"os_version":      o.OsVersion,
		"fault_tolerance": o.FaultTolerance,
		"id":              o.Id.String(),
		"name":            o.Name,
		"ir_type":         o.IrType,
		"flavor":          o.Flavor,
		"state":           o.State,
		"state_resize":    o.StateResize,
		"ip":              o.Ip,
		"step":            o.Step,
		"user":            o.User,
	}

	if o.Disk != 0 {
		serverMap["disk"] = o.Disk
	}

	if o.PublicSshName != "" {
		serverMap["public_ssh_name"] = o.PublicSshName
		serverMap["public_ssh"] = o.PublicSsh
	}

	if o.Hdd.StorageType != "" {
		serverMap["hdd"] = map[string]interface{}{"size": o.Disk, "storage_type": o.Hdd.StorageType}
	} else {
		serverMap["hdd"] = map[string]int{"size": o.Disk}
	}
	return serverMap
}

func (o *Server) Serialize() ([]byte, error) {
	serverMap := o.ToMap()
	delete(serverMap, "id")
	delete(serverMap, "name")
	//delete(serverMap, "ir_type")
	// delete(serverMap, "flavor")
	// delete(serverMap, "cpu")
	// delete(serverMap, "ram")
	delete(serverMap, "state")
	delete(serverMap, "state_resize")
	delete(serverMap, "ip")
	delete(serverMap, "step")
	delete(serverMap, "user")

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
	//var defaultSecurityGroup string
	//outputSecurityGroups := serverMap["outputs"].(map[string]interface{})["server_security_groups"].([]interface{})
	//if len(outputSecurityGroups) == 1 {
	//	defaultSecurityGroup = outputSecurityGroups[0].(map[string]interface{})["id"].(string)
	//}
	o.Name = serverMap["name"].(string)
	o.ServiceName = serverMap["service_name"].(string)
	description, ok := serverMap["comment"]
	if ok && description != nil {
		o.Comment = serverMap["comment"].(string)
	}
	zone, ok := serverMap["zone"]
	if ok && zone != nil {
		o.Zone = serverMap["zone"].(string)
	}
	//o.Zone = serverMap["zone"].(string)
	o.OsName = serverMap["os_name"].(string)
	o.OsVersion = serverMap["os_version"].(string)
	o.Virtualization = serverMap["virtualization"].(string)
	o.IrGroup = serverMap["ir_group"].(string)
	faultTolerance, ok := serverMap["fault_tolerance"]
	if ok && faultTolerance != nil {
		o.FaultTolerance = faultTolerance.(string)
	}
	//o.FaultTolerance = serverMap["fault_tolerance"].(string)
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
	groupId, ok := serverMap["group_id"]
	if ok && groupId != nil {
		o.GroupId = uuid.MustParse(groupId.(string))
	}
	projectId, ok := serverMap["project_id"]
	if ok && projectId != nil {
		o.ProjectId = uuid.MustParse(projectId.(string))
	}

	networkUuid, ok := serverMap["network_uuid"]
	if ok && networkUuid != "" && networkUuid != nil {
		o.NetworkUuid = uuid.MustParse(serverMap["network_uuid"].(string))
	}

	publicSshName, ok := serverMap["public_ssh_name"]
	if ok && publicSshName != "" {
		o.PublicSshName = publicSshName.(string)
	}
	hdd, ok := serverMap["hdd"]
	if ok && hdd != "" {
		o.Hdd.Size = hdd.(map[string]interface{})["size"].(int)
		o.Hdd.StorageType = hdd.(map[string]interface{})["storage_type"].(string)
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
	o.SecurityGroups = make([]uuid.UUID, 0)
	securityGroups, ok := serverMap["security_groups"]

	if ok && securityGroups != nil {
		for _, v := range securityGroups.([]interface{}) {
			if v.(map[string]interface{})["group_name"].(string) != "default" {
				id, err := uuid.Parse(v.(map[string]interface{})["security_group_id"].(string))
				if err != nil {
					log.Println(err)
				} else {
					o.SecurityGroups = append(o.SecurityGroups, id)
				}
			}
		}
	}
	o.Object.OnDeserialize(serverMap, o)

	return nil
}

func (o *Server) DeserializeSecurityGroups(data []byte) error {
	response := make(map[string]Server)
	err := json.Unmarshal(data, &response)
	if err != nil {
		return err
	}
	o.ResSecurityGroups = response["server"].ResSecurityGroups
	return nil
}

func (o *Server) CreateDI(data []byte) ([]byte, error) {
	return Api.NewRequestCreate(o.Object.Urls("create"), data)
}

func (o *Server) ReadDI() ([]byte, error) {
	return Api.NewRequestRead(fmt.Sprintf("servers/%s", o.Id))
}

func (o *Server) ReadSIStatusCode() ([]byte, int, error) {
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

func (o *Server) DescriptionAdd(description string) ([]byte, error) {
	request := map[string]map[string]string{
		"server": {
			"comment": description,
		},
	}
	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	return Api.NewRequestUpdate(fmt.Sprintf(o.Object.Urls("update"), o.Id), data)
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

func (o *Server) SecurityGroupVM(securityGroupId string, state string) ([]byte, error) {
	request := map[string]map[string]string{
		"security_group": {
			"state":             state,
			"security_group_id": securityGroupId,
		},
	}
	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	return Api.NewRequestUpdate(fmt.Sprintf("servers/%s/action", o.Id), data)
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

			responseBytes, responseStatusCode, err := o.ReadSIStatusCode()

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

func (o *Server) StateSecurityGroupChange(res *schema.ResourceData) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Timeout:      res.Timeout(schema.TimeoutCreate),
		PollInterval: 15 * time.Second,
		Pending:      []string{"Attaching", "Detaching"},
		Target:       []string{"Attached"},
		Refresh: func() (interface{}, string, error) {
			responseBytes, err := o.ReadDI()
			if err != nil {
				return nil, "error", err
			}

			err = o.DeserializeSecurityGroups(responseBytes)
			if err != nil {
				return nil, "error", err
			}

			log.Printf("[DEBUG] Refresh state for [%s]: %s", o.Id.String(), o.StateResize)

			for _, i := range o.ResSecurityGroups {
				for _, i2 := range i.AttachedToServer {
					if o.Id == i2.ServerUUID {
						if i2.Status == "attaching" {
							return o, "Attaching", nil
						} else if i2.Status == "detaching" {
							return o, "Detaching", nil
						}
					}
				}
			}
			return o, "Attached", nil
		},
	}
}
