package group

import "github.com/dwarvesf/smithy/backend/domain"

// Service interface for project service
type Service interface {
	Create(p *domain.Group) error
	Update(p *domain.Group) (*domain.Group, error)
	Find(p *domain.Group) (*domain.Group, error)
	FindAll() ([]domain.Group, error)
	FindByUser(p *domain.User) ([]domain.Group, error)
	GetPermission(p *domain.Group, dbName string, tableName string) ([]domain.Permission, error)
	Delete(p *domain.Group) error
}
