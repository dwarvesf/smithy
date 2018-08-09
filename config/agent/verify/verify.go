package verify

import (
	agentConfig "github.com/dwarvesf/smithy/config/agent"
	"github.com/dwarvesf/smithy/config/database"
)

// Verifier interface for verify agent model_list config
type Verifier interface {
	MissingColumns(models []database.Model) ([]agentConfig.MissingColumns, error)
	Verify(modelList []database.Model) error
}
