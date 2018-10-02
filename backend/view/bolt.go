package view

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
)

type boltImpl struct {
	bucket              string
	db                  *bolt.DB
	persistenceFileName string
}

// NewBoltWriteReadDeleter Peristent Bolt
func NewBoltWriteReadDeleter(persistenceFileName string) WriteReadDeleter {
	return boltImpl{
		bucket:              "View",
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

func (b boltImpl) Read(sqlID int) (*View, error) {
	var err error
	b.db, err = b.openConnection()
	if err != nil {
		return nil, err
	}
	defer b.db.Close()

	var view *View
	err = b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(b.bucket))
		if bucket == nil {
			return nil
		}
		v := bucket.Get([]byte(strconv.Itoa(sqlID)))
		view = &View{}
		return json.Unmarshal(v, view)
	})

	if err != nil {
		return nil, err
	}

	return view, nil
}

func (b boltImpl) Write(sql *View) error {
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
		sql.ID = int(id)

		buff, err := json.Marshal(sql)
		if err != nil {
			return err
		}

		return bucket.Put([]byte(strconv.Itoa(sql.ID)), buff)
	})
	return err
}

func (b boltImpl) Delete(sqlID int) error {
	var err error
	b.db, err = b.openConnection()
	if err != nil {
		return err
	}
	defer b.db.Close()

	err = b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(b.bucket))
		if bucket == nil {
			return nil
		}

		return bucket.Delete([]byte(strconv.Itoa(sqlID)))
	})

	return err
}

func (b boltImpl) ListCommands() ([]*View, error) {
	var err error
	b.db, err = b.openConnection()
	if err != nil {
		return nil, err
	}
	defer b.db.Close()

	views := make([]*View, 0)
	err = b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(b.bucket))

		if bucket == nil {
			return nil
		}

		c := bucket.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			sql := &View{}
			err := json.Unmarshal(v, sql)
			if err != nil {
				return err
			}

			views = append(views, sql)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return views, nil
}

func (b boltImpl) ListCommandsByDBName(databaseName string) ([]*View, error) {
	var err error
	b.db, err = b.openConnection()
	if err != nil {
		return nil, err
	}
	defer b.db.Close()

	views := make([]*View, 0)
	err = b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(b.bucket))

		if bucket == nil {
			return nil
		}

		c := bucket.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			sql := &View{}
			err := json.Unmarshal(v, sql)
			if err != nil {
				return err
			}

			if sql.DatabaseName == databaseName {
				views = append(views, sql)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return views, nil
}
