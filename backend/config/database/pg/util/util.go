package util

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/dwarvesf/smithy/backend/domain"
	"github.com/jinzhu/gorm"
)

const (
	dbHost     = "localhost"
	dbPort     = 5434
	dbUser     = "postgres"
	dbPassword = "example"
	dbName     = "test"
)

// CreateTestDatabase will create a test-database and test-schema
func CreateTestDatabase(t *testing.T) (*gorm.DB, string, func()) {
	testingHost := dbHost
	testingPort := fmt.Sprintf("%d", dbPort)
	if os.Getenv("POSTGRES_TESTING_HOST") != "" {
		testingHost = os.Getenv("POSTGRES_TESTING_HOST")
	}
	if os.Getenv("POSTGRES_TESTING_PORT") != "" {
		testingPort = os.Getenv("POSTGRES_TESTING_PORT")
	}
	connectionString := fmt.Sprintf("host = %s port=%s user=%s password=%s dbname=%s sslmode=disable", testingHost, testingPort, dbUser, dbPassword, dbName)
	db, dbErr := gorm.Open("postgres", connectionString)
	if dbErr != nil {
		t.Fatalf("Fail to create database. %s", dbErr.Error())
	}

	rand.Seed(time.Now().UnixNano())
	schemaName := "test" + strconv.FormatInt(rand.Int63(), 10)

	err := db.Exec("CREATE SCHEMA " + schemaName).Error
	if err != nil {
		t.Fatalf("Fail to create schema. %s", err.Error())
	}

	// set schema for current db connection
	err = db.Exec("SET search_path TO " + schemaName).Error
	if err != nil {
		t.Fatalf("Fail to set search_path to created schema. %s", err.Error())
	}

	return db, schemaName, func() {
		err := db.Exec("DROP SCHEMA " + schemaName + " CASCADE").Error
		if err != nil {
			t.Fatalf("Fail to drop database. %s", err.Error())
		}
	}
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
				Model:    domain.Model{ID: seedUserID},
				Username: "aaa",
				Role:     "admin",
				Password: "cccccc1@!a",
				Email:    "aaa@gmai.com",
			},
		},
	}

	group2 := domain.Group{
		Name:        "user",
		Description: "normal users",
		Role:        "user",
		Users: []domain.User{
			{
				Username: "bbb",
				Role:     "user",
				Password: "cccccc1@!a",
				Email:    "aaa@gmai.com",
			},
		},
	}

	p1 := domain.Permission{
		Model:        domain.Model{ID: seedPermission1ID},
		DatabaseName: "test1",
		TableName:    "users",
		Insert:       true,
		Select:       true,
		Update:       true,
		Delete:       false,
		GroupID:      seedGroup1ID,
	}

	p2 := domain.Permission{
		Model:        domain.Model{ID: seedPermission2ID},
		DatabaseName: "test1",
		TableName:    "users",
		Insert:       true,
		Select:       false,
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
