package service

import (
	backendConfig "github.com/dwarvesf/smithy/backend/config"
)

// Service ...
type Service struct {
	Config *backendConfig.Wrapper
}

// NewService new dashboard handler
func NewService(cfg *backendConfig.Config) Service {
	return Service{backendConfig.NewWrapper(cfg)}
}
