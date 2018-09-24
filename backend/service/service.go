package service

import (
	"github.com/dwarvesf/smithy/backend"
	backendConfig "github.com/dwarvesf/smithy/backend/config"
	"github.com/dwarvesf/smithy/backend/sqlmapper"
	"github.com/dwarvesf/smithy/backend/view"
)

// Service ...
type Service struct {
	*backendConfig.Wrapper
	sqlmapper.Mapper
	WriterReaderDeleter view.WriterReaderDeleter
}

// NewService new dashboard handler
func NewService(cfg *backendConfig.Config) (Service, error) {
	mapper, err := backend.NewSQLMapper(cfg)
	if err != nil {
		return Service{}, err
	}

	sqlWriterReaderDeleter := view.NewBoltWriterReaderDeleter(cfg.PersistenceFileName)

	return Service{
		Wrapper:             backendConfig.NewWrapper(cfg),
		Mapper:              mapper,
		WriterReaderDeleter: sqlWriterReaderDeleter,
	}, nil
}
