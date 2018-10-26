package pg

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"github.com/dwarvesf/smithy/backend/domain"
)

const (
	dbUsername      = "postgres"
	dbName          = "smithy"
	dbSSLModeOption = "disable"
	dbPassword      = "example"
	dbHostname      = "localhost"
	dbPort          = "5433"
)

// NewPG new a pg connection
func NewPG() (*gorm.DB, func(), error) {
	connectionString := fmt.Sprintf("user=%s dbname=%s sslmode=%s password=%s host=%s port=%s",
		dbUsername,
		dbName,
		dbSSLModeOption,
		dbPassword,
		dbHostname,
		dbPort)
	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		return nil, nil, err
	}

	return db, func() {
		db.Close()
	}, nil
}

// SeedCreateTable .
func SeedCreateTable(db *gorm.DB) error {
	err := db.AutoMigrate(
		domain.User{},
		domain.Group{},
		domain.Permission{},
	).Error

	if err != nil {
		return err
	}

	if err := db.Model(&domain.Permission{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		return err
	}
	if err := db.Model(&domain.Permission{}).AddForeignKey("group_id", "groups(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		return err
	}

	if err := db.Model(&domain.UserGroup{}).AddForeignKey("group_id", "groups(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		return err
	}
	if err := db.Model(&domain.UserGroup{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		return err
	}

	seedGroup1ID := domain.MustGetUUIDFromString("18a5cac3-1024-43c8-acbe-22be76006040")
	seedUserID := domain.MustGetUUIDFromString("18a5cac3-1024-43c8-acbe-22be76006041")
	seedPermission1ID := domain.MustGetUUIDFromString("18a5cac3-1024-43c8-acbe-22be76006042")
	seedPermission2ID := domain.MustGetUUIDFromString("18a5cac3-1024-43c8-acbe-22be76006043")

	group1 := domain.Group{
		Model:       domain.Model{ID: seedGroup1ID},
		Name:        "admin",
		Description: "general admins",
		Role:        "admin",
		Users: []domain.User{
			{
				Model:          domain.Model{ID: seedUserID},
				Username:       "admin",
				Role:           "admin",
				Password:       "admin",
				Email:          "aaa@gmai.com",
				IsEmailAccount: false,
			},
		},
	}

	group2 := domain.Group{
		Name:        "user",
		Description: "normal users",
	}

	p1 := domain.Permission{
		Model:        domain.Model{ID: seedPermission1ID},
		DatabaseName: "fortress",
		TableName:    "users",
		Insert:       true,
		Select:       true,
		Update:       true,
		Delete:       true,
		GroupID:      seedGroup1ID,
	}

	p2 := domain.Permission{
		Model:        domain.Model{ID: seedPermission2ID},
		DatabaseName: "fortress",
		TableName:    "users",
		Insert:       true,
		Select:       true,
		Update:       true,
		Delete:       true,
		UserID:       seedUserID,
	}

	if err := db.Where(domain.Group{Name: "admin"}).Assign(group1).FirstOrCreate(&group1).Error; err != nil {
		return err
	}

	if err := db.Where(domain.Group{Name: "user"}).Assign(group2).FirstOrCreate(&group2).Error; err != nil {
		return err
	}

	if err := db.Where(p1).Assign(p1).FirstOrCreate(&p1).Error; err != nil {
		return err
	}

	if err := db.Where(p2).Assign(p2).FirstOrCreate(&p2).Error; err != nil {
		return err
	}

	return nil
}
