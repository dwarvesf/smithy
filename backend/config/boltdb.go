package config

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/boltdb/bolt"

	"github.com/dwarvesf/smithy/common/database"
)

type boltImpl struct {
	bucket              string
	versionID           int
	db                  *bolt.DB
	persistenceFileName string
}

// NewBoltPersistent Peristent Bolt
func NewBoltPersistent(persistenceFileName string, versionID int) ReaderWriterQuerier {
	return boltImpl{
		bucket:              "ConfigVersion",
		versionID:           versionID,
		persistenceFileName: persistenceFileName,
	}
}

func (b boltImpl) openConnection() (*bolt.DB, error) {
	db, err := bolt.Open(b.persistenceFileName, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (b boltImpl) Read() (*Config, error) {
	var err error
	b.db, err = b.openConnection()
	if err != nil {
		return nil, err
	}
	defer b.db.Close()

	var cfg *Config
	err = b.db.View(func(tx *bolt.Tx) error {
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

	// init model map for prevent access nil map
	cfg.ModelMap = make(map[string]database.Model)

	return cfg, nil
}

func (b boltImpl) Write(cfg *Config) error {
	var err error
	b.db, err = b.openConnection()
	if err != nil {
		return err
	}
	defer b.db.Close()

	err = b.db.Update(func(tx *bolt.Tx) error {
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
	var err error
	b.db, err = b.openConnection()
	if err != nil {
		return nil, err
	}
	defer b.db.Close()

	versions := make([]Version, 0)
	err = b.db.View(func(tx *bolt.Tx) error {
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
	var err error
	b.db, err = b.openConnection()
	if err != nil {
		return nil, err
	}
	defer b.db.Close()

	cfg := &Config{}
	err = b.db.View(func(tx *bolt.Tx) error {
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
