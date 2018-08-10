package pg

import (
	"strings"

	"github.com/jinzhu/gorm"

	"github.com/dwarvesf/smithy/backend/sqlmapper"
	"github.com/dwarvesf/smithy/common/database"
)

type querier struct {
	db        *gorm.DB
	TableName string
	Columns   []database.Column
}

func (q *querier) columnNames() string {
	return strings.Join(database.Columns(q.Columns).Names(), ",")
}

// NewQuerier .
func NewQuerier(db *gorm.DB, tableName string, columns []database.Column) sqlmapper.Mapper {
	return &querier{db, tableName, columns}
}
