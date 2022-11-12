package views

import (
	"base.sw.sbc.space/pid/terraform-provider-si/models"
	"base.sw.sbc.space/pid/terraform-provider-si/utils"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

//var (
//Inventory *inventory_yaml.Inventory
//)

func CreateResource(o models.DIResource) schema.CreateContextFunc {
	f := func(ctx context.Context, res *schema.ResourceData, m interface{}) diag.Diagnostics {
		var diags diag.Diagnostics

		newObj := o.NewObj()

		server := models.Server{Object: newObj}
		// read from TF state
		server.ReadTF(res)

		err := server.GetPubKey()
		if err != nil {
			return diag.FromErr(err)
		}

		// serialize for request
		// requestBytes, err := server.Object.Serialize(&server)
		requestBytes, err := server.Serialize()
		if err != nil {
			return diag.FromErr(err)
		}

		// bb := bytes.Buffer{}
		// json.Indent(&bb, requestBytes, "", "\t")
		// log.Println(bb.String())
		// return nil

		// send request to DI (new server)
		// responseBytes, err := server.DI("create", requestBytes, 201)
		responseBytes, err := server.CreateDI(requestBytes)
		if err != nil {
			return diag.FromErr(err)
		}

		// get id from response
		err = server.ParseIdFromCreateResponse(responseBytes)
		if err != nil {
			return diag.FromErr(err)
		}

		// send request to DI (get info)
		responseBytes, err = server.ReadDI()
		if err != nil {
			return diag.FromErr(err)
		}

		// deserialize response to obj
		// err = server.Object.Deserialize(responseBytes, &server)
		err = server.Deserialize(responseBytes)
		if err != nil {
			return diag.FromErr(err)
		}
		server.WriteTF(res)

		// wait for running status
		_, err = server.StateChange(res).WaitForStateContext(ctx)
		if err != nil {
			log.Printf("[INFO] timeout on create for instance (%s), save current state: %s", server.Id.String(), server.State)
		}
		// if err != nil {
		// 	diags = append(diags, diag.Diagnostic{
		// 		Severity: diag.Warning,
		// 		Detail:   fmt.Sprintf("timeout on create for instance (%s), save current state: %s", server.Id, server.State),
		// 	})
		// 	return diag.FromErr(err)
		// }

		if server.State == "damaged" {
			err = server.DeleteDI()
			if err != nil {
				return diag.FromErr(err)
			}
			res.SetId("")
			return diag.Errorf(
				"server state: %s, remove server: %s %s, err_msg: %s",
				server.State,
				server.Id,
				server.ServiceName,
				server.ErrMsg,
			)
		}

		//groupName := newObj.GetType()[3:]
		//group := Inventory.All.GetGroup(groupName)
		//if group == nil {
		//	group = &inventory_yaml.Group{Name: groupName}
		//}

		//subgroupName := server.GetGroup()
		//subgroup := group.GetGroup(subgroupName)
		//if subgroup == nil {
		//	subgroup = &inventory_yaml.Group{Name: subgroupName}
		//	group.AddGroup(subgroup)
		//}
		//
		//if subgroup.Vars == nil {
		//	subgroup.Vars = make(map[string]interface{})
		//}
		//subgroup.Vars["service_name_en"] = utils.Reformat(server.ServiceName)
		//subgroup.Vars["service_name_ru"] = server.ServiceName

		//host := &inventory_yaml.Host{
		//// Name: server.Id.String(),
		//Name: server.DNSName,
		//Vars: newObj.HostVars(&server),
		//}
		//subgroup.AddHost(host)
		//Inventory.All.AddGroup(group)
		//err = Inventory.Save()
		//if err != nil {
		//	return diag.FromErr(err)
		//}
		//err = Inventory.ToBIN()
		//if err != nil {
		//	return diag.FromErr(err)
		//}

		// attach tags
		if res.HasChange("tag_ids") {
			_, tagIds := res.GetChange("tag_ids")
			tagSet := tagIds.(*schema.Set)
			for _, v := range tagSet.List() {
				_, err := server.TagAttachDI(v.(string))
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}
		return diags
	}
	return f
}

var IsImport bool

func ReadResource(obj models.DIResource) schema.ReadContextFunc {
	f := func(ctx context.Context, res *schema.ResourceData, m interface{}) diag.Diagnostics {
		var diags diag.Diagnostics
		// server := models.Server{Object: &models.VM{}}
		newObj := obj.NewObj()
		server := models.Server{Object: newObj, IsImport: IsImport}
		server.ReadTF(res)

		err := server.GetPubKey()
		if err != nil {
			return diag.FromErr(err)
		}

		responseBytes, err := server.ReadDI()
		if err != nil {
			return diag.FromErr(err)
		}
		err = server.Deserialize(responseBytes)
		if err != nil {
			return diag.FromErr(err)
		}
		server.WriteTF(res)
		// if obj.State == "damaged" {
		// 	err = obj.Delete()
		// 	if err != nil {
		// 		return diag.FromErr(err)
		// 	}
		// 	res.SetId("")
		// }

		//groupName := obj.GetType()[3:]
		//group := Inventory.All.GetGroup(groupName)
		//if group == nil {
		//	group = &inventory_yaml.Group{Name: groupName}
		//}
		//
		//subgroupName := server.GetGroup()
		//// subgroupName := obj.GetGroup()
		//// subgroupName := utils.Reformat(server.ServiceName)
		//subgroup := group.GetGroup(subgroupName)
		//if subgroup == nil {
		//	subgroup = &inventory_yaml.Group{Name: subgroupName}
		//	group.AddGroup(subgroup)
		//}
		//
		//if subgroup.Vars == nil {
		//	subgroup.Vars = make(map[string]interface{})
		//}
		//subgroup.Vars["service_name_en"] = utils.Reformat(server.ServiceName)
		//subgroup.Vars["service_name_ru"] = server.ServiceName
		//
		//host := &inventory_yaml.Host{
		//	// Name: server.Id.String(),
		//	Name: server.DNSName,
		//	Vars: newObj.HostVars(&server),
		//}
		//subgroup.AddHost(host)
		//Inventory.All.AddGroup(group)
		//err = Inventory.Save()
		//if err != nil {
		//	return diag.FromErr(err)
		//}
		//err = Inventory.ToBIN()
		//if err != nil {
		//	return diag.FromErr(err)
		//}

		return diags
	}
	return f
}

func UpdateResource(obj models.DIResource) schema.UpdateContextFunc {
	f := func(ctx context.Context, res *schema.ResourceData, m interface{}) diag.Diagnostics {
		var diags diag.Diagnostics

		server := models.Server{Object: obj}
		server.ReadTF(res)

		if res.Get("state") == "creating" || res.Get("state_resize") == "resizing" {
			return diag.FromErr(errors.New("can't update 'creating' or 'resizing' instance"))
		}
		if res.HasChange("disk") {
			return diag.FromErr(errors.New("main disk change not implemented in api"))
		}

		if res.HasChange("service_name") {
			_, service_name := res.GetChange("service_name")
			changes := map[string]map[string]string{
				"server": {
					"service_name": service_name.(string),
				},
			}
			objBytes, err := json.Marshal(changes)
			if err != nil {
				return diag.FromErr(err)
			}
			responseBytes, err := server.UpdateDI(objBytes)
			if err != nil {
				return diag.FromErr(err)
			}
			err = server.Deserialize(responseBytes)
			if err != nil {
				return diag.FromErr(err)
			}
			server.WriteTF(res)
		}

		if res.HasChange("project_id") {
			_, projectId := res.GetChange("project_id")
			changes := map[string]interface{}{
				"uuids":        []string{server.Id.String()},
				"project_uuid": projectId.(string),
			}
			objBytes, err := json.Marshal(changes)
			if err != nil {
				return diag.FromErr(err)
			}
			_, err = server.MoveDI(objBytes)
			if err != nil {
				return diag.FromErr(err)
			}
			server.ProjectId = uuid.MustParse(projectId.(string))
			server.WriteTF(res)
		}

		if res.HasChange("flavor") {
			_, flavor := res.GetChange("flavor")
			changes := map[string]interface{}{
				"resize": map[string]string{
					"flavor": flavor.(string),
				},
				"kvr_timeout": 1,
			}
			objBytes, err := json.Marshal(changes)
			if err != nil {
				return diag.FromErr(err)
			}

			_, err = server.ResizeDI(objBytes)
			if err != nil {
				return diag.FromErr(err)
			}
			server.StateResize = "resizing"
			server.WriteTF(res)

			_, err = server.StateResizeChange(res).WaitForStateContext(ctx)
			if err != nil {
				log.Printf("[INFO] timeout on resize for instance (%s), save current state: %s", server.Id.String(), server.StateResize)
			}
		}

		// openshift only
		//if server.Object.GetType() == "di_openshift" && (res.HasChange("ram") || res.HasChange("cpu")) {
		//	_, ram := res.GetChange("ram")
		//	_, cpu := res.GetChange("cpu")
		//	changes := map[string]map[string]int{
		//		"resize": {
		//			"ram": ram.(int),
		//			"cpu": cpu.(int),
		//		},
		//	}
		//	objBytes, err := json.Marshal(changes)
		//	if err != nil {
		//		return diag.FromErr(err)
		//	}
		//	_, err = server.ResizeDI(objBytes)
		//	if err != nil {
		//		return diag.FromErr(err)
		//	}
		//	server.StateResize = "resizing"
		//	server.WriteTF(res)
		//
		//	_, err = server.StateResizeChange(res).WaitForStateContext(ctx)
		//	if err != nil {
		//		log.Printf("[INFO] timeout on resize for instance (%s), save current state: %s", server.Id.String(), server.StateResize)
		//	}
		//}

		if res.HasChange("volume") {
			v1, v2 := res.GetChange("volume")
			volSet1 := v1.(*schema.Set)
			volSet2 := v2.(*schema.Set)
			removeVolSet := volSet1.Difference(volSet2)
			addVolSet := volSet2.Difference(volSet1)

			// add
			if addVolSet.Len() > 0 {
				changes := map[string][]interface{}{
					"volumes": addVolSet.List(),
				}
				objBytes, err := json.Marshal(changes)
				if err != nil {
					return diag.FromErr(err)
				}
				responseBytes, err := server.VolumeCreateDI(objBytes)
				if err != nil {
					return diag.FromErr(err)
				}
				err = server.Deserialize(responseBytes)
				if err != nil {
					return diag.FromErr(err)
				}
				server.WriteTF(res)
				_, err = server.StateResizeChange(res).WaitForStateContext(ctx)
				if err != nil {
					log.Printf("[INFO] timeout on add volume for instance (%s), save current state: %s", server.Id.String(), server.StateResize)
				}
			}

			// remove
			if removeVolSet.Len() > 0 {
				for _, v := range removeVolSet.List() {
					vol := v.(map[string]interface{})
					change := map[string]string{
						"volume_uuid": vol["volume_id"].(string),
					}
					objBytes, err := json.Marshal(change)
					if err != nil {
						return diag.FromErr(err)
					}
					_, err = server.VolumeRemoveDI(objBytes)
					if err != nil {
						return diag.FromErr(err)
					}
					_, err = server.StateResizeChange(res).WaitForStateContext(ctx)
				}
				server.WriteTF(res)
			}
		}

		if res.HasChange("tag_ids") {
			v1, v2 := res.GetChange("tag_ids")
			tagSet1 := v1.(*schema.Set)
			tagSet2 := v2.(*schema.Set)
			l1 := utils.ArrInterfaceToArrStr(tagSet1.List())
			l2 := utils.ArrInterfaceToArrStr(tagSet2.List())

			for _, v := range l1 {
				if !utils.ArrContainsStr(l2, v) {
					err := server.TagDetachDI(v)
					if err != nil {
						return diag.FromErr(err)
					}
				}
			}
			for _, v := range l2 {
				if !utils.ArrContainsStr(l1, v) {
					_, err := server.TagAttachDI(v)
					if err != nil {
						return diag.FromErr(err)
					}
				}
			}
			server.WriteTF(res)
		}
		return diags
	}
	return f
}

func DeleteResource(obj models.DIResource) schema.DeleteContextFunc {
	f := func(ctx context.Context, res *schema.ResourceData, m interface{}) diag.Diagnostics {
		var diags diag.Diagnostics
		server := models.Server{Object: obj}
		server.ReadTF(res)

		if res.Get("state") == "creating" || res.Get("state_resize") == "resizing" {
			return diag.FromErr(errors.New("can't delete 'creating' or 'resizing' instance"))
		}
		//server.Id = uuid.MustParse(res.Id())
		err := server.DeleteVM()
		//os.Exit(3)
		//responseBytes, err := server.ReadDI()
		//objRes := models.Server{Object: obj}
		//err = json.Unmarshal(responseBytes, &objRes)
		//if err != nil {
		//	return diag.FromErr(err)
		//}
		//objRes.Id = uuid.MustParse(res.Id())

		//return diags
		//_, err := objRes.StateChange(res).WaitForStateContext(ctx)
		_, err = server.StateChange(res).WaitForStateContext(ctx)
		if err != nil {
			log.Printf("[INFO] timeout on remove for instance (%s), save current state: %s", server.Id.String(), server.State)
		}
		//os.Exit(3)
		if err != nil {
			return diag.FromErr(err)
		}
		res.SetId("")

		//groupName := obj.GetType()[3:]

		//group := Inventory.All.GetGroup(groupName)
		//if group == nil {
		//	group = &inventory_yaml.Group{Name: groupName}
		//}
		//subgroupName := server.GetGroup()
		//// subgroupName := utils.Reformat(server.ServiceName)
		//subgroup := group.GetGroup(subgroupName)
		//if subgroup != nil {
		//	if len(subgroup.Hosts) == 1 {
		//		group.RmGroup(subgroup.Name)
		//	} else {
		//		subgroup.RmHost(server.DNSName)
		//	}
		//}
		//if len(group.Hosts) == 0 && len(group.Children) == 0 {
		//	Inventory.All.RmGroup(group.Name)
		//}
		//
		//err = Inventory.Save()
		//if err != nil {
		//	return diag.FromErr(err)
		//}
		//err = Inventory.ToBIN()
		//if err != nil {
		//	return diag.FromErr(err)
		//}

		return diags
	}
	return f
}

func ImportResource(obj models.DIResource) schema.StateContextFunc {
	return func(ctx context.Context, res *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
		// state := res.State()
		IsImport = true
		server := &models.Server{Object: obj, Id: uuid.MustParse(res.Id())}
		err := server.GetPubKey()
		if err != nil {
			return nil, err
		}
		responseBytes, err := server.ReadDI()
		if err != nil {
			return nil, err
		}
		err = server.Deserialize(responseBytes)
		if err != nil {
			return nil, err
		}
		// log.Println(pp.Sprint(server))

		hclRoot := server.GetHCLRoot()
		hclRoot.Resources.Type = obj.GetType()
		hclRoot.Resources.AppParams = obj.HCLAppParams()
		hclRoot.Resources.Volumes = obj.HCLVolumes()
		objBytes := server.GetHCLRootBytes(hclRoot)
		index := bytes.IndexByte(objBytes, byte('{'))
		firstString := objBytes[:index+1]
		if firstString[0] == byte(0x0a) {
			firstString = firstString[1:]
		}
		appParamsFind := []byte("app_params {")
		appParamsReplace := []byte("app_params = {")
		objBytes = bytes.Replace(objBytes, appParamsFind, appParamsReplace, 1)
		// log.Println(string(objBytes))

		// TODO: file names
		//fileName := "imports.tf"
		//fileBytes, err := ioutil.ReadFile(fileName)
		//if err != nil {
		//	return nil, err
		//}

		//toReplace := []byte(fmt.Sprintf("%s}", firstString))
		//newBytes := bytes.Replace(fileBytes, toReplace, objBytes, -1)
		//err = ioutil.WriteFile(fileName, newBytes, 0600)
		//if err != nil {
		//	return nil, err
		//}
		if len(server.TagIds) > 0 {
			var tags []string
			for _, v := range server.TagIds {
				tags = append(tags, v.String())
			}
			err = res.Set("tag_ids", tags)
			if err != nil {
				return nil, err
			}
		}
		// return nil, nil
		return []*schema.ResourceData{res}, nil
	}
}
