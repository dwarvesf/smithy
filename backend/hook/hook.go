package hook

import (
	"fmt"
	"log"
	"strings"

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
	All(tableName string, condition string) ([]map[interface{}]interface{}, error)
	Create(tableName string, data map[interface{}]interface{}) (map[interface{}]interface{}, error)
	Update(tableName string, primaryKey interface{}, data map[interface{}]interface{}) (map[interface{}]interface{}, error)
	Delete(tableName string, primaryKey interface{}) error
}

type pgLibImpl struct {
	db       *gorm.DB
	modelMap map[string]database.Model
}

// NewPGLib dblib implement by postgres
func NewPGLib(db *gorm.DB, modelMap map[string]database.Model) DBLib {
	return &pgLibImpl{
		db:       db,
		modelMap: modelMap,
	}
}

func (s *pgLibImpl) First(tableName string, condition string) (map[interface{}]interface{}, error) {
	model, ok := s.modelMap[tableName]
	if !ok {
		return nil, fmt.Errorf("uknown table_name %s", tableName)
	}
	cols := database.Columns(model.Columns).Names()
	colNames := strings.Join(cols, ",")
	rows, err := s.db.Table(tableName).Select(colNames).Where(condition).Limit(1).Rows()
	if err != nil {
		return nil, err
	}

	data, err := sqlmapper.SQLRowsToRows(rows, len(cols))
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, nil
	}

	first := data[0].([]interface{})
	res := make(map[interface{}]interface{})
	for i := range first {
		res[cols[i]] = first[i]
	}

	return res, nil
}
func (s *pgLibImpl) All(tableName string, condition string) ([]map[interface{}]interface{}, error) {
	return nil, nil
}
func (s *pgLibImpl) Create(tableName string, data map[interface{}]interface{}) (map[interface{}]interface{}, error) {
	return nil, nil
}
func (s *pgLibImpl) Update(tableName string, primaryKey interface{}, data map[interface{}]interface{}) (map[interface{}]interface{}, error) {
	return nil, nil
}
func (s *pgLibImpl) Delete(tableName string, primaryKey interface{}) error {
	return nil
}

func defineAnkoDBLib(env *vm.Env, lib DBLib) error {
	err := env.Define("db_first", lib.First)
	if err != nil {
		return err
	}
	err = env.Define("db_all", lib.All)
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
func NewAnkoScriptEngine(db *gorm.DB, modelMap map[string]database.Model) ScriptEngine {
	env := vm.NewEnv()
	err := env.Define("println", fmt.Println) // TODO: REMOVE THIS LATTER
	if err != nil {
		log.Fatalf("define error: %v\n", err)
	}

	lib := NewPGLib(db, modelMap)
	defineAnkoDBLib(env, lib)

	return &ankoScriptEngine{
		engine: env,
		dblib:  lib,
	}
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
