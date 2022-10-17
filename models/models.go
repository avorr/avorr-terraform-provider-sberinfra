package models

import (
	"base.sw.sbc.space/pid/terraform-provider-si/client"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	Api *client.Api
)

type DIDataResource interface {
	NewObj() DIDataResource
	ReadTF(*schema.ResourceData)
	WriteTF(*schema.ResourceData)
	ReadDI() ([]byte, error)
	Deserialize([]byte) error
	ReadAll() ([]byte, error)
	DeserializeAll([]byte) ([]DIDataResource, error)
	// ToHCL() *HCLResource
	// ToHCL2() *HCLDataResource
	GetId() string
	GetResType() string
	GetResName() string
	GetDomainId() uuid.UUID
	GetOutput() (string, string)
	SetResFields()
}

type DIResource interface {
	NewObj() DIResource
	OnSerialize(map[string]interface{}, *Server) map[string]interface{}
	OnDeserialize(map[string]interface{}, *Server)
	Urls(string) string
	OnReadTF(*schema.ResourceData, *Server)
	OnWriteTF(*schema.ResourceData, *Server)
	GetType() string
	GetGroup() string
	HostVars(*Server) map[string]interface{}
	// ToHCL(*Server) ([]byte, error)
	HCLAppParams() *HCLAppParams
	HCLVolumes() []*HCLVolume
}
