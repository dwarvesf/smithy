package hook

import (
	"fmt"
	"log"

	"github.com/mattn/anko/vm"
)

// ScriptEngine interface for running script
type ScriptEngine interface {
	Exec(content string) error
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

func (e *ankoScriptEngine) Exec(content string) error {
	// TODO: implement string processor
	_, err := e.engine.Execute(content)
	return err
}
