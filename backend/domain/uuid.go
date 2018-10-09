package domain

import (
	"database/sql/driver"
	"errors"

	uuid "github.com/satori/go.uuid"
)

// UUID implement for go-pg convert uuid
type UUID [16]byte

// NewUUID create new UUID with V4
func NewUUID() UUID {
	return UUID(uuid.NewV4())
}

// UUIDFromString convert string to UUID
func UUIDFromString(s string) (UUID, error) {
	id, err := uuid.FromString(s)
	return UUID(id), err
}

// IsZero check uuid is zero
func (u *UUID) IsZero() bool {
	if u == nil {
		return true
	}
	for _, c := range u {
		if c != 0 {
			return false
		}
	}
	return true
}

func (u UUID) String() string {
	return uuid.UUID(u).String()
}

// MarshalJSON implement for json encoding
func (u UUID) MarshalJSON() ([]byte, error) {
	if len(u) == 0 {
		return []byte(`""`), nil
	}
	return []byte(`"` + u.String() + `"`), nil
}

// UnmarshalJSON implement for json decoding
func (u *UUID) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == `""` {
		return nil
	}

	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return errors.New("invalid UUID format")
	}
	data = data[1 : len(data)-1]
	uu, err := uuid.FromString(string(data))
	if err != nil {
		return errors.New("invalid UUID format")
	}
	*u = UUID(uu)
	return nil
}

// Value .
func (u UUID) Value() (driver.Value, error) {
	if u.IsZero() {
		return nil, nil
	}
	return uuid.UUID(u).String(), nil
}

// Scan .
func (u *UUID) Scan(b interface{}) error {
	if b == nil {
		for i := range u {
			u[i] = 0
		}
		return nil
	}

	// postgres store DB as a string
	id, err := uuid.FromString(string(b.([]byte)))
	if err != nil {
		return err
	}

	for i, c := range id {
		u[i] = c
	}

	return nil
}

// MustGetUUIDFromString get uuid from string if failed throw panic,
// CAUTION: IT ONLY USE FOR TESTING
func MustGetUUIDFromString(s string) UUID {
	id, err := uuid.FromString(s)
	if err != nil {
		panic(err)
	}
	return UUID(id)
}
