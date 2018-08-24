package backend

import (
	"errors"

	jwtAuth "github.com/dwarvesf/smithy/backend/auth"
	backendConfig "github.com/dwarvesf/smithy/backend/config"
	"github.com/dwarvesf/smithy/backend/sqlmapper"
	sqlmapperDrv "github.com/dwarvesf/smithy/backend/sqlmapper/drivers"
	"github.com/dwarvesf/smithy/common/database"
)

// NewConfig check dashboard config is correct
func NewConfig(r backendConfig.Reader) (*backendConfig.Config, error) {
	cfg, err := r.Read()
	if err != nil {
		return nil, err
	}
	err = checkConfig(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// checkConfig check agent config is correct
func checkConfig(c *backendConfig.Config) error {
	// TODO: implement dashboard config checking
	return nil
}

// NewSQLMapper create new new sqlmapper to working with request query
func NewSQLMapper(c *backendConfig.Config, tableName string, columns []database.Column) (sqlmapper.Mapper, error) {
	switch c.DBType {
	case "postgres":
		return sqlmapperDrv.NewPGHookStore(
			sqlmapperDrv.NewPGStore(c.DB(), tableName, columns, c.ModelList),
		), nil
	default:
		return nil, errors.New("Uknown DB Driver")
	}
}

func SyncPersistent(c *backendConfig.Config) error {
	switch c.PersistenceSupport {
	case "boltdb":
		boltIO := backendConfig.NewBoltPersistent(c.PersistenceDB, 0)
		lastCfg, err := boltIO.LastestVersion()
		if err != nil {
			return err
		}

		if lastCfg == nil {
			return nil
		}
		return c.UpdateConfig(lastCfg)
	default:
		return errors.New("Uknown Persistent DB")
	}
}

//NewAuthenticate New JWT Authenticain
func NewAuthenticate(c *backendConfig.Config, username string, rule string) *jwtAuth.JWT {
	return jwtAuth.New(c.Authentication.SerectKey, username, rule)
}
