package domain

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// Model .
type Model struct {
	ID        UUID       `sql:",type:uuid" json:"id"`
	CreatedAt time.Time  `sql:"default:now()" json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// BeforeCreate prepare data before create data
func (m *Model) BeforeCreate(scope *gorm.Scope) error {
	if m.ID.IsZero() {
		if err := scope.SetColumn("ID", uuid.NewV4()); err != nil {
			return err
		}
	}

	if err := scope.SetColumn("CreatedAt", time.Now()); err != nil {
		return err
	}

	return nil
}
