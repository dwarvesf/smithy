package config

import (
	"github.com/dwarvesf/smithy/backend/config/dashboard"
)

// NewDashboardConfig check dashboard config is correct
func NewDashboardConfig(r dashboard.ConfigReader) (*dashboard.Config, error) {
	cfg, err := r.Read()
	if err != nil {
		return nil, err
	}
	err = checkDashBoardConfig(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// checkDashBoardConfig check agent config is correct
func checkDashBoardConfig(c *dashboard.Config) error {
	// TODO: implement dashboard config checking
	return nil
}
