package models

//
//import (
//	"encoding/json"
//	"fmt"
//	"io/ioutil"
//	"log"
//	"os"
//	"strconv"
//	"strings"
//
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
//	vault "github.com/sosedoff/ansible-vault-go"
//)
//
//type Patroni struct {
//	AppParams map[string]interface{} `json:"app_params"`
//	Volumes   []*Volume              `json:"volumes"`
//}
//
//func (o *Patroni) GetType() string {
//	return "di_patroni"
//}
//
//func (o *Patroni) NewObj() DIClusterResource {
//	return &Patroni{}
//}
//
//func (o *Patroni) Urls(action string) string {
//	urls := map[string]string{
//		"create":     "servers",
//		"read":       "servers/clusters/%s",
//		"update":     "servers/clusters/%s",
//		"delete":     "servers/clusters/%s",
//		"resize":     "servers/clusters/%s/resize",
//		"move":       "servers/moving_vms",
//		"tag_attach": "servers/%s/tags",
//		"tag_detach": "servers/%s/tags/%s",
//	}
//	return urls[action]
//}
//
//func (o *Patroni) OnSerialize(serverData map[string]interface{}, cluster *Cluster) map[string]interface{} {
//	delete(serverData, "cpu")
//	delete(serverData, "ram")
//	serialized := map[string]map[string]interface{}{
//		"version":          {"value": o.AppParams["version"].(string)},
//		"fault_tolerance":  {"value": o.AppParams["fault_tolerance"].(string)},
//		"box_server_count": {"value": o.AppParams["box_server_count"].(int)},
//		"dc_quantity":      {"value": o.AppParams["dc_quantity"].(int)},
//		"max_connections":  {"value": o.AppParams["max_connections"].(string)},
//		"postgres_db_password": {
//			"value": o.AppParams["postgres_db_password"].(string),
//		},
//		"postgres_dbs": {"value": strings.Split(o.AppParams["postgres_dbs"].(string), ",")},
//		// "postgres_dbs":     {"value": []string{o.AppParams["postgres_dbs"].(string)}},
//	}
//
//	joindomain, ok := o.AppParams["joindomain"]
//	if ok {
//		serialized["joindomain"] = map[string]interface{}{
//			"value": joindomain.(string),
//		}
//	}
//	serverData["app_params"] = serialized
//	serverData["hdd"] = map[string]int{
//		"size": cluster.Disk,
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
//func (o *Patroni) OnDeserialize(serverData map[string]interface{}, cluster *Cluster) {
//	servers := serverData["servers"].([]interface{})
//	server := servers[0].(map[string]interface{})
//	params := server["app_params"].(map[string]interface{})
//	if o.AppParams == nil {
//		o.AppParams = make(map[string]interface{})
//	}
//
//	dbs := []string{}
//	for _, v := range params["postgres_dbs"].(map[string]interface{})["value"].([]interface{}) {
//		dbs = append(dbs, v.(string))
//	}
//	dbsString := strings.Join(dbs, ",")
//	params["postgres_dbs"] = map[string]interface{}{
//		"value": dbsString,
//	}
//
//	newParams := make(map[string]interface{})
//	for k, v := range params {
//		if k == "endpoint" {
//			continue
//		}
//		val := v.(map[string]interface{})
//		newParams[k] = val["value"]
//	}
//
//	paramsBytes, err := json.Marshal(newParams)
//	if err != nil {
//		log.Println(err)
//	}
//	err = json.Unmarshal(paramsBytes, &o.AppParams)
//	if err != nil {
//		log.Println(err)
//	}
//	if cluster.ResAppParams == nil {
//		cluster.ResAppParams = &HCLAppParams{}
//	}
//	err = json.Unmarshal(paramsBytes, &cluster.ResAppParams)
//	if err != nil {
//		log.Println(err)
//	}
//
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
//					cluster.StateResize = "resizing"
//				}
//			}
//		}
//	}
//}
//
//func (o *Patroni) OnReadTF(res *schema.ResourceData, cluster *Cluster) {
//	if o.AppParams == nil {
//		o.AppParams = make(map[string]interface{})
//	}
//	o.AppParams = res.Get("app_params").(map[string]interface{})
//	if o.AppParams["box_server_count"] != nil {
//		o.AppParams["box_server_count"], _ = strconv.Atoi(o.AppParams["box_server_count"].(string))
//	}
//	if o.AppParams["dc_quantity"] != nil {
//		o.AppParams["dc_quantity"], _ = strconv.Atoi(o.AppParams["dc_quantity"].(string))
//	}
//	volumes, ok := res.GetOk("volume")
//	if ok {
//		for _, v := range volumes.([]interface{}) {
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
//	}
//
//	password := o.AppParams["postgres_db_password"]
//	if password == nil {
//		return
//	}
//
//	vaultPasswordFileLocation := os.Getenv("DI_VAULT_PASSWORD_FILE")
//	vaultPasswordFileBytes, err := ioutil.ReadFile(vaultPasswordFileLocation)
//	if err != nil {
//		log.Println(err)
//		return
//	}
//	// if last byte is '\n'- remove it
//	if vaultPasswordFileBytes[len(vaultPasswordFileBytes)-1] == 0x0a {
//		vaultPasswordFileBytes = vaultPasswordFileBytes[:len(vaultPasswordFileBytes)-1]
//	}
//
//	passwordEncrypted := o.AppParams["postgres_db_password"].(string)
//	passwordDecrypted, err := vault.Decrypt(passwordEncrypted, string(vaultPasswordFileBytes))
//	if err != nil {
//		log.Println(err)
//		return
//	}
//	o.AppParams["postgres_db_password"] = passwordDecrypted
//	o.AppParams["postgres_db_password_ansible"] = passwordEncrypted
//
//}
//
//func (o *Patroni) OnWriteTF(res *schema.ResourceData, cluster *Cluster) {
//	vaultPassword := o.AppParams["postgres_db_password_ansible"]
//	if vaultPassword != nil {
//		o.AppParams["postgres_db_password"] = vaultPassword
//		delete(o.AppParams, "postgres_db_password_ansible")
//	}
//	// res.Set("app_params", o.AppParams)
//
//	params := map[string]string{
//		"version":         o.AppParams["version"].(string),
//		"fault_tolerance": o.AppParams["fault_tolerance"].(string),
//		"max_connections": o.AppParams["max_connections"].(string),
//		"postgres_dbs":    o.AppParams["postgres_dbs"].(string),
//		// "postgres_db_password": o.AppParams["postgres_db_password"].(string),
//	}
//
//	password := o.AppParams["postgres_db_password"]
//	if password != nil {
//		params["postgres_db_password"] = o.AppParams["postgres_db_password"].(string)
//		// params["postgres_db_password"] = o.AppParams["postgres_db_password"].(string)
//	}
//
//	// releaseType, ok := o.AppParams["release_type"]
//	// if ok {
//	// 	params["release_type"] = releaseType.(string)
//	// }
//	switch o.AppParams["box_server_count"].(type) {
//	case int:
//		params["box_server_count"] = fmt.Sprintf("%d", o.AppParams["box_server_count"].(int))
//	case float64:
//		params["box_server_count"] = fmt.Sprintf("%.0f", o.AppParams["box_server_count"].(float64))
//	}
//	switch o.AppParams["dc_quantity"].(type) {
//	case int:
//		params["dc_quantity"] = fmt.Sprintf("%d", o.AppParams["dc_quantity"].(int))
//	case float64:
//		params["dc_quantity"] = fmt.Sprintf("%.0f", o.AppParams["dc_quantity"].(float64))
//	}
//	joindomain, ok := o.AppParams["joindomain"]
//	if ok {
//		params["joindomain"] = joindomain.(string)
//	}
//	res.Set("app_params", params)
//
//	// if o.Volumes != nil && len(o.Volumes) > 0 {
//	// 	volumes := make([]map[string]interface{}, 0)
//	// 	volumeBytes, err := json.Marshal(o.Volumes)
//	// 	if err != nil {
//	// 		log.Println(err)
//	// 		return
//	// 	}
//	// 	err = json.Unmarshal(volumeBytes, &volumes)
//	// 	if err != nil {
//	// 		log.Println(err)
//	// 		return
//	// 	}
//	// 	err = res.Set("volume", volumes)
//	// 	if err != nil {
//	// 		log.Println(err)
//	// 	}
//	// }
//
//}
//
//func (o *Patroni) HostVars(server *Server) map[string]interface{} {
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
//func (o *Patroni) GroupVars(cluster *Cluster) map[string]interface{} {
//	return map[string]interface{}{}
//}
