// +build integration

package drivers

import (
	"math"
	"reflect"
	"testing"

	"github.com/jinzhu/gorm"

	"github.com/dwarvesf/smithy/backend/sqlmapper"
	"github.com/dwarvesf/smithy/common/database"
	utilDB "github.com/dwarvesf/smithy/common/utils/database/pg"
	utilTest "github.com/dwarvesf/smithy/common/utils/database/pg/test/set1"
)

func Test_pgStore_Query(t *testing.T) {
	t.Parallel()
	cfg, clearDB := utilTest.CreateConfig(t)
	defer clearDB()

	// migrate tables
	err := utilTest.MigrateTables(cfg.DB())
	if err != nil {
		t.Fatalf("Failed to migrate table by error %v", err)
	}

	//create sample data
	users, err := utilTest.CreateUserSampleData(cfg.DB())
	if err != nil {
		t.Fatalf("Failed to create sample data by error %v", err)
	}

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
			name: "Valid test case",
			args: &sqlmapper.Query{
				SourceTable: "users",
				Fields:      []string{"id", "name"},
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
			name: "empty table",
			args: &sqlmapper.Query{
				SourceTable: "users",
				Fields:      []string{"id", "name"},
			},
			want:              []string{"id", "name"},
			want1:             []utilDB.User{},
			testForEmptyTable: true,
		},
		{
			name: "Find all",
			args: &sqlmapper.Query{
				SourceTable: "users",
				Fields:      []string{"id", "name"},
			},
			want:  []string{"id", "name"},
			want1: users,
		},
		{
			name: "id not exists",
			args: &sqlmapper.Query{
				SourceTable: "users",
				Fields:      []string{"id", "name"},
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
			name: "id too long",
			args: &sqlmapper.Query{
				SourceTable: "users",
				Fields:      []string{"id", "name"},
				Filter: sqlmapper.Filter{
					Operator:   "=",
					ColumnName: "id",
					Value:      math.MaxInt32 + 1,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid column name",
			args: &sqlmapper.Query{
				SourceTable: "users",
				Fields:      []string{"id", "name"},
				Filter: sqlmapper.Filter{
					Operator:   "=",
					ColumnName: "namemmm",
					Value:      "hieudeptrai1",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid table name",
			args: &sqlmapper.Query{
				SourceTable: "usershandsome",
				Fields:      []string{"id", "name"},
				Filter: sqlmapper.Filter{
					Operator:   "=",
					ColumnName: "name",
					Value:      "hieudeptrai1",
				},
			},
			wantErr: true,
		},
		{
			name: "offset=3 limit=2",
			args: &sqlmapper.Query{
				SourceTable: "users",
				Fields:      []string{"id", "name"},
				Offset:      3,
				Limit:       2,
			},
			want:  []string{"id", "name"},
			want1: users[3:5],
		},
		{
			name: "offset=3 limit=0",
			args: &sqlmapper.Query{
				SourceTable: "users",
				Fields:      []string{"id", "name"},
				Offset:      3,
			},
			want:  []string{"id", "name"},
			want1: users[3:],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s sqlmapper.Mapper
			if tt.testForEmptyTable {
				cfgEmpty, clearDB := utilTest.CreateConfig(t)
				defer clearDB()

				// migrate tables
				err := utilTest.MigrateTables(cfgEmpty.DB())
				if err != nil {
					t.Fatalf("Failed to migrate table by error %v", err)
				}
				s = NewPGStore(cfgEmpty.DB(), database.Models(cfgEmpty.ModelList).GroupByName())
			} else {
				s = NewPGStore(cfg.DB(), database.Models(cfg.ModelList).GroupByName())
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
	cfg, clearDB := utilTest.CreateConfig(t)
	defer clearDB()

	// migrate tables
	err := utilTest.MigrateTables(cfg.DB())
	if err != nil {
		t.Fatalf("Failed to migrate table by error %v", err)
	}

	//create sample data
	_, err = utilTest.CreateUserSampleData(cfg.DB())
	if err != nil {
		t.Fatalf("Failed to create sample data by error %v", err)
	}

	type args struct {
		tableName string
		fields    []interface{}
		data      []interface{}
	}
	tests := []struct {
		name              string
		tableName         string
		args              *args
		wantErr           bool
		testForEmptyTable bool
	}{
		{
			name:      "Valid test case: id",
			tableName: "users",
			args: &args{
				tableName: "users",
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
			name:      "Valid test case: id name",
			tableName: "users",
			args: &args{
				tableName: "users",
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
			name:      "Invalid testcase",
			tableName: "users",
			args: &args{
				tableName: "users",
				fields: []interface{}{
					"iddfdf",
				},
				data: []interface{}{
					"1",
				},
			},
			wantErr:           true,
			testForEmptyTable: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s sqlmapper.Mapper
			if tt.testForEmptyTable {
				cfgEmpty, clearDB := utilTest.CreateConfig(t)
				defer clearDB()

				// migrate tables
				err := utilTest.MigrateTables(cfgEmpty.DB())
				if err != nil {
					t.Fatalf("Failed to migrate table by error %v", err)
				}
				s = NewPGStore(cfgEmpty.DB(), database.Models(cfgEmpty.ModelList).GroupByName())
			} else {
				s = NewPGStore(cfg.DB(), database.Models(cfg.ModelList).GroupByName())
			}

			err := s.Delete(tt.tableName, tt.args.fields, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("pgStore.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_pgStore_Create(t *testing.T) {
	t.Parallel()
	cfg, clearDB := utilTest.CreateConfig(t)
	defer clearDB()

	// migrate tables
	err := utilTest.MigrateTables(cfg.DB())
	if err != nil {
		t.Fatalf("Failed to migrate table by error %v", err)
	}

	type args struct {
		data sqlmapper.RowData
	}
	tests := []struct {
		name      string
		tableName string
		args      args
		want      sqlmapper.RowData
		wantErr   bool
	}{
		{
			name:      "valid user",
			tableName: "users",
			args: args{
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
			name:      "empty input",
			tableName: "users",
			args: args{
				data: sqlmapper.RowData{},
			},
			wantErr: true,
		},
		{
			name:      "invalid column name",
			tableName: "users",
			args: args{
				data: sqlmapper.RowData{
					"namenmce": sqlmapper.ColData{
						Data: "hieudeptrai",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewPGStore(cfg.DB(), database.Models(cfg.ModelList).GroupByName())
			got, err := s.Create(tt.tableName, tt.args.data)
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
	cfg, clearDB := utilTest.CreateConfig(t)
	defer clearDB()
	// migrate tables
	err := utilTest.MigrateTables(cfg.DB())
	if err != nil {
		t.Fatalf("Failed to migrate table by error %v", err)
	}
	//create sample data
	users, err := utilTest.CreateUserSampleData(cfg.DB())
	if err != nil {
		t.Fatalf("Failed to create sample data by error %v", err)
	}
	type args struct {
		d  sqlmapper.RowData
		id int
	}
	tests := []struct {
		name      string
		tableName string
		args      args
		want      sqlmapper.RowData
		wantErr   bool
	}{
		{
			name:      "success",
			tableName: "users",
			args: args{
				d: sqlmapper.RowData{
					"name": sqlmapper.ColData{
						Data: "demo",
					},
				},
				id: users[0].ID,
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
				d: sqlmapper.RowData{
					"name": sqlmapper.ColData{
						Data: "demo",
					},
				},
				id: -1,
			},
			wantErr: true,
		},
		{
			name:      "primary key is duplicated",
			tableName: "users",
			args: args{
				d: sqlmapper.RowData{
					"name": sqlmapper.ColData{
						Data: "demo",
					},
					"id": sqlmapper.ColData{
						Data: "1",
					},
				},
				id: users[0].ID,
			},
			want: sqlmapper.RowData{
				"name": sqlmapper.ColData{
					Data: "demo",
				},
			},
			wantErr: false,
		},
		{
			name:      "rowData is empty",
			tableName: "users",
			args: args{
				id: users[0].ID,
			},
			wantErr: true,
		},
		{
			name:      "invalid column name",
			tableName: "users",
			args: args{
				d: sqlmapper.RowData{
					"blabla": sqlmapper.ColData{
						Data: "anmt",
					},
				},
				id: users[0].ID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewPGStore(cfg.DB(), database.Models(cfg.ModelList).GroupByName())
			got, err := s.Update(tt.tableName, tt.args.d, tt.args.id)
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
