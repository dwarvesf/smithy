package hook

import (
	"fmt"
	"log"

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
}

// NewAnkoScriptEngine engine for running a engine
func NewAnkoScriptEngine() ScriptEngine {
	env := vm.NewEnv()
	err := env.Define("println", fmt.Println) // TODO: REMOVE THIS LATTER
	if err != nil {
		log.Fatalf("Define error: %v\n", err)
	}
	return &ankoScriptEngine{engine: env}
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
