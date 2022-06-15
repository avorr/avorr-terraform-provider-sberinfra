package inventory_yaml

import (
	"sync"
)

type Group struct {
	mu       sync.Mutex
	Name     string                 `json:"name" yaml:"-"`
	Hosts    map[string]*Host       `json:"hosts" yaml:"hosts,omitempty"`
	Vars     map[string]interface{} `json:"vars" yaml:"vars,omitempty"`
	Children map[string]*Group      `json:"children" yaml:"children,omitempty"`
}

func (o *Group) AddGroup(group *Group) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.Children == nil {
		o.Children = make(map[string]*Group)
	}
	o.Children[group.Name] = group
}

func (o *Group) GetGroup(name string) *Group {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.Children[name]
}

func (o *Group) RmGroup(name string) {
	o.mu.Lock()
	defer o.mu.Unlock()
	delete(o.Children, name)
}

func (o *Group) AddHost(host *Host) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.Hosts == nil {
		o.Hosts = make(map[string]*Host)
	}
	o.Hosts[host.Name] = host
}

func (o *Group) RmHost(name string) {
	o.mu.Lock()
	defer o.mu.Unlock()
	delete(o.Hosts, name)
}
