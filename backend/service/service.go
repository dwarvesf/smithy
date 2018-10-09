package service

import (
	"github.com/dwarvesf/smithy/backend"
	backendConfig "github.com/dwarvesf/smithy/backend/config"
	groupSrv "github.com/dwarvesf/smithy/backend/service/group"
	permissionSrv "github.com/dwarvesf/smithy/backend/service/permission"
	userSrv "github.com/dwarvesf/smithy/backend/service/user"
	"github.com/dwarvesf/smithy/backend/sqlmapper"
	"github.com/dwarvesf/smithy/backend/view"
	"github.com/jinzhu/gorm"
)

// Service ...
type Service struct {
	*backendConfig.Wrapper
	sqlmapper.Mapper
	WriteReadDeleter  view.WriteReadDeleter
	UserService       userSrv.Service
	GroupService      groupSrv.Service
	PermissionService permissionSrv.Service
}

// NewService new dashboard handler
func NewService(cfg *backendConfig.Config, db *gorm.DB) (Service, error) {
	mapper, err := backend.NewSQLMapper(cfg)
	if err != nil {
		return Service{}, err
	}

	sqlWriteReadDeleter := view.NewBoltWriteReadDeleter(cfg.PersistenceFileName)

	return Service{
		Wrapper:           backendConfig.NewWrapper(cfg),
		Mapper:            mapper,
		WriteReadDeleter:  sqlWriteReadDeleter,
		UserService:       userSrv.NewPGService(db),
		GroupService:      groupSrv.NewPGService(db),
		PermissionService: permissionSrv.NewPGService(db),
	}, nil
}
