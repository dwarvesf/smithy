package permission

import (
	"github.com/jinzhu/gorm"

	"github.com/dwarvesf/smithy/backend/domain"
)

type pgService struct {
	db *gorm.DB
}

// NewPGService .
func NewPGService(db *gorm.DB) Service {
	return &pgService{
		db: db,
	}
}

// Update implement Update for Group service
func (s *pgService) Update(p *domain.Permission) (*domain.Permission, error) {
	old := domain.Permission{Model: domain.Model{ID: p.ID}}
	if err := s.db.Find(&old).Error; err != nil {
		return nil, err
	}

	old.Select = p.Select
	old.Insert = p.Insert
	old.Update = p.Update
	old.Delete = p.Delete

	return &old, s.db.Save(&old).Error
}

// FindByUser implement get permission for a user
func (s *pgService) FindByUser(p *domain.User) ([]domain.Permission, error) {
	pers := []domain.Permission{}
	if err := s.db.Where("user_id = ?", p.ID).Find(&pers).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}

	return pers, nil
}

// FindByGroup implement get permission for a user
func (s *pgService) FindByGroup(p *domain.Group) ([]domain.Permission, error) {
	pers := []domain.Permission{}
	if err := s.db.Where("group_id = ?", p.ID).Find(&pers).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}

	return pers, nil
}
