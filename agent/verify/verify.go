package verify

import (
	agentConfig "github.com/dwarvesf/smithy/agent/config"
	"github.com/dwarvesf/smithy/common/database"
)

// Verifier interface for verify agent model_list config
type Verifier interface {
	MissingColumns(models []database.Model) ([]agentConfig.MissingColumns, error)
	Verify(modelList []database.Model) error
}
