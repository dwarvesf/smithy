package agent

import (
	"fmt"

	"github.com/jinzhu/gorm"

	agentConfig "github.com/dwarvesf/smithy/agent/config"
	"github.com/dwarvesf/smithy/agent/dbtool/drivers"
	"github.com/dwarvesf/smithy/common/database"
	utilPGDb "github.com/dwarvesf/smithy/common/utils/database/pg"
)

const (
	pgDriver = "postgres"
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
	case pgDriver:
		return checkModelListPG(c)
	default:
		return fmt.Errorf("using not support database type %v", c.DBType)
	}
}

func checkModelListPG(c *agentConfig.Config) error {
	for _, dbase := range c.Databases {
		db, err := gorm.Open("postgres", c.DBConnectionString(dbase.DBName))
		if err != nil {
			return err
		}
		defer db.Close()

		err = drivers.NewPGStore(dbase.DBName, c.DBSchemaName, db).Verify(dbase.ModelList)
		if err != nil {
			return err
		}
	}

	return nil
}

// CreateUserWithACL using config to auto migrate missing columns and table
func CreateUserWithACL(cfg *agentConfig.Config, forceCreate bool) (*database.User, error) {
	switch cfg.DBType {
	case pgDriver:
		return createUserWithACLPG(cfg, forceCreate)
	default:
		return nil, fmt.Errorf("using not support database type: %s", cfg.DBType)
	}
}

// createUserWithACLPG create user with access list in model
func createUserWithACLPG(cfg *agentConfig.Config, forceCreate bool) (*database.User, error) {
	user := &database.User{
		Username: cfg.UserWithACL.Username,
		Password: cfg.UserWithACL.Password,
	}

	// delete user permission
	force := forceCreate || cfg.ForceRecreate
	if force {
		for _, dbase := range cfg.Databases {
			db, err := gorm.Open("postgres", cfg.DBConnectionString(dbase.DBName))
			if err != nil {
				return nil, err
			}
			defer db.Close()

			s := drivers.NewPGStore(cfg.DBName, cfg.DBSchemaName, db)
			err = s.RemoveACLUser(user.Username)
			if err != nil {
				return nil, err
			}
		}
	}

	// create user & grant permision
	isCreateUser := false
	for _, dbase := range cfg.Databases {
		db, err := gorm.Open("postgres", cfg.DBConnectionString(dbase.DBName))
		if err != nil {
			return nil, err
		}
		defer db.Close()

		s := drivers.NewPGStore(cfg.DBName, cfg.DBSchemaName, db)

		// priority passing argument than config file
		if !isCreateUser {
			if err = s.CreateACLUser(user, force); err != nil {
				return nil, err
			}
			isCreateUser = true
		}

		// gran permission
		err = s.CreateUserWithACL(dbase.ModelList, user, true)

		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

// AutoMigrate using config to auto migrate missing columns and table
func AutoMigrate(cfg *agentConfig.Config) error {
	switch cfg.DBType {
	case pgDriver:
		return autoMigrationPG(cfg)
	default:
		return fmt.Errorf("using not support database type: %s", cfg.DBType)
	}
}

func autoMigrationPG(cfg *agentConfig.Config) error {
	for _, d := range cfg.Databases {
		err := utilPGDb.CreatePGDatabase(cfg.DBPort, d.DBName)
		if err != nil {
			return err
		}

		db, err := gorm.Open("postgres", cfg.DBConnectionString(d.DBName))
		if err != nil {
			return err
		}
		defer db.Close()

		models := []database.Model{}
		for _, m := range d.ModelList {
			if m.AutoMigration {
				models = append(models, m)
			}
		}

		s := drivers.NewPGStore(d.DBName, cfg.DBSchemaName, db)
		missmap, err := s.MissingColumns(models)
		if err != nil {
			return err
		}

		err = s.AutoMigrate(missmap)
		if err != nil {
			return err
		}
	}

	return nil
}
