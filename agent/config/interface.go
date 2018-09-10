package config

import (
	"fmt"

	"github.com/dwarvesf/smithy/common/database"
)

// Reader interface for reading config for agent
type Reader interface {
	Read() (*Config, error)
}

// Writer interface for reading config for agent
type Writer interface {
	Write(res *Config) error
}

// Config contain config for agent
type Config struct {
	SerectKey               string `yaml:"serect_key" json:"-"`
	VerifyConfig            bool   `yaml:"verify_config" json:"-"`
	database.ConnectionInfo `yaml:"database_connection_info" json:"database_connection_info"`
	ForceRecreate           bool      `yaml:"force_recreate" json:"force_recreate"`
	Databases               Databases `yaml:"databases_list" json:"databases_list"`
}

// Databases contains list of databases
type Databases struct {
	Name      string           `yaml:"db_name" json:"db_name"`
	ModelList []database.Model `yaml:"model_list" json:"model_list"`
}

// DBConnectionString get pg connection string
func (c Config) DBConnectionString() string {
	switch c.DBType {
	case "postgres":
		return c.pgConnectionString()
	default:
		return ""
	}
}

// PGConnectionString get pg connection string
func (c Config) pgConnectionString() string {
	return fmt.Sprintf("user=%s dbname=%s sslmode=%s password=%s host=%s port=%s",
		c.DBUsername,
		c.DBName,
		c.DBSSLModeOption,
		c.DBPassword,
		c.DBHostname,
		c.DBPort,
	)
}

// ExistingColumnByTableName list of existing column_schema in database
type ExistingColumnByTableName map[string][]ColumnSchema

// ColumnSchema define of a column by database schema
type ColumnSchema struct {
	ColumnName             string
	UdtName                string // data_type name
	CharacterMaximumLength string
	IsNullable             string
	Order                  string
	ColumnDefault          string
	IsPrimary              bool
}

// ColumnSchemas array of ColumnSchema
type ColumnSchemas []ColumnSchema

// GroupByColumnName group column schema by column name
func (cols ColumnSchemas) GroupByColumnName() map[string][]ColumnSchema {
	res := make(map[string][]ColumnSchema)
	for _, col := range cols {
		if _, ok := res[col.ColumnName]; ok {
			res[col.ColumnName] = append(res[col.ColumnName], []ColumnSchema{col}...)
		} else {
			res[col.ColumnName] = []ColumnSchema{col}
		}
	}

	return res
}

// UpdateByColumnDefinition update column by column definition
func (c *ColumnSchema) UpdateByColumnDefinition(col database.Column) {
	c.ColumnName = col.Name
	c.IsNullable = "YES"
	if !col.IsNullable {
		c.IsNullable = "NO"
	}
	switch col.Type {
	case "int":
		c.UdtName = "int4"
	case "string":
		c.UdtName = "text"
	case "timestamp":
		c.UdtName = "timestamptz"
	default:
		c.UdtName = "text"
	}
	c.ColumnDefault = col.DefaultValue
	c.IsPrimary = col.IsPrimary
}

// MissingColumns define missing columns and groupted by table name
type MissingColumns struct {
	TableName string
	Columns   []ColumnSchema
	IsCreate  bool
}

// IsNeedMigrate check missing column is needed to make a migrate
func (mcs MissingColumns) IsNeedMigrate() bool {
	return len(mcs.Columns) > 0
}
