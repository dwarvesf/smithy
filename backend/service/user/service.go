package user

import "github.com/dwarvesf/smithy/backend/domain"

// Service interface for project service
type Service interface {
	Create(p *domain.User) error
	Update(p *domain.User) (*domain.User, error)
	Find(p *domain.User) (*domain.User, error)
	FindAll() ([]domain.User, error)
	GetPermission(p *domain.User, dbName string, tableName string) ([]domain.Permission, error)
	GetPermissionUserAndGroup(p *domain.User, dbName string, tableName string) (*domain.Permission, error)
	Delete(p *domain.User) error
}
