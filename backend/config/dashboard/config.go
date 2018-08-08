package dashboard

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/boltdb/bolt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/dwarvesf/smithy/backend/config/agent"
	"github.com/dwarvesf/smithy/backend/config/database"
)

// Config contain config for agent
type Config struct {
	SerectKey           string `yaml:"agent_serect_key"`
	AgentURL            string `yaml:"agent_url"`
	PersistenceSupport  string `yaml:"persistence_support"`
	PersistenceFileName string `yaml:"persistence_file_name"`

	// TODO: extend for using mutiple persistence DB
	PersistenceDB           *bolt.DB
	database.ConnectionInfo `yaml:"-"`
	ModelList               []database.Model `yaml:"-"`
	db                      *serviceDB
}

// GetDB get db connection from config
func (c Config) GetDB() *gorm.DB {
	c.db.Lock()
	defer c.db.Unlock()

	return c.db.DB
}

// UpdateConfigFromAgent update configuration from agent
func (c Config) UpdateConfigFromAgent() error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", c.AgentURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", c.SerectKey)

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	agentCfg := agent.Config{}
	err = json.NewDecoder(res.Body).Decode(&agentCfg)
	if err != nil {
		return err
	}

	c.ConnectionInfo = agentCfg.ConnectionInfo
	c.ModelList = agentCfg.ModelList
	err = c.updateDB()
	if err != nil {
		return err
	}

	return nil
}

// updateDB update db connection
func (c Config) updateDB() error {
	c.db.Lock()
	defer c.db.Unlock()

	newDB, err := c.openNewDBConnection()
	if err != nil {
		// TODO: add nicer error
		return err
	}
	c.db.DB = newDB

	return nil
}

// TODO: extend for using mutiple DB
func (c Config) openNewDBConnection() (*gorm.DB, error) {
	sslmode := "disable"
	if c.DBSSLModeOption == "enable" {
		sslmode = "require"
	}

	dbstring := fmt.Sprintf("user=%s dbname=%s sslmode=%s password=%s host=%s port=%s",
		c.DBUsername,
		c.DBName,
		sslmode,
		c.DBPassword,
		c.DBHostname,
		c.DBPort,
	)

	return gorm.Open("postgres", dbstring)
}

type serviceDB struct {
	sync.Mutex
	*gorm.DB
}

// ConfigReader interface for reading config for agent
type ConfigReader interface {
	Read() (*Config, error)
}
