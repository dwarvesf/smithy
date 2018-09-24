// +build integration

package drivers

import (
	"math"
	"reflect"
	"sort"
	"testing"

	"github.com/jinzhu/gorm"

	"github.com/dwarvesf/smithy/backend/sqlmapper"
	"github.com/dwarvesf/smithy/common/database"
	utilDB "github.com/dwarvesf/smithy/common/utils/database/pg"
	utilTest "github.com/dwarvesf/smithy/common/utils/database/pg/test/set1"
)

func Test_pgStore_Query(t *testing.T) {
	t.Parallel()
	// create config & create database with DOCKER SDK
	cfg, clearDB := utilTest.CreateConfig(t)
	defer clearDB()

	users := []utilDB.User{}
	for _, dbase := range cfg.Databases {
		// migrate tables
		err := utilTest.MigrateTables(cfg.DB(dbase.DBName))
		if err != nil {
			t.Fatalf("Failed to migrate table by error %v", err)
		}

		// create sample data
		users, err = utilTest.CreateUserSampleData(cfg.DB(dbase.DBName))
		if err != nil {
			t.Fatalf("Failed to create sample data by error %v", err)
		}
	}

	// sort for sample data
	// sort users by name in ascending order
	ascUserName := make([]utilDB.User, len(users))
	copy(ascUserName, users)
	sort.Slice(ascUserName, func(i, j int) bool { return ascUserName[i].Name < ascUserName[j].Name })

	// sort users by name in descending order
	descUserName := make([]utilDB.User, len(users))
	copy(descUserName, users)
	sort.Slice(descUserName, func(i, j int) bool { return descUserName[i].Name > descUserName[j].Name })

	// sort users by ID in ascending order
	ascUserID := make([]utilDB.User, len(users))
	copy(ascUserID, users)
	sort.Slice(users, func(i, j int) bool { return users[i].ID < users[j].ID })

	// sort users by ID in descending order
	descUserID := make([]utilDB.User, len(users))
	copy(descUserID, users)
	sort.Slice(descUserID, func(i, j int) bool { return descUserID[i].ID > descUserID[j].ID })

	dbTest := []string{"test1", "test2"}

	type fields struct {
		db       *gorm.DB
		modelMap map[string]database.Model
	}
	tests := []struct {
		name              string
		args              *sqlmapper.Query
		want              []string
		want1             []utilDB.User
		wantErr           bool
		testForEmptyTable bool
	}{
		{
			name: "Query an exist user in db",
			args: &sqlmapper.Query{
				SourceDatabase: dbTest[0],
				SourceTable:    "users",
				Fields:         []string{"id", "name"},
				Filter: sqlmapper.Filter{
					Operator:   "=",
					ColumnName: "id",
					Value:      1,
				},
			},
			want:  []string{"id", "name"},
			want1: users[0:1],
		},
		{
			name: "Query in an empty table",
			args: &sqlmapper.Query{
				SourceDatabase: dbTest[0],
				SourceTable:    "users",
				Fields:         []string{"id", "name"},
			},
			want:              []string{"id", "name"},
			want1:             []utilDB.User{},
			testForEmptyTable: true,
		},
		{
			name: "Query all records in db",
			args: &sqlmapper.Query{
				SourceDatabase: dbTest[0],
				SourceTable:    "users",
				Fields:         []string{"id", "name"},
			},
			want:  []string{"id", "name"},
			want1: users,
		},
		{
			name: "Query by an id not exists",
			args: &sqlmapper.Query{
				SourceDatabase: dbTest[0],
				SourceTable:    "users",
				Fields:         []string{"id", "name"},
				Filter: sqlmapper.Filter{
					Operator:   "=",
					ColumnName: "id",
					Value:      111,
				},
			},
			want:    []string{"id", "name"},
			wantErr: true,
		},
		{
			name: "Query a user by too long id",
			args: &sqlmapper.Query{
				SourceDatabase: dbTest[0],
				SourceTable:    "users",
				Fields:         []string{"id", "name"},
				Filter: sqlmapper.Filter{
					Operator:   "=",
					ColumnName: "id",
					Value:      math.MaxInt32 + 1,
				},
			},
			wantErr: true,
		},
		{
			name: "Query a record invalid column name",
			args: &sqlmapper.Query{
				SourceDatabase: dbTest[0],
				SourceTable:    "users",
				Fields:         []string{"id", "name"},
				Filter: sqlmapper.Filter{
					Operator:   "=",
					ColumnName: "namemmm",
					Value:      "hieudeptrai1",
				},
			},
			wantErr: true,
		},
		{
			name: "Query in an invalid table name",
			args: &sqlmapper.Query{
				SourceDatabase: dbTest[0],
				SourceTable:    "usershandsome",
				Fields:         []string{"id", "name"},
				Filter: sqlmapper.Filter{
					Operator:   "=",
					ColumnName: "name",
					Value:      "hieudeptrai1",
				},
			},
			wantErr: true,
		},
		{
			name: "Query with offset=3 limit=2",
			args: &sqlmapper.Query{
				SourceDatabase: dbTest[0],
				SourceTable:    "users",
				Fields:         []string{"id", "name"},
				Offset:         3,
				Limit:          2,
			},
			want:  []string{"id", "name"},
			want1: users[3:5],
		},
		{
			name: "Query with offset=3 limit=0",
			args: &sqlmapper.Query{
				SourceDatabase: dbTest[0],
				SourceTable:    "users",
				Fields:         []string{"id", "name"},
				Offset:         3,
			},
			want:  []string{"id", "name"},
			want1: users[3:],
		},
		{
			name: "Query and sort by name in descending order",
			args: &sqlmapper.Query{
				SourceDatabase: dbTest[0],
				SourceTable:    "users",
				Fields:         []string{"id", "name"},
				Order:          []string{"name", "desc"},
			},
			want:    []string{"id", "name"},
			want1:   descUserName,
			wantErr: false,
		},
		{
			name: "Query and sort by name in ascending order",
			args: &sqlmapper.Query{
				SourceDatabase: dbTest[0],
				SourceTable:    "users",
				Fields:         []string{"id", "name"},
				Order:          []string{"name", "asc"},
			},
			want:    []string{"id", "name"},
			want1:   ascUserName,
			wantErr: false,
		},
		{
			name: "Query and sort by ID in descending order",
			args: &sqlmapper.Query{
				SourceDatabase: dbTest[0],
				SourceTable:    "users",
				Fields:         []string{"id", "name"},
				Order:          []string{"id", "desc"},
			},
			want:    []string{"id", "name"},
			want1:   descUserID,
			wantErr: false,
		},
		{
			name: "Query and sort by ID in ascending order",
			args: &sqlmapper.Query{
				SourceDatabase: dbTest[0],
				SourceTable:    "users",
				Fields:         []string{"id", "name"},
				Order:          []string{"id", "asc"},
			},
			want:    []string{"id", "name"},
			want1:   ascUserID,
			wantErr: false,
		},
		{
			name: "Query without sort",
			args: &sqlmapper.Query{
				SourceDatabase: dbTest[0],
				SourceTable:    "users",
				Fields:         []string{"id", "name"},
				Order:          []string{},
			},
			want:    []string{"id", "name"},
			want1:   users,
			wantErr: false,
		},
		{
			name: "Query and sort with missing argument of Order",
			args: &sqlmapper.Query{
				SourceDatabase: dbTest[0],
				SourceTable:    "users",
				Fields:         []string{"id", "name"},
				Order:          []string{"name"},
			},
			wantErr: true,
		},
		{
			name: "Query and sort with wrong argument format for Order",
			args: &sqlmapper.Query{
				SourceDatabase: dbTest[0],
				SourceTable:    "users",
				Fields:         []string{"id", "name"},
				Order:          []string{"name", "ascending"},
			},
			wantErr: true,
		},
		{
			name: "Query and sort with column name doesn't exist",
			args: &sqlmapper.Query{
				SourceDatabase: dbTest[0],
				SourceTable:    "users",
				Fields:         []string{"id", "name"},
				Order:          []string{"age", "asc"},
			},
			wantErr: true,
		},
		{
			name: "Query user in other database",
			args: &sqlmapper.Query{
				SourceDatabase: dbTest[1],
				SourceTable:    "users",
				Fields:         []string{"id", "name"},
			},
			want:  []string{"id", "name"},
			want1: users,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s sqlmapper.Mapper
			if tt.testForEmptyTable {
				// create config & create database with DOCKER SDK
				cfgEmpty, clearDB := utilTest.CreateConfig(t)
				defer clearDB()

				// migrate tables
				for _, dbase := range cfgEmpty.Databases {
					err := utilTest.MigrateTables(cfgEmpty.DB(dbase.DBName))
					if err != nil {
						t.Fatalf("Failed to migrate table by error %v", err)
					}
				}
				s = NewPGStore(cfgEmpty.DBs(), cfgEmpty.ModelMap)
			} else {
				s = NewPGStore(cfg.DBs(), cfg.ModelMap)
			}

			got, got1, err := s.Query(*tt.args)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("pgStore.Query() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pgStore.Query() got = %v, want %v", got, tt.want)
			}

			if len(got1) != len(tt.want1) {
				t.Errorf("len(got1)=%v != len(tt.want1)=%v", len(got1), len(tt.want1))
				return
			}

			for i := 0; i < len(got1); i++ {
				u := got1[i].([]interface{})

				// convert data
				id := int(u[0].(int64))
				name := u[1].(string)

				if id != tt.want1[i].ID ||
					name != tt.want1[i].Name {
					t.Errorf("pgStore.FindByColumnName() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_pgStore_Delete(t *testing.T) {
	t.Parallel()
	// create config & create database with DOCKER SDK
	cfg, clearDB := utilTest.CreateConfig(t)
	defer clearDB()

	for _, dbase := range cfg.Databases {
		// migrate tables
		err := utilTest.MigrateTables(cfg.DB(dbase.DBName))
		if err != nil {
			t.Fatalf("Failed to migrate table by error %v", err)
		}

		// create sample data
		_, err = utilTest.CreateUserSampleData(cfg.DB(dbase.DBName))
		if err != nil {
			t.Fatalf("Failed to create sample data by error %v", err)
		}
	}

	dbTest := []string{"test1", "test2"}

	type args struct {
		databaseName string
		tableName    string
		fields       []interface{}
		data         []interface{}
	}
	tests := []struct {
		name              string
		tableName         string
		args              *args
		wantErr           bool
		testForEmptyTable bool
	}{
		{
			name:      "Delete an exist user by id",
			tableName: "users",
			args: &args{
				databaseName: dbTest[0],
				tableName:    "users",
				fields: []interface{}{
					"id",
				},
				data: []interface{}{
					"1",
				},
			},
			wantErr:           false,
			testForEmptyTable: false,
		},
		{
			name:      "Delete an exist user by name",
			tableName: "users",
			args: &args{
				databaseName: dbTest[0],
				tableName:    "users",
				fields: []interface{}{
					"id",
					"name",
				},
				data: []interface{}{
					"1",
					"hieudeptrai1",
				},
			},
			wantErr:           false,
			testForEmptyTable: false,
		},
		{
			name:      "Delete user missing primary key",
			tableName: "users",
			args: &args{
				databaseName: dbTest[0],
				tableName:    "users",
				fields: []interface{}{
					"minh dep trai chet di duoc",
				},
				data: []interface{}{
					"1",
				},
			},
			wantErr:           true,
			testForEmptyTable: false,
		},
		{
			name:      "Delete record in a empty table",
			tableName: "users",
			args: &args{
				databaseName: dbTest[0],
				tableName:    "users",
				fields: []interface{}{
					"iddfdf",
				},
				data: []interface{}{
					"1",
				},
			},
			wantErr:           true,
			testForEmptyTable: true,
		},
		{
			name:      "Delete record by id in a other database",
			tableName: "users",
			args: &args{
				databaseName: dbTest[1],
				tableName:    "users",
				fields: []interface{}{
					"id",
				},
				data: []interface{}{
					"1",
				},
			},
			wantErr:           false,
			testForEmptyTable: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s sqlmapper.Mapper
			if tt.testForEmptyTable {
				// create config & create database with DOCKER SDK
				cfgEmpty, clearDB := utilTest.CreateConfig(t)
				defer clearDB()

				for _, dbase := range cfgEmpty.Databases {
					// migrate tables
					err := utilTest.MigrateTables(cfgEmpty.DB(dbase.DBName))
					if err != nil {
						t.Fatalf("Failed to migrate table by error %v", err)
					}
				}
				s = NewPGStore(cfgEmpty.DBs(), cfgEmpty.ModelMap)
			} else {
				s = NewPGStore(cfg.DBs(), cfg.ModelMap)
			}

			err := s.Delete(tt.args.databaseName, tt.tableName, tt.args.fields, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("pgStore.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_pgStore_Create(t *testing.T) {
	t.Parallel()
	// create config & create database with DOCKER SDK
	cfg, clearDB := utilTest.CreateConfig(t)
	defer clearDB()

	for _, dbase := range cfg.Databases {
		// migrate tables
		err := utilTest.MigrateTables(cfg.DB(dbase.DBName))
		if err != nil {
			t.Fatalf("Failed to migrate table by error %v", err)
		}
	}

	dbTest := []string{"test1", "test2"}

	type args struct {
		data         sqlmapper.RowData
		databaseName string
	}
	tests := []struct {
		name      string
		tableName string
		args      args
		want      sqlmapper.RowData
		wantErr   bool
	}{
		{
			name:      "create a valid record",
			tableName: "users",
			args: args{
				databaseName: dbTest[0],
				data: sqlmapper.RowData{
					"name": sqlmapper.ColData{
						Data: "hieudeptrai",
					},
				},
			},
			want: sqlmapper.RowData{
				"id": sqlmapper.ColData{
					Data: 1,
				},
				"name": sqlmapper.ColData{
					Data: "hieudeptrai",
				},
			},
		},
		{
			name:      "create record missing data",
			tableName: "users",
			args: args{
				databaseName: dbTest[0],
				data:         sqlmapper.RowData{},
			},
			wantErr: true,
		},
		{
			name:      "invalid column name",
			tableName: "users",
			args: args{
				databaseName: dbTest[0],
				data: sqlmapper.RowData{
					"namenmce": sqlmapper.ColData{
						Data: "hieudeptrai",
					},
				},
			},
			wantErr: true,
		},
		{
			name:      "create a valid record in other database",
			tableName: "users",
			args: args{
				databaseName: dbTest[1],
				data: sqlmapper.RowData{
					"name": sqlmapper.ColData{
						Data: "hieudeptrai",
					},
				},
			},
			want: sqlmapper.RowData{
				"id": sqlmapper.ColData{
					Data: 1,
				},
				"name": sqlmapper.ColData{
					Data: "hieudeptrai",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewPGStore(cfg.DBs(), cfg.ModelMap)
			got, err := s.Create(tt.args.databaseName, tt.tableName, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("pgStore.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pgStore.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pgStore_Update(t *testing.T) {
	t.Parallel()
	// create config & create database with DOCKER SDK
	cfg, clearDB := utilTest.CreateConfig(t)
	defer clearDB()

	dbTest := []string{"test1", "test2"}

	users := []utilDB.User{}
	for _, dbase := range cfg.Databases {
		// migrate tables
		err := utilTest.MigrateTables(cfg.DB(dbase.DBName))
		if err != nil {
			t.Fatalf("Failed to migrate table by error %v", err)
		}

		//create sample data
		users, err = utilTest.CreateUserSampleData(cfg.DB(dbase.DBName))
		if err != nil {
			t.Fatalf("Failed to create sample data by error %v", err)
		}
	}

	type args struct {
		d            sqlmapper.RowData
		id           int
		databaseName string
	}
	tests := []struct {
		name      string
		tableName string
		args      args
		want      sqlmapper.RowData
		wantErr   bool
	}{
		{
			name:      "Update an valid record",
			tableName: "users",
			args: args{
				databaseName: dbTest[0],
				d: sqlmapper.RowData{
					"id": sqlmapper.ColData{
						Data: users[0].ID,
					},
					"name": sqlmapper.ColData{
						Data: "demo",
					},
				},
			},
			want: sqlmapper.RowData{
				"name": sqlmapper.ColData{
					Data: "demo",
				},
			},
			wantErr: false,
		},
		{
			name:      "primary key isn't exist",
			tableName: "users",
			args: args{
				databaseName: dbTest[0],
				d: sqlmapper.RowData{
					"id": sqlmapper.ColData{
						Data: -1,
					},
				},
			},
			wantErr: true,
		},
		{
			name:      "rowData is empty",
			tableName: "users",
			args: args{
				databaseName: dbTest[0],
			},
			wantErr: true,
		},
		{
			name:      "invalid column name",
			tableName: "users",
			args: args{
				databaseName: dbTest[0],
				d: sqlmapper.RowData{
					"id": sqlmapper.ColData{
						Data: users[0].ID,
					},
					"blabla": sqlmapper.ColData{
						Data: "anmt",
					},
				},
			},
			wantErr: true,
		},
		{
			name:      "Update a valid record in other database",
			tableName: "users",
			args: args{
				databaseName: dbTest[1],
				d: sqlmapper.RowData{
					"id": sqlmapper.ColData{
						Data: users[0].ID,
					},
					"name": sqlmapper.ColData{
						Data: "demo",
					},
				},
			},
			want: sqlmapper.RowData{
				"name": sqlmapper.ColData{
					Data: "demo",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewPGStore(cfg.DBs(), cfg.ModelMap)
			got, err := s.Update(tt.args.databaseName, tt.tableName, tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("pgStore.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pgStore.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pgStore_RawQuery(t *testing.T) {
	t.Parallel()
	// create config & create database with DOCKER SDK
	cfg, clearDB := utilTest.CreateConfig(t)
	defer clearDB()

	users := []utilDB.User{}
	for _, dbase := range cfg.Databases {
		// migrate tables
		err := utilTest.MigrateTables(cfg.DB(dbase.DBName))
		if err != nil {
			t.Fatalf("Failed to migrate table by error %v", err)
		}

		// create sample data
		users, err = utilTest.CreateUserSampleData(cfg.DB(dbase.DBName))
		if err != nil {
			t.Fatalf("Failed to create sample data by error %v", err)
		}
	}

	dbTest := []string{"test1", "test2"}

	type args struct {
		dbName string
		sql    string
	}
	tests := []struct {
		name    string
		args    args
		want    []utilDB.User
		wantErr bool
	}{
		{
			name: "in valid view",
			args: args{
				dbName: dbTest[0],
				sql:    "SELECTS * FROM users WHERE id = 1",
			},
			wantErr: true,
		},
		{
			name: "select user by id",
			args: args{
				dbName: dbTest[0],
				sql:    "SELECT * FROM users WHERE id = 1",
			},
			want: users[0:1],
		},
		{
			name: "select all user",
			args: args{
				dbName: dbTest[0],
				sql:    "SELECT * FROM users",
			},
			want: users,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewPGStore(cfg.DBs(), cfg.ModelMap)
			_, _, got, err := s.RawQuery(tt.args.dbName, tt.args.sql)
			if (err != nil) != tt.wantErr {
				t.Errorf("pgStore.RawQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("pgStore.RawQuery() got = %v, want %v", got, tt.want)
				return
			}

			for i := range got {
				u := got[i].([]interface{})
				id := int(u[0].(int64))
				name := u[1].(string)
				if id != tt.want[i].ID || name != tt.want[i].Name {
					t.Errorf("pgStore.RawQuery() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
