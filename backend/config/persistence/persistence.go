package persistence

import (
	"github.com/dwarvesf/smithy/backend/config"
)

type Persistence interface {
	Read(version string) (*config.Config, error)
	Write(cfg *config.Config) error
	ListVersion() []string
}
