package hook

import (
	"fmt"

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
	First(tableName string, condition string) (map[interface{}]interface{}, error)
	Where(tableName string, condition string) ([]map[interface{}]interface{}, error)
	Create(tableName string, data map[interface{}]interface{}) (map[interface{}]interface{}, error)
	Update(tableName string, primaryKey interface{}, data map[interface{}]interface{}) (map[interface{}]interface{}, error)
	Delete(tableName string, primaryKey interface{}) error
}

func toRowData(data map[interface{}]interface{}) sqlmapper.RowData {
	res := make(map[string]sqlmapper.ColData)
	for k, v := range data {
		res[k.(string)] = sqlmapper.ColData{Data: v}
	}

	return res
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
func NewAnkoScriptEngine(db *gorm.DB, modelMap map[string]database.Model) (ScriptEngine, error) {
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
