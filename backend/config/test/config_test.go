// +build integration

package test

import (
	"testing"

	_ "github.com/jinzhu/gorm/dialects/postgres"

	backendConfig "github.com/dwarvesf/smithy/backend/config"
)

func TestStoreConfigInPersistent_1(t *testing.T) {
	bCfg, clearDB := CreateBackendConfig(t)
	defer clearDB()

	agentCfg := CreateAgentConfig(t)

	t.Run("Store config in persistent", func(t *testing.T) {
		err := bCfg.UpdateConfigFromAgentConfig(agentCfg)
		if err != nil {
			t.Fatalf("Fail to update agent config error = %v", err)
		}

		sCfg, err := backendConfig.NewBoltPersistent(bCfg.PersistenceDB, bCfg.Version.VersionNumber).Read()
		if err != nil {
			t.Fatalf("Fail to read config from persistent error = %v", err)
		}

		if bCfg.Version.Checksum != sCfg.Version.Checksum {
			t.Errorf("Config.UpdateConfigFromAgentConfig() error = %v", err)
		}
	})
}

func TestStoreConfigInPersistent_2(t *testing.T) {
	bCfg, clearDB := CreateBackendConfig(t)
	defer clearDB()

	agentCfg_1 := CreateAgentConfig(t)
	agentCfg_2 := CreateAgentConfig(t)

	t.Run("Store config in persistent", func(t *testing.T) {
		// update config 1
		err := bCfg.UpdateConfigFromAgentConfig(agentCfg_1)
		if err != nil {
			t.Fatalf("Fail to update agent config error = %v", err)
		}
		v1 := bCfg.Version

		// update config 2
		err = bCfg.UpdateConfigFromAgentConfig(agentCfg_2)
		if err != nil {
			t.Fatalf("Fail to update agent config error = %v", err)
		}
		v2 := bCfg.Version

		// read from persistent by ver num
		sCfg_1, err := backendConfig.NewBoltPersistent(bCfg.PersistenceDB, v1.VersionNumber).Read()
		if err != nil {
			t.Fatalf("Fail to read config from persistent error = %v", err)
		}

		// read the lastest config from persistent
		sCfg_2, err := backendConfig.NewBoltPersistent(bCfg.PersistenceDB, 0).LastestVersion()
		if err != nil {
			t.Fatalf("Fail to read config from persistent error = %v", err)
		}

		if v1.Checksum != sCfg_1.Version.Checksum ||
			v2.Checksum != sCfg_2.Version.Checksum {
			t.Errorf("Config.UpdateConfigFromAgentConfig() error = %v", err)
		}
	})
}
