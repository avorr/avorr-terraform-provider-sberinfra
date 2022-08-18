package models

//
//import (
//	"fmt"
//	"net/url"
//	"strings"
//
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
//)
//
//type Openshift struct {
//	AppParams map[string]interface{} `json:"app_params"`
//}
//
//func (o *Openshift) GetType() string {
//	return "di_openshift"
//}
//
//func (o *Openshift) NewObj() DIResource {
//	return &Openshift{}
//}
//
//func (o *Openshift) Urls(action string) string {
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
//func (o *Openshift) OnSerialize(serverData map[string]interface{}, server *Server) map[string]interface{} {
//	delete(serverData, "disk")
//	delete(serverData, "os_version")
//	delete(serverData, "flavor")
//
//	serialized := map[string]map[string]interface{}{
//		"admin_user": {
//			"value": o.AppParams["admin_user"].(string),
//		},
//		"name_project": {
//			"value": o.AppParams["name_project"].(string),
//		},
//	}
//	joindomain, ok := o.AppParams["joindomain"]
//	if ok {
//		serialized["joindomain"] = map[string]interface{}{
//			"value": joindomain.(string),
//		}
//	}
//
//	serverData["hdd"] = map[string]int{
//		"size": 0,
//	}
//	serverData["greenfield"] = false
//	serverData["app_params"] = serialized
//	return serverData
//}
//
//func (o *Openshift) OnDeserialize(serverData map[string]interface{}, server *Server) {
//	params := serverData["app_params"].(map[string]interface{})
//	if o.AppParams == nil {
//		o.AppParams = make(map[string]interface{})
//	}
//	if params["admin_user"] != nil {
//		o.AppParams["admin_user"] = params["admin_user"].(map[string]interface{})["value"]
//	}
//	if params["name_project"] != nil {
//		o.AppParams["name_project"] = params["name_project"].(map[string]interface{})["value"]
//	}
//	joindomain, ok := params["joindomain"]
//	if ok {
//		o.AppParams["joindomain"] = joindomain.(map[string]interface{})["value"]
//	}
//	endpoint := params["endpoint"]
//	if endpoint != nil {
//		o.AppParams["endpoint"] = endpoint.(string)
//		o.apiHost()
//	}
//	outputs := serverData["outputs"]
//	if outputs != nil {
//		data := serverData["outputs"].(map[string]interface{})
//		o.AppParams["project"] = data["openshift_project_name"].(string)
//	}
//}
//
//func (o *Openshift) OnReadTF(res *schema.ResourceData, server *Server) {
//	o.AppParams = res.Get("app_params").(map[string]interface{})
//}
//
//func (o *Openshift) OnWriteTF(res *schema.ResourceData, server *Server) {
//	exclude := []string{"endpoint", "host", "project"}
//	params := make(map[string]interface{})
//	for k, v := range o.AppParams {
//		setKey := true
//		for _, field := range exclude {
//			if k == field {
//				setKey = false
//			}
//		}
//		if setKey {
//			params[k] = v
//		}
//	}
//	res.Set("app_params", params)
//}
//
//func (o *Openshift) HostVars(server *Server) map[string]interface{} {
//	m := map[string]interface{}{
//		"endpoint":           o.AppParams["endpoint"],
//		"dns_name":           server.DNSName,
//		"name":               server.Name,
//		"ansible_connection": "local",
//		// "service_name": server.ServiceName,
//	}
//	user := o.AppParams["user"]
//	if user != nil {
//		m["user"] = user
//	}
//	project := o.AppParams["project"]
//	if project != nil {
//		m["project"] = project
//	}
//	name_project := o.AppParams["name_project"]
//	if project != nil {
//		m["name_project"] = name_project
//	}
//	host := o.AppParams["host"]
//	if host != nil {
//		m["host"] = host
//	}
//	admin := o.AppParams["admin_user"]
//	if admin != nil {
//		m["admin_user"] = admin
//	}
//	password := o.AppParams["admin_password"]
//	if password != nil {
//		m["admin_password"] = password
//	}
//	return m
//}
//
//func (o *Openshift) apiHost() error {
//	if o.AppParams["endpoint"] == nil {
//		return fmt.Errorf("no endpoint")
//	}
//	port := 6443
//	urlObj, err := url.Parse(o.AppParams["endpoint"].(string))
//	if err != nil {
//		return err
//	}
//	urlObj.Host = fmt.Sprintf("%s:%d", urlObj.Host, port)
//	urlObj.Host = strings.Replace(urlObj.Host, "console-openshift-console.apps", "api", 1)
//	urlObj.Path = ""
//	o.AppParams["host"] = urlObj.String()
//	return nil
//}
//
//func (o *Openshift) GetGroup() string {
//	return o.AppParams["name_project"].(string)
//}
//
//func (o *Openshift) HCLAppParams() *HCLAppParams {
//	return &HCLAppParams{}
//}
//
//func (o *Openshift) HCLVolumes() []*HCLVolume {
//	return nil
//}
