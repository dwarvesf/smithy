package hook

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	"github.com/mattn/anko/vm"
)

// HookType for hooks
const (
	BeforeCreate = "BeforeCreate"
	AfterCreate  = "AfterCreate"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
)

// HookType for hooks
var (
	HookTypes = []string{
		BeforeCreate,
		AfterCreate,
		BeforeUpdate,
		AfterUpdate,
		BeforeDelete,
		AfterDelete,
	}
)

// IsAHookType check hook type is correct
func IsAHookType(hookType string) bool {
	res := false
	for _, t := range HookTypes {
		if hookType == t {
			res = true
		}
	}

	return res
}

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
	db *gorm.DB
}

// NewPGLib dblib implement by postgres
func NewPGLib(db *gorm.DB) DBLib {
	return &pgLibImpl{db: db}
}

func (s *pgLibImpl) First(tableName string, condition string) (map[interface{}]interface{}, error) {
	return nil, nil
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
func NewAnkoScriptEngine(db *gorm.DB) ScriptEngine {
	env := vm.NewEnv()
	err := env.Define("println", fmt.Println) // TODO: REMOVE THIS LATTER
	if err != nil {
		log.Fatalf("define error: %v\n", err)
	}

	lib := NewPGLib(db)
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
