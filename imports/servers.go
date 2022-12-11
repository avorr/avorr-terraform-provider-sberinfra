package imports

import (
	"encoding/json"
	"fmt"

	"base.sw.sbc.space/pid/terraform-provider-si/client"
	"base.sw.sbc.space/pid/terraform-provider-si/models"
)

type Servers struct {
	Servers []*models.Server `json:"servers"`
	//Clusters   []*models.Cluster `json:"clusters"`
	ClusterIds map[string]bool
	NonCluster []*models.Server
	Project    *models.Vdc
	Api        *client.Api
	Meta       map[string]interface{} `json:"meta"`
}

func (o *Servers) Urls(action string) string {
	urls := map[string]string{
		"servers": fmt.Sprintf("servers?project_id=%s", o.Project.ID.String()),
	}
	return urls[action]
}

func (o *Servers) Read() error {
	responseBytes, err := o.read(o.Urls("servers"))
	if err != nil {
		return err
	}
	err = o.deserialize(responseBytes)
	if err != nil {
		return err
	}

	//err = o.filter()
	//if err != nil {
	//	return err
	//}

	return nil
}

func (o *Servers) deserialize(data []byte) error {
	return json.Unmarshal(data, &o)
}

func (o *Servers) read(url string) ([]byte, error) {
	responseBytes, err := o.Api.NewRequestRead(url)
	if err != nil {
		return nil, err
	}
	return responseBytes, nil
}
