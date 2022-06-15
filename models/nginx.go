package models

import (
	"encoding/json"
	"log"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Nginx struct {
	AppParams map[string]interface{} `json:"app_params"`
	Volumes   []*Volume              `json:"volumes"`
}

func (o *Nginx) GetType() string {
	return "di_nginx"
}

func (o *Nginx) NewObj() DIResource {
	return &Nginx{}
}

func (o *Nginx) Urls(action string) string {
	urls := map[string]string{
		"create":        "servers",
		"read":          "servers/%s",
		"update":        "servers/%s",
		"delete":        "servers/%s",
		"resize":        "servers/%s/resize",
		"move":          "servers/moving_vms",
		"tag_attach":    "servers/%s/tags",
		"tag_detach":    "servers/%s/tags/%s",
		"volume_create": "servers/%s/volume_attachments",
		"volume_remove": "servers/%s/volume_detachments",
	}
	return urls[action]
}

func (o *Nginx) OnSerialize(serverData map[string]interface{}, server *Server) map[string]interface{} {
	delete(serverData, "cpu")
	delete(serverData, "ram")
	serialized := map[string]map[string]interface{}{
		"version": {
			"value": o.AppParams["version"].(string),
		},
		"nginx_geoip": {
			"value": o.AppParams["nginx_geoip"].(string),
		},
		"nginx_brotli": {
			"value": o.AppParams["nginx_brotli"].(string),
		},
	}
	joindomain, ok := o.AppParams["joindomain"]
	if ok {
		serialized["joindomain"] = map[string]interface{}{
			"value": joindomain.(string),
		}
	}
	serverData["app_params"] = serialized
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

func (o *Nginx) OnDeserialize(serverData map[string]interface{}, server *Server) {
	params := serverData["app_params"].(map[string]interface{})
	// delete(params, "endpoint")
	// for k, v := range params {
	// 	o.AppParams[k] = v.(map[string]interface{})["value"]
	// delete(v.(map[string]interface{}), "title")
	// }
	if o.AppParams == nil {
		o.AppParams = make(map[string]interface{})
	}
	o.AppParams["version"] = params["version"].(map[string]interface{})["value"]
	o.AppParams["nginx_geoip"] = params["nginx_geoip"].(map[string]interface{})["value"]
	o.AppParams["nginx_brotli"] = params["nginx_brotli"].(map[string]interface{})["value"]
	joindomain, ok := params["joindomain"]
	if ok {
		o.AppParams["joindomain"] = joindomain.(map[string]interface{})["value"]
	}
	// o.AppParams = params
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

func (o *Nginx) OnReadTF(res *schema.ResourceData, server *Server) {
	// if o.AppParams == nil {
	// 	o.AppParams = make(map[string]interface{})
	// }
	//
	// params := res.Get("app_params").([]interface{})
	// o.AppParams = params[0].(map[string]interface{})
	o.AppParams = res.Get("app_params").(map[string]interface{})
	volumes, ok := res.GetOk("volume")
	sort.Sort(ByPath(o.Volumes))
	if ok {
		volumeSet := volumes.(*schema.Set)
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
	}
}

func (o *Nginx) OnWriteTF(res *schema.ResourceData, server *Server) {
	// paramsList := make([]map[string]interface{}, 1)
	// paramsList[0] = o.AppParams
	res.Set("app_params", o.AppParams)
	if o.Volumes != nil && len(o.Volumes) > 0 {
		volumes := make([]map[string]interface{}, 0)
		sort.Sort(ByPath(o.Volumes))
		for _, v := range o.Volumes {
			volume := map[string]interface{}{
				"size":      v.Size,
				"path":      v.Path,
				"volume_id": v.VolumeId,
			}
			volumes = append(volumes, volume)
		}
		err := res.Set("volume", volumes)
		if err != nil {
			log.Println(err)
		}
	}
}

func (o *Nginx) HostVars(server *Server) map[string]interface{} {
	return map[string]interface{}{
		"ansible_host": server.Ip,
		"ansible_user": server.User,
		"dns_name":     server.DNSName,
		"name":         server.Name,
		"service_name": server.ServiceName,
		"state":        server.State,
	}
}

func (o *Nginx) GetGroup() string {
	return ""
}
