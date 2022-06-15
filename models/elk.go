package models

import (
	"encoding/json"
	"log"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ELK struct {
	AppParams map[string]interface{} `json:"app_params"`
	Volumes   []*Volume              `json:"volumes"`
}

func (o *ELK) GetType() string {
	return "di_elk"
}

func (o *ELK) Urls(action string) string {
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

func (o *ELK) NewObj() DIResource {
	return &ELK{}
}

func (o *ELK) OnSerialize(serverData map[string]interface{}, server *Server) map[string]interface{} {
	delete(serverData, "cpu")
	delete(serverData, "ram")
	serialized := map[string]map[string]interface{}{
		"version": {
			"value": o.AppParams["version"].(string),
		},
		"java_version": {
			"value": o.AppParams["java_version"].(string),
		},
		"elk_set": {
			"value": o.AppParams["elk_set"].(string),
		},
		"joindomain": {
			"value": o.AppParams["joindomain"].(string),
		},
	}
	serverData["app_params"] = serialized
	// log.Printf("OnSerialize: %v", pp.Sprint(o.Volumes))
	serverData["hdd"] = map[string]int{
		"size": server.Disk,
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

func (o *ELK) OnDeserialize(serverData map[string]interface{}, server *Server) {
	params, ok := serverData["app_params"].(map[string]interface{})
	if !ok {
		return
	}
	if o.AppParams == nil {
		o.AppParams = make(map[string]interface{})
	}
	o.AppParams["version"] = params["version"].(map[string]interface{})["value"]
	o.AppParams["java_version"] = params["java_version"].(map[string]interface{})["value"]
	o.AppParams["elk_set"] = params["elk_set"].(map[string]interface{})["value"]
	o.AppParams["joindomain"] = params["joindomain"].(map[string]interface{})["value"]
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
	// log.Println(pp.Sprint(volumes))
	// log.Printf("OnDeserialize: %v", pp.Sprint(o.Volumes))
}

func (o *ELK) OnReadTF(res *schema.ResourceData, server *Server) {
	params, ok := res.GetOk("app_params")
	if ok {
		o.AppParams = params.(map[string]interface{})
	}
	volumes, ok := res.GetOk("volume")
	if ok {
		volumeSet := volumes.(*schema.Set)
		// log.Printf("OnReadTF: %v", pp.Sprint(volumeSet))
		// 	log.Printf("OnReadTF: %v", pp.Sprint(v.List()))
		for _, v := range volumeSet.List() {
			values := v.(map[string]interface{})
			volume := &Volume{
				Size: values["size"].(int),
				Path: values["path"].(string),
				// Name:        values["name"].(string),
				// FSType:      values["fs_type"].(string),
				// StorageType: values["storage_type"].(string),
			}
			o.Volumes = append(o.Volumes, volume)
		}
		sort.Sort(ByPath(o.Volumes))
		// log.Printf("OnReadTF: %v", pp.Sprint(o.Volumes))
	}
}

func (o *ELK) OnWriteTF(res *schema.ResourceData, server *Server) {
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
				// "id":      v.Id,
				"size": v.Size,
				// "name":    v.Name,
				"path": v.Path,
				// "status":  v.Status,
				// "fs_type": v.FSType,
			}
			volumes = append(volumes, volume)
		}

		// volumeBytes, err := json.Marshal(o.Volumes)
		// if err != nil {
		// 	log.Println(err)
		// 	return
		// }
		// err = json.Unmarshal(volumeBytes, &volumes)
		// if err != nil {
		// 	log.Println(err)
		// 	return
		// }
		// sort.Sort(ByPath(o.Volumes))
		// log.Println(pp.Sprint(o.Volumes))
		// log.Printf("OnWriteTF: %v", pp.Sprint(volumes))
		err := res.Set("volume", volumes)
		if err != nil {
			log.Println(err)
		}
	}
}

func (o *ELK) HostVars(server *Server) map[string]interface{} {
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

func (o *ELK) GetGroup() string {
	return ""
}
