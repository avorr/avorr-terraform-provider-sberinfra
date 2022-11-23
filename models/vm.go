package models

import (
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type VM struct {
	AppParams map[string]interface{} `json:"app_params"`
	Volumes   []*Volume              `json:"volumes"`
}

//	pc, _, _, _ := runtime.Caller(0)
//	log.Println(pp.Sprintf("RUN FUNC %s", runtime.FuncForPC(pc).Name()))

func (o *VM) GetType() string {
	return "si_vm"
}

func (o *VM) GetFile() string {
	return "vm.tf"
}

func (o *VM) Urls(action string) string {
	urls := map[string]string{
		"create":         "servers",
		"read":           "servers/%s",
		"update":         "servers/%s",
		"delete":         "servers/%s",
		"resize":         "servers/%s/resize",
		"move":           "servers/moving_vms",
		"volume_create":  "servers/%s/volume_attachments",
		"volume_remove":  "servers/%s/volume_detachments",
		"tag_attach":     "servers/%s/tags",
		"tag_detach":     "servers/%s/tags/%s",
		"security_group": "servers/%s/action",
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
			//"size": server.Disk,
			"size": server.Disk,
		}
	}
	//if server.Region == "" {
	//	delete(serverData, "region")
	//}
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
		sort.Sort(ByPath(o.Volumes))
	}
}

func (o *VM) OnReadTF(res *schema.ResourceData, server *Server) {
	volumes, ok := res.GetOk("volume")
	if ok {
		volumeSet := volumes.(*schema.Set)
		for _, v := range volumeSet.List() {
			values := v.(map[string]interface{})
			if values["storage_type"].(string) != "__DEFAULT__" {
				volume := &Volume{
					Size: values["size"].(int),
					//Path:        values["path"].(string),
					StorageType: values["storage_type"].(string),
				}
				o.Volumes = append(o.Volumes, volume)

			} else {
				volume := &Volume{Size: values["size"].(int)}
				o.Volumes = append(o.Volumes, volume)
			}
		}
		sort.Sort(ByPath(o.Volumes))
	}
}

func (o *VM) OnWriteTF(res *schema.ResourceData, server *Server) {
	if o.Volumes != nil && len(o.Volumes) > 0 {
		volumes := make([]map[string]interface{}, 0)
		sort.Sort(ByPath(o.Volumes))
		for _, v := range o.Volumes {
			if v.StorageType != "__DEFAULT__" {
				volume := map[string]interface{}{
					"size": v.Size,
					//"path":         v.Path,
					"storage_type": v.StorageType,
					"volume_id":    v.VolumeId,
				}
				volumes = append(volumes, volume)
			} else {
				volume := map[string]interface{}{
					"size": v.Size,
					//"path": v.Path,
					//"storage_type": v.StorageType,
					"volume_id": v.VolumeId,
				}
				volumes = append(volumes, volume)
			}
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
		"name":         server.Name,
	}
	return m
}

func (o *VM) GetGroup() string {
	return ""
}

func (o *VM) HCLAppParams() *HCLAppParams {
	return &HCLAppParams{
		JoinDomain: "",
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
