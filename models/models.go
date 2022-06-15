package models

import (
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"stash.sigma.sbrf.ru/sddevops/terraform-provider-di/client"
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
}

type DIClusterResource interface {
	NewObj() DIClusterResource
	OnSerialize(map[string]interface{}, *Cluster) map[string]interface{}
	OnDeserialize(map[string]interface{}, *Cluster)
	Urls(string) string
	OnReadTF(*schema.ResourceData, *Cluster)
	OnWriteTF(*schema.ResourceData, *Cluster)
	GetType() string
	HostVars(*Server) map[string]interface{}
	GroupVars(*Cluster) map[string]interface{}
}
