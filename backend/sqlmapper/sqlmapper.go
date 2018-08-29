package sqlmapper

import (
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/dwarvesf/smithy/common/database"
)

// Mapper interface for mapping query from sql to corresponding database engine
type Mapper interface {
	Create(tableName string, d RowData) (RowData, error)
	Update(tableName string, d RowData, id int) (RowData, error)
	Delete(tableName string, id int) error
	Query(Query) ([]string, []interface{}, error)
}

// Query containt query data for a query request
type Query struct {
	SourceTable string   `json:"source_table"`
	Fields      []string `json:"fields"`
	Filter      Filter   `json:"filter"`
	Offset      int      `json:"offset"`
	Limit       int      `json:"limit"`
}

// ColumnNames return columns name in query
func (q *Query) ColumnNames() string {
	return strings.Join(q.Fields, ", ")
}

// Columns return columsn from query
func (q *Query) Columns() []string {
	return q.Fields
}

// ColumnMetadata .
func (q *Query) ColumnMetadata([]database.Column) []database.Column {
	return nil
}

// Filter containt filter
type Filter struct {
	Operator   string      `json:"operator"`    // "="
	ColumnName string      `json:"column_name"` // TODO: extend filter type
	Value      interface{} `json:"value"`
}

// IsZero check filter is empty
func (f *Filter) IsZero() bool {
	return f.Operator == ""
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
	Name     string      `json:"name"`
	Data     interface{} `json:"data"`
	DataType string      `json:"data_type"`
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

// SQLRowsToRows return rows from sql.Rows
func SQLRowsToRows(rows *sql.Rows, colNum int) ([]interface{}, error) {
	var res []interface{}
	for rows.Next() {
		columns := make([]interface{}, colNum)
		columnPointers := make([]interface{}, colNum)
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}
		tmp := []interface{}{}
		for i := range columnPointers {
			val := columnPointers[i].(*interface{})
			tmp = append(tmp, val)
		}

		res = append(res, tmp)
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
