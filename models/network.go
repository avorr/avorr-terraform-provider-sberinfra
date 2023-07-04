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
