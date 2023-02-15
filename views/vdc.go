package views

import (
	"gitlab.gos-tech.xyz/pid/iac/terraform-provider-sberinfra/models"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

func VdcCreate(ctx context.Context, res *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	obj := models.Vdc{}

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

	//objRes := models.Vdc{}

	err = obj.ParseIdFromCreateResponse(responseBytes)

	_, err = obj.StateChange(res).WaitForStateContext(ctx)

	if err != nil {
		log.Printf("[INFO] timeout on create for instance (%s), save current state: %s", obj.ID, obj.State)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	objRes2 := models.ResVdc{}
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

type limitsImport struct {
	Limits struct {
		RamGbAmount     int `json:"ram_gb_amount"`
		CoresVcpuCount  int `json:"cores_vcpu_count"`
		StorageGbAmount int `json:"storage_gb_amount"`
	} `json:"limits"`
}

func VdcRead(ctx context.Context, res *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	//obj := models.Vdc{}
	obj := models.ResVdc{}
	obj.ReadTFRes(res)

	responseBytes, err := obj.ReadDIRes()
	if err != nil {
		return diag.FromErr(err)
	}

	err = obj.DeserializeRead(responseBytes)
	if err != nil {
		return diag.FromErr(err)
	}

	bytes, err := obj.GetVdcQuota()

	limits := map[string]*limitsImport{}
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

func VdcUpdate(ctx context.Context, res *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	obj := models.Vdc{}
	//obj := models.ResVdc{}

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
		requestBytes, err := json.Marshal(map[string]map[string]string{
			"project": {
				"action": "change_name",
				"name":   obj.Name,
			},
		})
		responseBytes, err := obj.UpdateVdcName(requestBytes)
		if err != nil {
			return diag.Errorf(err.Error(), string(responseBytes))
		}
	}

	if res.HasChange("description") {
		requestBytes, err := json.Marshal(map[string]map[string]string{
			"project": {
				"action": "change_desc",
				"desc":   obj.Desc,
			},
		})
		responseBytes, err := obj.UpdateVdcDesc(requestBytes)
		if err != nil {
			return diag.Errorf(err.Error(), string(responseBytes))
		}
	}

	if res.HasChange("limits") {
		type updateVdcLimits struct {
			GroupID uuid.UUID `json:"group_id"`
			Limits  struct {
				CoresVcpuCount  int `json:"cores_vcpu_count"`
				StorageGbAmount int `json:"storage_gb_amount"`
				RamGbAmount     int `json:"ram_gb_amount"`
			} `json:"limits"`
		}

		objVdcLimits := updateVdcLimits{}
		objVdcLimits.GroupID = obj.GroupID
		objVdcLimits.Limits.CoresVcpuCount = obj.Limits.CoresVcpuCount
		objVdcLimits.Limits.StorageGbAmount = obj.Limits.StorageGbAmount
		objVdcLimits.Limits.RamGbAmount = obj.Limits.RamGbAmount

		requestBytes, err := json.Marshal(objVdcLimits)
		responseBytes, err := obj.UpdateVdcLimits(requestBytes)

		if err != nil {
			return diag.Errorf(err.Error(), string(responseBytes))
		}
	}

	objRes := models.ResVdc{}
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

func VdcDelete(ctx context.Context, res *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	obj := models.Vdc{}
	obj.ReadTF(res)

	err := obj.DeleteDI()
	if err != nil {
		return diag.FromErr(err)
	}
	res.SetId("")
	return diags
}

func VdcImport(ctx context.Context, res *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	//obj := models.SIProject{GroupId: uuid.MustParse(res.Id())}

	obj := models.Vdc{}
	obj.ID = uuid.MustParse(res.Id())
	responseBytes, err := obj.ReadDI()
	if err != nil {
		return nil, err
	}
	err = obj.Deserialize(responseBytes)

	bytes, err := obj.GetVdcQuota()

	limits := map[string]*limitsImport{}
	err = json.Unmarshal(bytes, &limits)

	if err != nil {
		return nil, err
	}

	obj.Limits.CoresVcpuCount = limits["data"].Limits.CoresVcpuCount
	obj.Limits.RamGbAmount = limits["data"].Limits.RamGbAmount
	obj.Limits.StorageGbAmount = limits["data"].Limits.StorageGbAmount

	obj.WriteTF(res)

	//obj := models.ResVdc{}
	//obj.Vdc.ID = uuid.MustParse(res.Id())
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
