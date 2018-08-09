package automigrate

import (
	"fmt"

	"github.com/jinzhu/gorm"

	agentConfig "github.com/dwarvesf/smithy/config/agent"
	"github.com/dwarvesf/smithy/config/agent/verify"
	"github.com/dwarvesf/smithy/config/database"
)

// AutoMigrater interface for automigrate implement
type AutoMigrater interface {
	Migrate([]agentConfig.MissingColumns) error
}

// AutoMigrate .
func AutoMigrate(cfg *agentConfig.Config) error {
	switch cfg.DBType {
	case "postgres":
		return autoMigrationPG(cfg)
	default:
		return fmt.Errorf("using not support database type: %s", cfg.DBType)
	}
}

func autoMigrationPG(cfg *agentConfig.Config) error {
	db, err := gorm.Open("postgres", cfg.DBConnectionString())
	if err != nil {
		return err
	}

	models := []database.Model{}
	for _, m := range cfg.ModelList {
		if m.AutoMigration {
			models = append(models, m)
		}
	}

	verifyStore := verify.NewPGStore(cfg.DBName, cfg.DBSchemaName, db)
	missmap, err := verifyStore.MissingColumns(models)
	if err != nil {
		return err
	}

	s := NewPGStore(cfg.DBName, cfg.DBSchemaName, db)
	err = s.Migrate(missmap)
	if err != nil {
		return err
	}

	return nil
}
