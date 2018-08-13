package drivers

import (
	"database/sql"
	"encoding/json"
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
	TableName string
	Columns   []database.Column
	ModelList []database.Model
}

func (s *pgStore) columnNames() string {
	return strings.Join(database.Columns(s.Columns).Names(), ",")
}

// NewPGStore .
func NewPGStore(db *gorm.DB, tableName string, columns []database.Column, modelList []database.Model) sqlmapper.Mapper {
	return &pgStore{db, tableName, columns, modelList}
}

func (s *pgStore) FindAll(offset int, limit int) ([]sqlmapper.RowData, error) {
	db := s.db.Table(s.TableName).
		Select(s.columnNames()).
		Offset(offset)

	var (
		rows *sql.Rows
		err  error
	)
	if limit <= 0 {
		rows, err = db.Rows()
	} else {
		rows, err = db.Limit(limit).Rows()
	}

	defer rows.Close()
	if err != nil {
		return nil, err
	}
	return sqlmapper.RowsToQueryResults(rows, s.Columns)
}

func (s *pgStore) FindByID(id int) (sqlmapper.RowData, error) {
	row := s.db.Table(s.TableName).Select(s.columnNames()).Where("id = ?", id).Row()
	res, err := sqlmapper.RowToQueryResult(row, s.Columns)
	if err != nil {
		return nil, err
	}

	return sqlmapper.RowData(res), nil
}

func (s *pgStore) Create(d sqlmapper.RowData) (sqlmapper.RowData, error) {
	if err := verifyCreate(&d, s.TableName, s.ModelList); err != nil {
		return nil, err
	}

	db := s.db.DB()
	cols, data := d.ColumnsAndData()

	phs := strmangle.Placeholders(true, len(cols), 1, 1)

	execQuery := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING id;",
		s.TableName,
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

func verifyCreate(d *sqlmapper.RowData, tableName string, modelList []database.Model) error {
	// name data_type nullable primary_key
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
					obligateErr := true
					for _, name := range cols {
						if name == column.Name {
							obligateErr = false
						}
					}

					if obligateErr {
						errMess := fmt.Sprintf("%s Obligate", column.Name)
						return errors.New(errMess)
					}
				}
			}

			// field invalid
			for _, name := range cols {
				err := true
				for _, column := range table.Columns {
					if name == column.Name {
						err = false
					}
				}

				if err {
					errMess := fmt.Sprintf("field %s in valid", name)
					return errors.New(errMess)
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

func (s *pgStore) FindByColumnName(columnName string, value string, offset int, limit int) ([]sqlmapper.RowData, error) {
	// TODO: check sql injection
	db := s.db.Table(s.TableName).
		Select(s.columnNames()).
		Where(columnName+" LIKE ?", "%"+value+"%").
		Offset(offset)

	var (
		rows *sql.Rows
		err  error
	)
	if limit <= 0 {
		rows, err = db.Rows()
	} else {
		rows, err = db.Limit(limit).Rows()
	}

	if err != nil {
		return nil, err
	}

	return sqlmapper.RowsToQueryResults(rows, s.Columns)
}

func (s *pgStore) Update(d sqlmapper.RowData, rowID string) ([]byte, error) {
	// TODO: verify column in data-set is correct, check rowData is empty, check primary key is not exist
	db := s.db.DB()
	cols, data := d.ColumnsAndData()

	rowQuery := make([]string, len(cols))

	for i := 0; i < len(cols); i++ {
		rowQuery[i] = fmt.Sprintf("%s = $%d", cols[i], i+1)
	}

	execQuery := fmt.Sprintf("UPDATE %s SET %s WHERE id = %s RETURNING id;",
		s.TableName,
		strings.Join(rowQuery, ","),
		rowID)

	res := db.QueryRow(execQuery, data...)

	var id int
	err := res.Scan(&id)
	if err != nil {
		return nil, err
	}

	// update id if create success
	d["id"] = sqlmapper.ColData{Data: id}

	return json.Marshal(d)
}
