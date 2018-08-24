package boil

import (
	"io"
	"os"
	"time"
)

var (
	// currentDB is a global database handle for the package
	currentDB Executor
	// timestampLocation is the timezone used for the
	// automated setting of created_at/updated_at columns
	timestampLocation = time.UTC
)

// DebugMode is a flag controlling whether generated sql statements and
// debug information is outputted to the DebugWriter handle
//
// NOTE: This should be disabled in production to avoid leaking sensitive data
var DebugMode = false

// DebugWriter is where the debug output will be sent if DebugMode is true
var DebugWriter io.Writer = os.Stdout

// SetDB initializes the database handle for all template db interactions
func SetDB(db Executor) {
	currentDB = db
}

// GetDB retrieves the global state database handle
func GetDB() Executor {
	return currentDB
}

// SetLocation sets the global timestamp Location.
// This is the timezone used by the generated package for the
// automated setting of created_at and updated_at columns.
// If the package was generated with the --no-auto-timestamps flag
// then this function has no effect.
func SetLocation(loc *time.Location) {
	timestampLocation = loc
}

// GetLocation retrieves the global timestamp Location.
// This is the timezone used by the generated package for the
// automated setting of created_at and updated_at columns
// if the package was not generated with the --no-auto-timestamps flag.
func GetLocation() *time.Location {
	return timestampLocation
}
