package test

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	backendConfig "github.com/dwarvesf/smithy/backend/config"
	"github.com/dwarvesf/smithy/common/database"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq" // postgresql driver
)

const (
	dbHost     = "localhost"
	dbPort     = "5439"
	dbUser     = "postgres"
	dbPassword = "example"
	dbName     = "test"
)

// MigrateTables migrate db with tables base by domain model
func MigrateTables(db *gorm.DB) error {
	type User struct {
		Id   int `sql:"primary_key"`
		Name string
	}
	return db.AutoMigrate(
		User{},
	).Error
}

func CreateModelList() []database.Model {
	dm := []database.Model{
		{
			TableName: "users",
			Columns: []database.Column{
				{
					Name:      "id",
					Type:      "int",
					IsPrimary: true,
				},
				{
					Name:       "name",
					Type:       "string",
					IsNullable: true,
				},
				{
					Name:       "title",
					Type:       "string",
					IsNullable: true,
				},
				{
					Name:       "description",
					Type:       "string",
					IsNullable: true,
				},
				{
					Name:       "age",
					Type:       "int",
					IsNullable: true,
				},
			},
		},
	}

	return dm
}

func CreateConfig(t *testing.T) (*backendConfig.Config, func()) {
	cfg := &backendConfig.Config{
		SerectKey: "fb0bc76a-dbb1-4944-bcf7-aaef0d9d6e95",
		AgentURL:  "http://localhost:3000/agent",
		ModelList: CreateModelList(),
		ConnectionInfo: database.ConnectionInfo{
			DBType:          "postgres",
			DBUsername:      dbUser,
			DBPassword:      dbPassword,
			DBName:          dbName,
			DBPort:          dbPort,
			DBHostname:      dbHost,
			DBSSLModeOption: "false",
		},
	}
	cfg.UpdateDB()

	rand.Seed(time.Now().UnixNano())
	schemaName := "test" + strconv.FormatInt(rand.Int63(), 10)

	err := cfg.DB().Exec("CREATE SCHEMA " + schemaName).Error
	if err != nil {
		t.Fatalf("Fail to create schema. %s", err.Error())
	}

	// set schema for current db connection
	err = cfg.DB().Exec("SET search_path TO " + schemaName).Error
	if err != nil {
		t.Fatalf("Fail to set search_path to created schema. %s", err.Error())
	}

	// set schema name to config
	cfg.DBSchemaName = schemaName

	return cfg, func() {
		err := cfg.DB().Exec("DROP SCHEMA " + schemaName + " CASCADE").Error
		if err != nil {
			t.Fatalf("Fail to drop database. %s", err.Error())
		}
	}
}
