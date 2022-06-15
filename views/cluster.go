package views

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strconv"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	inventory_yaml "stash.sigma.sbrf.ru/sddevops/terraform-provider-di/inventory-yaml"
	"stash.sigma.sbrf.ru/sddevops/terraform-provider-di/models"
	"stash.sigma.sbrf.ru/sddevops/terraform-provider-di/utils"
)

func CreateClusterResource(o models.DIClusterResource) schema.CreateContextFunc {
	f := func(ctx context.Context, res *schema.ResourceData, m interface{}) diag.Diagnostics {
		var diags diag.Diagnostics

		objNew := o.NewObj()
		cluster := models.Cluster{Object: objNew}
		cluster.ReadTF(res)

		err := cluster.GetPubKey()
		if err != nil {
			return diag.FromErr(err)
		}

		requestBytes, err := cluster.Serialize()
		if err != nil {
			return diag.FromErr(err)
		}

		// bb := bytes.Buffer{}
		// json.Indent(&bb, requestBytes, "", "\t")
		// log.Println(bb.String())
		// return nil

		responseBytes, err := cluster.CreateDI(requestBytes)
		if err != nil {
			return diag.FromErr(err)
		}

		serverIdString, err := cluster.ParseIdFromCreateResponse(responseBytes)
		if err != nil {
			return diag.FromErr(err)
		}

		server := models.Server{
			Id:     uuid.MustParse(serverIdString),
			Object: &models.VM{},
		}

		_, err = server.StateClusterChange(res).WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("timeout on read cluster server (%s)", server.Id)
		}

		cluster.Id = server.ClusterUuid

		// responseBytes, err = cluster.ReadDI()
		// if err != nil {
		// 	return diag.FromErr(err)
		// }
		// err = cluster.Deserialize(responseBytes)
		// if err != nil {
		// 	return diag.FromErr(err)
		// }
		// cluster.WriteTF(res)

		_, err = cluster.StateChange(res).WaitForStateContext(ctx)
		if err != nil {
			log.Printf("timeout on create cluster (%s)", cluster.Id)
		}

		groupName := objNew.GetType()[3:]
		group := Inventory.All.GetGroup(groupName)
		if group == nil {
			group = &inventory_yaml.Group{Name: groupName}
		}

		subgroupName := utils.Reformat(cluster.ServiceName)
		subgroup := group.GetGroup(subgroupName)
		if subgroup == nil {
			subgroup = &inventory_yaml.Group{Name: subgroupName}
			group.AddGroup(subgroup)
		}

		if subgroup.Vars == nil {
			subgroup.Vars = make(map[string]interface{})
		}
		subgroup.Vars["service_name_en"] = subgroupName
		subgroup.Vars["service_name_ru"] = cluster.ServiceName
		for k, v := range cluster.Object.GroupVars(&cluster) {
			subgroup.Vars[k] = v
		}

		for _, server := range cluster.Servers {
			host := &inventory_yaml.Host{
				Name: server.DNSName,
				Vars: objNew.HostVars(server),
			}
			subgroup.AddHost(host)
		}

		Inventory.All.AddGroup(group)
		err = Inventory.Save()
		if err != nil {
			return diag.FromErr(err)
		}
		err = Inventory.ToBIN()
		if err != nil {
			return diag.FromErr(err)
		}

		// attach tags
		if res.HasChange("tag_ids") {
			_, tagIds := res.GetChange("tag_ids")
			tagSet := tagIds.(*schema.Set)
			for _, v := range tagSet.List() {
				for _, clusterServer := range cluster.Servers {
					_, err := clusterServer.TagAttachDI(v.(string))
					if err != nil {
						return diag.FromErr(err)
					}
				}
			}
		}

		return diags
	}
	return f
}

func ReadClusterResource(obj models.DIClusterResource) schema.ReadContextFunc {
	f := func(ctx context.Context, res *schema.ResourceData, m interface{}) diag.Diagnostics {
		var diags diag.Diagnostics

		objNew := obj.NewObj()
		cluster := models.Cluster{Object: objNew}
		cluster.ReadTF(res)

		responseBytes, err := cluster.ReadDI()
		if err != nil {
			return diag.FromErr(err)
		}
		err = cluster.Deserialize(responseBytes)
		if err != nil {
			return diag.FromErr(err)
		}
		cluster.WriteTF(res)
		// if obj.State == "damaged" {
		// 	err = obj.Delete()
		// 	if err != nil {
		// 		return diag.FromErr(err)
		// 	}
		// 	res.SetId("")
		// }

		groupName := objNew.GetType()[3:]
		group := Inventory.All.GetGroup(groupName)
		if group == nil {
			group = &inventory_yaml.Group{Name: groupName}
		}

		subgroupName := utils.Reformat(cluster.ServiceName)
		subgroup := group.GetGroup(subgroupName)
		if subgroup == nil {
			subgroup = &inventory_yaml.Group{Name: subgroupName}
			group.AddGroup(subgroup)
		}

		if subgroup.Vars == nil {
			subgroup.Vars = make(map[string]interface{})
		}
		subgroup.Vars["service_name_en"] = subgroupName
		subgroup.Vars["service_name_ru"] = cluster.ServiceName
		for k, v := range cluster.Object.GroupVars(&cluster) {
			subgroup.Vars[k] = v
		}

		for _, server := range cluster.Servers {
			host := &inventory_yaml.Host{
				Name: server.DNSName,
				Vars: objNew.HostVars(server),
			}
			subgroup.AddHost(host)
		}

		Inventory.All.AddGroup(group)
		err = Inventory.Save()
		if err != nil {
			return diag.FromErr(err)
		}
		err = Inventory.ToBIN()
		if err != nil {
			return diag.FromErr(err)
		}

		return diags
	}
	return f
}

func UpdateClusterResource(obj models.DIClusterResource) schema.UpdateContextFunc {
	f := func(ctx context.Context, res *schema.ResourceData, m interface{}) diag.Diagnostics {
		var diags diag.Diagnostics

		objNew := obj.NewObj()
		cluster := models.Cluster{Object: objNew}
		cluster.ReadTF(res)

		if res.HasChange("disk") {
			return diag.FromErr(errors.New("main disk change not implemented in api"))
		}

		if res.Get("state") == "creating" || res.Get("state_resize") == "resizing" {
			return diag.FromErr(errors.New("can't update 'creating' or 'resizing' instance"))
		}

		if res.HasChange("app_params") {
			v1, v2 := res.GetChange("app_params")
			var oldBSC, newBSC int
			v1BoxServerCount := v1.(map[string]interface{})["box_server_count"]
			switch v1BoxServerCount.(type) {
			case string:
				v1BoxServerCountInt, err := strconv.Atoi(v1BoxServerCount.(string))
				if err != nil {
					return diag.FromErr(err)
				}
				oldBSC = v1BoxServerCountInt
			case int:
				oldBSC = v1BoxServerCount.(int)
			}
			v2BoxServerCount := v2.(map[string]interface{})["box_server_count"]
			switch v2BoxServerCount.(type) {
			case string:
				v2BoxServerCountInt, err := strconv.Atoi(v2BoxServerCount.(string))
				if err != nil {
					return diag.FromErr(err)
				}
				newBSC = v2BoxServerCountInt
			case int:
				newBSC = v2BoxServerCount.(int)
			}

			responseBytes, err := cluster.ReadDI()
			if err != nil {
				return diag.FromErr(err)
			}
			err = cluster.Deserialize(responseBytes)
			if err != nil {
				return diag.FromErr(err)
			}

			if oldBSC < newBSC {
				diff := newBSC - oldBSC
				resizeMap := map[string]interface{}{
					"number_of_nodes": diff,
					"cluster_uuid":    cluster.Id.String(),
				}
				objBytes, err := json.Marshal(resizeMap)
				if err != nil {
					return diag.FromErr(err)
				}
				_, err = cluster.UpScaleDI(objBytes)
				if err != nil {
					return diag.FromErr(err)
				}
				cluster.StateResize = "resizing"
				cluster.WriteTF(res)

				_, err = cluster.StateResizeChange(res).WaitForStateContext(ctx)
				if err != nil {
					log.Printf("[INFO] timeout on resize for instance (%s), save current state: %s", cluster.Id.String(), cluster.StateResize)
				}

			}
			if oldBSC > newBSC {
				diff := oldBSC - newBSC
				for i := 1; i <= diff; i++ {
					vm := cluster.Servers[len(cluster.Servers)-i]
					log.Printf("delete %s %s", vm.Id, vm.Name)
					err = vm.DeleteDI()
					if err != nil {
						return diag.FromErr(err)
					}
				}
				cluster.Servers = cluster.Servers[:len(cluster.Servers)-diff]
				// res.Set("app_params", v2)
				responseBytes, err = cluster.ReadDI()
				if err != nil {
					return diag.FromErr(err)
				}
				err = cluster.Deserialize(responseBytes)
				if err != nil {
					return diag.FromErr(err)
				}
				cluster.WriteTF(res)
			}
			// return diag.FromErr(errors.New("test cluster update"))
		}

		if res.HasChange("service_name") {
			_, service_name := res.GetChange("service_name")
			changes := map[string]map[string]string{
				"cluster": {
					"service_name": service_name.(string),
				},
			}
			objBytes, err := json.Marshal(changes)
			if err != nil {
				return diag.FromErr(err)
			}
			responseBytes, err := cluster.UpdateDI(objBytes)
			if err != nil {
				return diag.FromErr(err)
			}
			err = cluster.Deserialize(responseBytes)
			if err != nil {
				return diag.FromErr(err)
			}
			cluster.ServiceName = service_name.(string)
			cluster.WriteTF(res)
		}

		if res.HasChange("project_id") {
			_, projectId := res.GetChange("project_id")
			changes := map[string]interface{}{
				"uuids":        []string{cluster.Id.String()},
				"project_uuid": projectId.(string),
			}
			objBytes, err := json.Marshal(changes)
			if err != nil {
				return diag.FromErr(err)
			}
			_, err = cluster.MoveDI(objBytes)
			if err != nil {
				return diag.FromErr(err)
			}
			cluster.ProjectId = uuid.MustParse(projectId.(string))
			responseBytes, err := cluster.ReadDI()
			if err != nil {
				return diag.FromErr(err)
			}
			err = cluster.Deserialize(responseBytes)
			if err != nil {
				return diag.FromErr(err)
			}
			cluster.WriteTF(res)
		}

		// if res.HasChange("ram") || res.HasChange("cpu") {
		// 	_, ram := res.GetChange("ram")
		// 	_, cpu := res.GetChange("cpu")
		// 	changes := map[string]map[string]int{
		// 		"resize": {
		// 			"ram": ram.(int),
		// 			"cpu": cpu.(int),
		// 		},
		// 	}
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

			_, err = cluster.ResizeDI(objBytes)
			if err != nil {
				return diag.FromErr(err)
			}
			cluster.StateResize = "resizing"
			cluster.WriteTF(res)

			_, err = cluster.StateResizeChange(res).WaitForStateContext(ctx)
			if err != nil {
				log.Printf("[INFO] timeout on resize for instance (%s), save current state: %s", cluster.Id.String(), cluster.StateResize)
			}
		}

		if res.HasChange("tag_ids") {
			v1, v2 := res.GetChange("tag_ids")
			tagSet1 := v1.(*schema.Set)
			tagSet2 := v2.(*schema.Set)
			l1 := utils.ArrInterfaceToArrStr(tagSet1.List())
			l2 := utils.ArrInterfaceToArrStr(tagSet2.List())

			for _, server := range cluster.Servers {
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
			}
			// cluster.WriteTF(res)
		}

		cluster.WriteTF(res)
		return diags
	}
	return f
}

func DeleteClusterResource(obj models.DIClusterResource) schema.DeleteContextFunc {
	f := func(ctx context.Context, res *schema.ResourceData, m interface{}) diag.Diagnostics {
		var diags diag.Diagnostics

		objNew := obj.NewObj()
		cluster := models.Cluster{Object: objNew}
		cluster.ReadTF(res)

		if res.Get("state") == "creating" || res.Get("state_resize") == "resizing" {
			return diag.FromErr(errors.New("can't delete 'creating' or 'resizing' instance"))
		}

		err := cluster.DeleteDI()
		if err != nil {
			return diag.FromErr(err)
		}
		res.SetId("")

		groupName := obj.GetType()[3:]
		group := Inventory.All.GetGroup(groupName)
		if group == nil {
			group = &inventory_yaml.Group{Name: groupName}
		}
		subgroupName := utils.Reformat(cluster.ServiceName)
		subgroup := group.GetGroup(subgroupName)
		if subgroup != nil {
			group.RmGroup(subgroup.Name)
		}
		if len(group.Hosts) == 0 && len(group.Children) == 0 {
			Inventory.All.RmGroup(group.Name)
		}
		err = Inventory.Save()
		if err != nil {
			return diag.FromErr(err)
		}
		err = Inventory.ToBIN()
		if err != nil {
			return diag.FromErr(err)
		}

		return diags
	}
	return f
}
