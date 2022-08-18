package models

//
//import (
//	"encoding/json"
//	"log"
//	"sort"
//
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
//)
//
//type Sowa struct {
//	AppParams map[string]interface{} `json:"app_params"`
//	Volumes   []*Volume              `json:"volumes"`
//}
//
//func (o *Sowa) GetType() string {
//	return "di_sowa"
//}
//
//func (o *Sowa) NewObj() DIResource {
//	return &Sowa{}
//}
//
//func (o *Sowa) Urls(action string) string {
//	urls := map[string]string{
//		"create":     "servers",
//		"read":       "servers/%s",
//		"update":     "servers/%s",
//		"delete":     "servers/%s",
//		"resize":     "servers/%s/resize",
//		"move":       "servers/moving_vms",
//		"tag_attach": "servers/%s/tags",
//		"tag_detach": "servers/%s/tags/%s",
//	}
//	return urls[action]
//}
//
//func (o *Sowa) OnSerialize(serverData map[string]interface{}, server *Server) map[string]interface{} {
//	delete(serverData, "cpu")
//	delete(serverData, "ram")
//	serialized := map[string]map[string]interface{}{
//		"version": {
//			"value": o.AppParams["version"].(string),
//		},
//	}
//	joindomain, ok := o.AppParams["joindomain"]
//	if ok {
//		serialized["joindomain"] = map[string]interface{}{
//			"value": joindomain.(string),
//		}
//	}
//	serverData["app_params"] = serialized
//	serverData["hdd"] = map[string]int{
//		"size": server.Disk,
//	}
//	if o.Volumes != nil && len(o.Volumes) > 0 {
//		volumes := make([]map[string]interface{}, 0)
//		volumeBytes, err := json.Marshal(o.Volumes)
//		if err != nil {
//			log.Println(err)
//			return serverData
//		}
//		err = json.Unmarshal(volumeBytes, &volumes)
//		if err != nil {
//			log.Println(err)
//			return serverData
//		}
//		serverData["volumes"] = volumes
//	}
//	return serverData
//}
//
//func (o *Sowa) OnDeserialize(serverData map[string]interface{}, server *Server) {
//	log.Println("OnDeserialize")
//	params := serverData["app_params"].(map[string]interface{})
//	joindomain, ok := params["joindomain"]
//	if o.AppParams == nil {
//		o.AppParams = make(map[string]interface{})
//	}
//	if ok {
//		o.AppParams["joindomain"] = joindomain.(map[string]interface{})["value"]
//	}
//	o.AppParams["version"] = params["version"].(map[string]interface{})["value"]
//	volumes, ok := serverData["volumes"].([]interface{})
//	if ok {
//		if len(volumes) > 0 {
//			o.Volumes = make([]*Volume, 0)
//			volumeBytes, err := json.Marshal(volumes)
//			if err != nil {
//				log.Println(err)
//				return
//			}
//			err = json.Unmarshal(volumeBytes, &o.Volumes)
//			if err != nil {
//				log.Println(err)
//				return
//			}
//			for _, v := range o.Volumes {
//				if v.Status == "creating" {
//					server.StateResize = "resizing"
//				}
//			}
//		}
//		sort.Sort(ByPath(o.Volumes))
//	}
//}
//
//func (o *Sowa) OnReadTF(res *schema.ResourceData, server *Server) {
//	log.Println("OnReadTF")
//	o.AppParams = res.Get("app_params").(map[string]interface{})
//	volumes, ok := res.GetOk("volume")
//	if ok {
//		volumeSet := volumes.(*schema.Set)
//		for _, v := range volumeSet.List() {
//			values := v.(map[string]interface{})
//			volume := &Volume{
//				Size: values["size"].(int),
//				Path: values["path"].(string),
//			}
//			o.Volumes = append(o.Volumes, volume)
//		}
//		sort.Sort(ByPath(o.Volumes))
//	}
//}
//
//func (o *Sowa) OnWriteTF(res *schema.ResourceData, server *Server) {
//	log.Println("OnWriteTF")
//	res.Set("app_params", o.AppParams)
//	if o.Volumes != nil && len(o.Volumes) > 0 {
//		volumes := make([]map[string]interface{}, 0)
//		sort.Sort(ByPath(o.Volumes))
//		for _, v := range o.Volumes {
//			volume := map[string]interface{}{
//				"size": v.Size,
//				"path": v.Path,
//			}
//			volumes = append(volumes, volume)
//		}
//		err := res.Set("volume", volumes)
//		if err != nil {
//			log.Println(err)
//		}
//	}
//}
//
//func (o *Sowa) HostVars(server *Server) map[string]interface{} {
//	return map[string]interface{}{
//		"ansible_host": server.Ip,
//		"ansible_user": server.User,
//		"dns_name":     server.DNSName,
//		"name":         server.Name,
//		"service_name": server.ServiceName,
//		"state":        server.State,
//	}
//}
//
//func (o *Sowa) GetGroup() string {
//	return ""
//}
//
//func (o *Sowa) HCLAppParams() *HCLAppParams {
//	return &HCLAppParams{}
//}
//
//func (o *Sowa) HCLVolumes() []*HCLVolume {
//	return nil
//}
