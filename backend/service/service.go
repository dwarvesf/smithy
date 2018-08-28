package service

import (
	backendConfig "github.com/dwarvesf/smithy/backend/config"
	"github.com/dwarvesf/smithy/backend/sqlmapper"
)

// Service ...
type Service struct {
	*backendConfig.Wrapper
	sqlmapper.Mapper
}

// NewService new dashboard handler
func NewService(cfg *backendConfig.Config) Service {
	return Service{
		Wrapper: backendConfig.NewWrapper(cfg),
	}
}
