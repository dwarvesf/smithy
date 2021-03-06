// +build integration

package hook

import (
	"reflect"
	"testing"

	utilDB "github.com/dwarvesf/smithy/common/utils/database/pg"
	utilTest "github.com/dwarvesf/smithy/common/utils/database/pg/test/set1"
)

func Test_pgLibImpl_First(t *testing.T) {
	t.Parallel()
	cfg, clearDB := utilTest.CreateConfig(t)
	defer clearDB()

	for _, dbase := range cfg.Databases {
		// migrate tables
		err := utilTest.MigrateTables(cfg.DB(dbase.DBName))
		if err != nil {
			t.Fatalf("Failed to migrate table by error %v", err)
		}

		//create sample data
		_, err = utilTest.CreateUserSampleData(cfg.DB(dbase.DBName))
		if err != nil {
			t.Fatalf("Failed to create sample data by error %v", err)
		}
	}

	type args struct {
		databaseName string
		tableName    string
		condition    string
	}
	tests := []struct {
		name    string
		args    args
		want    map[interface{}]interface{}
		wantErr bool
	}{
		{
			name: "correct",
			args: args{
				databaseName: "test1",
				tableName:    "users",
				condition:    "id = 1",
			},
			want: map[interface{}]interface{}{
				"id":   int64(1),
				"name": "hieudeptrai0",
			},
			wantErr: false,
		},
		{
			name: "not found record",
			args: args{
				databaseName: "test1",
				tableName:    "users",
				condition:    "id = 100000",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "table not exist",
			args: args{
				databaseName: "test1",
				tableName:    "userss",
				condition:    "id = 1",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "wrong condittion",
			args: args{
				databaseName: "test1",
				tableName:    "users",
				condition:    "wrong_column = 1",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewPGLib(cfg.DBs(), cfg.ModelMap)

			got, err := s.First(tt.args.databaseName, tt.args.tableName, tt.args.condition)
			if (err != nil) != tt.wantErr {
				t.Errorf("pgLibImpl.First() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pgLibImpl.First() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pgLibImpl_Where(t *testing.T) {
	t.Parallel()
	cfg, clearDB := utilTest.CreateConfig(t)
	defer clearDB()

	for _, dbase := range cfg.Databases {
		// migrate tables
		err := utilTest.MigrateTables(cfg.DB(dbase.DBName))
		if err != nil {
			t.Fatalf("Failed to migrate table by error %v", err)
		}

		//create sample data
		_, err = utilTest.CreateUserSampleData(cfg.DB(dbase.DBName))
		if err != nil {
			t.Fatalf("Failed to create sample data by error %v", err)
		}
	}

	type args struct {
		databaseName string
		tableName    string
		condition    string
	}
	tests := []struct {
		name    string
		args    args
		want    []map[interface{}]interface{}
		wantErr bool
	}{
		{
			name: "correct",
			args: args{
				databaseName: "test1",
				tableName:    "users",
				condition:    "id = 1",
			},
			want: []map[interface{}]interface{}{
				{
					"id":   int64(1),
					"name": "hieudeptrai0",
				},
			},
			wantErr: false,
		},
		{
			name: "not found record",
			args: args{
				databaseName: "test1",
				tableName:    "users",
				condition:    "id = 100000",
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "table not exist",
			args: args{
				databaseName: "test1",
				tableName:    "userss",
				condition:    "id = 1",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "wrong condittion",
			args: args{
				databaseName: "test1",
				tableName:    "users",
				condition:    "wrong_column = 1",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewPGLib(cfg.DBs(), cfg.ModelMap)

			got, err := s.Where(tt.args.databaseName, tt.args.tableName, tt.args.condition)
			if (err != nil) != tt.wantErr {
				t.Errorf("pgLibImpl.Where() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pgLibImpl.Where() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pgLibImpl_Create(t *testing.T) {
	t.Parallel()
	cfg, clearDB := utilTest.CreateConfig(t)
	defer clearDB()

	for _, dbase := range cfg.Databases {
		// migrate tables
		err := utilTest.MigrateTables(cfg.DB(dbase.DBName))
		if err != nil {
			t.Fatalf("Failed to migrate table by error %v", err)
		}
	}

	type args struct {
		databaseName string
		tableName    string
		d            map[interface{}]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    map[interface{}]interface{}
		wantErr bool
	}{
		{
			name: "correct",
			args: args{
				databaseName: "test1",
				tableName:    "users",
				d: map[interface{}]interface{}{
					"name": "a_user_name",
				},
			},
			want: map[interface{}]interface{}{
				"id":   int64(1),
				"name": "a_user_name",
			},
			wantErr: false,
		},
		{
			name: "table not exist",
			args: args{
				databaseName: "test1",
				tableName:    "userss",
				d: map[interface{}]interface{}{
					"name": "a_user_name",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "wrong column",
			args: args{
				databaseName: "test1",
				tableName:    "users",
				d: map[interface{}]interface{}{
					"wrong_column": "a_user_name",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewPGLib(cfg.DBs(), cfg.ModelMap)

			got, err := s.Create(tt.args.databaseName, tt.args.tableName, tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("pgLibImpl.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pgLibImpl.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pgLibImpl_Update(t *testing.T) {
	t.Parallel()
	cfg, clearDB := utilTest.CreateConfig(t)
	defer clearDB()

	for _, dbase := range cfg.Databases {
		// migrate tables
		err := utilTest.MigrateTables(cfg.DB(dbase.DBName))
		if err != nil {
			t.Fatalf("Failed to migrate table by error %v", err)
		}

		//create sample data
		_, err = utilTest.CreateUserSampleData(cfg.DB(dbase.DBName))
		if err != nil {
			t.Fatalf("Failed to create sample data by error %v", err)
		}
	}

	type args struct {
		databaseName string
		tableName    string
		d            map[interface{}]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    map[interface{}]interface{}
		wantErr bool
	}{
		{
			name: "correct",
			args: args{
				databaseName: "test1",
				tableName:    "users",
				d: map[interface{}]interface{}{
					"id":   1,
					"name": "changed user name",
				},
			},
			want: map[interface{}]interface{}{
				"name": "changed user name",
			},
			wantErr: false,
		},
		{
			name: "primary key is not exist",
			args: args{
				databaseName: "test1",
				tableName:    "users",
				d: map[interface{}]interface{}{
					"id":   10000,
					"name": "changed user name",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "table not exist",
			args: args{
				databaseName: "test1",
				tableName:    "userss",
				d: map[interface{}]interface{}{
					"name": "a_user_name",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "wrong column",
			args: args{
				databaseName: "test1",
				tableName:    "users",
				d: map[interface{}]interface{}{
					"wrong_column": "a_user_name",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewPGLib(cfg.DBs(), cfg.ModelMap)
			got, err := s.Update(tt.args.databaseName, tt.args.tableName, tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("pgLibImpl.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pgLibImpl.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pgLibImpl_Delete(t *testing.T) {
	t.Parallel()
	cfg, clearDB := utilTest.CreateConfig(t)
	defer clearDB()

	for _, dbase := range cfg.Databases {
		// migrate tables
		err := utilTest.MigrateTables(cfg.DB(dbase.DBName))
		if err != nil {
			t.Fatalf("Failed to migrate table by error %v", err)
		}

		//create sample data
		_, err = utilTest.CreateUserSampleData(cfg.DB(dbase.DBName))
		if err != nil {
			t.Fatalf("Failed to create sample data by error %v", err)
		}
	}

	type args struct {
		databaseName string
		tableName    string
		fields       []interface{}
		data         []interface{}
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		testCorrect bool
	}{
		{
			name: "correct",
			args: args{
				databaseName: "test1",
				tableName:    "users",
				fields: []interface{}{
					"id",
				},
				data: []interface{}{
					"1",
				},
			},
			wantErr:     false,
			testCorrect: true,
		},
		{
			name: "id and name exist",
			args: args{
				databaseName: "test1",
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
			wantErr: false,
		},
		{
			name: "table not exist",
			args: args{
				databaseName: "test1",
				tableName:    "userss",
				fields: []interface{}{
					"id",
				},
				data: []interface{}{
					"1",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewPGLib(cfg.DBs(), cfg.ModelMap)

			if err := s.Delete(tt.args.databaseName, tt.args.tableName, tt.args.fields, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("pgLibImpl.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.testCorrect {
				_, err := s.First(utilDB.DBName, "users", "id = 1") // check record users already deleted in database
				if err == nil {
					t.Error("pgLibImpl.Delete() not delete record in database ")
				}
			}
		})
	}
}
