// +build integration

package test

import (
	"testing"

	_ "github.com/jinzhu/gorm/dialects/postgres"

	backendConfig "github.com/dwarvesf/smithy/backend/config"
)

// Test step:
// - Create random config
// - Save it in persistent
// - Read config in persistent
// - Compare 2 config
func TestStoreConfigInPersistent_1(t *testing.T) {
	bCfg, clearDB := CreateBackendConfig(t)
	defer clearDB()

	agentCfg := CreateAgentConfig(t)

	t.Run("Store config in persistent", func(t *testing.T) {
		err := bCfg.UpdateConfigFromAgentConfig(agentCfg)
		if err != nil {
			t.Fatalf("Fail to update agent config error = %v", err)
		}

		sCfg, err := backendConfig.NewBoltPersistent(bCfg.PersistenceFileName, bCfg.Version.ID).Read()
		if err != nil {
			t.Fatalf("Fail to read config from persistent error = %v", err)
		}

		if bCfg.Version.Checksum != sCfg.Version.Checksum {
			t.Errorf("Config.UpdateConfigFromAgentConfig() error = %v", err)
		}
	})
}

// Test step:
// - Create 2 random config
// - Save configs in persistent
// - Get the lastest version of config
// - Compare
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
		sCfg_1, err := backendConfig.NewBoltPersistent(bCfg.PersistenceFileName, v1.ID).Read()
		if err != nil {
			t.Fatalf("Fail to read config from persistent error = %v", err)
		}

		// read the lastest config from persistent
		sCfg_2, err := backendConfig.NewBoltPersistent(bCfg.PersistenceFileName, 0).LastestVersion()
		if err != nil {
			t.Fatalf("Fail to read lastest config from persistent error = %v", err)
		}

		if v1.Checksum != sCfg_1.Version.Checksum ||
			v2.Checksum != sCfg_2.Version.Checksum {
			t.Errorf("Config.UpdateConfigFromAgentConfig() error = %v", err)
		}
	})
}

// Test step:
// - Add 2 random config to persistent
// - Revert version
// - Add a new config
// - Check 3th config is the lastest config
func TestRevertConfig(t *testing.T) {
	bCfg, clearDB := CreateBackendConfig(t)
	defer clearDB()

	agentCfg_1 := CreateAgentConfig(t)
	agentCfg_2 := CreateAgentConfig(t)
	agentCfg_3 := CreateAgentConfig(t)

	t.Run("Revert config", func(t *testing.T) {
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

		// change version
		err = bCfg.ChangeVersion(v1.ID)
		if err != nil {
			t.Fatalf("Fail to change version of config error = %v", err)
		}

		if v1.Checksum != bCfg.Version.Checksum {
			t.Errorf("Config.UpdateConfigFromAgentConfig() error = %v", err)
		}

		// update config 3 => check if after revert version, adding new config
		err = bCfg.UpdateConfigFromAgentConfig(agentCfg_3)
		if err != nil {
			t.Fatalf("Fail to update agent config error = %v", err)
		}

		// get the lastest version config
		lastVerConfig, err := backendConfig.NewBoltPersistent(bCfg.PersistenceFileName, 0).LastestVersion()
		if err != nil {
			t.Fatalf("Fail to read lastest config from persistent error = %v", err)
		}

		//get 2th config
		v2Config, err := backendConfig.NewBoltPersistent(bCfg.PersistenceFileName, v2.ID).Read()
		if err != nil {
			t.Fatalf("Fail to read lastest config from persistent error = %v", err)
		}

		if v2Config == nil ||
			v2Config.Version.Checksum != v2.Checksum ||
			lastVerConfig.Version.Checksum != bCfg.Version.Checksum {
			t.Errorf("Config.UpdateConfigFromAgentConfig() error = %v", err)
		}
	})
}
