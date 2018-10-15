package service

import (
	"os"

	"github.com/dwarvesf/smithy/backend"
	"github.com/dwarvesf/smithy/backend/auth/gplus"
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
	*gplus.Provider
}

// NewService new dashboard handler
func NewService(cfg *backendConfig.Config, db *gorm.DB) (Service, error) {
	mapper, err := backend.NewSQLMapper(cfg)
	if err != nil {
		return Service{}, err
	}

	sqlWriteReadDeleter := view.NewBoltWriteReadDeleter(cfg.PersistenceFileName)
	provider := gplus.NewProvider(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"))

	return Service{
		Wrapper:           backendConfig.NewWrapper(cfg),
		Mapper:            mapper,
		WriteReadDeleter:  sqlWriteReadDeleter,
		UserService:       userSrv.NewPGService(db),
		GroupService:      groupSrv.NewPGService(db),
		PermissionService: permissionSrv.NewPGService(db),
		Provider:          provider,
	}, nil
}
