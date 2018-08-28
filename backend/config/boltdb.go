package config

import (
	"encoding/json"
	"strconv"

	"github.com/boltdb/bolt"
)

type boltImpl struct {
	bucket    string
	versionID int
	db        *bolt.DB
}

// NewBoltPersistent Peristent Bolt
func NewBoltPersistent(db *bolt.DB, versionID int) ReaderWriterQuerier {
	return boltImpl{
		bucket:    "ConfigVersion",
		versionID: versionID,
		db:        db,
	}
}

func (b boltImpl) Read() (*Config, error) {
	var cfg *Config
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(b.bucket))
		if bucket == nil {
			return nil
		}
		v := bucket.Get([]byte(strconv.Itoa(b.versionID)))
		cfg = &Config{}
		return json.Unmarshal(v, cfg)
	})

	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (b boltImpl) Write(cfg *Config) error {
	err := b.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(b.bucket))
		if err != nil {
			return err
		}

		id, _ := bucket.NextSequence()
		cfg.Version.ID = int(id)

		buff, err := json.Marshal(cfg)
		if err != nil {
			return err
		}

		return bucket.Put([]byte(strconv.Itoa(cfg.Version.ID)), buff)
	})
	return err
}

func (b boltImpl) ListVersion() ([]Version, error) {
	versions := make([]Version, 0)
	err := b.db.View(func(tx *bolt.Tx) error {
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

	if err != nil {
		return nil, err
	}

	return versions, nil
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
