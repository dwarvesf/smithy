package config

import (
	"fmt"

	"github.com/dwarvesf/smithy/backend/config/agent"
	"github.com/dwarvesf/smithy/backend/config/agent/verify"
	"github.com/jinzhu/gorm"
)

// NewAgentConfig get agent config
func NewAgentConfig(r agent.ConfigReader) (*agent.Config, error) {
	cfg, err := r.Read()
	if err != nil {
		return nil, err
	}

	if cfg.VerifyConfig {
		err = checkAgentConfig(cfg)
		if err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

// checkAgentConfig check agent config is correct
func checkAgentConfig(c *agent.Config) error {
	return checkModelList(c)
}

func checkModelList(c *agent.Config) error {
	switch c.DBType {
	case "postgres":
		return checkModelListPG(c)
	default:
		return fmt.Errorf("using not support database type %v", c.DBType)
	}
}

func checkModelListPG(c *agent.Config) error {
	db, err := gorm.Open("postgres", c.DBConnectionString())
	if err != nil {
		return err
	}

	return verify.NewPGStore(c.DBName, c.DBSchemaName, db).Verify(c.ModelList)
}
