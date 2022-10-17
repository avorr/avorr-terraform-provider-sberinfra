package models

import "github.com/google/uuid"

type SSHKey struct {
	Id            uuid.UUID `json:"uuid"`
	PublicSshName string    `json:"public_ssh_name"`
	PublicSsh     string    `json:"public_ssh"`
	Default       bool      `json:"default"`
}
