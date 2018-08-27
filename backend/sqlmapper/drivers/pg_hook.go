package drivers

import (
	"github.com/dwarvesf/smithy/backend/hook"
	"github.com/dwarvesf/smithy/backend/sqlmapper"
)

type pgHookStore struct {
	pgStore    sqlmapper.Mapper
	hookEngine hook.ScriptEngine
}

// NewPGHookStore new pg implement for hook
func NewPGHookStore(store sqlmapper.Mapper) sqlmapper.Mapper {
	return &pgHookStore{
		pgStore:    store,
		hookEngine: hook.NewAnkoScriptEngine(),
	}
}

func (s *pgHookStore) Create(d sqlmapper.RowData) (sqlmapper.RowData, error) {
	err := s.hookEngine.Eval(`println("before create hook")`)
	if err != nil {
		return nil, err
	}

	res, err := s.pgStore.Create(d)
	if err != nil {
		return nil, err
	}

	err = s.hookEngine.Eval(`println("after create hook")`)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *pgHookStore) Delete(id int) error {
	err := s.hookEngine.Eval(`println("before delete hook")`)
	if err != nil {
		return err
	}

	res := s.pgStore.Delete(id)

	err = s.hookEngine.Eval(`println("after delete hook")`)
	if err != nil {
		return err
	}

	return res
}

func (s *pgHookStore) FindByID(id int) (sqlmapper.RowData, error) {
	return s.pgStore.FindByID(id)
}

func (s *pgHookStore) Update(d sqlmapper.RowData, id int) (sqlmapper.RowData, error) {
	err := s.hookEngine.Eval(`println("before update hook")`)
	if err != nil {
		return nil, err
	}

	res, err := s.pgStore.Update(d, id)
	if err != nil {
		return nil, err
	}

	err = s.hookEngine.Eval(`println("after update hook")`)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *pgHookStore) FindByColumnName(columnName string, value string, offset int, limit int) ([]sqlmapper.RowData, error) {
	return s.pgStore.FindByColumnName(columnName, value, offset, limit)
}

func (s *pgHookStore) FindAll(offset int, limit int) ([]sqlmapper.RowData, error) {
	return s.pgStore.FindAll(offset, limit)
}
