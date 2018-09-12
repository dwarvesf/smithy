package hook

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/mattn/anko/vm"

	"github.com/dwarvesf/smithy/backend/sqlmapper"
	"github.com/dwarvesf/smithy/common/database"
)

// ScriptEngine interface for running script
type ScriptEngine interface {
	Eval(ctx map[string]interface{}, content string) error
}

type ankoScriptEngine struct {
	engine *vm.Env
	dblib  DBLib
}

// DBLib interface for lib in db
type DBLib interface {
	First(dbName string, tableName string, condition string) (map[interface{}]interface{}, error)
	Where(dbName string, tableName string, condition string) ([]map[interface{}]interface{}, error)
	Create(dbName string, tableName string, data map[interface{}]interface{}) (map[interface{}]interface{}, error)
	Update(dbName string, tableName string, primaryKey interface{}, data map[interface{}]interface{}) (map[interface{}]interface{}, error)
	Delete(dbName string, tableName string, fields, data []interface{}) error
}

func toRowData(data map[interface{}]interface{}) sqlmapper.RowData {
	res := make(map[string]sqlmapper.ColData)
	for k, v := range data {
		res[k.(string)] = sqlmapper.ColData{Data: v}
	}

	return res
}

func getJSONAPI(headers map[interface{}]interface{}, url string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// add json header
	req.Header.Set("Content-Type", "application/json")
	err = addHeader(headers, req)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	res := map[string]interface{}{}
	err = json.Unmarshal(buf, &res)
	return res, err
}

func postJSONAPI(data map[interface{}]interface{}, headers map[interface{}]interface{}, url string) (map[string]interface{}, error) {
	buf, err := convertAnkoMapToJSON(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}

	// add json header
	req.Header.Set("Content-Type", "application/json")
	err = addHeader(headers, req)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	res := map[string]interface{}{}
	err = json.Unmarshal(buf, &res)
	return res, err
}

func addHeader(headers map[interface{}]interface{}, req *http.Request) error {
	for k, v := range headers {
		header, ok := k.(string)
		if !ok {
			return fmt.Errorf("format of header was wrong")
		}
		value, ok := v.(string)
		if !ok {
			return fmt.Errorf("format of header key was wrong")
		}

		req.Header.Set(header, value)
	}

	return nil
}

func convertAnkoMapToJSON(m map[interface{}]interface{}) ([]byte, error) {
	tmp := map[string]interface{}{}
	for k, v := range m {
		var key string
		switch val := k.(type) {
		case int:
			key = strconv.Itoa(val)
		case string:
			key = val
		default:
			return nil, errors.New("unknown json format")
		}

		tmp[key] = v
	}

	return json.Marshal(tmp)
}

func defineAPICallLib(env *vm.Env) error {
	err := env.Define("json_post", postJSONAPI)
	if err != nil {
		return err
	}

	return env.Define("json_get", getJSONAPI)
}

func defineAnkoDBLib(env *vm.Env, lib DBLib) error {
	err := env.Define("db_first", lib.First)
	if err != nil {
		return err
	}
	err = env.Define("db_where", lib.Where)
	if err != nil {
		return err
	}
	err = env.Define("db_create", lib.Create)
	if err != nil {
		return err
	}
	err = env.Define("db_update", lib.Update)
	if err != nil {
		return err
	}

	return env.Define("db_delete", lib.Delete)
}

// NewAnkoScriptEngine engine for running a engine
func NewAnkoScriptEngine(db map[string]*gorm.DB, modelMap map[string]map[string]database.Model) (ScriptEngine, error) {
	env := vm.NewEnv()
	err := env.Define("println", fmt.Println) // TODO: REMOVE THIS LATTER
	if err != nil {
		return nil, fmt.Errorf("define error: %v", err)
	}

	lib := NewPGLib(db, modelMap)
	err = defineAnkoDBLib(env, lib)
	if err != nil {
		return nil, fmt.Errorf("define error: %v", err)
	}

	err = defineAPICallLib(env)
	if err != nil {
		return nil, fmt.Errorf("define error: %v", err)
	}

	return &ankoScriptEngine{
		engine: env,
		dblib:  lib,
	}, nil
}

type libCtx struct {
	data map[string]interface{}
}

func (l *libCtx) Ctx() map[string]interface{} {
	return l.data
}

func (e *ankoScriptEngine) Eval(ctx map[string]interface{}, content string) error {
	l := libCtx{ctx}
	env := e.engine.NewEnv()
	err := env.Define("ctx", l.Ctx)
	if err != nil {
		return err
	}

	// TODO: implement string processor
	_, err = env.Execute(content)
	return err
}
