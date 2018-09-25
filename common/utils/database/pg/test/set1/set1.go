package set1

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/dwarvesf/smithy/common/database"
	"github.com/jinzhu/gorm"

	backendConfig "github.com/dwarvesf/smithy/backend/config"
	utilDB "github.com/dwarvesf/smithy/common/utils/database/pg"
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

func ACLUsersTable(t *testing.T, cfg *backendConfig.Config) func() {
	db := cfg.DB(utilDB.DBName)

	var (
		schemaName = cfg.DBSchemaName
		username   = cfg.DBUsername
		password   = cfg.DBPassword
	)

	if err := db.Exec(fmt.Sprintf("CREATE ROLE %s LOGIN PASSWORD '%s';", username, password)).Error; err != nil {
		t.Fatalf("Fail to create ROLE. %s", err.Error())
	}
	if err := db.Exec(fmt.Sprintf("GRANT USAGE ON SCHEMA %s TO %s ;", schemaName, username)).Error; err != nil {
		t.Fatalf("Fail to create ROLE. %s", err.Error())
	}
	if err := db.Exec(fmt.Sprintf("GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA %s TO %s;", schemaName, username)).Error; err != nil {
		t.Fatalf("Fail to GRANT SCHEMA. %s", err.Error())
	}
	if err := db.Exec(fmt.Sprintf("GRANT SELECT, INSERT ON users TO %s", username)).Error; err != nil {
		t.Fatalf("Fail to GRANT SELECT, CREATE ON users. %s", err.Error())
	}
	return func() {
		db.Exec(fmt.Sprintf("REASSIGN OWNED BY %s TO postgres;", username))
		db.Exec(fmt.Sprintf("DROP OWNED BY %s;", username))

		err := db.Exec(fmt.Sprintf("DROP ROLE IF EXISTS %s;", username)).Error
		if err != nil {
			t.Fatalf("Fail to drop ROLE. %s", err.Error())
		}
	}
}

// CreateModelList fake model list from agent
func CreateModelList() []database.Model {
	dm := []database.Model{
		{
			TableName: "users",
			ACL:       "cr",
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
	dm := []database.Database{
		{
			DBName:    "test1",
			ModelList: CreateModelList(),
		},
		{
			DBName:    "test2",
			ModelList: CreateModelList(),
		},
	}

	return dm
}

// CreateConfig fake config for test
func CreateConfig(t *testing.T) (*backendConfig.Config, func()) {
	cfg := &backendConfig.Config{
		Databases: CreateDatabaseList(),
		ModelMap:  make(map[string]map[string]database.Model),
		ConnectionInfo: database.ConnectionInfo{
			DBType:          "postgres",
			DBUsername:      utilDB.DBUser,
			DBPassword:      utilDB.DBPassword,
			DBName:          utilDB.DBName,
			DBPort:          utilDB.DBPort,
			DBHostname:      utilDB.DBHost,
			DBSSLModeOption: "disable",
		},
		Authentication: backendConfig.Authentication{
			SerectKey: "lalala",
			UserList: []backendConfig.User{
				{
					Username: "aaa",
					Password: "abc",
					Role:     "admin",
					DatabaseList: []backendConfig.Database{
						{
							DBName: "test",
							Tables: []backendConfig.Table{
								{
									TableName: "users",
									ACL:       "cru",
								},
							},
						},
					},
				},
				{
					Username: "bbb",
					Password: "abc",
					Role:     "user",
				},
			},
		},
	}

	// migrate database
	for _, dbase := range cfg.Databases {
		err := utilDB.CreatePGDatabase(cfg.DBPort, dbase.DBName)
		if err != nil {
			t.Fatalf("Failed to migrate database by error %v", err)
		}
	}

	err := cfg.UpdateConfig(cfg) // update table map
	if err != nil {
		t.Fatal(err)
	}

	clearDBs := utilDB.CreateDatabase(t, cfg)

	return cfg, clearDBs
}
