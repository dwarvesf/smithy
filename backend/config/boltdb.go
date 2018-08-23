package config

import (
	"encoding/json"
	"strconv"

	"github.com/boltdb/bolt"
)

type boltImpl struct {
	bucket  string
	version int
	db      *bolt.DB
}

func NewBoltIO(db *bolt.DB, version int) ReaderWriterQuerier {
	return boltImpl{
		bucket:  "ConfigVersion",
		version: version,
		db:      db,
	}
}

func (b boltImpl) Read() (*Config, error) {
	cfg := &Config{}
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(b.bucket))
		v := bucket.Get([]byte(strconv.Itoa(b.version)))
		return json.Unmarshal(v, cfg)
	})

	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (b boltImpl) Write(cfg *Config) error {
	buff, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	err = b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(b.bucket))
		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte(b.bucket))
			if err != nil {
				return err
			}
		}
		return bucket.Put([]byte(strconv.Itoa(cfg.Version.VersionNumber)), buff)
	})
	return err
}

func (b boltImpl) ListVersion() []Version {
	versions := make([]Version, 0)
	b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(b.bucket))

		if bucket == nil {
			return nil
		}

		c := bucket.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			cfg := &Config{}
			err := json.Unmarshal(v, cfg)
			if err != nil {
				return err
			}

			versions = append(versions, cfg.Version)
		}

		return nil
	})

	return versions
}

func (b boltImpl) LastestVersion() (*Config, error) {
	cfg := &Config{}
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(b.bucket))

		if bucket == nil {
			cfg = nil
			return nil
		}

		c := bucket.Cursor()
		_, v := c.Last()

		return json.Unmarshal(v, cfg)
	})

	return cfg, err
}
