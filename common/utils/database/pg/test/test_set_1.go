package test

import (
	"testing"

	backendConfig "github.com/dwarvesf/smithy/backend/config"
	"github.com/dwarvesf/smithy/common/database"
	utilDB "github.com/dwarvesf/smithy/common/utils/database/pg"
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

	clearDB := utilDB.CreateDatabase(t, cfg)

	return cfg, clearDB
}
