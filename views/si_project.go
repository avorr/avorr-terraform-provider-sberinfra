package views

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"stash.sigma.sbrf.ru/sddevops/terraform-provider-di/models"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func SiProjectCreate(ctx context.Context, res *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	obj := models.SIProject{}

	obj.ReadTF(res)

	networks := res.Get("network").(*schema.Set).List()

	log.Printf("!!NET %v, %T", networks, networks)

	defaultNetworkCount := 0
	for _, v := range networks {
		if v.(map[string]interface{})["is_default"] == true {
			defaultNetworkCount += 1
			if defaultNetworkCount > 1 {
				return diag.Errorf("Default networks should not be more than one")
			}
		} else {
			log.Printf("!!%v, %T\n", v.(map[string]interface{})["is_default"], v)
		}
	}

	log.Println("DNC", defaultNetworkCount)

	//return diags
	requestBytes, err := obj.Serialize()

	log.Println("!!RB", requestBytes)
	log.Println("!!RBS", string(requestBytes))

	if err != nil {
		return diag.FromErr(err)
	}

	responseBytes, err := obj.CreateDI(requestBytes)

	if err != nil {
		return diag.FromErr(err)
	}

	objRes := models.SIProject{}

	err = objRes.ParseIdFromCreateResponse(responseBytes)

	_, err = objRes.StateChange(res).WaitForStateContext(ctx)
	if err != nil {
		log.Printf("[INFO] timeout on create for instance (%s), save current state: %s", objRes.Project.ID.String(), objRes.Project.State)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	responseBytes, err = objRes.ReadDI()

	if err != nil {
		return diag.FromErr(err)
	}

	//err = objRes.Deserialize(responseBytes)
	//if err != nil {
	//	return diag.FromErr(err)
	//}
	log.Println("RRES", res.Get("network"))
	objRes.WriteTF(res)
	return diags
}

func SiProjectRead(ctx context.Context, res *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("##", "read_func")
	log.Println("##", res.Id())

	var diags diag.Diagnostics

	obj := models.SIProject{}
	obj.ReadTF(res)

	responseBytes, err := obj.ReadDI()
	if err != nil {
		return diag.FromErr(err)
	}

	err = obj.Deserialize(responseBytes)
	if err != nil {
		return diag.FromErr(err)
	}

	obj.WriteTF(res)
	return diags
}

func SiProjectUpdate(ctx context.Context, res *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	obj := models.SIProject{}
	obj.ReadTF(res)

	//log.Println("HC-ir_group", res.HasChange("ir_group"))
	//log.Println("HC-type", res.HasChange("type"))
	//log.Println("HC-vdc_openstack", res.HasChange("vdc_openstack"))
	//log.Println("HC-openstack", res.HasChange("openstack"))
	//log.Println("HC-name", res.HasChange("name"))
	//log.Println("HC-group_id", res.HasChange("group_id"))
	//log.Println("HC-datacenter", res.HasChange("datacenter"))
	//log.Println("HC-jump_host", res.HasChange("jump_host"))
	//log.Println("HC-limits", res.HasChange("limits"))

	log.Println("HC-network", res.HasChange("network"))
	log.Println("")

	if res.HasChange("name") {
		type updateSIProjectName struct {
			Project struct {
				Action string `json:"action"`
				Name   string `json:"name"`
			} `json:"project"`
		}
		objSIProjectUpdate := updateSIProjectName{}
		objSIProjectUpdate.Project.Action = "change_name"
		objSIProjectUpdate.Project.Name = obj.Project.Name
		requestBytes, err := json.Marshal(objSIProjectUpdate)
		responseBytes, err := obj.UpdateSIProjectName(requestBytes)
		if err != nil {
			return diag.FromErr(err)
			log.Println(responseBytes)
		}
	}

	if res.HasChange("desc") {
		type updateSIProjectDesc struct {
			Project struct {
				Action string `json:"action"`
				Desc   string `json:"desc"`
			} `json:"project"`
		}
		objSIProjectUpdate := updateSIProjectDesc{}
		objSIProjectUpdate.Project.Action = "change_desc"
		objSIProjectUpdate.Project.Desc = obj.Project.Desc
		requestBytes, err := json.Marshal(objSIProjectUpdate)
		responseBytes, err := obj.UpdateSIProjectDesc(requestBytes)
		if err != nil {
			return diag.FromErr(err)
			log.Println(responseBytes)
		}
	}

	if res.HasChange("limits") {
		type updateSIProjectsLimits struct {
			GroupID uuid.UUID `json:"group_id"`
			Limits  struct {
				CoresVcpuCount  int `json:"cores_vcpu_count"`
				StorageGbAmount int `json:"storage_gb_amount"`
				RAMGbAmount     int `json:"ram_gb_amount"`
			} `json:"limits"`
		}

		objSIProjectLimits := updateSIProjectsLimits{}
		objSIProjectLimits.GroupID = obj.Project.GroupID
		objSIProjectLimits.Limits.CoresVcpuCount = obj.Project.Limits.CoresVcpuCount
		objSIProjectLimits.Limits.StorageGbAmount = obj.Project.Limits.StorageGbAmount
		objSIProjectLimits.Limits.RAMGbAmount = obj.Project.Limits.RAMGbAmount

		requestBytes, err := json.Marshal(objSIProjectLimits)
		responseBytes, err := obj.UpdateSIProjectLimits(requestBytes)

		if err != nil {
			return diag.FromErr(err)
			log.Println(responseBytes)
		}

		log.Println(res.HasChange("limits"))
		log.Println("LIM", res.Get("limits"))

	}

	//log.Println("RESID", res.Id())
	//log.Println("RESNAME", res.Get("name"))
	//log.Println("RESNAME", obj.Project.Name)
	//log.Println("RES", reflect.TypeOf(res))

	return diags

	requestBytes, err := obj.Serialize()
	if err != nil {
		return diag.FromErr(err)
	}

	responseBytes, err := obj.UpdateDI(requestBytes)
	if err != nil {
		return diag.FromErr(err)
	}

	err = obj.Deserialize(responseBytes)
	if err != nil {
		return diag.FromErr(err)
	}

	obj.WriteTF(res)
	return diags
}

func SiProjectDelete(ctx context.Context, res *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	obj := models.SIProject{}
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
