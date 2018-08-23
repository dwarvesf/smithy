package config

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	"github.com/jinzhu/gorm"

	agentConfig "github.com/dwarvesf/smithy/agent/config"
	"github.com/dwarvesf/smithy/common/database"
)

// Reader interface for reading config for agent
type UserConfigReader interface {
	Read() (*UserConfig, error)
}

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

	// TODO: extend for using mutiple persistence DB
	PersistenceDB           *bolt.DB
	database.ConnectionInfo `yaml:"-"`
	ModelList               []database.Model `yaml:"-"`
	Version                 Version          `yaml:"-" json:"version"`
	db                      *gorm.DB
	UserConfig              *UserConfig

	sync.Mutex
}

// Version version of backend config
type Version struct {
	Checksum      string    `json:"checksum"`
	VersionNumber int64     `json:"version_number"`
	SyncAt        time.Time `json:"sync_at"`
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

// CheckSum to checksum md5 when agent-sync check version
func (c *Config) CheckSum() (string, error) {
	buff, err := json.Marshal(c)
	if err != nil {
		return "", err
	}

	h := md5.New()
	return fmt.Sprintf("%x", h.Sum(buff)), nil
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

	// Copy config file into tempCfg
	tempCfg := Config{}

	tempCfg.ConnectionInfo = agentCfg.ConnectionInfo
	tempCfg.DBUsername = agentCfg.UserWithACL.Username
	tempCfg.DBPassword = agentCfg.UserWithACL.Password
	tempCfg.ModelList = agentCfg.ModelList

	// If available new version, update config then save it into persistence
	tmpVer := tempCfg.Version
	tempCfg.Version = Version{}
	checksum, err := tempCfg.CheckSum()
	if err != nil {
		return err
	}
	if checksum == c.Version.Checksum {
		return nil
	}
	tempCfg.Version = tmpVer

	err = c.UpdateConfig(&tempCfg)
	if err != nil {
		return err
	}

	c.Version.Checksum = checksum
	c.Version.SyncAt = time.Now()
	c.Version.VersionNumber = c.Version.SyncAt.Unix()

	wr := NewBoltPersistent(c.PersistenceDB, 0)
	return wr.Write(c)
}

// UpdateConfig update configuration
func (c *Config) UpdateConfig(cfg *Config) error {
	// check config was enable
	c.Lock()
	defer c.Unlock()

	c.ConnectionInfo = cfg.ConnectionInfo
	c.DBUsername = cfg.DBUsername
	c.DBPassword = cfg.DBPassword
	c.ModelList = cfg.ModelList
	c.Version = cfg.Version

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

// add user config
func (c *Config) AddUserConfigToConfig(userConfig *UserConfig) *Config {
	c.UserConfig = userConfig
	return c
}

type UserConfig struct {
	SerectKey string          `yaml:"secret_key" json:"secret_key"`
	Users     map[string]User `yaml:"users" json:"users"`
}

type User struct {
	Password string `yaml:"password" json:"password"`
	Role     string `yaml:"role" json:"role"`
	Acl      string `yaml:"acl" json:"acl"`
}
