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

	diags = obj.ReadTF(res)

	cores := obj.Limits.CoresVcpuCount
	ram := obj.Limits.RamGbAmount
	storage := obj.Limits.StorageGbAmount

	if diags.HasError() {
		return diags
	}

	networks := res.Get("network").(*schema.Set)

	additionalNetworks := make([]interface{}, 0)

	defaultNetworkCount := 0
	//networkNames := make([]string, 0)
	for _, v := range networks.List() {
		v := v.(map[string]interface{})
		//for _, name := range networkNames {
		//	if v["name"].(string) == name {
		//		return diag.Errorf("There mustn't be networks with the same name [%s, %s]", name, v["name"])
		//	}
		//}
		//networkNames = append(networkNames, v["name"].(string))

		if v["default"] == true {
			defaultNetworkCount += 1
			//if defaultNetworkCount > 1 {
			//	return diag.Errorf("Default networks should not be more than one")
			//}
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

	//objRes := models.Project{}

	err = obj.ParseIdFromCreateResponse(responseBytes)

	_, err = obj.StateChange(res).WaitForStateContext(ctx)

	if err != nil {
		log.Printf("[INFO] timeout on create for instance (%s), save current state: %s", obj.ID, obj.State)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	objRes2 := models.ResProject{}
	objRes2.ID = obj.ID
	obj.AddNetwork(ctx, res, additionalNetworks)
	responseBytes, err = objRes2.ReadDIRes()

	if err != nil {
		return diag.FromErr(err)
	}

	err = objRes2.DeserializeRead(responseBytes)

	if err != nil {
		return diag.FromErr(err)
	}

	objRes2.Limits.CoresVcpuCount = cores
	objRes2.Limits.RamGbAmount = ram
	objRes2.Limits.StorageGbAmount = storage

	//objRes2.Limits = obj.Limits

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

	err = obj.DeserializeRead(responseBytes)
	if err != nil {
		return diag.FromErr(err)
	}

	bytes, err := obj.GetProjectQuota()

	type LimitsImport struct {
		Limits struct {
			RamGbAmount     int `json:"ram_gb_amount"`
			CoresVcpuCount  int `json:"cores_vcpu_count"`
			StorageGbAmount int `json:"storage_gb_amount"`
		} `json:"limits"`
	}

	limits := map[string]*LimitsImport{}
	err = json.Unmarshal(bytes, &limits)

	if err != nil {
		return diag.FromErr(err)
	}

	obj.Limits.CoresVcpuCount = limits["data"].Limits.CoresVcpuCount
	obj.Limits.RamGbAmount = limits["data"].Limits.RamGbAmount
	obj.Limits.StorageGbAmount = limits["data"].Limits.StorageGbAmount

	//err = obj.Deserialize(responseBytes)

	obj.WriteTFRes(res)
	return diags
}

func ProjectUpdate(ctx context.Context, res *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	obj := models.Project{}
	//obj := models.ResProject{}

	diags = obj.ReadTF(res)

	if diags.HasError() {
		return diags
	}

	//return diags
	if res.HasChange("network") {
		v1, v2 := res.GetChange("network")
		netSet1 := v1.(*schema.Set)
		netSet2 := v2.(*schema.Set)

		for _, net2 := range netSet2.List() {
			net2 := net2.(map[string]interface{})
			for _, net1 := range netSet1.List() {
				net1 := net1.(map[string]interface{})
				if net1["name"] == net2["name"] {
					net2["id"] = net1["id"]
				}
			}
		}

		//for i1, net1 := range netSet2.List() {
		//net1 := net1.(map[string]interface{})
		//netIsDefault1 := net1["default"]
		//netName1 := net1["name"]
		//for i2, net2 := range netSet2.List() {
		//net2 := net2.(map[string]interface{})
		//netIsDefault2 := net2["default"]
		//if i2-i1 > 0 && netIsDefault2 == true && netIsDefault1 == true {
		//	return diag.Errorf("Default networks shouldn't be more than one")
		//}
		//netName2 := net2["name"]
		//if i2-i1 > 0 && netName1 == netName2 {
		//	return diag.Errorf("There mustn't be networks with the same name [%s, %s]", netName1, netName2)
		//}
		//}
		//}

		//return diags

		var existNetwork func([]interface{}, string) bool
		existNetwork = func(m []interface{}, networkName string) bool {
			for _, i := range m {
				if i.(map[string]interface{})["name"] == networkName {
					return true
				}
			}
			return false
		}

		addNetSet := make([]interface{}, 0)
		removeNetSet := make([]interface{}, 0)

		for _, net := range netSet2.List() {
			net := net.(map[string]interface{})
			if !existNetwork(netSet1.List(), net["name"].(string)) {
				addNetSet = append(addNetSet, net)
			}

			if net["default"].(bool) && net["id"] != res.Get("default_network") && net["id"] != "" {
				err := obj.SetDefaultNetwork(net["id"].(string))
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}

		for _, net := range netSet1.List() {
			net := net.(map[string]interface{})
			if !existNetwork(netSet2.List(), net["name"].(string)) {
				removeNetSet = append(removeNetSet, net)
			}
		}

		//removeNetSet := netSet1.Difference(netSet2)
		//addNetSet := netSet2.Difference(netSet1)

		// add
		if len(addNetSet) > 0 {
			obj.AddNetwork(ctx, res, addNetSet)
		}

		// remove
		if len(removeNetSet) > 0 {
			for _, v := range removeNetSet {
				vol := v.(map[string]interface{})
				err := obj.DeleteNetwork(vol["id"].(string))
				if err != nil {
					return diag.FromErr(err)
				}
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
		objProjectUpdate.Project.Name = obj.Name
		requestBytes, err := json.Marshal(objProjectUpdate)
		responseBytes, err := obj.UpdateProjectName(requestBytes)
		if err != nil {
			return diag.Errorf(err.Error(), string(responseBytes))
		}
	}

	if res.HasChange("description") {
		type updateProjectDesc struct {
			Project struct {
				Action string `json:"action"`
				Desc   string `json:"desc"`
			} `json:"project"`
		}
		objProjectUpdate := updateProjectDesc{}
		objProjectUpdate.Project.Action = "change_desc"
		objProjectUpdate.Project.Desc = obj.Desc
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
				RamGbAmount     int `json:"ram_gb_amount"`
			} `json:"limits"`
		}

		objProjectLimits := updateProjectLimits{}
		objProjectLimits.GroupID = obj.GroupID
		objProjectLimits.Limits.CoresVcpuCount = obj.Limits.CoresVcpuCount
		objProjectLimits.Limits.StorageGbAmount = obj.Limits.StorageGbAmount
		objProjectLimits.Limits.RamGbAmount = obj.Limits.RamGbAmount

		requestBytes, err := json.Marshal(objProjectLimits)
		responseBytes, err := obj.UpdateProjectLimits(requestBytes)

		if err != nil {
			return diag.Errorf(err.Error(), string(responseBytes))
		}
	}

	objRes := models.ResProject{}
	objRes.ID = obj.ID
	objRes.Limits = obj.Limits
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

func ProjectImport(ctx context.Context, res *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	//obj := models.SIProject{GroupId: uuid.MustParse(res.Id())}

	obj := models.Project{}
	obj.ID = uuid.MustParse(res.Id())
	responseBytes, err := obj.ReadDI()
	if err != nil {
		return nil, err
	}
	err = obj.Deserialize(responseBytes)

	bytes, err := obj.GetProjectQuota()

	type LimitsImport struct {
		Limits struct {
			RamGbAmount     int `json:"ram_gb_amount"`
			CoresVcpuCount  int `json:"cores_vcpu_count"`
			StorageGbAmount int `json:"storage_gb_amount"`
		} `json:"limits"`
	}

	limits := map[string]*LimitsImport{}
	err = json.Unmarshal(bytes, &limits)

	if err != nil {
		return nil, err
	}

	obj.Limits.CoresVcpuCount = limits["data"].Limits.CoresVcpuCount
	obj.Limits.RamGbAmount = limits["data"].Limits.RamGbAmount
	obj.Limits.StorageGbAmount = limits["data"].Limits.StorageGbAmount

	obj.WriteTF(res)

	//obj := models.ResProject{}
	//obj.Project.ID = uuid.MustParse(res.Id())
	//responseBytes, err := obj.ReadDIRes()
	//if err != nil {
	//	return nil, err
	//}
	//err = obj.DeserializeRead(responseBytes)
	//if err != nil {
	//	return nil, err
	//}
	//obj.WriteTFRes(res)

	//
	//objBytes, _ := obj.ToHCL(nil)
	//log.Println(string(objBytes))
	//
	//index := bytes.IndexByte(objBytes, byte('{'))
	//
	//firstString := objBytes[:index+1]
	//
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
