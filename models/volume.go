package models

type Volume struct {
	VolumeId string `json:"volume_id"`
	Size     int    `json:"size" hcl:"size"`
	Name     string `json:"name,omitempty"`
	//Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
	StorageType string `json:"storage_type,omitempty" hcl:"storage_type,optional"`
}

type ByPath []*Volume

func (o ByPath) Len() int { return len(o) }

func (o ByPath) Swap(i, j int) { o[i], o[j] = o[j], o[i] }
