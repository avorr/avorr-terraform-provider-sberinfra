package models

type Network struct {
	Cidr           string   `json:"cidr"`
	DnsNameservers []string `json:"dns_nameservers"`
	EnableDhcp     bool     `json:"enable_dhcp"`
	IsDefault      bool     `json:"is_default"`
	NetworkName    string   `json:"network_name"`
	NetworkUuid    string   `json:"network_uuid" hcl:"network_uuid"`
}

type ByNetworkPath []*Network

//func (o ByNetworkPath) Len() int           { return len(o) }
//func (o ByNetworkPath) Less(i, j int) bool { return o[i].Path < o[j].Path }
//func (o ByNetworkPath) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
//
