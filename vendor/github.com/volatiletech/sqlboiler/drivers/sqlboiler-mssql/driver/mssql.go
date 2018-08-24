package driver

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"

	// Side effect import go-mssqldb
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/drivers"
	"github.com/volatiletech/sqlboiler/importers"
	"github.com/volatiletech/sqlboiler/strmangle"
)

func init() {
	drivers.RegisterFromInit("mssql", &MSSQLDriver{})
}

//go:generate go-bindata -nometadata -pkg driver -prefix override override/...

// Assemble is more useful for calling into the library so you don't
// have to instantiate an empty type.
func Assemble(config drivers.Config) (dbinfo *drivers.DBInfo, err error) {
	driver := MSSQLDriver{}
	return driver.Assemble(config)
}

// MSSQLDriver holds the database connection string and a handle
// to the database connection.
type MSSQLDriver struct {
	connStr string
	conn    *sql.DB
}

// Templates that should be added/overridden
func (MSSQLDriver) Templates() (map[string]string, error) {
	names := AssetNames()
	tpls := make(map[string]string)
	for _, n := range names {
		b, err := Asset(n)
		if err != nil {
			return nil, err
		}

		tpls[n] = base64.StdEncoding.EncodeToString(b)
	}

	return tpls, nil
}

// Assemble all the information we need to provide back to the driver
func (m *MSSQLDriver) Assemble(config drivers.Config) (dbinfo *drivers.DBInfo, err error) {
	defer func() {
		if r := recover(); r != nil && err == nil {
			dbinfo = nil
			err = r.(error)
		}
	}()

	user := config.MustString(drivers.ConfigUser)
	pass, _ := config.String(drivers.ConfigPass)
	dbname := config.MustString(drivers.ConfigDBName)
	host := config.MustString(drivers.ConfigHost)
	port := config.DefaultInt(drivers.ConfigPort, 1433)
	sslmode := config.DefaultString(drivers.ConfigSSLMode, "true")

	schema := config.DefaultString(drivers.ConfigSchema, "dbo")
	whitelist, _ := config.StringSlice(drivers.ConfigWhitelist)
	blacklist, _ := config.StringSlice(drivers.ConfigBlacklist)

	m.connStr = MSSQLBuildQueryString(user, pass, dbname, host, port, sslmode)
	m.conn, err = sql.Open("mssql", m.connStr)
	if err != nil {
		return nil, errors.Wrap(err, "sqlboiler-mssql failed to connect to database")
	}

	defer func() {
		if e := m.conn.Close(); e != nil {
			dbinfo = nil
			err = e
		}
	}()

	dbinfo = &drivers.DBInfo{
		Schema: schema,
		Dialect: drivers.Dialect{
			LQ: '[',
			RQ: ']',

			UseIndexPlaceholders: true,
			UseSchema:            true,
			UseDefaultKeyword:    true,

			UseAutoColumns:          true,
			UseTopClause:            true,
			UseOutputClause:         true,
			UseCaseWhenExistsClause: true,
		},
	}
	dbinfo.Tables, err = drivers.Tables(m, schema, whitelist, blacklist)
	if err != nil {
		return nil, err
	}

	return dbinfo, err
}

// MSSQLBuildQueryString builds a query string for MSSQL.
func MSSQLBuildQueryString(user, pass, dbname, host string, port int, sslmode string) string {
	query := url.Values{}
	query.Add("database", dbname)
	query.Add("encrypt", sslmode)

	u := &url.URL{
		Scheme: "sqlserver",
		User:   url.UserPassword(user, pass),
		Host:   fmt.Sprintf("%s:%d", host, port),
		// Path:  instance, // if connecting to an instance instead of a port
		RawQuery: query.Encode(),
	}

	return u.String()
}

// TableNames connects to the postgres database and
// retrieves all table names from the information_schema where the
// table schema is schema. It uses a whitelist and blacklist.
func (m *MSSQLDriver) TableNames(schema string, whitelist, blacklist []string) ([]string, error) {
	var names []string

	query := `
		SELECT table_name
		FROM   information_schema.tables
		WHERE  table_schema = ? AND table_type = 'BASE TABLE'`

	args := []interface{}{schema}
	if len(whitelist) > 0 {
		tables := drivers.TablesFromList(whitelist)
		if len(tables) > 0 {
			query += fmt.Sprintf(" AND table_name IN (%s);", strings.Repeat(",?", len(tables))[1:])
			for _, w := range tables {
				args = append(args, w)
			}
		}
	} else if len(blacklist) > 0 {
		tables := drivers.TablesFromList(whitelist)
		if len(tables) > 0 {
			query += fmt.Sprintf(" AND table_name not IN (%s);", strings.Repeat(",?", len(tables))[1:])
			for _, b := range tables {
				args = append(args, b)
			}
		}
	}

	rows, err := m.conn.Query(query, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}

	return names, nil
}

// Columns takes a table name and attempts to retrieve the table information
// from the database information_schema.columns. It retrieves the column names
// and column types and returns those as a []Column after TranslateColumnType()
// converts the SQL types to Go types, for example: "varchar" to "string"
func (m *MSSQLDriver) Columns(schema, tableName string, whitelist, blacklist []string) ([]drivers.Column, error) {
	var columns []drivers.Column
	args := []interface{}{schema, tableName}
	query := `
	SELECT column_name,
       CASE
         WHEN character_maximum_length IS NULL THEN data_type
         ELSE data_type + '(' + CAST(character_maximum_length AS VARCHAR) + ')'
       END AS full_type,
       data_type,
	   column_default,
       CASE
         WHEN is_nullable = 'YES' THEN 1
         ELSE 0
       END AS is_nullable,
       CASE
         WHEN EXISTS (SELECT c.column_name
                      FROM information_schema.table_constraints tc
                        INNER JOIN information_schema.key_column_usage kcu
                                ON tc.constraint_name = kcu.constraint_name
                               AND tc.table_name = kcu.table_name
                               AND tc.table_schema = kcu.table_schema
                      WHERE c.column_name = kcu.column_name
                      AND   tc.table_name = c.table_name
                      AND   (tc.constraint_type = 'PRIMARY KEY' OR tc.constraint_type = 'UNIQUE')
                      AND   (SELECT COUNT(*)
                             FROM information_schema.key_column_usage
                             WHERE table_schema = kcu.table_schema
                             AND   table_name = tc.table_name
                             AND   constraint_name = tc.constraint_name) = 1) THEN 1
         ELSE 0
       END AS is_unique,
	   COLUMNPROPERTY(object_id($1 + '.' + $2), c.column_name, 'IsIdentity') as is_identity
	FROM information_schema.columns c
	WHERE table_schema = $1 AND table_name = $2`

	if len(whitelist) > 0 {
		cols := drivers.ColumnsFromList(whitelist, tableName)
		if len(cols) > 0 {
			query += fmt.Sprintf(" and c.column_name in (%s)", strmangle.Placeholders(true, len(cols), 3, 1))
			for _, w := range cols {
				args = append(args, w)
			}
		}
	} else if len(blacklist) > 0 {
		cols := drivers.ColumnsFromList(blacklist, tableName)
		if len(cols) > 0 {
			query += fmt.Sprintf(" and c.column_name not in (%s)", strmangle.Placeholders(true, len(cols), 3, 1))
			for _, w := range cols {
				args = append(args, w)
			}
		}
	}

	rows, err := m.conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var colName, colType, colFullType string
		var nullable, unique, identity, auto bool
		var defaultValue *string
		if err := rows.Scan(&colName, &colFullType, &colType, &defaultValue, &nullable, &unique, &identity); err != nil {
			return nil, errors.Wrapf(err, "unable to scan for table %s", tableName)
		}

		auto = strings.EqualFold(colType, "timestamp") || strings.EqualFold(colType, "rowversion")

		column := drivers.Column{
			Name:          colName,
			FullDBType:    colFullType,
			DBType:        colType,
			Nullable:      nullable,
			Unique:        unique,
			AutoGenerated: auto,
		}

		if defaultValue != nil && *defaultValue != "NULL" {
			column.Default = *defaultValue
		} else if identity || auto {
			column.Default = "auto"
		}
		columns = append(columns, column)
	}

	return columns, nil
}

// PrimaryKeyInfo looks up the primary key for a table.
func (m *MSSQLDriver) PrimaryKeyInfo(schema, tableName string) (*drivers.PrimaryKey, error) {
	pkey := &drivers.PrimaryKey{}
	var err error

	query := `
	SELECT constraint_name
	FROM   information_schema.table_constraints
	WHERE  table_name = ? AND constraint_type = 'PRIMARY KEY' AND table_schema = ?;`

	row := m.conn.QueryRow(query, tableName, schema)
	if err = row.Scan(&pkey.Name); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	queryColumns := `
	SELECT column_name
	FROM   information_schema.key_column_usage
	WHERE  table_name = ? AND constraint_name = ? AND table_schema = ?;`

	var rows *sql.Rows
	if rows, err = m.conn.Query(queryColumns, tableName, pkey.Name, schema); err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var column string

		err = rows.Scan(&column)
		if err != nil {
			return nil, err
		}

		columns = append(columns, column)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	pkey.Columns = columns

	return pkey, nil
}

// ForeignKeyInfo retrieves the foreign keys for a given table name.
func (m *MSSQLDriver) ForeignKeyInfo(schema, tableName string) ([]drivers.ForeignKey, error) {
	var fkeys []drivers.ForeignKey

	query := `
	SELECT ccu.constraint_name ,
		ccu.table_name AS local_table ,
		ccu.column_name AS local_column ,
		kcu.table_name AS foreign_table ,
		kcu.column_name AS foreign_column
	FROM information_schema.constraint_column_usage ccu
	INNER JOIN information_schema.referential_constraints rc ON ccu.constraint_name = rc.constraint_name
	INNER JOIN information_schema.key_column_usage kcu ON kcu.constraint_name = rc.unique_constraint_name
	WHERE ccu.table_schema = ?
	  AND ccu.constraint_schema = ?
	  AND ccu.table_name = ?
	`

	var rows *sql.Rows
	var err error
	if rows, err = m.conn.Query(query, schema, schema, tableName); err != nil {
		return nil, err
	}

	for rows.Next() {
		var fkey drivers.ForeignKey
		var sourceTable string

		fkey.Table = tableName
		err = rows.Scan(&fkey.Name, &sourceTable, &fkey.Column, &fkey.ForeignTable, &fkey.ForeignColumn)
		if err != nil {
			return nil, err
		}

		fkeys = append(fkeys, fkey)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return fkeys, nil
}

// TranslateColumnType converts postgres database types to Go types, for example
// "varchar" to "string" and "bigint" to "int64". It returns this parsed data
// as a Column object.
func (m *MSSQLDriver) TranslateColumnType(c drivers.Column) drivers.Column {
	if c.Nullable {
		switch c.DBType {
		case "tinyint":
			c.Type = "null.Int8"
		case "smallint":
			c.Type = "null.Int16"
		case "mediumint":
			c.Type = "null.Int32"
		case "int":
			c.Type = "null.Int"
		case "bigint":
			c.Type = "null.Int64"
		case "real":
			c.Type = "null.Float32"
		case "float":
			c.Type = "null.Float64"
		case "boolean", "bool", "bit":
			c.Type = "null.Bool"
		case "date", "datetime", "datetime2", "smalldatetime", "time":
			c.Type = "null.Time"
		case "binary", "varbinary":
			c.Type = "null.Bytes"
		case "timestamp", "rowversion":
			c.Type = "null.Bytes"
		case "xml":
			c.Type = "null.String"
		case "uniqueidentifier":
			c.Type = "null.String"
			c.DBType = "uuid"
		case "numeric", "decimal", "dec":
			c.Type = "types.NullDecimal"
		default:
			c.Type = "null.String"
		}
	} else {
		switch c.DBType {
		case "tinyint":
			c.Type = "int8"
		case "smallint":
			c.Type = "int16"
		case "mediumint":
			c.Type = "int32"
		case "int":
			c.Type = "int"
		case "bigint":
			c.Type = "int64"
		case "real":
			c.Type = "float32"
		case "float":
			c.Type = "float64"
		case "boolean", "bool", "bit":
			c.Type = "bool"
		case "date", "datetime", "datetime2", "smalldatetime", "time":
			c.Type = "time.Time"
		case "binary", "varbinary":
			c.Type = "[]byte"
		case "timestamp", "rowversion":
			c.Type = "[]byte"
		case "xml":
			c.Type = "string"
		case "uniqueidentifier":
			c.Type = "string"
			c.DBType = "uuid"
		case "numeric", "decimal", "dec":
			c.Type = "types.Decimal"
		default:
			c.Type = "string"
		}
	}

	return c
}

// Imports returns important imports for the driver
func (MSSQLDriver) Imports() (col importers.Collection, err error) {
	col.All = importers.Set{
		Standard: importers.List{
			`"strconv"`,
		},
	}
	col.Singleton = importers.Map{
		"mssql_upsert": {
			Standard: importers.List{
				`"fmt"`,
				`"strings"`,
			},
			ThirdParty: importers.List{
				`"github.com/volatiletech/sqlboiler/strmangle"`,
				`"github.com/volatiletech/sqlboiler/drivers"`,
			},
		},
	}
	col.TestSingleton = importers.Map{
		"mssql_suites_test": {
			Standard: importers.List{
				`"testing"`,
			},
		},
		"mssql_main_test": {
			Standard: importers.List{
				`"bytes"`,
				`"database/sql"`,
				`"fmt"`,
				`"os"`,
				`"os/exec"`,
				`"regexp"`,
				`"strings"`,
			},
			ThirdParty: importers.List{
				`"github.com/kat-co/vala"`,
				`"github.com/pkg/errors"`,
				`"github.com/spf13/viper"`,
				`"github.com/volatiletech/sqlboiler/drivers/sqlboiler-mssql/driver"`,
				`"github.com/volatiletech/sqlboiler/randomize"`,
				`_ "github.com/denisenkom/go-mssqldb"`,
			},
		},
	}

	col.BasedOnType = importers.Map{
		"null.Float32": {
			ThirdParty: importers.List{`"github.com/volatiletech/null"`},
		},
		"null.Float64": {
			ThirdParty: importers.List{`"github.com/volatiletech/null"`},
		},
		"null.Int": {
			ThirdParty: importers.List{`"github.com/volatiletech/null"`},
		},
		"null.Int8": {
			ThirdParty: importers.List{`"github.com/volatiletech/null"`},
		},
		"null.Int16": {
			ThirdParty: importers.List{`"github.com/volatiletech/null"`},
		},
		"null.Int32": {
			ThirdParty: importers.List{`"github.com/volatiletech/null"`},
		},
		"null.Int64": {
			ThirdParty: importers.List{`"github.com/volatiletech/null"`},
		},
		"null.Uint": {
			ThirdParty: importers.List{`"github.com/volatiletech/null"`},
		},
		"null.Uint8": {
			ThirdParty: importers.List{`"github.com/volatiletech/null"`},
		},
		"null.Uint16": {
			ThirdParty: importers.List{`"github.com/volatiletech/null"`},
		},
		"null.Uint32": {
			ThirdParty: importers.List{`"github.com/volatiletech/null"`},
		},
		"null.Uint64": {
			ThirdParty: importers.List{`"github.com/volatiletech/null"`},
		},
		"null.String": {
			ThirdParty: importers.List{`"github.com/volatiletech/null"`},
		},
		"null.Bool": {
			ThirdParty: importers.List{`"github.com/volatiletech/null"`},
		},
		"null.Time": {
			ThirdParty: importers.List{`"github.com/volatiletech/null"`},
		},
		"null.Bytes": {
			ThirdParty: importers.List{`"github.com/volatiletech/null"`},
		},
		"time.Time": {
			Standard: importers.List{`"time"`},
		},
		"types.Decimal": {
			Standard: importers.List{`"github.com/volatiletech/sqlboiler/types"`},
		},
		"types.NullDecimal": {
			Standard: importers.List{`"github.com/volatiletech/sqlboiler/types"`},
		},
	}
	return col, err
}
