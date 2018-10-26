package user

import (
	"errors"

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

// Create implement Create for User service
func (s *pgService) Create(p *domain.User) error {
	return s.db.Create(p).Error
}

// Update implement Update for User service
func (s *pgService) Update(p *domain.User) (*domain.User, error) {
	old := &domain.User{
		Model:          domain.Model{ID: p.ID},
		Username:       p.Username,
		Email:          p.Email,
		IsEmailAccount: p.IsEmailAccount,
	}

	//if don't have id, find by username
	if old.ID.IsZero() {
		if !old.IsEmailAccount && old.Username != "" {
			if err := s.db.Where("username = ?", old.Username).First(old).Error; err != nil {
				return nil, err
			}
		} else {
			if err := s.db.Where("email = ?", old.Email).First(old).Error; err != nil {
				return nil, err
			}
		}
	} else {
		if err := s.db.Find(old).Error; err != nil {
			return nil, err
		}
	}

	if p.Username != "" {
		old.Username = p.Username
	}

	if p.Password != "" {
		old.Password = p.Password
	}

	if p.Role != "" {
		old.Role = p.Role
	}

	if p.Email != "" {
		old.Email = p.Email
	}

	if p.ConfirmCode != "" {
		old.ConfirmCode = p.ConfirmCode
	}

	return old, s.db.Save(old).Error
}

// Find implement Find for User service
func (s *pgService) Find(p *domain.User) (*domain.User, error) {
	res := p
	if p.IsEmailAccount && p.Username != "" {
		return nil, errors.New("username is invalid")
	}

	if res.ID.IsZero() {
		if !p.IsEmailAccount || p.Username != "" {
			if err := s.db.Where("username = ?", p.Username).First(p).Error; err != nil {
				return nil, err
			}
		} else {
			if err := s.db.Where("email = ?", p.Email).First(p).Error; err != nil {
				return nil, err
			}

		}
	} else {
		if err := s.db.Find(&res).Error; err != nil {
			return nil, err
		}
	}

	return res, nil
}

// FindAll implement FindAll for User service
func (s *pgService) FindAll() ([]domain.User, error) {
	res := []domain.User{}
	return res, s.db.Find(&res).Error
}

// Delete implement Delete for User service
func (s *pgService) Delete(p *domain.User) error {
	old := domain.User{Model: domain.Model{ID: p.ID}}
	if err := s.db.Find(&old).Error; err != nil {
		return err
	}
	return s.db.Delete(old).Error
}

// GetPermission implement get permission for User service
func (s *pgService) GetPermission(p *domain.User, dbName string, tableName string) ([]domain.Permission, error) {
	pers := []domain.Permission{}
	if err := s.db.Where("user_id = ? AND database_name = ? AND table_name = ?", p.ID, dbName, tableName).Find(&pers).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}

	return pers, nil
}

// GetPermissionUserAndGroup implement get permission (include group permission) for User service
func (s *pgService) GetPermissionUserAndGroup(p *domain.User, dbName string, tableName string) (*domain.Permission, error) {
	if p.ID.IsZero() {
		if !p.IsEmailAccount {
			if err := s.db.Where("username = ?", p.Username).First(p).Error; err != nil {
				return nil, err
			}
		} else {
			if err := s.db.Where("email = ?", p.Email).First(p).Error; err != nil {
				return nil, err
			}
		}
	}

	// get user permission
	userPermission := &domain.Permission{}
	if err := s.db.Where("user_id = ? AND database_name = ? AND table_name = ?", p.ID, dbName, tableName).First(userPermission).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}

	//get group
	groups := []domain.Group{}
	s.db.Model(p).Association("Groups").Find(&groups)

	// get group permission
	groupPermissions := []domain.Permission{}
	for _, group := range groups {
		g := domain.Permission{}
		if err := s.db.Where("group_id = ? AND database_name = ? AND table_name = ?", group.ID, dbName, tableName).First(&g).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				continue
			} else {
				return nil, err
			}
		}
		groupPermissions = append(groupPermissions, g)
	}

	// AND permissions
	for _, groupPermission := range groupPermissions {
		u := userPermission.AND(groupPermission)
		userPermission = &u
	}

	return userPermission, nil
}
