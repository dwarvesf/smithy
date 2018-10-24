package permission

import "github.com/dwarvesf/smithy/backend/domain"

// Service interface for project service
type Service interface {
	Update(p *domain.Permission) (*domain.Permission, error)
	FindByUser(p *domain.User) ([]domain.Permission, error)
	FindByGroup(p *domain.Group) ([]domain.Permission, error)
}
