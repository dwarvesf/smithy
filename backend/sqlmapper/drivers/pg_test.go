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
	cfg, clearDB := utilTest.CreateConfig(t)
	defer clearDB()

	// migrate tables
	err := utilTest.MigrateTables(cfg.DB())
	if err != nil {
		t.Fatalf("Failed to migrate table by error %v", err)
	}

	// create sample data
	users, err := utilTest.CreateUserSampleData(cfg.DB())
	if err != nil {
		t.Fatalf("Failed to create sample data by error %v", err)
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
		{
			name: "success sort by name in descending order",
			args: &sqlmapper.Query{
				SourceTable: "users",
				Fields:      []string{"id", "name"},
				Order:       []string{"name", "desc"},
			},
			want:    []string{"id", "name"},
			want1:   descUserName,
			wantErr: false,
		},
		{
			name: "success sort by name in ascending order",
			args: &sqlmapper.Query{
				SourceTable: "users",
				Fields:      []string{"id", "name"},
				Order:       []string{"name", "asc"},
			},
			want:    []string{"id", "name"},
			want1:   ascUserName,
			wantErr: false,
		},
		{
			name: "success sort by ID in descending order",
			args: &sqlmapper.Query{
				SourceTable: "users",
				Fields:      []string{"id", "name"},
				Order:       []string{"id", "desc"},
			},
			want:    []string{"id", "name"},
			want1:   descUserID,
			wantErr: false,
		},
		{
			name: "success sort by ID in ascending order",
			args: &sqlmapper.Query{
				SourceTable: "users",
				Fields:      []string{"id", "name"},
				Order:       []string{"id", "asc"},
			},
			want:    []string{"id", "name"},
			want1:   ascUserID,
			wantErr: false,
		},
		{
			name: "success for none sort",
			args: &sqlmapper.Query{
				SourceTable: "users",
				Fields:      []string{"id", "name"},
				Order:       []string{},
			},
			want:    []string{"id", "name"},
			want1:   users,
			wantErr: false,
		},
		{
			name: "fail sort because missing argument of Order",
			args: &sqlmapper.Query{
				SourceTable: "users",
				Fields:      []string{"id", "name"},
				Order:       []string{"name"},
			},
			wantErr: true,
		},
		{
			name: "fail sort because wrong argument format for Order",
			args: &sqlmapper.Query{
				SourceTable: "users",
				Fields:      []string{"id", "name"},
				Order:       []string{"name", "ascending"},
			},
			wantErr: true,
		},
		{
			name: "fail sort because column name doesn't exist",
			args: &sqlmapper.Query{
				SourceTable: "users",
				Fields:      []string{"id", "name"},
				Order:       []string{"age", "asc"},
			},
			wantErr: true,
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
		id int
	}
	tests := []struct {
		name              string
		tableName         string
		args              *args
		wantErr           bool
		testForEmptyTable bool
	}{
		{
			name:      "Valid test case",
			tableName: "users",
			args: &args{
				id: 1,
			},
		},
		{
			name:      "id not exists",
			tableName: "users",
			args: &args{
				id: 111,
			},
			wantErr: true,
		},
		{
			name:      "id too long",
			tableName: "users",
			args: &args{
				id: math.MaxInt32 + 1,
			},
			wantErr: true,
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

			err := s.Delete(tt.tableName, tt.args.id)
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
