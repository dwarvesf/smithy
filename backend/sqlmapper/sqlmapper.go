package sqlmapper

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dwarvesf/smithy/common/database"
)

// Mapper interface for mapping query from sql to corresponding database engine
type Mapper interface {
	Create(tableName string, d RowData) (RowData, error)
	Update(tableName string, d RowData, id int) (RowData, error)
	Delete(tableName string, id int) error
	Query(Query) ([]string, []interface{}, error)
	ColumnMetadata(Query) ([]database.Column, error)
}

// Query contain query data for a query request
type Query struct {
	SourceTable string   `json:"-"`
	Fields      []string `json:"fields"`
	Filter      Filter   `json:"filter"`
	Offset      int      `json:"offset"`
	Limit       int      `json:"limit"`
	Order       []string `json:"order"` // 2 elements: "columnName" and "asc" if ascending order, "desc" if descending order
}

// ColumnNames return columns name in query
func (q *Query) ColumnNames() string {
	return strings.Join(q.Fields, ", ")
}

// OrderSequence return order sequence of columns in query
func (q *Query) OrderSequence() string {
	return strings.Join(q.Order, " ")
}

// Columns return columsn from query
func (q *Query) Columns() []string {
	return q.Fields
}

// ColumnMetadata convert query to column spec
func (q *Query) ColumnMetadata(columns []database.Column) ([]database.Column, error) {
	res := []database.Column{}
	colMap := database.Columns(columns).GroupByName()
	for _, field := range q.Fields {
		cols, ok := colMap[field]
		if !ok {
			return nil, fmt.Errorf("unknown field %s ", field)
		}

		res = append(res, cols[0]) // expect all cols is a same column, if dupplicate happened
	}

	return res, nil
}

// Filter containt filter
type Filter struct {
	Operator   string      `json:"operator"` // "="
	ColumnName string      `json:"column_name"`
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

// Ctx for current sqlmapper data context
type Ctx map[string]interface{}

// ToCtx return map[string]interface{} from a RowData
func (r RowData) ToCtx() Ctx {
	res := make(map[string]interface{})
	for k, v := range r {
		res[k] = v.Data
	}

	return res
}

// ToRowData convert Ctx back to RowData
func (c Ctx) ToRowData() RowData {
	res := make(map[string]ColData)
	for k, v := range c {
		res[k] = ColData{Data: v}
	}

	return res
}

// ColData hold data of a column
type ColData struct {
	Name     string      `json:"name"`
	Data     interface{} `json:"data"`
	DataType string      `json:"data_type"`
}

// MakeRowData make row_data from fields([]string) and data([]{}interface)
func MakeRowData(fields []string, data []interface{}) (RowData, error) {
	if len(fields) != len(data) {
		return nil, errors.New("length of fields and data is not the same")
	}
	res := make(map[string]ColData)

	for i := range fields {
		res[fields[i]] = ColData{Data: data[i], Name: fields[i]}
	}

	return res, nil
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
func SQLRowsToRows(rows *sql.Rows) ([]interface{}, error) {
	var res []interface{}
	columns, _ := rows.Columns()
	for rows.Next() {
		row := make([]interface{}, len(columns))
		for idx := range columns {
			row[idx] = new(metalScanner)
		}

		err := rows.Scan(row...)
		if err != nil {
			fmt.Println(err)
		}

		tmp := []interface{}{}
		for idx := range columns {
			var scanner = row[idx].(*metalScanner)
			tmp = append(tmp, scanner.value)
		}

		res = append(res, tmp)
	}

	return res, nil
}

type metalScanner struct {
	valid bool
	value interface{}
}

func (scanner *metalScanner) getBytes(src interface{}) []byte {
	if a, ok := src.([]uint8); ok {
		return a
	}
	return nil
}

func (scanner *metalScanner) Scan(src interface{}) error {
	switch src.(type) {
	case int64:
		if value, ok := src.(int64); ok {
			scanner.value = value
			scanner.valid = true
		}
	case float64:
		if value, ok := src.(float64); ok {
			scanner.value = value
			scanner.valid = true
		}
	case bool:
		if value, ok := src.(bool); ok {
			scanner.value = value
			scanner.valid = true
		}
	case string:
		scanner.value = src
		scanner.valid = true
	case []byte:
		value := scanner.getBytes(src)
		scanner.value = value
		scanner.valid = true
	case time.Time:
		if value, ok := src.(time.Time); ok {
			scanner.value = value
			scanner.valid = true
		}
	case nil:
		scanner.value = nil
		scanner.valid = true
	}
	return nil
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
