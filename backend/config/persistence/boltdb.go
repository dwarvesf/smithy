package persistence

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/dwarvesf/smithy/backend/config"
)

type boltPersistence struct {
	bucket string
	file   string
	db     *bolt.DB
}

func NewBoltPersistence(file string, db *bolt.DB) Persistence {
	return boltPersistence{
		bucket: "ConfigVersion",
		file:   file,
		db:     db,
	}
}

func (b boltPersistence) Read(version string) (*config.Config, error) {
	cfg := &config.Config{}
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(b.bucket))
		v := bucket.Get([]byte(version))
		return json.Unmarshal(v, cfg)
	})

	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (b boltPersistence) Write(cfg *config.Config) error {
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

func (b boltPersistence) ListVersion() []string {
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
