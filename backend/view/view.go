package view

import (
	"errors"
	"strings"
	"time"

	"github.com/dwarvesf/smithy/backend/sqlmapper"
)

// View .
type View struct {
	ID           int       `json:"id"`
	SQL          string    `json:"sql"`
	DatabaseName string    `json:"database_name"`
	CreatedAt    time.Time `json:"created_at"`
}

// Validate validate a view
func (s *View) Validate(m sqlmapper.Mapper) (interface{}, error) {
	// just get first view
	SQLs := strings.Split(s.SQL, ";")
	s.SQL = strings.TrimSpace(SQLs[0])

	if strings.Index(s.SQL, "SELECT") != 0 {
		return nil, errors.New(`view must be begining with "SELECT"`)
	}

	q, err := m.Explain(s.DatabaseName, s.SQL)
	if err != nil {
		return nil, err
	}
	return q, nil
}

// Reader .
type Reader interface {
	ListCommands() ([]*View, error)
	ListCommandsByDBName(databaseName string) ([]*View, error)
	Read(sqlID int) (*View, error)
}

// Writer .
type Writer interface {
	Write(sql *View) error
}

// Deleter .
type Deleter interface {
	Delete(sqlID int) error
}

// WriterReaderDeleter .
type WriterReaderDeleter interface {
	Reader
	Writer
	Deleter
}
