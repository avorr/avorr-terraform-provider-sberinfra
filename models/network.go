package models

type Network1 struct {
	// Id             string `json:"volume_id"`
	// ToscaRequestId string `json:"tosca_request_id"`
	// Name           string `json:"name,omitempty" hcl:"name,optional"`
	// FSType         string `json:"fs_type,omitempty" hcl:"fs_type,optional"`
	NetworkName    string `json:"volume_id,omitempty" hcl:"-"`
	StorageType    string `json:"storage_type,omitempty" hcl:"storage_type,optional"`
	NetworkUuid    int    `json:"size" hcl:"size"`
	Cidr           string `json:"path,omitempty" hcl:"path,optional"`
	DnsNameServers string `json:"status,omitempty"`
}

type Network struct {
	Cidr           string   `json:"cidr"`
	DnsNameservers []string `json:"dns_nameservers"`
	EnableDhcp     bool     `json:"enable_dhcp"`
	IsDefault      bool     `json:"is_default"`
	NetworkName    string   `json:"network_name"`
	NetworkUuid    string   `json:"network_uuid"`
}

//					"network_name": {Type: schema.TypeString, Required: true},
//					"network_uuid": {Type: schema.TypeString, Computed: true},
//					"cidr":         {Type: schema.TypeString, Required: true},
//					"dns_nameservers": {
//						Type:     schema.TypeSet,
//						Required: true,
//						Elem:     &schema.Schema{Type: schema.TypeString},
//					},
//					"enable_dhcp": {Type: schema.TypeBool, Required: true},
//					"is_default":  {Type: schema.TypeBool, Optional: true},

//type HCLVolume struct {
//	Size        int    `json:"size" hcl:"size"`
//	Path        string `json:"path,omitempty" hcl:"path,optional"`
//	StorageType string `json:"storage_type,omitempty" hcl:"storage_type,optional"`
//}

type ByNetworkPath []*Network

//func (o ByNetworkPath) Len() int           { return len(o) }
//func (o ByNetworkPath) Less(i, j int) bool { return o[i].Path < o[j].Path }
//func (o ByNetworkPath) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
//
