example:
```go
package main

import (
	"log"
	"os"
	"path/filepath"

	inventory_yaml "stash.sigma.sbrf.ru/sddevops/inventory-yaml.git"
)

func main() {

	inventory_yaml.FLAG_RM_EMPTY_DICT_VAL = true
	rootgr := &inventory_yaml.Group{
		Name:  "rootgroup",
		Hosts: nil,
		Vars: map[string]interface{}{
			"group_all_var": "value",
		},
		Children: nil,
	}
	group1 := &inventory_yaml.Group{
		Name:     "testgr1",
		Hosts:    nil,
		Vars:     nil,
		Children: nil,
	}
	group2 := &inventory_yaml.Group{
		Name:  "testgr2",
		Hosts: nil,
		Vars: map[string]interface{}{
			"groupvar1": 1,
			"groupvar2": "two",
		},
		Children: nil,
	}
	group3 := &inventory_yaml.Group{
		Name:  "testgr3",
		Hosts: nil,
		Vars: map[string]interface{}{
			"groupvar31": 31,
			"groupvar32": "trirty two",
		},
		Children: nil,
	}

	group2.AddGroup(group3)

	host1 := &inventory_yaml.Host{
		Name: "test-di",
	}
	host2 := &inventory_yaml.Host{Name: "testhost2", Vars: map[string]interface{}{"host_var2": "value2"}}
	host3 := &inventory_yaml.Host{
		Name: "testhost3",
		Vars: map[string]interface{}{
			"ansible_host": "127.0.0.1",
			"hostvar3":     33,
		},
	}

	rootgr.AddHost(host1)
	group1.AddHost(host2)
	group2.AddHost(host3)
	rootgr.AddGroup(group1)
	rootgr.AddGroup(group2)

	inv := &inventory_yaml.Inventory{All: rootgr}
	err := inv.Save()
	if err != nil {
		log.Println(err)
		return
	}

}
```

```bash
$ cat inventory.yml
---
all:
  hosts:
    test-di:
  vars:
    group_all_var: value
  children:
    testgr1:
      hosts:
        testhost2:
          host_var2: value2
    testgr2:
      hosts:
        testhost3:
          ansible_host: 127.0.0.1
          hostvar3: 33
      vars:
        groupvar1: 1
        groupvar2: two
      children:
        testgr3:
          vars:
            groupvar31: 31
            groupvar32: trirty two
```