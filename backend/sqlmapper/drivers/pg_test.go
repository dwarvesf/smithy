// +build integration

package drivers

import (
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

func Test_pgStore_FindByID(t *testing.T) {
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
		id int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    sqlmapper.RowData
		wantErr bool
	}{
		{
			name:      "Valid test case",
			tableName: "users",
			args: &args{
				id: 1,
			},
			want: users,
		},
		{
			name:      "id not exists",
			tableName: "users",
			args: &args{
				id: 111
			},
			wantErr: true,
		}, 
		{
			name: "id too long",
			tableName: "users",
			args: &args {
				id: 11111111111111111111111111111111111111111,
			}
			wantErr: true
		},
		{
			name: "id empty",
			tableName: "users",
			args: &args{
				id: "",
			}
			wantErr: true
		}
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

			got, err := s.FindByColumnName(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("pgStore.FindByID() error = %v, wantErr %v", err, tt.wantErr)
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
					t.Errorf("pgStore.FindByID() = %v, want %v", got, tt.want)
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
		id int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:      "Valid test case",
			tableName: "users",
			args: &args{
				id: 1,
			},
			want: users,
		},
		{
			name:      "id not exists",
			tableName: "users",
			args: &args{
				id: 111
			},
			wantErr: true,
		}, 
		{
			name: "id too long",
			tableName: "users",
			args: &args {
				id: 11111111111111111111111111111111111111111,
			}
			wantErr: true
		},
		{
			name: "id empty",
			tableName: "users",
			args: &args{
				id: "",
			}
			wantErr: true
		}
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

			got, err := s.FindByColumnName(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("pgStore.Delete() error = %v, wantErr %v", err, tt.wantErr)
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
					t.Errorf("pgStore.Delete() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}