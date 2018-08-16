// +build integration

package drivers

import (
	"reflect"
	"testing"

	"github.com/dwarvesf/smithy/backend/sqlmapper"
	"github.com/dwarvesf/smithy/common/database"
	utilPg "github.com/dwarvesf/smithy/common/utils/database/pg"
	utilTest "github.com/dwarvesf/smithy/common/utils/database/pg/test/set1"
	utilReflect "github.com/dwarvesf/smithy/common/utils/reflect"
)

/**
 * use for integration test with AGENT later
func Test_SyncAgent(t *testing.T) {
	// read config from agent
	cfg, err := backendConfig.ReadYAML("../../../example_dashboard_config.yaml").Read()
	if err != nil {
		t.Errorf("Can't read config file %v", err)
	}

	err = cfg.UpdateConfigFromAgent()
	if err != nil {
		t.Errorf("Can't sync from agent %v", err)
	}
}
*/

func Test_pgStore_FindAll(t *testing.T) {
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

	cols := []database.Column{
		{
			Name: "id",
			Type: "int",
		},
		{
			Name: "name",
			Type: "string",
		},
	}

	type args struct {
		Offset int
		Limit  int
	}

	tests := []struct {
		name              string
		tableName         string
		args              *args
		want              []utilPg.User
		wantErr           bool
		testForEmptyTable bool
	}{
		{
			name:      "Valid test case",
			tableName: "users",
			args: &args{
				Offset: 0,
				Limit:  0,
			},
			want: users,
		},
		{
			name:      "empty table",
			tableName: "users",
			args: &args{
				Offset: 0,
				Limit:  0,
			},
			want:              []utilPg.User{},
			testForEmptyTable: true,
		},
		{
			name:      "offset = 2 limit = 10",
			tableName: "users",
			args: &args{
				Offset: 2,
				Limit:  10,
			},
			want: users[2:12],
		},
		{
			name:      "offset = 2 limit = 0",
			tableName: "users",
			args: &args{
				Offset: 2,
				Limit:  0,
			},
			want: users[2:],
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
				s = NewPGStore(cfgEmpty.DB(), tt.tableName, cols, cfgEmpty.ModelList)
			} else {
				s = NewPGStore(cfg.DB(), tt.tableName, cols, cfg.ModelList)
			}

			got, err := s.FindAll(tt.args.Offset, tt.args.Limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("pgStore.FindAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("len(got)=%v != len(tt.want)=%v", len(got), len(tt.want))
				return
			}

			for i := 0; i < len(got); i++ {
				// convert data
				iId, err := utilReflect.ConvertFromInterfacePtr(got[i]["id"].Data)
				if err != nil {
					t.Fatal(err)
				}
				id := iId.(int)
				iName, err := utilReflect.ConvertFromInterfacePtr(got[i]["name"].Data)
				if err != nil {
					t.Fatal(err)
				}
				name := iName.(string)

				if len(got) != len(tt.want) ||
					id != tt.want[i].Id ||
					name != tt.want[i].Name {
					t.Errorf("pgStore.FindByColumnName() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_pgStore_FindByColumnName(t *testing.T) {
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

	cols := []database.Column{
		{
			Name: "id",
			Type: "int",
		},
		{
			Name: "name",
			Type: "string",
		},
	}

	type args struct {
		ColumnName string
		Value      string
		Offset     int
		Limit      int
	}

	tests := []struct {
		name              string
		tableName         string
		args              *args
		want              []utilPg.User
		wantErr           bool
		testForEmptyTable bool
	}{
		{
			name:      "Valid test case",
			tableName: "users",
			args: &args{
				ColumnName: "name",
				Value:      "hieudeptrai",
				Offset:     0,
				Limit:      0,
			},
			want: users,
		},
		{
			name:      "empty table",
			tableName: "users",
			args: &args{
				ColumnName: "name",
				Value:      "hieudeptrai",
				Offset:     0,
				Limit:      0,
			},
			want:              []utilPg.User{},
			testForEmptyTable: true,
		},
		{
			name:      "offset = 2 limit = 10",
			tableName: "users",
			args: &args{
				ColumnName: "name",
				Value:      "hieudeptrai",
				Offset:     2,
				Limit:      10,
			},
			want: users[2:12],
		},
		{
			name:      "offset = 2 limit = 0",
			tableName: "users",
			args: &args{
				ColumnName: "name",
				Value:      "hieudeptrai",
				Offset:     2,
				Limit:      0,
			},
			want: users[2:],
		},
		{
			name:      "invalid column name",
			tableName: "users",
			args: &args{
				ColumnName: "namexxx",
				Value:      "hieudeptrai",
				Offset:     2,
				Limit:      0,
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
				s = NewPGStore(cfgEmpty.DB(), tt.tableName, cols, cfgEmpty.ModelList)
			} else {
				s = NewPGStore(cfg.DB(), tt.tableName, cols, cfg.ModelList)
			}

			got, err := s.FindByColumnName(tt.args.ColumnName, tt.args.Value, tt.args.Offset, tt.args.Limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("pgStore.FindByColumnName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("len(got)=%v != len(tt.want)=%v", len(got), len(tt.want))
				return
			}

			for i := 0; i < len(got); i++ {
				// convert data
				iId, err := utilReflect.ConvertFromInterfacePtr(got[i]["id"].Data)
				if err != nil {
					t.Fatal(err)
				}
				id := iId.(int)
				iName, err := utilReflect.ConvertFromInterfacePtr(got[i]["name"].Data)
				if err != nil {
					t.Fatal(err)
				}
				name := iName.(string)

				if len(got) != len(tt.want) ||
					id != tt.want[i].Id ||
					name != tt.want[i].Name {
					t.Errorf("pgStore.FindByColumnName() = %v, want %v", got, tt.want)
				}
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
				id: users[0].Id,
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
				id: users[0].Id,
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
				id: users[0].Id,
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
				id: users[0].Id,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewPGStore(cfg.DB(), tt.tableName, []database.Column{}, cfg.ModelList)
			got, err := s.Update(tt.args.d, tt.args.id)
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
