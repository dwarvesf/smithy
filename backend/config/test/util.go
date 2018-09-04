package test

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	agentConfig "github.com/dwarvesf/smithy/agent/config"
	backendConfig "github.com/dwarvesf/smithy/backend/config"
	"github.com/dwarvesf/smithy/common/database"
	utilDB "github.com/dwarvesf/smithy/common/utils/database/bolt"
)

const (
	dbHost     = "localhost"
	dbPort     = "5439"
	dbUser     = "postgres"
	dbPassword = "example"
	dbName     = "test"
)

// CreateModelList fake model list from agent
func CreateModelList() []database.Model {
	rand.Seed(time.Now().UnixNano())
	randTableName := strconv.FormatInt(rand.Int63(), 10)

	dm := []database.Model{
		{
			TableName: randTableName,
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

// CreateBackendConfig fake config for test
func CreateBackendConfig(t *testing.T) (*backendConfig.Config, func()) {
	cfg := &backendConfig.Config{
		ModelList: CreateModelList(),
		ModelMap:  make(map[string]database.Model),
		ConnectionInfo: database.ConnectionInfo{
			DBType:          "postgres",
			DBUsername:      dbUser,
			DBPassword:      dbPassword,
			DBName:          dbName,
			DBPort:          dbPort,
			DBHostname:      dbHost,
			DBSSLModeOption: "false",
			UserWithACL: database.User{
				Username: dbUser,
				Password: dbPassword,
			},
		},
	}

	persistenceFileName, clearDB := utilDB.CreateDatabase(t)

	cfg.PersistenceFileName = persistenceFileName

	return cfg, clearDB
}

// CreateAgentConfig fake config for test
func CreateAgentConfig(t *testing.T) *agentConfig.Config {
	cfg := &agentConfig.Config{
		ModelList: CreateModelList(),
		ConnectionInfo: database.ConnectionInfo{
			DBType:          "postgres",
			DBUsername:      dbUser,
			DBPassword:      dbPassword,
			DBName:          dbName,
			DBPort:          dbPort,
			DBHostname:      dbHost,
			DBSSLModeOption: "false",
			UserWithACL: database.User{
				Username: dbUser,
				Password: dbPassword,
			},
		},
	}

	return cfg
}
