package drivers

import "github.com/dwarvesf/smithy/backend/sqlmapper"

type pgHookStore struct {
	BeforeHook func() error
	pgStore    sqlmapper.Mapper
	AfterHook  func() error
}

// NewPGHookStore new pg implement for hook
func NewPGHookStore(store sqlmapper.Mapper) sqlmapper.Mapper {
	return &pgHookStore{
		BeforeHook: func() error { return nil },
		AfterHook:  func() error { return nil },
		pgStore:    store,
	}
}

func (s *pgHookStore) Create(d sqlmapper.RowData) (sqlmapper.RowData, error) {
	err := s.BeforeHook()
	if err != nil {
		return nil, err
	}

	res, err := s.pgStore.Create(d)
	if err != nil {
		return nil, err
	}

	err = s.AfterHook()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *pgHookStore) Delete(id int) error {
	err := s.BeforeHook()
	if err != nil {
		return err
	}

	res := s.pgStore.Delete(id)

	err = s.AfterHook()
	if err != nil {
		return err
	}

	return res
}

func (s *pgHookStore) FindByID(id int) (sqlmapper.RowData, error) {
	return s.pgStore.FindByID(id)
}

func (s *pgHookStore) Update(d sqlmapper.RowData, id int) (sqlmapper.RowData, error) {
	err := s.BeforeHook()
	if err != nil {
		return nil, err
	}

	res, err := s.pgStore.Update(d, id)
	if err != nil {
		return nil, err
	}

	err = s.AfterHook()
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
