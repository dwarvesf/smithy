package drivers

import (
	"errors"
	"fmt"
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

func (s *pgStore) Create(tableName string, d sqlmapper.RowData) (sqlmapper.RowData, error) {
	if err := verifyInput(d, tableName, s.modelMap); err != nil {
		return nil, err
	}

	db := s.db.DB()

	cols, data := d.ColumnsAndData()

	phs := strmangle.Placeholders(true, len(cols), 1, 1)

	execQuery := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING id;",
		tableName,
		strings.Join(cols, ","),
		phs)

	res := db.QueryRow(execQuery, data...)

	var id int
	err := res.Scan(&id)
	if err != nil {
		return nil, err
	}

	// update id if create success
	d["id"] = sqlmapper.ColData{Data: id}

	return d, nil
}

func (s *pgStore) Delete(tableName string, fields, data []interface{}) error {
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

	if _, err := s.db.DB().Exec(exec); err != nil {
		return fmt.Errorf("%v", err.Error())
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

func (s *pgStore) Update(tableName string, d sqlmapper.RowData, id int) (sqlmapper.RowData, error) {
	if notExist, _ := s.isIDNotExist(tableName, id); !notExist {
		return nil, errors.New("primary key is not exist")
	}

	if err := verifyInput(d, tableName, s.modelMap); err != nil {
		return nil, err
	}

	db := s.db.DB()
	cols, data := d.ColumnsAndData()

	rowQuery := make([]string, len(cols))

	for i := 0; i < len(cols); i++ {
		rowQuery[i] = fmt.Sprintf("%s = $%d", cols[i], i+1)
	}

	execQuery := fmt.Sprintf("UPDATE %s SET %s WHERE id = %d",
		tableName,
		strings.Join(rowQuery, ","),
		id)

	if _, err := db.Exec(execQuery, data...); err != nil {
		return nil, err
	}
	return d, nil
}
func (s *pgStore) isIDNotExist(tableName string, id int) (bool, error) {
	data := struct {
		Result bool
	}{}

	execQuery := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id = %d) as result", tableName, id)

	return data.Result, s.db.Raw(execQuery).Scan(&data).Error
}
