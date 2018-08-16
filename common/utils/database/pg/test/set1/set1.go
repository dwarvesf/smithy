package set1

import (
	"strconv"
	"testing"

	"github.com/dwarvesf/smithy/common/database"
	"github.com/jinzhu/gorm"

	backendConfig "github.com/dwarvesf/smithy/backend/config"
	utilDB "github.com/dwarvesf/smithy/common/utils/database/pg"
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
	return db.Exec(`CREATE TABLE "users" (
		"id" int NOT NULL,
		"name" text,
		CONSTRAINT "users_pkey" PRIMARY KEY ("id")
	  ) WITH (oids = false);`).Error
}

func CreateUserSampleData(db *gorm.DB) ([]utilDB.User, error) {
	users := make([]utilDB.User, 0)

	for i := 0; i < 15; i++ {
		user := utilDB.User{
			Id:   i + 1,
			Name: "hieudeptrai" + strconv.Itoa(i),
		}
		err := db.Create(&user).Error

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
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
