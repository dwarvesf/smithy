package bolt

import (
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

// CreateDatabase create bolt db for test
func CreateDatabase(t *testing.T) (string, func()) {
	var err error

	rand.Seed(time.Now().UnixNano())
	persistenceFileName := "test" + strconv.FormatInt(rand.Int63(), 10)

	if err != nil {
		t.Fatalf("Fail to create persistence test db. %s", err.Error())
	}

	return persistenceFileName, func() {
		err := os.Remove(persistenceFileName)
		if err != nil {
			t.Fatalf("Fail to drop BoltDB. %s", err.Error())
		}
	}
}
