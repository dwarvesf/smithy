package config

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/jinzhu/gorm"

	agentConfig "github.com/dwarvesf/smithy/agent/config"
	"github.com/dwarvesf/smithy/common/database"
)

// Reader interface for reading config for agent
type Reader interface {
	Read() (*Config, error)
}

// Writer interface for reading config for agent
type Writer interface {
	Write(cfg *Config) error
}

// Querier interface for reading config for agent
type Querier interface {
	ListVersion() ([]Version, error)
	LastestVersion() (*Config, error)
}

// ReaderWriterQuerier compose interface for read/write/query config for agent
type ReaderWriterQuerier interface {
	Reader
	Writer
	Querier
}

// Config contain config for dashboard
type Config struct {
	SerectKey           string `yaml:"agent_serect_key"`
	AgentURL            string `yaml:"agent_url"`
	PersistenceSupport  string `yaml:"persistence_support"`
	PersistenceFileName string `yaml:"persistence_file_name"`

	database.ConnectionInfo `yaml:"-"`
	Databases               []database.Database                  `yaml:"databases_list" json:"databases_list"`
	ModelMap                map[string]map[string]database.Model `yaml:"-" json:"-"`
	Version                 Version                              `yaml:"-" json:"version"`
	db                      map[string]*gorm.DB
	Authentication          Authentication `yaml:"authentication" json:"authentication"`

	sync.Mutex
}

// Version version of backend config
type Version struct {
	Checksum string    `json:"checksum"`
	ID       int       `json:"id"`
	SyncAt   time.Time `json:"sync_at"`
}

// Wrapper use to hide detail of a config
type Wrapper struct {
	cfg *Config
}

// NewWrapper .
func NewWrapper(cfg *Config) *Wrapper {
	return &Wrapper{cfg}
}

// SyncConfig get synchronized config from wrapper
func (w *Wrapper) SyncConfig() *Config {
	w.cfg.Lock()
	defer w.cfg.Unlock()
	return w.cfg
}

// DB get db connection from config
func (c *Config) DB(dbName string) *gorm.DB {
	return c.db[dbName]
}

// DBs get db connection from config
func (c *Config) DBs() map[string]*gorm.DB {
	return c.db
}

// CheckSum to checksum md5 when agent-sync check version
func (c *Config) CheckSum() (string, error) {
	buff, err := json.Marshal(c)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", md5.Sum(buff)), nil
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

	agentCfg := &agentConfig.Config{}
	err = json.NewDecoder(res.Body).Decode(agentCfg)
	if err != nil {
		return err
	}

	return c.UpdateConfigFromAgentConfig(agentCfg)
}

// UpdateConfigFromAgentConfig update config from AgentConfig
func (c *Config) UpdateConfigFromAgentConfig(agentCfg *agentConfig.Config) error {
	// Copy config file into tempCfg
	tempCfg := Config{}

	tempCfg.ConnectionInfo = agentCfg.ConnectionInfo
	tempCfg.DBUsername = agentCfg.UserWithACL.Username
	tempCfg.DBPassword = agentCfg.UserWithACL.Password
	tempCfg.Databases = agentCfg.Databases

	// If available new version, update config then save it into persistence
	checksum, err := tempCfg.CheckSum()
	if err != nil {
		return err
	}
	if checksum == c.Version.Checksum {
		return nil
	}

	err = c.UpdateConfig(&tempCfg)
	if err != nil {
		return err
	}

	c.Version.Checksum = checksum
	c.Version.SyncAt = time.Now()

	wr := NewBoltPersistent(c.PersistenceFileName, 0)
	return wr.Write(c)
}

// AddHook add hook to configuration
func (c *Config) AddHook(tableName, hookType, content string) error {
	for i := range c.Databases {
		models := c.Databases[i].ModelList
		for i := range models {
			if models[i].TableName == tableName {
				err := models[i].AddHook(hookType, content)
				if err != nil {
					return err
				}
				return nil
			}
		}
	}

	return errors.New("table_name not exist")
}

// UpdateConfig update configuration
func (c *Config) UpdateConfig(cfg *Config) error {
	// check config was enable
	c.Lock()
	defer c.Unlock()

	c.ConnectionInfo = cfg.ConnectionInfo
	c.DBUsername = cfg.DBUsername
	c.DBPassword = cfg.DBPassword
	c.Databases = cfg.Databases
	c.Version = cfg.Version

	for k := range c.ModelMap {
		delete(c.ModelMap, k)
	}

	for _, db := range c.Databases {
		tmp := database.Models(db.ModelList).GroupByName()
		c.ModelMap[db.DBName] = make(map[string]database.Model)
		for k := range tmp {
			c.ModelMap[db.DBName][k] = tmp[k]
		}
	}

	return c.UpdateDB()
}

// ChangeVersion get config in persistent by version number
func (c *Config) ChangeVersion(id int) error {
	reader := NewBoltPersistent(c.PersistenceFileName, id)
	cfg, err := reader.Read()
	if err != nil {
		return err
	}

	return c.UpdateConfig(cfg)
}

// UpdateDB update db connection
func (c *Config) UpdateDB() error {
	c.db = make(map[string]*gorm.DB)
	for i := range c.Databases {
		newDB, err := c.openNewDBConnection(c.Databases[i].DBName)
		if err != nil {
			// TODO: add nicer error
			return err
		}
		c.db[c.Databases[i].DBName] = newDB
	}

	return nil
}

// TODO: extend for using mutiple DB
func (c *Config) openNewDBConnection(dbName string) (*gorm.DB, error) {
	sslmode := "disable"
	if c.DBSSLModeOption == "enable" {
		sslmode = "require"
	}

	dbstring := fmt.Sprintf("user=%s dbname=%s sslmode=%s password=%s host=%s port=%s",
		c.DBUsername,
		dbName,
		sslmode,
		c.DBPassword,
		c.DBHostname,
		c.DBPort,
	)

	return gorm.Open("postgres", dbstring)
}

// ConvertUserListToMap convert user list to user map
func (c *Config) ConvertUserListToMap() map[string]User {
	userMap := make(map[string]User)

	for _, user := range c.Authentication.UserList {
		userMap[user.Username] = user
	}
	return userMap
}

// Authentication use to authenticate
type Authentication struct {
	SerectKey string `yaml:"secret_key" json:"secret_key"`
	UserList  []User `yaml:"users" json:"users"`
}

// User .
type User struct {
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	Role     string `yaml:"role" json:"role"`
}
