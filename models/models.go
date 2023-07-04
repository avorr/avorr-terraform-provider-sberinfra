package models

import (
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gitlab.gos-tech.xyz/pid/iac/terraform-provider-sberinfra/client"
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
}
