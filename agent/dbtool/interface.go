package dbtool

import (
	agentConfig "github.com/dwarvesf/smithy/agent/config"
	"github.com/dwarvesf/smithy/common/database"
)

// DBTool interface for tooling when agent working with db
type DBTool interface {
	MissingColumns(models []database.Model) ([]agentConfig.MissingColumns, error)
	Verify(modelList []database.Model) error
	AutoMigrate([]agentConfig.MissingColumns) error
	CreateUserWithACL(models []database.Model, username string, password string, forceCreate bool) (*database.User, error)
}
