package drivers

import (
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
}

func (s *pgStore) columnNames() string {
	return strings.Join(database.Columns(s.Columns).Names(), ",")
}

// NewPGStore .
func NewPGStore(db *gorm.DB, tableName string, columns []database.Column) sqlmapper.Mapper {
	return &pgStore{db, tableName, columns}
}

func (s *pgStore) FindAll() ([]sqlmapper.RowData, error) {
	return s.executeFindAllQuery()
}

func (s *pgStore) executeFindAllQuery() (sqlmapper.QueryResults, error) {
	rows, err := s.db.Table(s.TableName).Select(s.columnNames()).Rows()
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	return sqlmapper.RowsToQueryResults(rows, s.Columns)
}

func (s *pgStore) FindByID(id int) (sqlmapper.RowData, error) {
	return s.executeFindByIDQuery(id)
}

func (s *pgStore) executeFindByIDQuery(id int) (sqlmapper.RowData, error) {
	row := s.db.Table(s.TableName).Select(s.columnNames()).Where("id = ?", id).Row()
	res, err := sqlmapper.RowToQueryResult(row, s.Columns)
	if err != nil {
		return nil, err
	}

	return sqlmapper.RowData(res), nil
}

func (s *pgStore) Create(d sqlmapper.RowData) (sqlmapper.RowData, error) {
	// TODO: verify column in data-set is correct, check rowData is empty, check primary key is not exist
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

func (s *pgStore) executeFindByColumnName(request sqlmapper.RequestFindBy) (sqlmapper.QueryResults, error) {
	// TODO: check sql injection
	rows, err := s.db.Table(s.TableName).
		Select(s.columnNames()).
		Where(request.ColumnName+" = ?", request.Value).
		Offset(request.Offset).
		Limit(request.Limit).
		Rows()
	if err != nil {
		return nil, err
	}

	return sqlmapper.RowsToQueryResults(rows, s.Columns)
}

func (s *pgStore) FindByColumnName(request sqlmapper.RequestFindBy) ([]sqlmapper.RowData, error) {
	return s.executeFindByColumnName(request)
}
