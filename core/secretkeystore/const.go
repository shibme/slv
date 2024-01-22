package secretkeystore

import "errors"

const (
	slvSecreKeyEnvarName      = "SLV_SECRET_KEY"
	slvAccessBindingEnvarName = "SLV_ACCESS_BINDING"
)

var (
	errEnvironmentAccessNotFound = errors.New("environment doesn't have access. please set one of the environment variables: " + slvSecreKeyEnvarName + " or " + slvAccessBindingEnvarName)
)