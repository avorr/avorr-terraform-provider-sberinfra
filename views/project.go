package views

import (
	"base.sw.sbc.space/pid/terraform-provider-si/models"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

func ProjectCreate(ctx context.Context, res *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	obj := models.Project{}

	obj.ReadTF(res)

	networks := res.Get("network").(*schema.Set).List()

	additionalNetworks := make([]interface{}, 0)

	defaultNetworkCount := 0
	networkNames := make([]string, 0)
	for _, v := range networks {
		v := v.(map[string]interface{})
		for _, name := range networkNames {
			if v["network_name"].(string) == name {
				return diag.Errorf("There mustn't be networks with the same name [%s, %s]", name, v["network_name"])
			}
		}
		networkNames = append(networkNames, v["network_name"].(string))

		if v["is_default"] == true {
			defaultNetworkCount += 1
			if defaultNetworkCount > 1 {
				return diag.Errorf("Default networks should not be more than one")
			}
		} else {
			additionalNetworks = append(additionalNetworks, v)
		}
	}

	requestBytes, err := obj.Serialize()

	if err != nil {
		return diag.FromErr(err)
	}
	responseBytes, err := obj.CreateDI(requestBytes)
	if err != nil {
		return diag.FromErr(err)
	}

	objRes := models.Project{}

	err = objRes.ParseIdFromCreateResponse(responseBytes)

	_, err = objRes.StateChange(res).WaitForStateContext(ctx)

	if err != nil {
		log.Printf("[INFO] timeout on create for instance (%s), save current state: %s", objRes.Project.ID, objRes.Project.State)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	objRes2 := models.ResProject{}
	objRes2.Project.ID = objRes.Project.ID
	objRes.AddNetwork(ctx, res, additionalNetworks)
	responseBytes, err = objRes2.ReadDIRes()

	if err != nil {
		return diag.FromErr(err)
	}

	err = objRes2.DeserializeRead(responseBytes)

	if err != nil {
		return diag.FromErr(err)
	}
	objRes2.WriteTFRes(res)
	return diags
}

func ProjectRead(ctx context.Context, res *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	//obj := models.Project{}
	obj := models.ResProject{}
	obj.ReadTFRes(res)

	responseBytes, err := obj.ReadDIRes()

	if err != nil {
		return diag.FromErr(err)
	}

	//err = obj.Deserialize(responseBytes)
	err = obj.DeserializeRead(responseBytes)

	if err != nil {
		return diag.FromErr(err)
	}

	obj.WriteTFRes(res)
	return diags
}

func ProjectUpdate(ctx context.Context, res *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	obj := models.Project{}
	//obj := models.ResProject{}
	obj.ReadTF(res)

	if res.HasChange("network") {
		v1, v2 := res.GetChange("network")
		netSet1 := v1.(*schema.Set)
		netSet2 := v2.(*schema.Set)

		for i1, net1 := range netSet2.List() {
			net1 := net1.(map[string]interface{})
			netIsDefault1 := net1["is_default"]
			netName1 := net1["network_name"]
			for i2, net2 := range netSet2.List() {
				net2 := net2.(map[string]interface{})
				netIsDefault2 := net2["is_default"]
				if i2-i1 > 0 && netIsDefault2 == netIsDefault1 {
					return diag.Errorf("Default networks shouldn't be more than one")
				}
				netName2 := net2["network_name"]
				if i2-i1 > 0 && netName1 == netName2 {
					return diag.Errorf("There mustn't be networks with the same name [%s, %s]", netName1, netName2)
				}
			}
		}

		removeNetSet := netSet1.Difference(netSet2)
		addNetSet := netSet2.Difference(netSet1)

		log.Println("NETSET2", netSet2)
		log.Println("NETSET1", netSet1)

		return diags

		// add
		if addNetSet.Len() > 0 {
			obj.AddNetwork(ctx, res, addNetSet.List())
		}

		// remove
		if removeNetSet.Len() > 0 {
			for _, v := range removeNetSet.List() {
				vol := v.(map[string]interface{})
				err := obj.DeleteNetwork(vol["network_uuid"].(string))
				log.Println(err.Error())

				//if err != nil {
				//	return diag.FromErr(err)
				//}
			}
		}
	}

	if res.HasChange("name") {
		type updateProjectName struct {
			Project struct {
				Action string `json:"action"`
				Name   string `json:"name"`
			} `json:"project"`
		}
		objProjectUpdate := updateProjectName{}
		objProjectUpdate.Project.Action = "change_name"
		objProjectUpdate.Project.Name = obj.Project.Name
		requestBytes, err := json.Marshal(objProjectUpdate)
		responseBytes, err := obj.UpdateProjectName(requestBytes)
		if err != nil {
			return diag.Errorf(err.Error(), string(responseBytes))
		}
	}

	if res.HasChange("desc") {
		type updateProjectDesc struct {
			Project struct {
				Action string `json:"action"`
				Desc   string `json:"desc"`
			} `json:"project"`
		}
		objProjectUpdate := updateProjectDesc{}
		objProjectUpdate.Project.Action = "change_desc"
		objProjectUpdate.Project.Desc = obj.Project.Desc
		requestBytes, err := json.Marshal(objProjectUpdate)
		responseBytes, err := obj.UpdateProjectDesc(requestBytes)
		if err != nil {
			return diag.Errorf(err.Error(), string(responseBytes))
		}
	}

	if res.HasChange("limits") {
		type updateProjectLimits struct {
			GroupID uuid.UUID `json:"group_id"`
			Limits  struct {
				CoresVcpuCount  int `json:"cores_vcpu_count"`
				StorageGbAmount int `json:"storage_gb_amount"`
				RAMGbAmount     int `json:"ram_gb_amount"`
			} `json:"limits"`
		}

		objProjectLimits := updateProjectLimits{}
		objProjectLimits.GroupID = obj.Project.GroupID
		objProjectLimits.Limits.CoresVcpuCount = obj.Project.Limits.CoresVcpuCount
		objProjectLimits.Limits.StorageGbAmount = obj.Project.Limits.StorageGbAmount
		objProjectLimits.Limits.RAMGbAmount = obj.Project.Limits.RAMGbAmount

		requestBytes, err := json.Marshal(objProjectLimits)
		responseBytes, err := obj.UpdateProjectLimits(requestBytes)

		if err != nil {
			return diag.Errorf(err.Error(), string(responseBytes))
		}
	}

	objRes := models.ResProject{}
	objRes.Project.ID = obj.Project.ID
	responseBytes, err := objRes.ReadDIRes()
	if err != nil {
		return diag.FromErr(err)
	}
	err = objRes.DeserializeRead(responseBytes)
	if err != nil {
		return diag.FromErr(err)
	}
	objRes.WriteTFRes(res)

	return diags

	//requestBytes, err := obj.Serialize()
	//if err != nil {
	//	return diag.FromErr(err)
	//}

	//responseBytes, err := obj.UpdateDI(requestBytes)
	//if err != nil {
	//	return diag.FromErr(err)
	//}

	//err = obj.Deserialize(responseBytes)
	//if err != nil {
	//	return diag.FromErr(err)
	//}

	//obj.WriteTF(res)
	//return diags
}

func ProjectDelete(ctx context.Context, res *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	obj := models.Project{}
	obj.ReadTF(res)

	err := obj.DeleteDI()
	if err != nil {
		return diag.FromErr(err)
	}
	res.SetId("")
	return diags
}

func SiProjectImport(ctx context.Context, res *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	//obj := models.SIProject{GroupId: uuid.MustParse(res.Id())}
	//responseBytes, err := obj.ReadDI()
	//if err != nil {
	//	return nil, err
	//}
	//err = obj.Deserialize(responseBytes)
	//if err != nil {
	//	return nil, err
	//}
	//obj.WriteTF(res)

	//objBytes, _ := obj.ToHCL(nil)
	// log.Println(string(objBytes))

	//index := bytes.IndexByte(objBytes, byte('{'))

	//firstString := objBytes[:index+1]

	//fileBytes, err := ioutil.ReadFile("project.tf")
	//if err != nil {
	//	return nil, err
	//}

	//toReplace := []byte(fmt.Sprintf("%s}", firstString))

	//newBytes := bytes.Replace(fileBytes, toReplace, objBytes, -1)

	//err = ioutil.WriteFile("project.tf", newBytes, 0600)
	//if err != nil {
	//	return nil, err
	//}

	return []*schema.ResourceData{res}, nil
}
