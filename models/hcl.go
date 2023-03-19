package models

import (
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"gitlab.gos-tech.xyz/pid/iac/terraform-provider-sberinfra/utils"
)

type HCLDataResource struct {
	ResType  string `hcl:"type,label"`
	ResName  string `hcl:"name,label"`
	Resource DIDataResource
}

type HCLDataRoot struct {
	Resources DIDataResource `hcl:"data,block"`
	// Variables []Variable `hcl:"var,block"`
}

type HCLResourceRoot struct {
	Resources DIResource `hcl:"resource,block"`
}

type HCLOutputRoot struct {
	Resources *HCLOutput `hcl:"output,block"`
}

type HCLOutput struct {
	ResName string `hcl:"name,label"`
	Value   string `hcl:"value"`
}

func ToHCLData(res DIDataResource) []byte {
	res.SetResFields()
	dataRoot := &HCLDataRoot{Resources: res}
	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(dataRoot, f.Body())
	return utils.Regexp(f.Bytes())
}

func ToHCLResource(res DIResource) []byte {
	// res.SetResFields()
	dataRoot := &HCLResourceRoot{Resources: res}
	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(dataRoot, f.Body())
	return utils.Regexp(f.Bytes())
}

func ToHCLOutput(res DIDataResource) []byte {
	res.SetResFields()
	name, value := res.GetOutput()
	dataRoot := &HCLOutputRoot{
		Resources: &HCLOutput{
			ResName: name,
			Value:   value,
		},
	}
	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(dataRoot, f.Body())
	return utils.Regexp(f.Bytes())
}

type HCLTags []string

type HCL struct {
	Type           string `hcl:"type,label"`
	Name           string `hcl:"name,label"`
	GroupId        string `hcl:"group_id"`
	ProjectId      string `hcl:"project_id"`
	ServiceName    string `hcl:"service_name"`
	IrGroup        string `hcl:"ir_group"`
	OsName         string `hcl:"os_name"`
	OsVersion      string `hcl:"os_version"`
	Virtualization string `hcl:"virtualization"`
	FaultTolerance string `hcl:"fault_tolerance"`
	//Region         string        `hcl:"region"`
	Disk          int           `hcl:"disk"`
	Flavor        string        `hcl:"flavor"`
	Zone          string        `hcl:"zone"`
	PublicSshName string        `hcl:"public_ssh_name,optional"`
	AppParams     *HCLAppParams `hcl:"app_params,block"`
	Volumes       []*HCLVolume  `hcl:"volume,block"`
	TagIds        *HCLTags      `hcl:"tag_ids,optional"`
}

type HCLRoot struct {
	Resources *HCL `hcl:"resource,block"`
}

type HCLAppParams struct {
	JoinDomain string `hcl:"joindomain"`
}
