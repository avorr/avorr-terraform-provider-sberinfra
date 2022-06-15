package inventory_yaml

type Host struct {
	Name string                 `json:"name" yaml:"-"`
	Vars map[string]interface{} `json:"vars" yaml:"vars,omitempty,inline"`
}

// func (o *Host) MarshalYAML() ([]byte, error) {
// }
