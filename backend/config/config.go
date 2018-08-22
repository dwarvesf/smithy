package config

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/boltdb/bolt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	yaml "gopkg.in/yaml.v2"

	agentConfig "github.com/dwarvesf/smithy/agent/config"
	"github.com/dwarvesf/smithy/common/database"
)

// Reader interface for reading config for agent
type Reader interface {
	Read() (*Config, error)
}

// Config contain config for dashboard
type Config struct {
	SerectKey           string `yaml:"agent_serect_key"`
	AgentURL            string `yaml:"agent_url"`
	PersistenceSupport  string `yaml:"persistence_support"`
	PersistenceFileName string `yaml:"persistence_file_name"`

	// TODO: extend for using mutiple persistence DB
	PersistenceDB           *bolt.DB
	database.ConnectionInfo `yaml:"-"`
	ModelList               []database.Model `yaml:"-"`
	Version                 string           `yaml:"-" json:"-"`
	db                      *gorm.DB

	sync.Mutex
}

// Wrapper use to hide detail of a config
type Wrapper struct {
	cfg *Config
}

// NewWrapper .
func NewWrapper(cfg *Config) *Wrapper {
	return &Wrapper{cfg}
}

// Config get sync config from wrapper
func (w *Wrapper) Config() *Config {
	w.cfg.Lock()
	defer w.cfg.Unlock()
	return w.cfg
}

// DB get db connection from config
func (c *Config) DB() *gorm.DB {
	return c.db
}

// CheckSum to checksum sha256 when agent-sync check version
func (c *Config) CheckSum() (string, error) {
	buff, err := yaml.Marshal(c)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(buff)
	fmt.Printf("%x", sum)
	return string(sum[:]), nil
}

// UpdateConfigFromAgent update configuration from agent
func (c *Config) UpdateConfigFromAgent() error {
	// check config was enable

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

	agentCfg := agentConfig.Config{}
	err = json.NewDecoder(res.Body).Decode(&agentCfg)
	if err != nil {
		return err
	}

	tempCfg := *c
	checksum, err := tempCfg.CheckSum()
	if checksum == c.Version {
		return nil
	}

	c.Lock()
	defer c.Unlock()

	c.ConnectionInfo = agentCfg.ConnectionInfo
	c.DBUsername = c.UserWithACL.Username
	c.DBPassword = c.UserWithACL.Password
	c.ModelList = agentCfg.ModelList
	c.Version = checksum
	err = c.UpdateDB()
	if err != nil {
		return err
	}

	return nil
}

// UpdateConfig update configuration
func (c *Config) UpdateConfig(cfg *Config) error {
	// check config was enable
	c.Lock()
	defer c.Unlock()

	c = cfg
	return c.UpdateDB()
}

// UpdateDB update db connection
func (c *Config) UpdateDB() error {
	newDB, err := c.openNewDBConnection()
	if err != nil {
		// TODO: add nicer error
		return err
	}
	c.db = newDB

	return nil
}

// TODO: extend for using mutiple DB
func (c *Config) openNewDBConnection() (*gorm.DB, error) {
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
