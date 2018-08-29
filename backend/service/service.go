package service

import (
	"github.com/dwarvesf/smithy/backend"
	backendConfig "github.com/dwarvesf/smithy/backend/config"
	"github.com/dwarvesf/smithy/backend/sqlmapper"
)

// Service ...
type Service struct {
	*backendConfig.Wrapper
	sqlmapper.Mapper
}

// NewService new dashboard handler
func NewService(cfg *backendConfig.Config) (Service, error) {
	mapper, err := backend.NewSQLMapper(cfg)
	if err != nil {
		return Service{}, err
	}

	return Service{
		Wrapper: backendConfig.NewWrapper(cfg),
		Mapper:  mapper,
	}, nil
}
