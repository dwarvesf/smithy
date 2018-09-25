package drivers

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/volatiletech/sqlboiler/strmangle"

	"github.com/dwarvesf/smithy/backend/sqlmapper"
	"github.com/dwarvesf/smithy/common/database"
)

type pgStore struct {
	db       map[string]*gorm.DB
	modelMap map[string]map[string]database.Model
}

// NewPGStore .
func NewPGStore(db map[string]*gorm.DB, modelMap map[string]map[string]database.Model) sqlmapper.Mapper {
	return &pgStore{
		db:       db,
		modelMap: modelMap,
	}
}

func (s *pgStore) addFilter(q sqlmapper.Query, db *gorm.DB) (*gorm.DB, error) {
	if q.Filter.IsZero() {
		return db, nil
	}

	switch q.Filter.Operator {
	case "=":
		return db.Where(q.Filter.ColumnName+" = ?", q.Filter.Value), nil
	default:
		return db, fmt.Errorf("unknown filter operator %s", q.Filter.Operator)
	}
}

func (s *pgStore) addLimitOffset(q sqlmapper.Query, db *gorm.DB) *gorm.DB {
	db = db.Offset(q.Offset)
	if q.Limit > 0 {
		db = db.Limit(q.Limit)
	}

	return db
}

func (s *pgStore) addOrder(q sqlmapper.Query, db *gorm.DB) (*gorm.DB, error) {
	if len(q.Order) == 0 {
		return db, nil
	}

	if len(q.Order) != 2 {
		return db, fmt.Errorf("error require 2 elements: column name and 'asc' if ascending order, 'desc' if descending order")
	}
	return db.Order(q.OrderSequence()), nil
}

func (s *pgStore) Query(q sqlmapper.Query) ([]string, []interface{}, error) {
	db := s.db[q.SourceDatabase].Table(q.SourceTable).
		Select(q.ColumnNames())

	db = s.addLimitOffset(q, db)
	db, err := s.addFilter(q, db)
	if err != nil {
		return nil, nil, err
	}

	db, err = s.addOrder(q, db)
	if err != nil {
		return nil, nil, err
	}

	rows, err := db.Rows()
	if err != nil {
		return nil, nil, err
	}

	defer func() {
		if err != nil {
			return
		}
		rows.Close()
	}()

	if err != nil {
		return nil, nil, err
	}
	data, err := sqlmapper.SQLRowsToRows(rows)

	return q.Columns(), data, err
}

func (s *pgStore) ColumnMetadata(q sqlmapper.Query) ([]database.Column, error) {
	m, ok := s.modelMap[q.SourceDatabase][q.SourceTable]
	if !ok {
		return nil, fmt.Errorf("uknown database_name/table_name %s/%s", q.SourceDatabase, q.SourceTable)
	}

	return q.ColumnMetadata(m.Columns)
}

func (s *pgStore) getRelationshipType(dbName string, tableName string, relateTableName string) (string, error) {
	model, ok := s.modelMap[dbName][tableName]
	if !ok {
		return "", fmt.Errorf("uknown table_name %s", tableName)
	}

	relationships := model.Relationship

	var relationshipType string
	for _, rel := range relationships {
		if rel.Table == relateTableName {
			relationshipType = rel.Type
		}
	}

	return relationshipType, nil
}

func (s *pgStore) getForeignKeyColumn(dbName string, tableName string, relateTableName string) (*database.Column, error) {
	m, ok := s.modelMap[dbName][relateTableName]
	if !ok {
		return nil, fmt.Errorf("uknown database_name/table_name %s/%s", dbName, relateTableName)
	}

	cs := m.Columns
	var c *database.Column
	for i := 0; i < len(cs); i++ {
		if cs[i].ForeignKey.Table == tableName {
			c = &cs[i]
			break
		}
	}

	if c == nil {
		return nil, errors.New("Can't find foreign key column")
	}

	return c, nil
}

func (s *pgStore) Create(dbName string, tableName string, row sqlmapper.RowData) (sqlmapper.RowData, error) {
	d, ok := s.modelMap[dbName]
	if !ok {
		return nil, fmt.Errorf("uknown database_name %s", dbName)
	}
	// clear primary columns
	if _, err := s.getPrimaryKeyMap(row, dbName, tableName); err != nil {
		return nil, err
	}

	if err := verifyInput(row, tableName, d); err != nil {
		return nil, err
	}

	tx, err := s.db[dbName].DB().Begin()
	if err != nil {
		return nil, err
	}

	cols, data := row.ColumnsAndData()

	phs := strmangle.Placeholders(true, len(cols), 1, 1)
	sqlQuery := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING id;",
		tableName,
		strings.Join(cols, ","),
		phs)

	stmt, err := tx.Prepare(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// update id if create success
	var id int
	err = stmt.QueryRow(data...).Scan(&id)
	if err != nil {
		err = tx.Rollback()
		return nil, err
	}

	row["id"] = sqlmapper.ColData{Data: id}
	// create relation data
	relateRowData := row.RelateData()
	for relateTableName, rows := range relateRowData {
		relationship, err := s.getRelationshipType(dbName, tableName, relateTableName)
		if err != nil {
			return nil, err
		}

		switch relationship {
		case "has_many":
			err = s.createWithHasMany(tx, dbName, tableName, id, relateTableName, rows)
			if err != nil {
				err = tx.Rollback()
				return nil, err
			}
		default:
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return row, nil
}

func (s *pgStore) createWithHasMany(tx *sql.Tx, dbName string, parentTableName string, parentID int, tableName string, datas []sqlmapper.RowData) error {
	for _, row := range datas {
		// find relate column
		c, err := s.getForeignKeyColumn(dbName, parentTableName, tableName)
		if c == nil {
			return err
		}

		cols, data := row.ColumnsAndData()
		cols = append(cols, c.Name)
		data = append(data, parentID)

		phs := strmangle.Placeholders(true, len(cols), 1, 1)
		sqlQuery := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING id;",
			tableName,
			strings.Join(cols, ","),
			phs)

		stmt, err := tx.Prepare(sqlQuery)
		if err != nil {
			return err
		}
		defer stmt.Close()

		// update id if create success
		var id int
		err = stmt.QueryRow(data...).Scan(&id)
		if err != nil {
			return err
		}

		row["id"] = sqlmapper.ColData{Data: id}
	}

	return nil
}

func (s *pgStore) Delete(dbName string, tableName string, fields, data []interface{}) error {
	d, ok := s.modelMap[dbName]
	if !ok {
		return fmt.Errorf("uknown database_name %s", dbName)
	}

	if !tableExisted(tableName, d) {
		return fmt.Errorf("Table not exists")
	}

	execPostfix := fmt.Sprintf("DELETE FROM %s WHERE", tableName)

	if len(fields) != len(data) {
		return errors.New("Fields and data isn't match")
	}

	param := []string{}

	numberOfParam := len(fields)
	for i := 0; i < numberOfParam; i++ {
		param = append(param, fmt.Sprintf("%v='%v'", fields[i], data[i]))
	}

	exec := fmt.Sprintf("%s %s", execPostfix, strings.Join(param, " AND "))

	if _, err := s.db[dbName].DB().Exec(exec); err != nil {
		return fmt.Errorf("%v", err)
	}
	return nil
}

func tableExisted(tableName string, modalList map[string]database.Model) bool {
	for _, table := range modalList {
		if table.TableName == tableName {
			return true
		}
	}
	return false
}

func verifyInput(d sqlmapper.RowData, tableName string, modelList map[string]database.Model) error {
	// name data_type nullable primary_key
	if len(d) <= 0 {
		return errors.New("rowData is empty")
	}

	tableNotExist := true
	for _, table := range modelList {
		if table.TableName == tableName {
			// table existed
			tableNotExist = false

			cols, _ := d.ColumnsAndData()

			// create field valid
			for _, column := range table.Columns {
				// file id database will auto generate, so bypass check id
				if column.IsPrimary {
					continue
				}

				if !column.IsNullable {
					if err := checkColumnFieldIsValid(cols, column.Name); err != nil {
						return err
					}
				}
			}
			//field invalid
			for _, name := range cols {
				agentColumns := database.Columns(table.Columns).Names()
				if err := checkColumnFieldIsValid(agentColumns, name); err != nil {
					return err
				}
			}

			break
		}
	}
	if tableNotExist {
		return fmt.Errorf("table %s doesn't exist", tableName)
	}

	return nil
}

func checkColumnFieldIsValid(inputColumns []string, colName string) error {
	err := true
	for _, name := range inputColumns {
		if name == colName {
			err = false
		}
	}

	if err {
		return fmt.Errorf("field %s in valid", colName)
	}

	return nil
}

func (s *pgStore) Update(dbName, tableName string, row sqlmapper.RowData) (sqlmapper.RowData, error) {
	d, ok := s.modelMap[dbName]
	if !ok {
		return nil, fmt.Errorf("uknown database_name %s", dbName)
	}

	primaryKeyMap, err := s.getPrimaryKeyMap(row, dbName, tableName)
	if err != nil {
		return nil, err
	}
	if err := verifyInput(row, tableName, d); err != nil {
		return nil, err
	}
	if exist, _ := s.isPrimaryKeyExist(dbName, tableName, primaryKeyMap); !exist {
		return nil, errors.New("primary key is not exist")
	}
	tx, err := s.db[dbName].DB().Begin()
	if err != nil {
		return nil, err
	}
	if err = s.handleUpdate(tx, row, primaryKeyMap, dbName, tableName); err != nil {
		return nil, err
	}
	relationalRowData := row.RelateData()
	for relateTableName, rows := range relationalRowData {
		relationship, err := s.getRelationshipType(dbName, tableName, relateTableName)
		if err != nil {
			return nil, err
		}

		switch relationship {
		case "has_many":
			err = s.updateWithHasMany(tx, dbName, relateTableName, rows, d)
			if err != nil {
				if errRoolBack := tx.Rollback(); errRoolBack != nil {
					return nil, errRoolBack
				}
				return nil, err
			}
		default:
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return row, nil
}
func (s *pgStore) isPrimaryKeyExist(dbName, tableName string, primaryKeyMap sqlmapper.RowData) (bool, error) {
	data := struct {
		Result bool
	}{}

	params := []string{}
	for colName, colData := range primaryKeyMap {
		params = append(params, fmt.Sprintf("%s = '%v'", colName, colData.Data))
	}
	execQuery := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE %s) as result", tableName, strings.Join(params, " AND "))

	return data.Result, s.db[dbName].Raw(execQuery).Scan(&data).Error
}

func (s *pgStore) isIDNotExist(tableName, colName string, id interface{}) (bool, error) {
	data := struct {
		Result bool
	}{}

	execQuery := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE %s = %v) as result", tableName, colName, id)

	return data.Result, s.db[tableName].Raw(execQuery).Scan(&data).Error
}

func (s *pgStore) handleUpdate(tx *sql.Tx, row, primaryKeyMap sqlmapper.RowData, dbName, tableName string) error {
	cols, data := row.ColumnsAndData()
	foreignColumns, err := s.getRelationalColumn(dbName, tableName)
	if err != nil {
		return err
	}
	if foreignColumns != nil {
		if err := s.isForeignKeyExist(cols, data, foreignColumns); err != nil {
			return err
		}
	}

	if err := s.execUpdateSQL(tx, primaryKeyMap, data, cols, tableName); err != nil {
		return err
	}

	return nil
}

func (s *pgStore) isPrimaryKey(dbName, colName, tableName string) bool {
	columns := s.modelMap[dbName][tableName].Columns
	for _, col := range columns {
		if col.IsPrimary && col.Name == colName {
			return true
		}
	}
	return false
}

func (s *pgStore) getRelationalColumn(dbName, tableName string) ([]database.Column, error) {
	m, ok := s.modelMap[dbName][tableName]
	if !ok {
		return nil, fmt.Errorf("uknown database_name/table_name %s/%s", dbName, tableName)
	}
	cs := m.Columns
	c := []database.Column{}
	for i := 0; i < len(cs); i++ {
		if cs[i].ForeignKey.Table != "" {
			c = append(c, cs[i])
		}
	}
	return c, nil
}

func (s *pgStore) getPrimaryKeyMap(row sqlmapper.RowData, dbName, tableName string) (sqlmapper.RowData, error) {
	if _, ok := s.modelMap[dbName][tableName]; !ok {
		return nil, fmt.Errorf("uknown database_name/table_name %s/%s", dbName, tableName)
	}
	primaryKeyMap := make(sqlmapper.RowData)
	for colName, colData := range row {
		if ok := s.isPrimaryKey(dbName, colName, tableName); ok {
			primaryKeyMap[colName] = colData
			delete(row, colName)
		}
	}
	return primaryKeyMap, nil
}

func (s *pgStore) isForeignKeyExist(cols []string, data []interface{}, foreignColumns []database.Column) error {
	for index, colName := range cols {
		for _, foreignColumn := range foreignColumns {
			if colName == foreignColumn.Name {
				if exist, err := s.isIDNotExist(foreignColumn.ForeignKey.Table, foreignColumn.ForeignKey.ForeignColumn, data[index]); !exist {
					return err
				}
			}
		}
	}
	return nil
}

func (s *pgStore) execUpdateSQL(tx *sql.Tx, primaryKeyMap sqlmapper.RowData, data []interface{}, cols []string, tableName string) error {
	rowQuery := make([]string, len(cols))
	for i := 0; i < len(cols); i++ {
		rowQuery[i] = fmt.Sprintf("%s = $%d", cols[i], i+1)
	}

	params := []string{}
	for colName, colData := range primaryKeyMap {
		params = append(params, fmt.Sprintf("%s = '%v'", colName, colData.Data))
	}

	execQuery := fmt.Sprintf("UPDATE %s SET %s WHERE %s",
		tableName,
		strings.Join(rowQuery, ","),
		strings.Join(params, " AND "))

	stmt, err := tx.Prepare(execQuery)
	if err != nil {
		return err
	}

	if _, err := stmt.Exec(data...); err != nil {
		return err
	}
	defer stmt.Close()

	return nil
}

func (s *pgStore) updateWithHasMany(tx *sql.Tx, dbName, tableName string, rows []sqlmapper.RowData, d map[string]database.Model) error {
	for _, row := range rows {
		primaryKeyMap, err := s.getPrimaryKeyMap(row, dbName, tableName)
		if err != nil {
			return err
		}
		if err := verifyInput(row, tableName, d); err != nil {
			return err
		}
		if exist, _ := s.isPrimaryKeyExist(dbName, tableName, primaryKeyMap); !exist {
			if _, err := s.Create(dbName, tableName, row); err != nil {
				return err
			}
		}
		if err := s.handleUpdate(tx, row, primaryKeyMap, dbName, tableName); err != nil {
			return err
		}
	}
	return nil
}
