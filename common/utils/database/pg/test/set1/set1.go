package set1

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

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
	return db.Exec(`CREATE SEQUENCE user_id_seq;
		CREATE TABLE "users" (
		"id" int NOT NULL DEFAULT nextval('user_id_seq'),
		"name" text,
		CONSTRAINT "users_pkey" PRIMARY KEY ("id")
	  ) WITH (oids = false);`).Error
}

// CreateUserSampleData sample data for test
func CreateUserSampleData(db *gorm.DB) ([]utilDB.User, error) {
	users := make([]utilDB.User, 0)

	for i := 0; i < 15; i++ {
		user := utilDB.User{
			ID:   i + 1,
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

// CreateModelList fake model list from agent
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
			},
		},
	}

	return dm
}

// CreateDatabaseList .
func CreateDatabaseList() []database.Database {
	rand.Seed(time.Now().UnixNano())
	randDBName := strconv.FormatInt(rand.Int63(), 10)

	dm := []database.Database{
		{
			DBName:    randDBName,
			ModelList: CreateModelList(),
		},
	}

	return dm
}

// CreateConfig fake config for test
func CreateConfig(t *testing.T) (*backendConfig.Config, func()) {
	cfg := &backendConfig.Config{
		Databases: CreateDatabaseList(),
		ModelMap:  make(map[string]database.Model),
		ConnectionInfo: database.ConnectionInfo{
			DBType:          "postgres",
			DBUsername:      dbUser,
			DBPassword:      dbPassword,
			DBName:          dbName,
			DBPort:          dbPort,
			DBHostname:      dbHost,
			DBSSLModeOption: "false",
		},
		Authentication: backendConfig.Authentication{
			SerectKey: "lalala",
			UserList: []backendConfig.User{
				{
					Username: "aaa",
					Password: "abc",
					Role:     "client",
				},
				{
					Username: "bbb",
					Password: "abc",
					Role:     "user",
				},
			},
		},
	}

	err := cfg.UpdateConfig(cfg) // update table map
	if err != nil {
		t.Fatal(err)
	}

	err = cfg.UpdateDB()
	if err != nil {
		t.Fatal(err)
	}

	clearDB := utilDB.CreateDatabase(t, cfg)

	return cfg, clearDB
}
