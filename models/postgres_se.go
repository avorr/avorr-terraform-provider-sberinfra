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
//type PostgresSE struct {
//	AppParams map[string]interface{} `json:"app_params"`
//	Volumes   []*Volume              `json:"volumes"`
//}
//
//func (o *PostgresSE) GetType() string {
//	return "di_postgres_se"
//}
//
//func (o *PostgresSE) NewObj() DIResource {
//	return &PostgresSE{}
//}
//
//func (o *PostgresSE) Urls(action string) string {
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
//func (o *PostgresSE) serializeAppParams() map[string]interface{} {
//	serialized := make(map[string]interface{})
//	for k, v := range o.AppParams {
//		serialized[k] = map[string]interface{}{
//			"value": v,
//		}
//	}
//	return serialized
//}
//
//func (o *PostgresSE) OnSerialize(serverData map[string]interface{}, server *Server) map[string]interface{} {
//	delete(serverData, "cpu")
//	delete(serverData, "ram")
//
//	serverData["app_params"] = o.serializeAppParams()
//
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
//func (o *PostgresSE) OnDeserialize(serverData map[string]interface{}, server *Server) {
//	params := serverData["app_params"].(map[string]interface{})
//	if o.AppParams == nil {
//		o.AppParams = make(map[string]interface{})
//	}
//	o.AppParams = params
//	// o.AppParams["version"] = params["version"].(map[string]interface{})["value"]
//	// o.AppParams["postgres_db_name"] = params["postgres_db_name"].(map[string]interface{})["value"]
//	// o.AppParams["postgres_db_user"] = params["postgres_db_user"].(map[string]interface{})["value"]
//	// o.AppParams["postgres_db_password"] = params["postgres_db_password"].(map[string]interface{})["value"]
//	// joindomain, ok := params["joindomain"]
//	// if ok {
//	// 	o.AppParams["joindomain"] = joindomain.(map[string]interface{})["value"]
//	// }
//	// maxconnections, ok := params["max_connections"]
//	// if ok {
//	// 	o.AppParams["max_connections"] = maxconnections.(map[string]interface{})["value"]
//	// }
//}
//
//func (o *PostgresSE) OnReadTF(res *schema.ResourceData, server *Server) {
//	o.AppParams = res.Get("app_params").(map[string]interface{})
//
//	volumes, ok := res.GetOk("volume")
//	if ok {
//		volumeSet := volumes.(*schema.Set)
//		// log.Printf("OnReadTF: %v", pp.Sprint(volumeSet))
//		// 	log.Printf("OnReadTF: %v", pp.Sprint(v.List()))
//		for _, v := range volumeSet.List() {
//			values := v.(map[string]interface{})
//			volume := &Volume{
//				Size: values["size"].(int),
//				Path: values["path"].(string),
//				// Name:        values["name"].(string),
//				// FSType:      values["fs_type"].(string),
//				// StorageType: values["storage_type"].(string),
//			}
//			o.Volumes = append(o.Volumes, volume)
//		}
//		sort.Sort(ByPath(o.Volumes))
//		// log.Printf("OnReadTF: %v", pp.Sprint(o.Volumes))
//	}
//
//	// password := o.AppParams["postgres_db_password"]
//	// if password == nil {
//	// 	return
//	// }
//	//
//	// vaultPasswordFileLocation := os.Getenv("DI_VAULT_PASSWORD_FILE")
//	// vaultPasswordFileBytes, err := ioutil.ReadFile(vaultPasswordFileLocation)
//	// if err != nil {
//	// 	log.Println(err)
//	// 	return
//	// }
//	// // if last byte is '\n'- remove it
//	// if vaultPasswordFileBytes[len(vaultPasswordFileBytes)-1] == 0x0a {
//	// 	vaultPasswordFileBytes = vaultPasswordFileBytes[:len(vaultPasswordFileBytes)-1]
//	// }
//	//
//	// passwordEncrypted := o.AppParams["postgres_db_password"].(string)
//	// passwordDecrypted, err := vault.Decrypt(passwordEncrypted, string(vaultPasswordFileBytes))
//	// if err != nil {
//	// 	log.Println(err)
//	// 	return
//	// }
//	// o.AppParams["postgres_db_password"] = passwordDecrypted
//	// o.AppParams["postgres_db_password_ansible"] = passwordEncrypted
//}
//
//func (o *PostgresSE) OnWriteTF(res *schema.ResourceData, server *Server) {
//	// vaultPassword := o.AppParams["postgres_db_password_ansible"]
//	// if vaultPassword != nil {
//	// 	o.AppParams["postgres_db_password"] = vaultPassword
//	// 	delete(o.AppParams, "postgres_db_password_ansible")
//	// }
//	res.Set("app_params", o.AppParams)
//}
//
//func (o *PostgresSE) HostVars(server *Server) map[string]interface{} {
//	return map[string]interface{}{
//		"ansible_host": server.Ip,
//		"ansible_user": server.User,
//		"dns_name":     server.DNSName,
//		"name":         server.Name,
//		// "postgres_db_name":     o.AppParams["postgres_db_name"],
//		// "postgres_db_user":     o.AppParams["postgres_db_user"],
//		// "postgres_db_password": o.AppParams["postgres_db_password"],
//	}
//}
//
//func (o *PostgresSE) GetGroup() string {
//	return ""
//}
//
//func (o *PostgresSE) HCLAppParams() *HCLAppParams {
//	return &HCLAppParams{}
//}
//
//func (o *PostgresSE) HCLVolumes() []*HCLVolume {
//	return nil
//}
