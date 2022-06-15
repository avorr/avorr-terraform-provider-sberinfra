package inventory_yaml

import (
	"log"
	"testing"
)

func TestNewInventory(t *testing.T) {
	i := NewInventory()
	if i.Yml != "inventory.yml" {
		t.Errorf("expected = %s, got = %s", "inventory.yml", i.Yml)
	}
}

func TestNewInventory2(t *testing.T) {

	inventory := NewInventory()
	rootgr := &Group{
		Name:  "rootgroup",
		Hosts: nil,
		Vars: map[string]interface{}{
			"group_all_var": "value",
		},
		Children: nil,
	}
	inventory.All = rootgr

	group1 := &Group{
		Name:     "testgr1",
		Hosts:    nil,
		Vars:     nil,
		Children: nil,
	}
	group2 := &Group{
		Name:  "testgr2",
		Hosts: nil,
		Vars: map[string]interface{}{
			"groupvar1": 1,
			"groupvar2": "two",
		},
		Children: nil,
	}
	group3 := &Group{
		Name:  "testgr3",
		Hosts: nil,
		Vars: map[string]interface{}{
			"groupvar31": 31,
			"groupvar32": "trirty two",
			"password2":  "65343365313832313861396565633431316233656636373235393563636233373735313264393239\n3961623634356337393038333065316139643234393063330a633331316562623138393166653937\n35373731666362393762343866356261616637636666633163616661373266396164366130393461\n6136376633306461660a363462656664613535376538653330623431343966616234336432353032\n3633\n",
		},
		Children: nil,
	}

	group2.AddGroup(group3)

	host1 := &Host{
		Name: "test1",
	}
	host2 := &Host{Name: "testhost2", Vars: map[string]interface{}{"host_var2": "value2"}}
	host3 := &Host{
		Name: "testhost3",
		Vars: map[string]interface{}{
			"ansible_host": "127.0.0.1",
			"hostvar3":     33,
			"password":     "$ANSIBLE_VAULT;1.1;AES256\n65343365313832313861396565633431316233656636373235393563636233373735313264393239\n3961623634356337393038333065316139643234393063330a633331316562623138393166653937\n35373731666362393762343866356261616637636666633163616661373266396164366130393461\n6136376633306461660a363462656664613535376538653330623431343966616234336432353032\n3633\n",
		},
	}

	rootgr.AddHost(host1)
	group1.AddHost(host2)
	group2.AddHost(host3)
	rootgr.AddGroup(group1)
	rootgr.AddGroup(group2)

	err := inventory.Save()
	if err != nil {
		log.Println(err)
		return
	}

	// err = inventory.ToBIN()
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// pp.Println(inventory)

}

func TestInventory_FromBIN(t *testing.T) {
	inventory := NewInventory()
	// pp.Println(inventory)
	err := inventory.FromBIN()
	if err != nil {
		log.Println(err)
		return
	}
	// pp.Println(inventory)
}
