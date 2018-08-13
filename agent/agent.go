package agent

import (
	"fmt"

	"github.com/jinzhu/gorm"

	agentConfig "github.com/dwarvesf/smithy/agent/config"
	"github.com/dwarvesf/smithy/agent/dbtool/drivers"
	"github.com/dwarvesf/smithy/common/database"
)

// NewConfig get agent config from reader
func NewConfig(r agentConfig.Reader) (*agentConfig.Config, error) {
	cfg, err := r.Read()
	if err != nil {
		return nil, err
	}

	if cfg.VerifyConfig {
		err = checkConfig(cfg)
		if err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

// checkConfig check agent config is correct
func checkConfig(c *agentConfig.Config) error {
	return checkModelList(c)
}

func checkModelList(c *agentConfig.Config) error {
	switch c.DBType {
	case "postgres":
		return checkModelListPG(c)
	default:
		return fmt.Errorf("using not support database type %v", c.DBType)
	}
}

func checkModelListPG(c *agentConfig.Config) error {
	db, err := gorm.Open("postgres", c.DBConnectionString())
	if err != nil {
		return err
	}

	return drivers.NewPGStore(c.DBName, c.DBSchemaName, db).Verify(c.ModelList)
}

// AutoMigrate using config to auto migrate missing columns and table
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

	s := drivers.NewPGStore(cfg.DBName, cfg.DBSchemaName, db)
	missmap, err := s.MissingColumns(models)
	if err != nil {
		return err
	}

	err = s.AutoMigrate(missmap)
	if err != nil {
		return err
	}

	return nil
}
