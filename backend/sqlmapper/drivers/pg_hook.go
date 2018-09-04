package drivers

import (
	"github.com/jinzhu/gorm"

	"github.com/dwarvesf/smithy/backend/hook"
	"github.com/dwarvesf/smithy/backend/sqlmapper"
	"github.com/dwarvesf/smithy/common/database"
)

type pgHookStore struct {
	pgStore    sqlmapper.Mapper
	hookEngine hook.ScriptEngine
	modelMap   map[string]database.Model
}

// NewPGHookStore new pg implement for hook
func NewPGHookStore(store sqlmapper.Mapper, modelMap map[string]database.Model, db *gorm.DB) sqlmapper.Mapper {
	return &pgHookStore{
		pgStore:    store,
		hookEngine: hook.NewAnkoScriptEngine(db, modelMap),
		modelMap:   modelMap,
	}
}

func (s *pgHookStore) Query(q sqlmapper.Query) ([]string, []interface{}, error) {
	return s.pgStore.Query(q)
}

func (s *pgHookStore) ColumnMetadata(q sqlmapper.Query) ([]database.Column, error) {
	return s.pgStore.ColumnMetadata(q)
}

func (s *pgHookStore) Create(tableName string, d sqlmapper.RowData) (sqlmapper.RowData, error) {
	ctx := d.ToCtx()

	model := s.modelMap[tableName]
	if model.IsBeforeCreateEnable() {
		err := s.hookEngine.Eval(ctx, model.Hooks.BeforeCreate.Content)
		if err != nil {
			return nil, err
		}
		d = sqlmapper.Ctx(ctx).ToRowData()
	}

	res, err := s.pgStore.Create(tableName, d)
	if err != nil {
		return nil, err
	}

	if model.IsAfterCreateEnable() {
		err := s.hookEngine.Eval(ctx, model.Hooks.AfterCreate.Content)
		if err != nil {
			return nil, err
		}

		d = sqlmapper.Ctx(ctx).ToRowData()
	}

	return res, nil
}

func (s *pgHookStore) Delete(tableName string, id int) error {
	model := s.modelMap[tableName]
	if model.IsBeforeDeleteEnable() {
		err := s.hookEngine.Eval(nil, model.Hooks.BeforeDelete.Content)
		if err != nil {
			return err
		}
	}

	res := s.pgStore.Delete(tableName, id)

	if model.IsAfterDeleteEnable() {
		err := s.hookEngine.Eval(nil, model.Hooks.AfterDelete.Content)
		if err != nil {
			return err
		}
	}

	return res
}

func (s *pgHookStore) Update(tableName string, d sqlmapper.RowData, id int) (sqlmapper.RowData, error) {
	model := s.modelMap[tableName]
	if model.IsBeforeUpdateEnable() {
		err := s.hookEngine.Eval(nil, model.Hooks.BeforeUpdate.Content)
		if err != nil {
			return nil, err
		}
	}

	res, err := s.pgStore.Update(tableName, d, id)
	if err != nil {
		return nil, err
	}

	if model.IsAfterUpdateEnable() {
		err := s.hookEngine.Eval(nil, model.Hooks.AfterUpdate.Content)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}
