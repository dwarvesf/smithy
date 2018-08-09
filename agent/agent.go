package agent

import (
	"fmt"

	"github.com/jinzhu/gorm"

	agentConfig "github.com/dwarvesf/smithy/agent/config"
	"github.com/dwarvesf/smithy/agent/verify"
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

	return verify.NewPGStore(c.DBName, c.DBSchemaName, db).Verify(c.ModelList)
}
