package group

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

// Create implement Create for Group service
func (s *pgService) Create(p *domain.Group) error {
	return s.db.Create(p).Error
}

// Update implement Update for Group service
func (s *pgService) Update(p *domain.Group) (*domain.Group, error) {
	old := domain.Group{Model: domain.Model{ID: p.ID}}
	if err := s.db.Find(&old).Error; err != nil {
		return nil, err
	}

	if p.Name != "" {
		old.Name = p.Name
	}

	if p.Description != "" {
		old.Description = p.Description
	}

	return &old, s.db.Save(&old).Error
}

// Find implement Find for Group service
func (s *pgService) Find(p *domain.Group) (*domain.Group, error) {
	res := p
	if err := s.db.Find(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func (s *pgService) FindByUser(p *domain.User) ([]domain.Group, error) {
	if p.ID.IsZero() {
		if err := s.db.Where("username = ?", p.Username).First(p).Error; err != nil {
			return nil, err
		}
	}

	groups := []domain.Group{}
	if err := s.db.Model(p).Association("Groups").Find(&groups).Error; err != nil {
		return nil, err
	}

	return groups, nil
}

// FindAll implement FindAll for Group service
func (s *pgService) FindAll() ([]domain.Group, error) {
	res := []domain.Group{}
	return res, s.db.Find(&res).Error
}

// Delete implement Delete for Group service
func (s *pgService) Delete(p *domain.Group) error {
	old := domain.Group{Model: domain.Model{ID: p.ID}}
	if err := s.db.Find(&old).Error; err != nil {
		return err
	}
	return s.db.Delete(old).Error
}

// GetPermission implement get permission for User service
func (s *pgService) GetPermission(p *domain.Group, dbName string, tableName string) ([]domain.Permission, error) {
	pers := []domain.Permission{}
	if err := s.db.Where("group_id = ? AND database_name = ? AND table_name = ?", p.ID, dbName, tableName).Find(&pers).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}

	return pers, nil
}
