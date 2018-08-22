package config

import (
	"encoding/json"

	"github.com/boltdb/bolt"
)

type BoltReaderWriterQuerierImpl struct {
	bucket  string
	version string
	db      *bolt.DB
}

func NewBoltReader(version string, db *bolt.DB) Reader {
	return BoltReaderWriterQuerierImpl{
		bucket:  "ConfigVersion",
		version: version,
		db:      db,
	}
}

func NewBoltQuerier(db *bolt.DB) Querier {
	return BoltReaderWriterQuerierImpl{
		bucket: "ConfigVersion",
		db:     db,
	}
}

func NewBoltWriter(db *bolt.DB) Writer {
	return BoltReaderWriterQuerierImpl{
		bucket: "ConfigVersion",
		db:     db,
	}
}

func (b BoltReaderWriterQuerierImpl) Read() (*Config, error) {
	cfg := &Config{}
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(b.bucket))
		v := bucket.Get([]byte(b.version))
		return json.Unmarshal(v, cfg)
	})

	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (b BoltReaderWriterQuerierImpl) Write(cfg *Config) error {
	buff, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	err = b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(b.bucket))
		return bucket.Put([]byte(cfg.Version), buff)
	})
	return err
}

func (b BoltReaderWriterQuerierImpl) ListVersion() []string {
	versions := make([]string, 0)
	b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(b.bucket))

		if bucket == nil {
			return nil
		}

		c := bucket.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			versions = append(versions, string(k))
		}

		return nil
	})

	return versions
}
