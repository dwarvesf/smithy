package drivers

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/volatiletech/sqlboiler/strmangle"

	"github.com/dwarvesf/smithy/backend/sqlmapper"
	"github.com/dwarvesf/smithy/common/database"
)

type pgStore struct {
	db       *gorm.DB
	modelMap map[string]database.Model
}

// NewPGStore .
func NewPGStore(db *gorm.DB, modelMap map[string]database.Model) sqlmapper.Mapper {
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
	db := s.db.Table(q.SourceTable).
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
	m, ok := s.modelMap[q.SourceTable]
	if !ok {
		return nil, fmt.Errorf("uknown source_table %s", q.SourceTable)
	}

	return q.ColumnMetadata(m.Columns)
}

func (s *pgStore) getRelationshipType(tableName string, relateTableName string) string {
	relationships := s.modelMap[tableName].Relationship
	var relationshipType string
	for _, rel := range relationships {
		if rel.Table == relateTableName {
			relationshipType = rel.Type
		}
	}

	return relationshipType
}

func (s *pgStore) getForeignKeyColumn(tableName string, relateTableName string) (*database.Column, error) {
	cs := s.modelMap[relateTableName].Columns
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

func (s *pgStore) Create(tableName string, row sqlmapper.RowData) (sqlmapper.RowData, error) {
	if err := verifyInput(row, tableName, s.modelMap); err != nil {
		return nil, err
	}

	tx, err := s.db.DB().Begin()
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
		relationship := s.getRelationshipType(tableName, relateTableName)

		switch relationship {
		case "has_many":
			err = s.createWithHasMany(tx, tableName, id, relateTableName, rows)
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

func (s *pgStore) createWithHasMany(tx *sql.Tx, parentTableName string, parentID int, tableName string, datas []sqlmapper.RowData) error {
	for _, row := range datas {
		// find relate column
		c, err := s.getForeignKeyColumn(parentTableName, tableName)
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

func (s *pgStore) Delete(tableName string, id int) error {
	if notExist, _ := s.isIDNotExist(tableName, id); !notExist {
		return errors.New("primary key is not exist")
	}

	exec := fmt.Sprintf("DELETE FROM %s WHERE %s=%v",
		tableName,
		"id",
		id)

	if _, err := s.db.DB().Exec(exec); err != nil {
		return errors.New("delete error")
	}
	return nil
}

func verifyInput(d sqlmapper.RowData, tableName string, modelList map[string]database.Model) error {
	// name data_type nullable primary_key
	d = filterRowData(d)
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
				if column.Name == "id" {
					continue
				}

				if !column.IsNullable || column.IsPrimary {
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

// remove id field out of rowdata because it duplicates with primary key
func filterRowData(d sqlmapper.RowData) sqlmapper.RowData {
	_, ok := d["id"]
	if ok {
		delete(d, "id")
	}
	return d
}

func (s *pgStore) Update(tableName string, row sqlmapper.RowData, id int) (sqlmapper.RowData, error) {
	if notExist, _ := s.isIDNotExist(tableName, id); !notExist {
		return nil, errors.New("primary key is not exist")
	}

	if err := verifyInput(row, tableName, s.modelMap); err != nil {
		return nil, err
	}

	tx, err := s.db.DB().Begin()
	if err != nil {
		return nil, err
	}

	cols, data := row.ColumnsAndData()
	foreignColumns := s.getRelationalColumn(tableName)
	if foreignColumns != nil {
		if err = s.isForeignKeyExist(cols, data, foreignColumns); err != nil {
			return nil, err
		}
	}

	if err = s.execUpdateSQL(tx, id, data, cols, tableName); err != nil {
		return nil, err
	}

	relationalRowData := row.RelateData()
	primaryKeyList := s.ignorePrimaryKey(tableName, relationalRowData)

	for relateTableName, rows := range relationalRowData {
		relationship := s.getRelationshipType(tableName, relateTableName)
		switch relationship {
		case "has_many":
			err = s.updateWithHasMany(tx, tableName, relateTableName, rows, primaryKeyList)
			if err != nil {
				if errRoolBack := tx.Rollback(); errRoolBack != nil {
					return nil, errRoolBack
				}
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
func (s *pgStore) isIDNotExist(tableName string, id int) (bool, error) {
	data := struct {
		Result bool
	}{}

	execQuery := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id = %d) as result", tableName, id)

	return data.Result, s.db.Raw(execQuery).Scan(&data).Error
}

func (s *pgStore) isPrimaryKey(colName, tableName string) bool {
	columns := s.modelMap[tableName].Columns
	for _, col := range columns {
		if col.IsPrimary && col.Name == colName {
			return true
		}
	}
	return false
}

func (s *pgStore) ignorePrimaryKey(tableName string, relateRowData map[string][]sqlmapper.RowData) []int {
	primaryKeyList := []int{}
	for _, rowList := range relateRowData {
		for _, rowData := range rowList {
			for _, colData := range rowData {
				if ok := s.isPrimaryKey(colData.Name, tableName); ok {
					primaryKeyList = append(primaryKeyList, int(colData.Data.(float64)))
					delete(rowData, colData.Name)
				}
			}
		}
	}
	return primaryKeyList
}

func (s *pgStore) getRelationalColumn(tableName string) []database.Column {
	cs := s.modelMap[tableName].Columns
	c := []database.Column{}
	for i := 0; i < len(cs); i++ {
		if cs[i].ForeignKey.Table != "" {
			c = append(c, cs[i])
		}
	}
	return c
}

func (s *pgStore) isForeignKeyExist(cols []string, data []interface{}, foreignColumns []database.Column) error {
	for index, colName := range cols {
		for _, foreignColumn := range foreignColumns {
			if colName == foreignColumn.Name {
				var (
					id  int
					err error
				)
				if id, err = strconv.Atoi(data[index].(string)); err != nil {
					return err
				}
				if notExist, _ := s.isIDNotExist(foreignColumn.ForeignKey.Table, id); !notExist {
					return fmt.Errorf("primary key of %s is not exist", colName)
				}
			}
		}
	}
	return nil
}

func (s *pgStore) execUpdateSQL(tx *sql.Tx, id int, data []interface{}, cols []string, tableName string) error {
	rowQuery := make([]string, len(cols))

	for j := 0; j < len(cols); j++ {
		rowQuery[j] = fmt.Sprintf("%s = $%d", cols[j], j+1)
	}

	execQuery := fmt.Sprintf("UPDATE %s SET %s WHERE id = %d",
		tableName,
		strings.Join(rowQuery, ","),
		id)

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

func (s *pgStore) updateWithHasMany(tx *sql.Tx, parentTableName string, tableName string, rows []sqlmapper.RowData, primaryKeyList []int) error {
	for i, row := range rows {
		// find relate column
		if notExist, _ := s.isIDNotExist(tableName, primaryKeyList[i]); !notExist {
			if _, err := s.Create(tableName, row); err != nil {
				return err
			}
		}

		if err := verifyInput(row, tableName, s.modelMap); err != nil {
			return err
		}

		c, err := s.getForeignKeyColumn(parentTableName, tableName)
		if c == nil {
			return err
		}

		cols, data := row.ColumnsAndData()
		for index, name := range cols {
			if name == c.Name {
				if notExist, _ := s.isIDNotExist(parentTableName, int(data[index].(float64))); !notExist {
					return fmt.Errorf("primary key of %s is not exist", name)
				}
			}
		}

		if err = s.execUpdateSQL(tx, primaryKeyList[i], data, cols, tableName); err != nil {
			return err
		}
	}

	return nil
}
