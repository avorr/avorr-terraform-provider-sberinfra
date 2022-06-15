package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vault "github.com/sosedoff/ansible-vault-go"
)

type Ignite struct {
	AppParams map[string]interface{} `json:"app_params"`
	Volumes   []*Volume              `json:"volumes"`
}

func (o *Ignite) GetType() string {
	return "di_ignite"
}

func (o *Ignite) NewObj() DIClusterResource {
	return &Ignite{}
}

func (o *Ignite) Urls(action string) string {
	urls := map[string]string{
		"create":     "servers",
		"read":       "servers/clusters/%s",
		"update":     "servers/clusters/%s",
		"delete":     "servers/clusters/%s",
		"resize":     "servers/clusters/%s/resize",
		"move":       "servers/moving_vms",
		"tag_attach": "servers/%s/tags",
		"tag_detach": "servers/%s/tags/%s",
	}
	return urls[action]
}

func (o *Ignite) OnSerialize(serverData map[string]interface{}, cluster *Cluster) map[string]interface{} {
	delete(serverData, "cpu")
	delete(serverData, "ram")
	serialized := map[string]map[string]interface{}{
		"ise_email":           {"value": o.AppParams["ise_email"].(string)},
		"ise_client_password": {"value": o.AppParams["ise_client_password"].(string)},
		"fault_tolerance":     {"value": o.AppParams["fault_tolerance"].(string)},
		"box_server_count":    {"value": o.AppParams["box_server_count"].(int)},
		"version":             {"value": o.AppParams["version"].(string)},
	}
	joindomain, ok := o.AppParams["joindomain"]
	if ok {
		serialized["joindomain"] = map[string]interface{}{
			"value": joindomain.(string),
		}
	}
	serverData["app_params"] = serialized
	serverData["hdd"] = map[string]int{
		"size": cluster.Disk,
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

func (o *Ignite) OnDeserialize(serverData map[string]interface{}, cluster *Cluster) {
	servers := serverData["servers"].([]interface{})
	server := servers[0].(map[string]interface{})
	params := server["app_params"].(map[string]interface{})
	if o.AppParams == nil {
		o.AppParams = make(map[string]interface{})
	}

	o.AppParams["version"] = params["version"].(map[string]interface{})["value"]
	o.AppParams["ise_email"] = params["ise_email"].(map[string]interface{})["value"]
	o.AppParams["fault_tolerance"] = params["fault_tolerance"].(map[string]interface{})["value"]
	o.AppParams["box_server_count"] = params["box_server_count"].(map[string]interface{})["value"]
	// o.AppParams["gg_client_password"] = params["gg_client_password"].(map[string]interface{})["value"]
	joindomain, ok := params["joindomain"]
	if ok {
		o.AppParams["joindomain"] = joindomain.(map[string]interface{})["value"]
	}
	// newParams := make(map[string]interface{})
	// for k, v := range params {
	// 	val := v.(map[string]interface{})
	// 	newParams[k] = val["value"]
	// }
	// paramsBytes, err := json.Marshal(newParams)
	// if err != nil {
	// 	log.Println(err)
	// }
	// err = json.Unmarshal(paramsBytes, &o.AppParams)
	// if err != nil {
	// 	log.Println(err)
	// }
	paramsBytes, err := json.Marshal(o.AppParams)
	if err != nil {
		log.Println(err)
	}
	if cluster.ResAppParams == nil {
		cluster.ResAppParams = &HCLAppParams{}
	}
	err = json.Unmarshal(paramsBytes, &cluster.ResAppParams)
	if err != nil {
		log.Println(err)
	}
}

func (o *Ignite) OnReadTF(res *schema.ResourceData, cluster *Cluster) {
	if o.AppParams == nil {
		o.AppParams = make(map[string]interface{})
	}
	o.AppParams = res.Get("app_params").(map[string]interface{})
	if o.AppParams["box_server_count"] != nil {
		o.AppParams["box_server_count"], _ = strconv.Atoi(o.AppParams["box_server_count"].(string))
	}
	volumes, ok := res.GetOk("volume")
	if ok {
		for _, v := range volumes.([]interface{}) {
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

	password := o.AppParams["ise_client_password"]
	if password == nil {
		return
	}

	vaultPasswordFileLocation := os.Getenv("DI_VAULT_PASSWORD_FILE")
	vaultPasswordFileBytes, err := ioutil.ReadFile(vaultPasswordFileLocation)
	if err != nil {
		log.Println(err)
		return
	}
	// if last byte is '\n'- remove it
	if vaultPasswordFileBytes[len(vaultPasswordFileBytes)-1] == 0x0a {
		vaultPasswordFileBytes = vaultPasswordFileBytes[:len(vaultPasswordFileBytes)-1]
	}

	passwordEncrypted := password.(string)
	passwordDecrypted, err := vault.Decrypt(passwordEncrypted, string(vaultPasswordFileBytes))
	if err != nil {
		log.Println(err)
		return
	}
	o.AppParams["ise_client_password"] = passwordDecrypted
	o.AppParams["ise_client_password_ansible"] = passwordEncrypted
}

func (o *Ignite) OnWriteTF(res *schema.ResourceData, cluster *Cluster) {
	params := map[string]interface{}{
		"ise_email":       o.AppParams["ise_email"].(string),
		"fault_tolerance": o.AppParams["fault_tolerance"].(string),
		"version":         o.AppParams["version"].(string),
		// "box_server_count": fmt.Sprintf("%.0f", o.AppParams["box_server_count"].(float64)),
		// "box_server_count": o.AppParams["box_server_count"].(int),
		// "box_server_count":   strconv.Itoa(o.AppParams["box_server_count"].(int)),
	}

	switch o.AppParams["box_server_count"].(type) {
	case float64:
		params["box_server_count"] = fmt.Sprintf("%.0f", o.AppParams["box_server_count"].(float64))
	case int:
		// params["box_server_count"] = o.AppParams["box_server_count"].(int)
		params["box_server_count"] = strconv.Itoa(o.AppParams["box_server_count"].(int))
	}

	password := o.AppParams["ise_client_password"]
	if password != nil {
		params["ise_client_password"] = o.AppParams["ise_client_password"].(string)
	}

	joindomain, ok := o.AppParams["joindomain"]
	if ok {
		params["joindomain"] = joindomain.(string)
	}
	vaultPassword := o.AppParams["ise_client_password_ansible"]
	if vaultPassword != nil {
		params["ise_client_password"] = vaultPassword.(string)
		o.AppParams["ise_client_password"] = vaultPassword.(string)
		delete(o.AppParams, "ise_client_password_ansible")
	}

	err := res.Set("app_params", params)
	if err != nil {
		log.Printf("SET ERROR: %v", err)
	}
}

func (o *Ignite) HostVars(server *Server) map[string]interface{} {
	m := map[string]interface{}{
		"ansible_host": server.Ip,
		"ansible_user": server.User,
		"dns_name":     server.DNSName,
		"name":         server.Name,
	}
	if server.Password != "" {
		vaultPasswordFileLocation := os.Getenv("DI_VAULT_PASSWORD_FILE")
		vaultPasswordFileBytes, err := ioutil.ReadFile(vaultPasswordFileLocation)
		// if last byte is '\n'- remove it
		if vaultPasswordFileBytes[len(vaultPasswordFileBytes)-1] == 0x0a {
			vaultPasswordFileBytes = vaultPasswordFileBytes[:len(vaultPasswordFileBytes)-1]
		}
		passwordEncrypted, err := vault.Encrypt(server.Password, string(vaultPasswordFileBytes))
		if err != nil {
			log.Println(err)
			return m
		}
		m["ansible_password"] = passwordEncrypted
	}
	return m
}

func (o *Ignite) GroupVars(cluster *Cluster) map[string]interface{} {
	m := make(map[string]interface{})
	password := o.AppParams["ise_client_password"]
	if password != nil {
		m["ise_client_password"] = o.AppParams["ise_client_password"].(string)
	}
	return m
}
