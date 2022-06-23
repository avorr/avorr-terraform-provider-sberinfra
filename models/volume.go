package models

type Volume struct {
	// Id             string `json:"volume_id"`
	// ToscaRequestId string `json:"tosca_request_id"`
	VolumeId string `json:"volume_id,omitempty" hcl:"-"`
	Size     int    `json:"size" hcl:"size"`
	Path     string `json:"path,omitempty" hcl:"path,optional"`
	Status   string `json:"status,omitempty"`
	// Name           string `json:"name,omitempty" hcl:"name,optional"`
	// FSType         string `json:"fs_type,omitempty" hcl:"fs_type,optional"`
	StorageType string `json:"storage_type,omitempty" hcl:"storage_type,optional"`
}

type HCLVolume struct {
	Size        int     `json:"size" hcl:"size"`
	Path        *string `json:"path,omitempty" hcl:"path,optional"`
	StorageType string  `json:"storage_type,omitempty" hcl:"storage_type,optional"`
}

type ByPath []*Volume

func (o ByPath) Len() int           { return len(o) }
func (o ByPath) Less(i, j int) bool { return o[i].Path < o[j].Path }
func (o ByPath) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
