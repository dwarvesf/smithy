package backend

import (
	backendConfig "github.com/dwarvesf/smithy/backend/config"
)

// NewConfig check dashboard config is correct
func NewConfig(r backendConfig.Reader) (*backendConfig.Config, error) {
	cfg, err := r.Read()
	if err != nil {
		return nil, err
	}
	err = checkConfig(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// checkConfig check agent config is correct
func checkConfig(c *backendConfig.Config) error {
	// TODO: implement dashboard config checking
	return nil
}
