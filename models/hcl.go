package models

import (
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"stash.sigma.sbrf.ru/sddevops/terraform-provider-di/utils"
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

type HCLAppParams struct {
	VersionJDK         *string `json:"version_jdk" hcl:"version_jdk,optional"`
	JDKVersion         *string `json:"jdk_version" hcl:"jdk_version,optional"`
	SowaVersion        *string `json:"sowa_version" hcl:"sowa_version,optional"`
	Version            *string `json:"version" hcl:"version,optional"`
	JoinDomain         *string `json:"join_domain" hcl:"joindomain,optional"`
	AdminUser          *string `json:"admin_user" hcl:"admin_user,optional"`
	NameProject        *string `json:"name_project" hcl:"name_project,optional"`
	NginxBrotli        *string `json:"nginx_brotli" hcl:"nginx_brotli,optional"`
	NginxGeoip         *string `json:"nginx_geoip" hcl:"nginx_geoip,optional"`
	PostgresDbName     *string `json:"postgres_db_name" hcl:"postgres_db_name,optional"`
	PostgresDbUser     *string `json:"postgres_db_user" hcl:"postgres_db_user,optional"`
	PostgresDbPassword *string `json:"postgres_db_password" hcl:"postgres_db_password,optional"`
	Security           *string `json:"security" hcl:"security,optional"`
	BoxServerCount     *int    `json:"box_server_count" hcl:"box_server_count,optional"`
	FaultTolerance     *string `json:"fault_tolerance" hcl:"fault_tolerance,optional"`
	IseEmail           *string `json:"ise_email" hcl:"ise_email,optional"`
	GGClientPassword   *string `json:"gg_client_password" hcl:"gg_client_password,optional"`
	IseClientPassword  *string `json:"ise_client_password" hcl:"ise_client_password,optional"`
	// Endpoint           *string `json:"endpoint" hcl:"endpoint,optional"`
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
