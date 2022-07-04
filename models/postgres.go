package models

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vault "github.com/sosedoff/ansible-vault-go"
)

type Postgres struct {
	AppParams map[string]interface{} `json:"app_params"`
	Volumes   []*Volume              `json:"volumes"`
}

func (o *Postgres) GetType() string {
	return "di_postgres"
}

func (o *Postgres) NewObj() DIResource {
	return &Postgres{}
}

func (o *Postgres) Urls(action string) string {
	urls := map[string]string{
		"create":     "servers",
		"read":       "servers/%s",
		"update":     "servers/%s",
		"delete":     "servers/%s",
		"resize":     "servers/%s/resize",
		"move":       "servers/moving_vms",
		"tag_attach": "servers/%s/tags",
		"tag_detach": "servers/%s/tags/%s",
	}
	return urls[action]
}

func (o *Postgres) OnSerialize(serverData map[string]interface{}, server *Server) map[string]interface{} {
	delete(serverData, "cpu")
	delete(serverData, "ram")
	serialized := map[string]map[string]interface{}{
		"version": {
			"value": o.AppParams["version"].(string),
		},
		"postgres_db_name": {
			"value": o.AppParams["postgres_db_name"].(string),
		},
		"postgres_db_user": {
			"value": o.AppParams["postgres_db_user"].(string),
		},
		"postgres_db_password": {
			"value": o.AppParams["postgres_db_password"].(string),
		},
	}
	joindomain, ok := o.AppParams["joindomain"]
	if ok {
		serialized["joindomain"] = map[string]interface{}{
			"value": joindomain.(string),
		}
	}
	maxconnections, ok := o.AppParams["max_connections"]
	if ok {
		serialized["max_connections"] = map[string]interface{}{
			"value": maxconnections.(string),
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

func (o *Postgres) OnDeserialize(serverData map[string]interface{}, server *Server) {
	params := serverData["app_params"].(map[string]interface{})
	if o.AppParams == nil {
		o.AppParams = make(map[string]interface{})
	}
	o.AppParams["version"] = params["version"].(map[string]interface{})["value"]
	o.AppParams["postgres_db_name"] = params["postgres_db_name"].(map[string]interface{})["value"]
	o.AppParams["postgres_db_user"] = params["postgres_db_user"].(map[string]interface{})["value"]
	// o.AppParams["postgres_db_password"] = params["postgres_db_password"].(map[string]interface{})["value"]
	joindomain, ok := params["joindomain"]
	if ok {
		o.AppParams["joindomain"] = joindomain.(map[string]interface{})["value"]
	}
	maxconnections, ok := params["max_connections"]
	if ok {
		o.AppParams["max_connections"] = maxconnections.(map[string]interface{})["value"]
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

func (o *Postgres) OnReadTF(res *schema.ResourceData, server *Server) {
	o.AppParams = res.Get("app_params").(map[string]interface{})
	volumes, ok := res.GetOk("volume")
	if ok {
		volumeSet := volumes.(*schema.Set)
		for _, v := range volumeSet.List() {
			values := v.(map[string]interface{})
			volume := &Volume{
				Size: values["size"].(int),
				Path: values["path"].(string),
			}
			o.Volumes = append(o.Volumes, volume)
		}
		sort.Sort(ByPath(o.Volumes))
	}

	password := o.AppParams["postgres_db_password"]
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

	passwordEncrypted := o.AppParams["postgres_db_password"].(string)
	passwordDecrypted, err := vault.Decrypt(passwordEncrypted, string(vaultPasswordFileBytes))
	if err != nil {
		log.Println(err)
		return
	}
	o.AppParams["postgres_db_password"] = passwordDecrypted
	o.AppParams["postgres_db_password_ansible"] = passwordEncrypted
}

func (o *Postgres) OnWriteTF(res *schema.ResourceData, server *Server) {
	vaultPassword := o.AppParams["postgres_db_password_ansible"]
	if vaultPassword != nil {
		o.AppParams["postgres_db_password"] = vaultPassword
		delete(o.AppParams, "postgres_db_password_ansible")
	}
	res.Set("app_params", o.AppParams)
}

func (o *Postgres) HostVars(server *Server) map[string]interface{} {
	return map[string]interface{}{
		"ansible_host":         server.Ip,
		"ansible_user":         server.User,
		"dns_name":             server.DNSName,
		"name":                 server.Name,
		"postgres_db_name":     o.AppParams["postgres_db_name"],
		"postgres_db_user":     o.AppParams["postgres_db_user"],
		"postgres_db_password": o.AppParams["postgres_db_password"],
	}
}

func (o *Postgres) GetGroup() string {
	return ""
}

func (o *Postgres) HCLAppParams() *HCLAppParams {
	return &HCLAppParams{}
}

func (o *Postgres) HCLVolumes() []*HCLVolume {
	return nil
}
