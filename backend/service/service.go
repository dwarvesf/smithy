package service

import (
	backendConfig "github.com/dwarvesf/smithy/backend/config"
)

// Service ...
type Service struct {
	Config    *backendConfig.Wrapper
	SecretKey string
}

// NewService new dashboard handler
func NewService(cfg *backendConfig.Config) Service {
	// We must load secrect key from persistant
	return Service{backendConfig.NewWrapper(cfg), "K8UeMDPyb9AwFkzS"}
}
