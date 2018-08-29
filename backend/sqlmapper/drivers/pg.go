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
	db        *gorm.DB
	ModelList []database.Model
}

// NewPGStore .
func NewPGStore(db *gorm.DB, modelList []database.Model) sqlmapper.Mapper {
	return &pgStore{db, modelList}
}

func (s *pgStore) Query(q sqlmapper.Query) ([]string, []interface{}, error) {
	db := s.db.Table(q.SourceTable).
		Select(q.ColumnNames()).
		Offset(q.Offset)

	var (
		rows *sql.Rows
		err  error
	)
	if q.Limit <= 0 {
		rows, err = db.Rows()
	} else {
		rows, err = db.Limit(q.Limit).Rows()
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
	data, err := sqlmapper.SQLRowsToRows(rows, len(q.Columns()))

	return q.Columns(), data, err
}

// func (s *pgStore) FindAll(q sqlmapper.Query) ([]sqlmapper.RowData, error) {
// 	db := s.db.Table(q.SourceTable).
// 		Select(q.ColumnNames()).
// 		Offset(q.Offset)

// 	var (
// 		rows *sql.Rows
// 		err  error
// 	)
// 	if q.Limit <= 0 {
// 		rows, err = db.Rows()
// 	} else {
// 		rows, err = db.Limit(q.Limit).Rows()
// 	}

// 	defer func() {
// 		if err != nil {
// 			return
// 		}
// 		rows.Close()
// 	}()

// 	if err != nil {
// 		return nil, err
// 	}

// 	return sqlmapper.RowsToQueryResults(rows, q.Fields)
// }

// func (s *pgStore) FindByID(q sqlmapper.Query) (sqlmapper.RowData, error) {
// 	row := s.db.Table(q.SourceTable).Select(q.ColumnNames()).Where("id = ?", q.Filter.Value).Row()
// 	res, err := sqlmapper.RowToQueryResult(row, q.Fields)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return sqlmapper.RowData(res), nil
// }

func (s *pgStore) Create(tableName string, d sqlmapper.RowData) (sqlmapper.RowData, error) {
	if err := verifyInput(d, tableName, s.ModelList); err != nil {
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

func verifyInput(d sqlmapper.RowData, tableName string, modelList []database.Model) error {
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
		return errors.New("Table have existed yet")
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

// func (s *pgStore) FindByColumnName(q sqlmapper.Query) ([]sqlmapper.RowData, error) {
// 	// TODO: check sql injection
// 	db := s.db.Table(q.SourceTable).
// 		Select(q.ColumnNames()).
// 		Where(q.Filter.ColName+" LIKE ?", "%"+fmt.Sprintf("%v", q.Filter.Value)+"%").
// 		Offset(q.Offset)

// 	var (
// 		rows *sql.Rows
// 		err  error
// 	)
// 	if q.Limit <= 0 {
// 		rows, err = db.Rows()
// 	} else {
// 		rows, err = db.Limit(q.Limit).Rows()
// 	}

// 	defer func() {
// 		if err != nil {
// 			return
// 		}
// 		rows.Close()
// 	}()

// 	if err != nil {
// 		return nil, err
// 	}

// 	return sqlmapper.RowsToQueryResults(rows, q.Fields)
// }

func (s *pgStore) Update(tableName string, d sqlmapper.RowData, id int) (sqlmapper.RowData, error) {
	if notExist, _ := s.isIDNotExist(tableName, id); !notExist {
		return nil, errors.New("primary key is not exist")
	}

	if err := verifyInput(d, tableName, s.ModelList); err != nil {
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
