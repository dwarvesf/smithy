package hook

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dwarvesf/smithy/backend/sqlmapper"
	"github.com/dwarvesf/smithy/common/database"
	"github.com/jinzhu/gorm"
	"github.com/volatiletech/sqlboiler/strmangle"
)

type pgLibImpl struct {
	db       map[string]*gorm.DB
	modelMap map[string]map[string]database.Model
}

// NewPGLib dblib implement by postgres
func NewPGLib(db map[string]*gorm.DB, modelMap map[string]map[string]database.Model) DBLib {
	return &pgLibImpl{
		db:       db,
		modelMap: modelMap,
	}
}

func (s *pgLibImpl) First(dbName string, tableName string, condition string) (map[interface{}]interface{}, error) {
	model, ok := s.modelMap[dbName][tableName]
	if !ok {
		return nil, fmt.Errorf("uknown database_name/table_name %s/%s", dbName, tableName)
	}
	cols := database.Columns(model.Columns).Names()
	colNames := strings.Join(cols, ",")
	rows, err := s.db[dbName].Table(tableName).Select(colNames).Where(condition).Limit(1).Rows()
	if err != nil {
		return nil, err
	}

	data, err := sqlmapper.SQLRowsToRows(rows)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, errors.New("record not found")
	}

	first := data[0].([]interface{}) // get only first element
	res := make(map[interface{}]interface{})
	for i := range first {
		res[cols[i]] = first[i]
	}

	return res, nil
}

func (s *pgLibImpl) Where(dbName string, tableName string, condition string) ([]map[interface{}]interface{}, error) {
	model, ok := s.modelMap[dbName][tableName]
	if !ok {
		return nil, fmt.Errorf("uknown database_name/table_name %s/%s", dbName, tableName)
	}
	cols := database.Columns(model.Columns).Names()
	colNames := strings.Join(cols, ",")
	rows, err := s.db[dbName].Table(tableName).Select(colNames).Where(condition).Rows()
	if err != nil {
		return nil, err
	}

	data, err := sqlmapper.SQLRowsToRows(rows)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, nil
	}

	res := []map[interface{}]interface{}{}
	for i := range data {
		tmp := make(map[interface{}]interface{})
		for j := range cols {
			tmp[cols[j]] = data[i].([]interface{})[j] // row is a []interface{}
		}

		res = append(res, tmp)
	}

	return res, nil
}

func (s *pgLibImpl) Create(dbName string, tableName string, d map[interface{}]interface{}) (map[interface{}]interface{}, error) {
	db := s.db[dbName].DB()
	row := toRowData(d)

	cols, data := row.ColumnsAndData()

	phs := strmangle.Placeholders(true, len(cols), 1, 1)

	execQuery := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING id;",
		tableName,
		strings.Join(cols, ","),
		phs)

	res := db.QueryRow(execQuery, data...)

	var id int64
	err := res.Scan(&id)
	if err != nil {
		return nil, err
	}

	// update id if create success
	d["id"] = id

	return d, nil
}

func (s *pgLibImpl) Update(dbName string, tableName string, primaryKey interface{}, d map[interface{}]interface{}) (map[interface{}]interface{}, error) {
	db := s.db[dbName].DB()
	if notExist, _ := s.isIDNotExist(dbName, tableName, primaryKey); !notExist {
		return nil, errors.New("primary key is not exist")
	}

	row := toRowData(d)
	cols, data := row.ColumnsAndData()

	rowQuery := make([]string, len(cols))

	for i := 0; i < len(cols); i++ {
		rowQuery[i] = fmt.Sprintf("%s = $%d", cols[i], i+1)
	}

	execQuery := fmt.Sprintf("UPDATE %s SET %s WHERE id = %d", // FIXME: set primary key could have dynamic name not only id
		tableName,
		strings.Join(rowQuery, ","),
		primaryKey)

	if _, err := db.Exec(execQuery, data...); err != nil {
		return nil, err
	}
	return d, nil
}
func (s *pgLibImpl) Delete(dbName string, tableName string, fields, data []interface{}) error {
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
		return errors.New("delete error")
	}
	return nil
}

func (s *pgLibImpl) isIDNotExist(dbName string, tableName string, primaryKey interface{}) (bool, error) {
	data := struct {
		Result bool
	}{}

	execQuery := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id = %v) as result", tableName, primaryKey)

	return data.Result, s.db[dbName].Raw(execQuery).Scan(&data).Error
}
