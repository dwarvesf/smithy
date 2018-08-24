package bolt

import (
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/boltdb/bolt"
)

// CreateDatabase create bolt db for test
func CreateDatabase(t *testing.T) (*bolt.DB, func()) {
	var err error

	rand.Seed(time.Now().UnixNano())
	persistenceFileName := "test" + strconv.FormatInt(rand.Int63(), 10)
	boltDB, err := bolt.Open(persistenceFileName, 0600, &bolt.Options{Timeout: 1 * time.Second})

	if err != nil {
		t.Fatalf("Fail to create persistence test db. %s", err.Error())
	}

	return boltDB, func() {
		if boltDB != nil {
			boltDB.Close()
		}

		err := os.Remove(persistenceFileName)
		if err != nil {
			t.Fatalf("Fail to drop BoltDB. %s", err.Error())
		}
	}
}
