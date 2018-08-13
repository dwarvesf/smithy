package sqlmapper

import (
	"database/sql"
	"encoding/json"

	"github.com/dwarvesf/smithy/common/database"
)

// Mapper interface for mapping query from sql to corresponding database engine
type Mapper interface {
	Create(d RowData) ([]byte, error)
	FindAll(request RequestFindAll) ([]byte, error)
	FindByID(id int) ([]byte, error)
	FindByColumnName(request RequestFindBy) ([]byte, error)
}

// Columns return columns listed in RowData
func (r RowData) Columns() []string {
	tmp := []string{}
	for colName := range r {
		tmp = append(tmp, colName)
	}

	return tmp
}

// ColumnsAndData return columns and data of a row_data
func (r RowData) ColumnsAndData() ([]string, []interface{}) {
	cols := []string{}
	data := []interface{}{}
	for k, v := range r {
		cols = append(cols, k)
		data = append(data, v.Data)
	}

	return cols, data
}

// ColData hold data of a column
type ColData struct {
	Name     string
	Data     interface{}
	DataType string
}

// RowData hold data of a row
type RowData map[string]ColData

// MarshalJSON encode json of a row
func (r RowData) MarshalJSON() ([]byte, error) {
	res := map[string]interface{}{}
	for k, v := range r {
		res[k] = v.Data
	}

	return json.Marshal(res)
}

func makeRowDataSet(columns []database.Column) RowData {
	res := map[string]ColData{}
	for _, col := range columns {
		res[col.Name] = ColData{DataType: col.Type, Name: col.Name}
	}

	return res
}

// RowsToQueryResults rows to query results
func RowsToQueryResults(rows *sql.Rows, coldefs []database.Column) (QueryResults, error) {
	cols := database.Columns(coldefs).Names()
	res := []RowData{}
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}
		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		rowData := makeRowDataSet(coldefs)
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			rowData[colName] = ColData{Data: val, DataType: rowData[colName].DataType}
		}

		res = append(res, rowData)
	}

	return res, nil
}

// RowToQueryResult rows to query result
func RowToQueryResult(row *sql.Row, colDefines []database.Column) (QueryResult, error) {
	cols := database.Columns(colDefines).Names()
	columns := make([]interface{}, len(cols))
	columnPointers := make([]interface{}, len(cols))
	for i := range columns {
		columnPointers[i] = &columns[i]
	}
	// Scan the result into the column pointers...
	if err := row.Scan(columnPointers...); err != nil {
		return nil, err
	}

	rowData := makeRowDataSet(colDefines)
	for i, colName := range cols {
		val := columnPointers[i].(*interface{})
		rowData[colName] = ColData{Data: val, DataType: rowData[colName].DataType}
	}

	return QueryResult(rowData), nil
}

// QueryResults hold data of a query more than 1 row
type QueryResults []RowData

// QueryResult hold data of a query have result is 1 row
type QueryResult RowData

// MarshalJSON encode json of a row
func (r QueryResult) MarshalJSON() ([]byte, error) {
	res := map[string]interface{}{}
	for k, v := range r {
		res[k] = v.Data
	}

	return json.Marshal(res)
}

type RequestFindBy struct {
	ColumnName string
	Value      string
	Offset     int `default:"0"`
	Limit      int `default:"0"`
}

type RequestFindAll struct {
	Offset int `default:"0"`
	Limit  int `default:"0"`
}
