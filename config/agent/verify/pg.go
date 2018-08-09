package verify

import (
	"fmt"

	agentConfig "github.com/dwarvesf/smithy/config/agent"
	"github.com/dwarvesf/smithy/config/database"
	"github.com/jinzhu/gorm"
)

type pgStore struct {
	databaseName string
	schemaName   string
	db           *gorm.DB
}

// NewPGStore return an implement for verifier inteface in postgres
func NewPGStore(databaseName, schemaName string, db *gorm.DB) Verifier {
	return &pgStore{databaseName, schemaName, db}
}

// Verify verify for agent model_list
func (s *pgStore) Verify(modelList []database.Model) error {
	missingColumns, err := s.MissingColumns(modelList)
	if err != nil {
		return err
	}
	cols := []agentConfig.ColumnSchema{}
	for _, mc := range missingColumns {
		cols = append(cols, mc.Columns...)
	}

	// TODO: change to return multiple errors
	if len(cols) > 0 {
		return fmt.Errorf("config was missing column %+v", cols)
	}

	return nil
}

func (s *pgStore) MissingColumns(tableDefinitions []database.Model) ([]agentConfig.MissingColumns, error) {
	existColumns, err := s.existColumnsByTableName()
	if err != nil {
		return nil, err
	}

	res := []agentConfig.MissingColumns{}
	colDefs := database.Models(tableDefinitions).ColumnsByTableName()
	for tblName, columns := range colDefs {
		// check not created table
		{
			if _, ok := existColumns[tblName]; !ok {
				missingColumns := []agentConfig.ColumnSchema{}
				for _, col := range columns {
					tmp := agentConfig.ColumnSchema{}
					tmp.UpdateByColumnDefinition(col)
					missingColumns = append(missingColumns, tmp)
				}
				res = append(res, agentConfig.MissingColumns{
					TableName: tblName,
					Columns:   missingColumns,
					IsCreate:  true,
				})
				continue
			}

		}

		// check created table
		{
			colDefs := database.Columns(columns).GroupByName()
			existCols := agentConfig.ColumnSchemas(existColumns[tblName]).GroupByColumnName()
			missingColumns := []agentConfig.ColumnSchema{}
			for colName, cols := range colDefs {
				if _, ok := existCols[colName]; !ok {
					for _, v := range cols {
						tmp := agentConfig.ColumnSchema{}
						tmp.UpdateByColumnDefinition(v)
						missingColumns = append(missingColumns, tmp)
					}
				}
			}

			res = append(res, agentConfig.MissingColumns{
				TableName: tblName,
				Columns:   missingColumns,
			})
		}
	}

	return res, nil
}

func (s *pgStore) existColumnsByTableName() (agentConfig.ExistingColumnByTableName, error) {
	tableNames, err := s.existTableNames()
	if err != nil {
		return nil, err
	}
	cmap := make(map[string][]agentConfig.ColumnSchema)
	for _, tn := range tableNames {
		cs, err := s.getSchemaOfTable(tn, s.databaseName)
		if err != nil {
			return nil, err
		}
		cmap[tn] = cs
	}

	return cmap, nil
}

func (s *pgStore) existTableNames() ([]string, error) {
	tmp := []struct {
		TableName string
	}{}
	err := s.setSearchPath("information_schema")
	if err != nil {
		return nil, err
	}

	err = s.db.Table("tables").
		Select("table_name").
		Where("table_catalog = ? AND table_schema = ? AND table_type = 'BASE TABLE'", s.databaseName, s.schemaName).
		Scan(&tmp).Error
	if err != nil {
		return nil, err
	}

	res := []string{}
	for _, t := range tmp {
		res = append(res, t.TableName)
	}

	return res, nil
}

func (s *pgStore) getSchemaOfTable(tableName, databaseName string) ([]agentConfig.ColumnSchema, error) {
	err := s.setSearchPath("information_schema")
	if err != nil {
		return nil, err
	}
	cs := []agentConfig.ColumnSchema{}
	return cs, s.db.Table("columns").
		Select("column_name, udt_name, is_nullable, character_maximum_length, ordinal_position as order, column_default").
		Where("table_name = ? AND table_catalog = ?", tableName, databaseName).
		Scan(&cs).Error
}

// setSearchPath set search path (schema name in term for postgres)
func (s *pgStore) setSearchPath(schemaName string) error {
	return s.db.Exec("SET search_path TO " + schemaName).Error
}
