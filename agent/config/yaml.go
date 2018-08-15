package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type yamlReaderImpl struct {
	file string
}

// ReadYAML .
func ReadYAML(file string) Reader {
	return yamlReaderImpl{file}
}

func (c yamlReaderImpl) Read() (*Config, error) {
	res := &Config{}
	buf, err := ioutil.ReadFile(c.file)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(buf, res)
	if err != nil {
		return nil, err
	}

	if os.Getenv("ENV") == "test" {
		res.DBUsername = "postgres"
		res.DBName = "test"
		res.DBSSLModeOption = "disable"
		res.DBPassword = "example"
		res.DBHostname = "localhost"
		res.DBPort = "5439"
	}

	return res, nil
}
