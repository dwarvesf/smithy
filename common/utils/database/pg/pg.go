package pg

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	backendConfig "github.com/dwarvesf/smithy/backend/config"
)

func CreateDatabase(t *testing.T, cfg *backendConfig.Config) func() {
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

	return func() {
		err := cfg.DB().Exec("DROP SCHEMA " + schemaName + " CASCADE").Error
		if err != nil {
			t.Fatalf("Fail to drop database. %s", err.Error())
		}
	}
}

type User struct {
	Id   int `sql:"primary_key"`
	Name string
}
