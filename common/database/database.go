package database

import (
	"errors"

	"github.com/dwarvesf/smithy/backend/hook"
)

// ConnectionInfo store information to connect to a database
type ConnectionInfo struct {
	DBType          string `yaml:"db_type" json:"db_type"`
	DBUsername      string `yaml:"db_username" json:"db_username"`
	DBPassword      string `yaml:"db_password" json:"db_password"`
	DBName          string `yaml:"db_name" json:"db_name"`
	DBSSLModeOption string `yaml:"db_ssl_mode_option" json:"db_ssl_mode_option"`
	DBHostname      string `yaml:"db_hostname" json:"db_hostname"`
	DBPort          string `yaml:"db_port" json:"db_port"`
	DBEnvironment   string `yaml:"db_environment" json:"db_environment"`
	DBSchemaName    string `yaml:"db_schema_name" json:"db_schema_name"`
	UserWithACL     User   `yaml:"user_with_acl" json:"user_with_acl"`
}

// Model store information of model can manage
type Model struct {
	ACL               string    `yaml:"acl" json:"acl"`
	ACLDetail         ACLDetail `yaml:"-" json:"-"`
	TableName         string    `yaml:"table_name" json:"table_name"`
	Columns           []Column  `yaml:"columns" json:"columns"`
	AutoMigration     bool      `yaml:"auto_migration" json:"auto_migration"` // auto_migration if table not exist or misisng column
	DisplayName       string    `yaml:"display_name" json:"display_name"`
	NameDisplayColumn string    `yaml:"name_display_column" json:"name_display_column"`
	Hooks             Hooks     `yaml:"hooks" json:"hooks"`
}

// AddHook add hook to model base on hookType
func (m *Model) AddHook(hookType, content string) error {
	h := Hook{Enable: true, Content: content} // TODO: add check hook content

	switch hookType {
	case hook.BeforeCreate:
		m.Hooks.BeforeCreate = h
	case hook.AfterCreate:
		m.Hooks.AfterCreate = h
	case hook.BeforeUpdate:
		m.Hooks.BeforeUpdate = h
	case hook.AfterUpdate:
		m.Hooks.AfterUpdate = h
	case hook.BeforeDelete:
		m.Hooks.BeforeDelete = h
	case hook.AfterDelete:
		m.Hooks.AfterDelete = h
	default:
		return errors.New("hook_type is not exist")

	}

	return nil
}

// IsBeforeCreateEnable check before_create is enable
func (m *Model) IsBeforeCreateEnable() bool {
	return m.Hooks.BeforeCreate.Enable
}

// IsAfterCreateEnable check after_create is enable
func (m *Model) IsAfterCreateEnable() bool {
	return m.Hooks.AfterCreate.Enable
}

// IsBeforeUpdateEnable check before_update is enable
func (m *Model) IsBeforeUpdateEnable() bool {
	return m.Hooks.BeforeUpdate.Enable
}

// IsAfterUpdateEnable check after_update is enable
func (m *Model) IsAfterUpdateEnable() bool {
	return m.Hooks.AfterUpdate.Enable
}

// IsBeforeDeleteEnable check before_delete is enable
func (m *Model) IsBeforeDeleteEnable() bool {
	return m.Hooks.BeforeDelete.Enable
}

// IsAfterDeleteEnable check after_delete is enable
func (m *Model) IsAfterDeleteEnable() bool {
	return m.Hooks.AfterDelete.Enable
}

// Hooks hook declaration for a model
type Hooks struct {
	BeforeCreate Hook `yaml:"before_create" json:"before_create"`
	AfterCreate  Hook `yaml:"after_create" json:"after_create"`
	BeforeUpdate Hook `yaml:"before_update" json:"before_update"`
	AfterUpdate  Hook `yaml:"after_update" json:"after_update"`
	BeforeDelete Hook `yaml:"before_delete" json:"before_delete"`
	AfterDelete  Hook `yaml:"after_delete" json:"after_delete"`
}

// Hook define a
type Hook struct {
	Enable  bool   `yaml:"enable" json:"enable"` // Is enable a hook
	Content string `yaml:"content" json:"content"`
}

// User detail of a user in database
type User struct {
	Username string
	Password string
}

// MakeACLDetailFromACL update access list detail
func (m *Model) MakeACLDetailFromACL() {
	ad := ACLDetail{}
	for _, r := range m.ACL {
		switch r {
		case 'C', 'c':
			ad.Insert = true
		case 'R', 'r':
			ad.Select = true
		case 'U', 'u':
			ad.Update = true
		case 'D', 'd':
			ad.Delete = true
		}
	}

	m.ACLDetail = ad
}

// ACLDetail .
type ACLDetail struct {
	Select bool
	Insert bool
	Update bool
	Delete bool
}

// Models array of model
type Models []Model

// GroupByName create map model by table name
func (ms Models) GroupByName() map[string]Model {
	res := make(map[string]Model)
	for _, m := range ms {
		res[m.TableName] = m
	}

	return res
}

// ColumnsByTableName create map columns by table name from array of column
func (ms Models) ColumnsByTableName() map[string][]Column {
	res := make(map[string][]Column)
	for _, m := range ms {
		if _, ok := res[m.TableName]; ok {
			res[m.TableName] = append(res[m.TableName], m.Columns...)
		} else {
			res[m.TableName] = m.Columns
		}
	}

	return res
}

// ModelByTableName create map model by table name
func (ms Models) ModelByTableName() map[string]Model {
	res := make(map[string]Model)
	for _, m := range ms {
		res[m.TableName] = m
	}

	return res
}

// Column store information of a column
type Column struct {
	Name         string `yaml:"name" json:"name"`
	Type         string `yaml:"type" json:"type"`
	Tags         string `yaml:"tags" json:"tags"`
	IsNullable   bool   `yaml:"is_nullable" json:"is_nullable"`
	IsPrimary    bool   `yaml:"is_primary" json:"is_primary"`
	DefaultValue string `yaml:"default_value" json:"default_value"`
}

// Columns array of column
type Columns []Column

// GroupByName group column by name
func (cols Columns) GroupByName() map[string][]Column {
	res := make(map[string][]Column)
	for _, col := range cols {
		if _, ok := res[col.Name]; ok {
			res[col.Name] = append(res[col.Name], col)
		} else {
			res[col.Name] = []Column{col}
		}
	}

	return res
}

// Names return names of all columns
func (cols Columns) Names() []string {
	res := []string{}
	for _, col := range cols {
		res = append(res, col.Name)

	}

	return res
}
