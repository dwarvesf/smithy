package drivers

import (
	"github.com/dwarvesf/smithy/backend/hook"
	"github.com/dwarvesf/smithy/backend/sqlmapper"
	"github.com/dwarvesf/smithy/common/database"
)

type pgHookStore struct {
	pgStore    sqlmapper.Mapper
	hookEngine hook.ScriptEngine
	models     []database.Model
}

// NewPGHookStore new pg implement for hook
func NewPGHookStore(store sqlmapper.Mapper, models []database.Model) sqlmapper.Mapper {
	return &pgHookStore{
		pgStore:    store,
		hookEngine: hook.NewAnkoScriptEngine(),
		models:     models,
	}
}

func (s *pgHookStore) Query(q sqlmapper.Query) ([]interface{}, error) {
	return s.pgStore.Query(q)
}

func (s *pgHookStore) Create(tableName string, d sqlmapper.RowData) (sqlmapper.RowData, error) {
	model := database.Models(s.models).ModelByTableName()[tableName]
	if model.IsBeforeCreateEnable() {
		err := s.hookEngine.Eval(model.Hooks.BeforeCreate.Content)
		if err != nil {
			return nil, err
		}
	}

	res, err := s.pgStore.Create(tableName, d)
	if err != nil {
		return nil, err
	}

	if model.IsAfterCreateEnable() {
		err := s.hookEngine.Eval(model.Hooks.AfterCreate.Content)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (s *pgHookStore) Delete(tableName string, id int) error {
	model := database.Models(s.models).ModelByTableName()[tableName]
	if model.IsBeforeDeleteEnable() {
		err := s.hookEngine.Eval(model.Hooks.BeforeDelete.Content)
		if err != nil {
			return err
		}
	}

	res := s.pgStore.Delete(tableName, id)

	if model.IsAfterDeleteEnable() {
		err := s.hookEngine.Eval(model.Hooks.AfterDelete.Content)
		if err != nil {
			return err
		}
	}

	return res
}

func (s *pgHookStore) FindByID(q sqlmapper.Query) (sqlmapper.RowData, error) {
	return s.pgStore.FindByID(q)
}

func (s *pgHookStore) Update(tableName string, d sqlmapper.RowData, id int) (sqlmapper.RowData, error) {
	model := database.Models(s.models).ModelByTableName()[tableName]
	if model.IsBeforeUpdateEnable() {
		err := s.hookEngine.Eval(model.Hooks.BeforeUpdate.Content)
		if err != nil {
			return nil, err
		}
	}

	res, err := s.pgStore.Update(tableName, d, id)
	if err != nil {
		return nil, err
	}

	if model.IsAfterUpdateEnable() {
		err := s.hookEngine.Eval(model.Hooks.AfterUpdate.Content)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (s *pgHookStore) FindByColumnName(q sqlmapper.Query) ([]sqlmapper.RowData, error) {
	return s.pgStore.FindByColumnName(q)
}

func (s *pgHookStore) FindAll(q sqlmapper.Query) ([]sqlmapper.RowData, error) {
	return s.pgStore.FindAll(q)
}
