package models

import (
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclwrite"

	"base.sw.sbc.space/pid/terraform-provider-si/utils"
)

type VaultGenericSecretRoot struct {
	Resources *VaultGenericSecret `hcl:"data,block"`
}

type VaultGenericSecret struct {
	ResType string `json:"-" hcl:"type,label"`
	ResName string `json:"-" hcl:"name,label"`
	Path    string `json:"-" hcl:"path"`
	Data    string `json:"-"`
	Field   string `json:"-"`
}

func (o *VaultGenericSecret) ToHCL() []byte {
	dataRoot := &VaultGenericSecretRoot{Resources: o}
	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(dataRoot, f.Body())
	return utils.Regexp(f.Bytes())
}

func (o *VaultGenericSecret) ToBash() []byte {
	data := fmt.Sprintf("echo -n '%s' | base64 -D | vault kv put \\\n %s \\\n %s=-\n\n",
		base64.StdEncoding.EncodeToString([]byte(o.Data)),
		o.Path,
		o.Field,
	)
	return []byte(data)
}
