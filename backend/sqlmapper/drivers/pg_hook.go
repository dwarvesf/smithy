package drivers

import (
	"github.com/dwarvesf/smithy/backend/hook"
	"github.com/dwarvesf/smithy/backend/sqlmapper"
	"github.com/dwarvesf/smithy/common/database"
)

type pgHookStore struct {
	pgStore    sqlmapper.Mapper
	hookEngine hook.ScriptEngine
	model      database.Model
}

func (s *pgHookStore) isBeforeCreateEnable() bool {
	return s.model.Hooks.BeforeCreate.Enable
}

func (s *pgHookStore) isAfterCreateEnable() bool {
	return s.model.Hooks.BeforeCreate.Enable
}

func (s *pgHookStore) isBeforeUpdateEnable() bool {
	return s.model.Hooks.BeforeCreate.Enable
}

func (s *pgHookStore) isAfterUpdateEnable() bool {
	return s.model.Hooks.BeforeCreate.Enable
}

func (s *pgHookStore) isBeforeDeleteEnable() bool {
	return s.model.Hooks.BeforeCreate.Enable
}

func (s *pgHookStore) isAfterDeleteEnable() bool {
	return s.model.Hooks.BeforeCreate.Enable
}

// NewPGHookStore new pg implement for hook
func NewPGHookStore(store sqlmapper.Mapper, model database.Model) sqlmapper.Mapper {
	return &pgHookStore{
		pgStore:    store,
		hookEngine: hook.NewAnkoScriptEngine(),
		model:      model,
	}
}

func (s *pgHookStore) Create(d sqlmapper.RowData) (sqlmapper.RowData, error) {
	if s.isBeforeCreateEnable() {
		err := s.hookEngine.Eval(s.model.Hooks.BeforeCreate.Content)
		if err != nil {
			return nil, err
		}
	}

	res, err := s.pgStore.Create(d)
	if err != nil {
		return nil, err
	}

	if s.isAfterCreateEnable() {
		err := s.hookEngine.Eval(s.model.Hooks.AfterCreate.Content)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (s *pgHookStore) Delete(id int) error {
	if s.isBeforeDeleteEnable() {
		err := s.hookEngine.Eval(s.model.Hooks.BeforeDelete.Content)
		if err != nil {
			return err
		}
	}

	res := s.pgStore.Delete(id)

	if s.isAfterDeleteEnable() {
		err := s.hookEngine.Eval(s.model.Hooks.AfterDelete.Content)
		if err != nil {
			return err
		}
	}

	return res
}

func (s *pgHookStore) FindByID(id int) (sqlmapper.RowData, error) {
	return s.pgStore.FindByID(id)
}

func (s *pgHookStore) Update(d sqlmapper.RowData, id int) (sqlmapper.RowData, error) {
	if s.isBeforeUpdateEnable() {
		err := s.hookEngine.Eval(s.model.Hooks.BeforeUpdate.Content)
		if err != nil {
			return nil, err
		}
	}

	res, err := s.pgStore.Update(d, id)
	if err != nil {
		return nil, err
	}

	if s.isAfterUpdateEnable() {
		err := s.hookEngine.Eval(s.model.Hooks.AfterUpdate.Content)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (s *pgHookStore) FindByColumnName(columnName string, value string, offset int, limit int) ([]sqlmapper.RowData, error) {
	return s.pgStore.FindByColumnName(columnName, value, offset, limit)
}

func (s *pgHookStore) FindAll(offset int, limit int) ([]sqlmapper.RowData, error) {
	return s.pgStore.FindAll(offset, limit)
}
