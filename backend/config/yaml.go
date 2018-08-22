package config

import (
	"io/ioutil"
	"time"

	"github.com/boltdb/bolt"
	"gopkg.in/yaml.v2"
)

type yamlReaderImpl struct {
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

	if res.PersistenceSupport == "boltdb" {
		if res.PersistenceFileName != "" {
			res.PersistenceDB, err = bolt.Open(res.PersistenceFileName, 0600, &bolt.Options{Timeout: 1 * time.Second})
			if err != nil {
				return nil, err
			}
		}
	}

	return res, nil
}
func (c yamlReaderImpl) ReadToken() (*TokenInfo, error) {
	res := &TokenInfo{}
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
