package pg

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dwarvesf/smithy/backend/sqlmapper"
)

func (s *querier) FindAll() ([]byte, error) {
	rs, err := s.executeFindAllQuery()
	if err != nil {
		return nil, err
	}

	buf, err := json.Marshal(rs)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (s *querier) executeFindAllQuery() (sqlmapper.QueryResults, error) {
	rows, err := s.db.Table(s.TableName).Select(s.columnNames()).Rows()
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	return sqlmapper.RowsToQueryResults(rows, s.Columns)
}

func (s *querier) FindByID(id int) ([]byte, error) {
	rs, err := s.executeFindByIDQuery(id)
	if err != nil {
		return nil, err
	}

	buf, err := json.Marshal(rs)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (s *querier) executeFindByIDQuery(id int) (sqlmapper.QueryResult, error) {
	row := s.db.Table(s.TableName).Select(s.columnNames()).Where("id = ?", id).Row()
	return sqlmapper.RowToQueryResult(row, s.Columns)
}

func (s *querier) Create(d sqlmapper.RowData) ([]byte, error) {
	// TODO: verify column in data-set is correct, check rowData is empty, check primary key is not exist
	db := s.db.DB()

	qms := []string{}
	for i := 0; i < len(d.Data()); i++ {
		qms = append(qms, fmt.Sprintf("$%d", i+1))
	}

	execQuery := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING id;",
		s.TableName,
		d.ColumnsString(),
		strings.Join(qms, ","))

	data := d.Data()
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
