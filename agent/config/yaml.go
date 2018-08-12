package config

import (
	"io/ioutil"

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

	return res, nil
}
