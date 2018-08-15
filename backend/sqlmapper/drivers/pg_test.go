// +build integration

package drivers

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/dwarvesf/smithy/backend/sqlmapper"
	"github.com/dwarvesf/smithy/common/database"
	utilTest "github.com/dwarvesf/smithy/common/utils/database/pg/test"
	utilReflect "github.com/dwarvesf/smithy/common/utils/reflect"
<<<<<<< HEAD
=======
	utilTest "github.com/dwarvesf/smithy/common/utils/test"
	"github.com/jinzhu/gorm"
>>>>>>> Unit test
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
	sampleData := make(map[string]sqlmapper.ColData)
	sampleData["id"] = sqlmapper.ColData{
		Data:     999,
		DataType: "int",
	}
	sampleData["name"] = sqlmapper.ColData{
		Data:     "hieudeptrai",
		DataType: "string",
	}

	tableName := "users"
	cfg.DB().Exec(fmt.Sprintf("DELETE FROM %s;", tableName))
	execQuery := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING id;",
		tableName,
		strings.Join([]string{"id", "name"}, ","),
		"?,?")
	cfg.DB().Exec(execQuery, sampleData["id"].Data, sampleData["name"].Data)

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
		name    string
		args    *args
		want    []sqlmapper.RowData
		wantErr bool
	}{
		{
			name: "Valid test case",
			args: &args{
				Offset: 0,
				Limit:  10,
			},
			want: []sqlmapper.RowData{sampleData},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewPGStore(cfg.DB(), tableName, cols, cfg.ModelList)
			got, err := s.FindAll(tt.args.Offset, tt.args.Limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("pgStore.FindAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// convert data
			iId, err := utilReflect.ConvertFromInterfacePtr(got[0]["id"].Data)
			if err != nil {
				t.Fatal(err)
			}
			id := iId.(int)
			iName, err := utilReflect.ConvertFromInterfacePtr(got[0]["name"].Data)
			if err != nil {
				t.Fatal(err)
			}
			name := iName.(string)

			if len(got) != len(tt.want) ||
				id != tt.want[0]["id"].Data ||
				name != tt.want[0]["name"].Data {
				t.Errorf("pgStore.FindAll() = %v, want %v", got, tt.want)
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
	sampleData := make(map[string]sqlmapper.ColData)
	sampleData["id"] = sqlmapper.ColData{
		Data:     999,
		DataType: "int",
	}
	sampleData["name"] = sqlmapper.ColData{
		Data:     "hieudeptrai",
		DataType: "string",
	}

	tableName := "users"
	cfg.DB().Exec(fmt.Sprintf("DELETE FROM %s;", tableName))
	execQuery := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING id;",
		tableName,
		strings.Join([]string{"id", "name"}, ","),
		"?,?")
	cfg.DB().Exec(execQuery, sampleData["id"].Data, sampleData["name"].Data)

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
		name    string
		args    *args
		want    []sqlmapper.RowData
		wantErr bool
	}{
		{
			name: "Valid test case",
			args: &args{
				ColumnName: "name",
				Value:      "hieudeptrai",
				Offset:     0,
				Limit:      10,
			},
			want: []sqlmapper.RowData{sampleData},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewPGStore(cfg.DB(), tableName, cols, cfg.ModelList)
			got, err := s.FindByColumnName(tt.args.ColumnName, tt.args.Value, tt.args.Offset, tt.args.Limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("pgStore.FindByColumnName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// convert data
			iId, err := utilReflect.ConvertFromInterfacePtr(got[0]["id"].Data)
			if err != nil {
				t.Fatal(err)
			}
			id := iId.(int)
			iName, err := utilReflect.ConvertFromInterfacePtr(got[0]["name"].Data)
			if err != nil {
				t.Fatal(err)
			}
			name := iName.(string)

			if len(got) != len(tt.want) ||
				id != tt.want[0]["id"].Data ||
				name != tt.want[0]["name"].Data {
				t.Errorf("pgStore.FindByColumnName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pgStore_FindByID(t *testing.T) {
	t.Parallel()
	cfg, clearDB := utilTest.CreateConfig(t)
	defer clearDB(
	
	// migrate tables
	err := utilTest.MigrateTables(cfg.DB())
	if err != nil {
		t.Fatalf("Failed to migrate table by error %v", err)
	}

	//create sample data
	sampleData := make(map[string]sqlmapper.ColData)
	sampleData["id"] = sqlmapper.ColData{
		Data:     123,
		DataType: "int",
	}
	sampleData["name"] = sqlmapper.ColData{
		Data:     "minh dep trai vl",
		DataType: "string",
	}

	tableName := "users"
	cfg.DB().Exec(fmt.Sprintf("DELETE FROM %s;", tableName))
	execQuery := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING id;",
		tableName,
		strings.Join([]string{"id", "name"}, ","),
		"?,?")
	cfg.DB().Exec(execQuery, sampleData["id"].Data, sampleData["name"].Data)

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
			name: "Valid test case",
			args: &args{
				id: 123
			},
			want: []sqlmapper.RowData{sampleData},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewPGStore(cfg.DB(), tableName, cols, cfg.ModelList)
			got, err := s.FindByColumnName(tt.args.ColumnName, tt.args.Value, tt.args.Offset, tt.args.Limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("pgStore.FindByColumnName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// convert data
			iId, err := utilReflect.ConvertFromInterfacePtr(got[0]["id"].Data)
			if err != nil {
				t.Fatal(err)
			}
			id := iId.(int)
			iName, err := utilReflect.ConvertFromInterfacePtr(got[0]["name"].Data)
			if err != nil {
				t.Fatal(err)
			}
			name := iName.(string)

			if len(got) != len(tt.want) ||
				id != tt.want[0]["id"].Data ||
				name != tt.want[0]["name"].Data {
				t.Errorf("pgStore.FindByID() = %v, want %v", got, tt.want)
			}
		})
	}
}
