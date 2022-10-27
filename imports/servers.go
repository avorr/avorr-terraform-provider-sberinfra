package imports

//
//import (
//	"encoding/json"
//	"fmt"
//
//	"base.sw.sbc.space/pid/terraform-provider-si/client"
//	"base.sw.sbc.space/pid/terraform-provider-si/models"
//)
//
//type Servers struct {
//	Servers    []*models.Server  `json:"servers"`
//	Clusters   []*models.Cluster `json:"clusters"`
//	ClusterIds map[string]bool
//	NonCluster []*models.Server
//	Project    *models.Project
//	Api        *client.Api
//	Meta       map[string]interface{} `json:"meta"`
//}
//
//func (o *Servers) Urls(action string) string {
//	urls := map[string]string{
//		"servers": fmt.Sprintf("servers?project_id=%s", o.Project.Id.String()),
//	}
//	return urls[action]
//}
//
//func (o *Servers) Read() error {
//	responseBytes, err := o.read(o.Urls("servers"))
//	if err != nil {
//		return err
//	}
//	err = o.deserialize(responseBytes)
//	if err != nil {
//		return err
//	}
//
//	err = o.filter()
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (o *Servers) ReadCluster(id string) error {
//	responseBytes, err := o.read(fmt.Sprintf("servers/clusters/%s", id))
//	if err != nil {
//		return err
//	}
//	type ClusterResponse struct {
//		Cluster *models.Cluster `json:"cluster"`
//	}
//	response := &ClusterResponse{}
//	err = json.Unmarshal(responseBytes, response)
//	if err != nil {
//		return err
//	}
//	o.Clusters = append(o.Clusters, response.Cluster)
//	response.Cluster.SetObject()
//	return nil
//}
//
//func (o *Servers) deserialize(data []byte) error {
//	return json.Unmarshal(data, &o)
//}
//
//func (o *Servers) filter() error {
//	for _, v := range o.Servers {
//		if v.ClusterUuid.ID() != 0 {
//			if o.ClusterIds == nil {
//				o.ClusterIds = make(map[string]bool)
//			}
//			o.ClusterIds[v.ClusterUuid.String()] = true
//		} else {
//			o.NonCluster = append(o.NonCluster, v)
//			v.SetObject()
//		}
//	}
//
//	for k, _ := range o.ClusterIds {
//		err := o.ReadCluster(k)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}
//
//func (o *Servers) read(url string) ([]byte, error) {
//	responseBytes, err := o.Api.NewRequestRead(url)
//	if err != nil {
//		return nil, err
//	}
//	return responseBytes, nil
//}
