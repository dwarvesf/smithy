package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/dwarvesf/smithy/common/database"
)

type yamlReaderImpl struct {
	file string
}

type yamlWriterImpl struct {
	file string
}

// ReadYAML reader dashboard config from front-end
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

	// init model map for prevent access nil map
	res.ModelMap = make(map[string]map[string]database.Model)

	return res, nil
}

// WriteYAML to write data to .yaml file
func WriteYAML(file string) Writer {
	return yamlWriterImpl{file: file}
}

func (c yamlWriterImpl) Write(res *Config) error {
	buff, err := yaml.Marshal(&res)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(c.file, buff, 0644)
	if err != nil {
		return err
	}

	return nil
}
