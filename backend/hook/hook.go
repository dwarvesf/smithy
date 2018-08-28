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
	Eval(content string) error
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

func (e *ankoScriptEngine) Eval(content string) error {
	// TODO: implement string processor
	_, err := e.engine.Execute(content)
	return err
}
