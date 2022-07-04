package models

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type VM struct {
	AppParams map[string]interface{} `json:"app_params"`
	Volumes   []*Volume              `json:"volumes"`
}

func (o *VM) GetType() string {
	return "di_vm"
}

func (o *VM) Urls(action string) string {
	urls := map[string]string{
		"create":        "servers",
		"read":          "servers/%s",
		"update":        "servers/%s",
		"delete":        "servers/%s",
		"resize":        "servers/%s/resize",
		"move":          "servers/moving_vms",
		"volume_create": "servers/%s/volume_attachments",
		"tag_attach":    "servers/%s/tags",
		"tag_detach":    "servers/%s/tags/%s",
	}
	return urls[action]
}

func (o *VM) NewObj() DIResource {
	return &VM{}
}

func (o *VM) OnSerialize(serverData map[string]interface{}, server *Server) map[string]interface{} {
	delete(serverData, "cpu")
	delete(serverData, "ram")
	// if o.AppParams == nil || len(o.AppParams) == 0 {
	// 	return serverData
	// }
	serialized := make(map[string]map[string]interface{})
	joindomain, ok := o.AppParams["joindomain"]
	if ok {
		serialized["joindomain"] = map[string]interface{}{
			"value": joindomain.(string),
		}
	}
	versionJDK, ok := o.AppParams["version_jdk"]
	if ok {
		serialized["version_jdk"] = map[string]interface{}{
			"value": versionJDK.(string),
		}
	}
	serverData["app_params"] = serialized
	serverData["hdd"] = map[string]int{
		"size": server.Disk,
	}
	if server.Region == "" {
		delete(serverData, "region")
	}
	if o.Volumes != nil && len(o.Volumes) > 0 {
		volumes := make([]map[string]interface{}, 0)
		volumeBytes, err := json.Marshal(o.Volumes)
		if err != nil {
			log.Println(err)
			return serverData
		}
		err = json.Unmarshal(volumeBytes, &volumes)
		if err != nil {
			log.Println(err)
			return serverData
		}
		serverData["volumes"] = volumes
	}

	return serverData
}

func (o *VM) OnDeserialize(serverData map[string]interface{}, server *Server) {
	params, ok := serverData["app_params"].(map[string]interface{})
	if !ok {
		return
	}
	joindomain, ok := params["joindomain"]
	if ok {
		if o.AppParams == nil {
			o.AppParams = make(map[string]interface{})
		}
		o.AppParams["joindomain"] = joindomain.(map[string]interface{})["value"]
	}
	versionJDK, ok := params["version_jdk"]
	if ok {
		versionJDKmap := versionJDK.(map[string]interface{})
		if versionJDKmap["value"] != "Не устанавливать" {
			if o.AppParams == nil {
				o.AppParams = make(map[string]interface{})
			}
			switch versionJDKmap["value"].(type) {
			case string:
				o.AppParams["version_jdk"] = versionJDKmap["value"].(string)
			case float64:
				o.AppParams["version_jdk"] = fmt.Sprintf("%.0f", versionJDKmap["value"].(float64))
			}
		}
	}
	volumes, ok := serverData["volumes"].([]interface{})
	if ok {
		if len(volumes) > 0 {
			o.Volumes = make([]*Volume, 0)
			volumeBytes, err := json.Marshal(volumes)
			if err != nil {
				log.Println(err)
				return
			}
			err = json.Unmarshal(volumeBytes, &o.Volumes)
			if err != nil {
				log.Println(err)
				return
			}
			for _, v := range o.Volumes {
				if v.Status == "creating" {
					server.StateResize = "resizing"
				}
			}
		}
		sort.Sort(ByPath(o.Volumes))
	}
}

func (o *VM) OnReadTF(res *schema.ResourceData, server *Server) {
	params, ok := res.GetOk("app_params")
	if ok {
		o.AppParams = params.(map[string]interface{})
	}
	volumes, ok := res.GetOk("volume")
	if ok {
		volumeSet := volumes.(*schema.Set)
		for _, v := range volumeSet.List() {
			values := v.(map[string]interface{})
			volume := &Volume{
				Size:        values["size"].(int),
				Path:        values["path"].(string),
				StorageType: values["storage_type"].(string),
			}
			o.Volumes = append(o.Volumes, volume)
		}
		sort.Sort(ByPath(o.Volumes))
	}
}

func (o *VM) OnWriteTF(res *schema.ResourceData, server *Server) {
	if o.AppParams != nil && len(o.AppParams) > 0 {
		// tmp := strconv.Itoa(o.AppParams["version_jdk"].(int))
		// o.AppParams["version_jdk"] = tmp
		err := res.Set("app_params", o.AppParams)
		if err != nil {
			log.Println(err)
		}
	}
	if o.Volumes != nil && len(o.Volumes) > 0 {
		volumes := make([]map[string]interface{}, 0)
		sort.Sort(ByPath(o.Volumes))
		for _, v := range o.Volumes {
			volume := map[string]interface{}{
				"size":         v.Size,
				"path":         v.Path,
				"storage_type": v.StorageType,
			}
			volumes = append(volumes, volume)
		}
		err := res.Set("volume", volumes)
		if err != nil {
			log.Println(err)
		}
	}
}

func (o *VM) HostVars(server *Server) map[string]interface{} {
	m := map[string]interface{}{
		"ansible_host": server.Ip,
		"ansible_user": server.User,
		"dns_name":     server.DNSName,
		"name":         server.Name,
	}
	passwordEncrypted, err := server.GetAnsibleVaultPassword()
	if err != nil {
		log.Println(err)
	} else {
		m["ansible_password"] = passwordEncrypted
	}
	return m
}

func (o *VM) GetGroup() string {
	return ""
}

func (o *VM) HCLAppParams() *HCLAppParams {
	return &HCLAppParams{
		JoinDomain: o.AppParams["joindomain"].(string),
	}
}

func (o *VM) HCLVolumes() []*HCLVolume {
	if len(o.Volumes) == 0 {
		return nil
	}
	hclVolumes := make([]*HCLVolume, 0)
	for _, v := range o.Volumes {
		vol := &HCLVolume{
			Size:        v.Size,
			Path:        v.Path,
			StorageType: v.StorageType,
		}
		hclVolumes = append(hclVolumes, vol)
	}
	return hclVolumes
}
