package automigrate

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"

	agentConfig "github.com/dwarvesf/smithy/agent/config"
)

type pgStore struct {
	databaseName string
	schemaName   string
	db           *gorm.DB
}

// NewPGStore implementer for automigrater inteface
func NewPGStore(databaseName, schemaName string, db *gorm.DB) AutoMigrater {
	return &pgStore{databaseName, schemaName, db}
}

// Migrate migreate missing column
func (s *pgStore) Migrate(ms []agentConfig.MissingColumns) error {
	if err := s.setSearchPath(s.schemaName); err != nil {
		return err
	}

	tx := s.db.Begin()
	var err error
	defer func() {
		if err == nil {
			tx.Commit()
			return
		}
		tx.Rollback()
	}()
	for _, m := range ms {
		if m.IsNeedMigrate() {
			sql, err := s.makeMigrateQuery(m)
			if err != nil {
				return err
			}
			err = tx.Exec(sql).Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *pgStore) makeMigrateQuery(m agentConfig.MissingColumns) (string, error) {
	execQuery := fmt.Sprintf("ALTER TABLE %s.%s %s", s.schemaName, m.TableName, s.makeUpdateQueries(m.Columns))
	if m.IsCreate {
		createQueries, err := s.makeCreateQueries(m.Columns)
		if err != nil {
			return "", fmt.Errorf("create table %s was failed because: %v", m.TableName, err)
		}
		execQuery = fmt.Sprintf("CREATE TABLE %s.%s ( %s );", s.schemaName, m.TableName, createQueries)
	}
	return execQuery, nil
}

func (s *pgStore) makeCreateQueries(cols []agentConfig.ColumnSchema) (string, error) {
	queries := []string{}
	havePrimaryKey := false
	for _, col := range cols {
		dataType := col.UdtName

		optional := ""
		if col.IsNullable == "NO" {
			optional += "NOT NULL"
		}

		// add serial primary key
		if col.IsPrimary && dataType == "int4" {
			havePrimaryKey = true
			queries = append(queries, fmt.Sprintf("%s SERIAL PRIMARY KEY", col.ColumnName))
			continue
		}

		queries = append(queries, fmt.Sprintf("%s %s %s", col.ColumnName, dataType, optional))
	}

	if !havePrimaryKey {
		return "", errors.New("missing primary key")
	}

	return s.groupCreateQueries(queries), nil

}

func (s *pgStore) groupCreateQueries(queries []string) string {
	res := ""
	for i, q := range queries {
		res += " " + q
		delimiter := ","
		if i == len(queries)-1 {
			delimiter = ""
		}
		res += delimiter

	}

	return res
}

func (s *pgStore) makeUpdateQueries(cols []agentConfig.ColumnSchema) string {
	queries := []string{}
	for _, col := range cols {
		dataType := col.UdtName
		optional := ""
		if col.IsNullable == "NO" {
			optional += "NOT NULL"
		}

		queries = append(queries, fmt.Sprintf("ADD COLUMN %s %s %s", col.ColumnName, dataType, optional))
	}

	return s.groupUpdateQueries(queries)
}

func (s *pgStore) groupUpdateQueries(queries []string) string {
	res := ""
	for i, q := range queries {
		res += " " + q
		delimiter := ","
		if i == len(queries)-1 {
			delimiter = ";"
		}
		res += delimiter

	}

	return res
}

// setSearchPath set search path (schema name in term for postgres)
func (s *pgStore) setSearchPath(schemaName string) error {
	return s.db.Exec("SET search_path TO " + schemaName).Error
}
