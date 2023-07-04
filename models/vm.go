package models

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

type VM struct {
	Volumes []*Volume `json:"volumes"`
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
		"volume_remove": "servers/%s/volume_detachments",
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
	if server.Hdd.StorageType != "" {
		serverData["hdd"] = map[string]interface{}{
			//"size": server.Disk,
			"size":         server.Disk,
			"storage_type": server.Hdd.StorageType,
		}
	} else {
		serverData["hdd"] = map[string]int{
			"size": server.Disk,
		}
	}
	if server.NetworkUuid == uuid.Nil {
		delete(serverData, "network_uuid")
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
				if v.Status == "removing" {
					server.StateResize = "resizing"
				}
			}
		}
	}
}

func (o *VM) OnReadTF(res *schema.ResourceData, server *Server) {
	volumes, ok := res.GetOk("volume")
	if ok {
		volumeSet := volumes.(*schema.Set)
		for _, v := range volumeSet.List() {
			values := v.(map[string]interface{})
			volume := &Volume{
				Size: values["size"].(int),
				Name: values["name"].(string),
				//Description: values["description"].(string),
				StorageType: values["storage_type"].(string),
			}
			o.Volumes = append(o.Volumes, volume)
		}
	}
}

func (o *VM) OnWriteTF(res *schema.ResourceData, server *Server) {
	if o.Volumes != nil && len(o.Volumes) > 0 {
		volumes := make([]map[string]interface{}, 0)
		for _, v := range o.Volumes {
			volume := map[string]interface{}{
				"size": v.Size,
				"name": v.Name,
				//"description":  v.Description,
				"storage_type": v.StorageType,
				"volume_id":    v.VolumeId,
			}
			volumes = append(volumes, volume)
		}
		err := res.Set("volume", volumes)
		if err != nil {
			log.Println(err)
		}
	}
}
