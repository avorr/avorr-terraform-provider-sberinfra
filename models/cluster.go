package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sort"
	"time"

	"base.sw.sbc.space/pid/terraform-provider-si/utils"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclwrite"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Cluster struct {
	Id               uuid.UUID         `json:"id"`
	Object           DIClusterResource `json:"object"`
	GroupId          uuid.UUID         `json:"group_id"`
	ProjectId        uuid.UUID         `json:"project_id"`
	State            string            `json:"state"`
	StateResize      string            `json:"state_resize"`
	Step             string            `json:"step"`
	Name             string            `json:"name"`
	ServiceName      string            `json:"service_name" hcl:"service_name"`
	IrGroup          string            `json:"ir_group" hcl:"ir_group"`
	Type             string            `json:"type"`
	OsName           string            `json:"os_name" hcl:"os_name"`
	OsVersion        string            `json:"os_version" hcl:"os_version"`
	Virtualization   string            `json:"virtualization" hcl:"virtualization"`
	FaultTolerance   string            `json:"fault_tolerance" hcl:"fault_tolerance"`
	Region           string            `json:"region" hcl:"region"`
	Cpu              int               `json:"cpu" hcl:"cpu"`
	Ram              int               `json:"ram" hcl:"ram"`
	Disk             int               `json:"disk" hcl:"disk"`
	Flavor           string            `json:"flavor" hcl:"flavor"`
	Zone             string            `json:"zone" hcl:"zone"`
	Servers          []*Server         `json:"servers"`
	PublicSshName    string            `json:"public_ssh_name,omitempty" hcl:"public_ssh_name,optional"`
	PublicSsh        string            `json:"public_ssh,omitempty"`
	Group            string            `json:"group,omitempty"`
	ResId            string            `json:"-" hcl:"id"`
	ResType          string            `json:"-" hcl:"type,label"`
	ResName          string            `json:"-" hcl:"name,label"`
	ResGroupIdUUID   string            `json:"-" hcl:"group_id_uuid"`
	ResGroupId       string            `json:"-" hcl:"group_id"`
	ResProjectIdUUID string            `json:"-" hcl:"project_id_uuid"`
	ResProjectId     string            `json:"-" hcl:"project_id"`
	ResAppParams     *HCLAppParams     `json:"-" hcl:"app_params,block"`
	TagIds           []uuid.UUID       `json:"tag_ids" hcl:"-"`
}

func (o *Cluster) GetPubKey() error {
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

func (o *Cluster) ReadTF(res *schema.ResourceData) {
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
	o.ServiceName = res.Get("service_name").(string)
	o.IrGroup = res.Get("ir_group").(string)
	o.OsName = res.Get("os_name").(string)
	o.OsVersion = res.Get("os_version").(string)
	o.FaultTolerance = res.Get("fault_tolerance").(string)
	o.Virtualization = res.Get("virtualization").(string)
	o.Region = res.Get("region").(string)
	o.Zone = res.Get("zone").(string)
	o.Cpu = res.Get("cpu").(int)
	o.Ram = res.Get("ram").(int)
	o.Disk = res.Get("disk").(int)
	o.Flavor = res.Get("flavor").(string)

	// _, ok := res.GetOk("cluster_uuid")
	// if ok {
	// 	o.ClusterUuid = uuid.MustParse(res.Get("cluster_uuid").(string))
	// }
	typeVar, ok := res.GetOk("type")
	if ok {
		o.Type = typeVar.(string)
	}
	// user, ok := res.GetOk("user")
	// if ok {
	// 	o.User = user.(string)
	// }
	// o.Flavor = res.Get("flavor").(string)
	o.Name = res.Get("name").(string)
	// o.DNS = res.Get("dns").(string)
	// o.DNSName = res.Get("dns_name").(string)
	// o.Ip = res.Get("ip").(string)
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
			id, err := uuid.Parse(v.(string))
			if err != nil {
				log.Println(err)
			}
			o.TagIds = append(o.TagIds, id)
		}
	}

	o.Object.OnReadTF(res, o)
}

func (o *Cluster) WriteTF(res *schema.ResourceData) {
	res.SetId(o.Id.String())

	res.Set("state", o.State)
	res.Set("state_resize", o.StateResize)
	res.Set("step", o.Step)
	res.Set("name", o.Name)
	res.Set("service_name", o.ServiceName)
	res.Set("ir_group", o.IrGroup)
	res.Set("type", o.Type)
	res.Set("os_name", o.OsName)
	res.Set("os_version", o.OsVersion)
	res.Set("fault_tolerance", o.FaultTolerance)
	res.Set("virtualization", o.Virtualization)
	// res.Set("flavor", o.Flavor)
	res.Set("cpu", o.Cpu)
	res.Set("ram", o.Ram)
	res.Set("disk", o.Disk)
	// res.Set("dns", o.DNS)
	// res.Set("dns_name", o.DNSName)
	// res.Set("ip", o.Ip)
	res.Set("zone", o.Zone)
	//res.Set("region", o.Region)
	res.Set("project_id", o.ProjectId.String())
	res.Set("group_id", o.GroupId.String())

	// if o.ClusterUuid.ID() != uint32(0) {
	// 	res.Set("cluster_uuid", o.ClusterUuid.String())
	// }
	servers := make([]map[string]interface{}, 0)
	for _, v := range o.Servers {
		servers = append(servers, v.ToMap())
	}
	res.Set("servers", servers)

	if o.PublicSshName != "" {
		res.Set("public_ssh_name", o.PublicSshName)
	}
	if o.Group != "" {
		res.Set("group", o.Group)
	}

	o.Object.OnWriteTF(res, o)
}

func (o *Cluster) Serialize() ([]byte, error) {
	serverMap := map[string]interface{}{
		"group_id":        o.GroupId,
		"project_id":      o.ProjectId,
		"service_name":    o.ServiceName,
		"ir_group":        o.IrGroup,
		"region":          o.Region,
		"zone":            o.Zone,
		"cpu":             o.Cpu,
		"ram":             o.Ram,
		"disk":            o.Disk,
		"virtualization":  o.Virtualization,
		"os_name":         o.OsName,
		"os_version":      o.OsVersion,
		"fault_tolerance": o.FaultTolerance,
		// "flavor":          o.Flavor,
	}
	if o.PublicSshName != "" {
		serverMap["public_ssh_name"] = o.PublicSshName
		serverMap["public_ssh"] = o.PublicSsh
	}
	if o.Flavor == "" {
		serverMap["flavor"] = nil
	} else {
		serverMap["flavor"] = o.Flavor
	}

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

func (o *Cluster) Deserialize(data []byte) error {
	response := make(map[string]map[string]interface{})
	err := json.Unmarshal(data, &response)
	if err != nil {
		return err
	}

	serverMap := response["cluster"]

	o.Name = serverMap["name"].(string)
	o.State = serverMap["state"].(string)
	o.Type = serverMap["type"].(string)
	o.IrGroup = serverMap["ir_group"].(string)
	o.ServiceName = serverMap["service_name"].(string)
	o.Zone = serverMap["zone"].(string)
	o.OsName = serverMap["os_name"].(string)
	o.OsVersion = serverMap["os_version"].(string)
	if serverMap["step"] != nil {
		o.Step = serverMap["step"].(string)
	}
	o.GroupId = uuid.MustParse(serverMap["group_id"].(string))
	o.Cpu = int(serverMap["cpu"].(float64))
	o.Ram = int(serverMap["ram"].(float64))
	o.Disk = int(serverMap["disk"].(float64))
	step, ok := serverMap["step"]
	if ok && step != nil {
		o.Step = step.(string)
	}

	o.StateResize = "stable"

	o.Servers = make([]*Server, 0)
	for _, v := range serverMap["servers"].([]interface{}) {
		serverInfo := v.(map[string]interface{})
		serverData := map[string]map[string]interface{}{
			"server": serverInfo,
		}
		serialized, err := json.Marshal(serverData)
		if err != nil {
			return err
		}
		server := &Server{Object: &VM{}, Id: uuid.MustParse(serverInfo["id"].(string))}
		err = server.Deserialize(serialized)
		if err != nil {
			return err
		}
		o.Servers = append(o.Servers, server)
		if server.StateResize == "resizing" {
			o.StateResize = "resizing"
		}
	}
	sort.Sort(ById(o.Servers))

	if len(o.Servers) > 0 {
		o.Virtualization = o.Servers[0].Virtualization
		o.FaultTolerance = o.Servers[0].FaultTolerance
		o.Region = o.Servers[0].Region
		o.PublicSshName = o.Servers[0].PublicSshName
		o.PublicSsh = o.Servers[0].PublicSsh
		o.ProjectId = o.Servers[0].ProjectId
		o.Flavor = o.Servers[0].Flavor
	}

	o.Object.OnDeserialize(serverMap, o)

	return nil
}

func (o *Cluster) CreateDI(data []byte) ([]byte, error) {
	return Api.NewRequestCreate(o.Object.Urls("create"), data)
}

func (o *Cluster) ReadDI() ([]byte, error) {
	return Api.NewRequestRead(fmt.Sprintf(o.Object.Urls("read"), o.Id))
}

func (o *Cluster) UpdateDI(data []byte) ([]byte, error) {
	return Api.NewRequestUpdate(fmt.Sprintf(o.Object.Urls("update"), o.Id), data)
}

func (o *Cluster) DeleteDI() error {
	return Api.NewRequestDelete(fmt.Sprintf(o.Object.Urls("delete"), o.Id), nil)
}

func (o *Cluster) ResizeDI(data []byte) ([]byte, error) {
	return Api.NewRequestResize(fmt.Sprintf(o.Object.Urls("resize"), o.Id), data)
}

func (o *Cluster) MoveDI(data []byte) ([]byte, error) {
	return Api.NewRequestMove(o.Object.Urls("move"), data)
}

func (o *Cluster) UpScaleDI(data []byte) ([]byte, error) {
	return Api.NewRequestUpScale(o.Object.Urls("add_nodes"), data)
}

func (o *Cluster) ParseIdFromCreateResponse(data []byte) (string, error) {
	response := make(map[string][]map[string]interface{})
	err := json.Unmarshal(data, &response)
	if err != nil {
		return "", err
	}
	if len(response["servers"]) == 0 {
		return "", errors.New("no server in response")
	}
	server := response["servers"][0]
	return server["id"].(string), nil
}

func (o *Cluster) StateChange(res *schema.ResourceData) *resource.StateChangeConf {
	log.Printf("schema timeout %s", schema.TimeoutCreate)
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
			o.WriteTF(res)
			log.Printf("[DEBUG] Refresh state for [%s]: state: %s, step: %s", o.Id.String(), o.State, o.Step)
			if o.State == "running" {
				// write to TF state
				// o.WriteTF(res)
				return o, "Running", nil
			}
			if o.State == "damaged" {
				return o, "Damaged", nil
			}
			return o, "Creating", nil
		},
	}
}

func (o *Cluster) StateResizeChange(res *schema.ResourceData) *resource.StateChangeConf {
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

func (o *Cluster) ToHCL() []byte {
	type HCLServerRoot struct {
		Resources *Cluster `hcl:"resource,block"`
	}
	root := &HCLServerRoot{Resources: o}
	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(root, f.Body())
	return utils.Regexp(f.Bytes())
}

type ById []*Server

func (o ById) Len() int           { return len(o) }
func (o ById) Less(i, j int) bool { return o[i].Id.String() < o[j].Id.String() }
func (o ById) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }

func (o *Cluster) SetObject() bool {
	if o.IrGroup == "" {
		return false
	}
	if o.Object != nil {
		return false
	}
	var obj DIClusterResource
	switch o.IrGroup {
	//case "kafka":
	//	obj = &Kafka{}
	//case "ignite":
	//	obj = &Ignite{}
	//case "patroni":
	//	obj = &Patroni{}
	}
	o.Object = obj
	return true
}

func (o *Cluster) HCLHeader() []byte {
	return []byte(fmt.Sprintf(
		"resource \"%s\" \"%s\" {}\n",
		o.Object.GetType(),
		utils.Reformat(fmt.Sprintf("%s_%s", o.ServiceName, o.Name)),
	))
}

func (o *Cluster) ImportCmd() []byte {
	return []byte(fmt.Sprintf(
		"terraform import %s.%s %s\n",
		o.Object.GetType(),
		utils.Reformat(fmt.Sprintf("%s_%s", o.ServiceName, o.Name)),
		o.Id.String(),
	))
}
