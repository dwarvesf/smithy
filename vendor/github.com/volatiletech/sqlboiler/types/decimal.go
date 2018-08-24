package types

import (
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/ericlagergren/decimal"
)

// Decimal is a DECIMAL in sql. Its zero value is valid for use with both
// Value and Scan.
//
// Although decimal can represent NaN and Infinity it will return an error
// if an attempt to store these values in the database is made.
//
// Because it cannot be nil, when Big is nil Value() will return "0"
// It will error if an attempt to Scan() a "null" value into it.
type Decimal struct {
	*decimal.Big
}

// NullDecimal is the same as Decimal, but allows the Big pointer to be nil.
// See docmentation for Decimal for more details.
//
// When going into a database, if Big is nil it's value will be "null".
type NullDecimal struct {
	*decimal.Big
}

// NewDecimal creates a new decimal from a decimal
func NewDecimal(d *decimal.Big) Decimal {
	return Decimal{Big: d}
}

// NewNullDecimal creates a new null decimal from a decimal
func NewNullDecimal(d *decimal.Big) NullDecimal {
	return NullDecimal{Big: d}
}

// Value implements driver.Valuer.
func (d Decimal) Value() (driver.Value, error) {
	return decimalValue(d.Big, false)
}

// Scan implements sql.Scanner.
func (d *Decimal) Scan(val interface{}) error {
	newD, err := decimalScan(d.Big, val, false)
	if err != nil {
		return err
	}

	d.Big = newD
	return nil
}

// Randomize implements sqlboiler's randomize interface
func (d *Decimal) Randomize(nextInt func() int64, fieldType string, shouldBeNull bool) {
	d.Big = randomDecimal(nextInt, fieldType, false)
}

// Value implements driver.Valuer.
func (n NullDecimal) Value() (driver.Value, error) {
	return decimalValue(n.Big, true)
}

// Scan implements sql.Scanner.
func (n *NullDecimal) Scan(val interface{}) error {
	newD, err := decimalScan(n.Big, val, true)
	if err != nil {
		return err
	}

	n.Big = newD
	return nil
}

// Randomize implements sqlboiler's randomize interface
func (n *NullDecimal) Randomize(nextInt func() int64, fieldType string, shouldBeNull bool) {
	n.Big = randomDecimal(nextInt, fieldType, shouldBeNull)
}

func randomDecimal(nextInt func() int64, fieldType string, shouldBeNull bool) *decimal.Big {
	if shouldBeNull {
		return nil
	}

	randVal := fmt.Sprintf("%d.%d", nextInt()%10, nextInt()%10)
	random, success := new(decimal.Big).SetString(randVal)
	if !success {
		panic("randVal could not be turned into a decimal")
	}

	return random
}

func decimalValue(d *decimal.Big, canNull bool) (driver.Value, error) {
	if canNull && d == nil {
		return nil, nil
	}

	if d.IsNaN(0) {
		return nil, errors.New("refusing to allow NaN into database")
	}
	if d.IsInf(0) {
		return nil, errors.New("refusing to allow infinity into database")
	}

	return d.String(), nil
}

func decimalScan(d *decimal.Big, val interface{}, canNull bool) (*decimal.Big, error) {
	if val == nil {
		if !canNull {
			return nil, errors.New("null cannot be scanned into decimal")
		}

		return nil, nil
	}

	if d == nil {
		d = new(decimal.Big)
	}

	switch t := val.(type) {
	case float64:
		d.SetFloat64(t)
		return d, nil
	case string:
		if _, ok := d.SetString(t); !ok {
			if err := d.Context.Err(); err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("invalid decimal syntax: %q", t)
		}
		return d, nil
	case []byte:
		if err := d.UnmarshalText(t); err != nil {
			return nil, err
		}
		return d, nil
	default:
		return nil, fmt.Errorf("cannot scan decimal value: %#v", val)
	}
}
