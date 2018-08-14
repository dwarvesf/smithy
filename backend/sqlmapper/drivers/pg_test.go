package drivers

import (
	"fmt"
	"strings"
	"testing"

	backendConfig "github.com/dwarvesf/smithy/backend/config"
	"github.com/dwarvesf/smithy/backend/sqlmapper"
	"github.com/dwarvesf/smithy/common/database"
	"github.com/dwarvesf/smithy/common/utils"
)

func Test_pgStore_FindAll(t *testing.T) {
	// read config from agent
	cfg, err := backendConfig.ReadYAML("../../../example_dashboard_config.yaml").Read()
	if err != nil {
		t.Fatalf("Can't read config file %v", err)
	}

	err = cfg.UpdateConfigFromAgent()
	if err != nil {
		t.Fatalf("Can't sync from agent %v", err)
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
			iId, err := utils.ConvertFromInterfacePtr(got[0]["id"].Data)
			if err != nil {
				t.Fatal(err)
			}
			id := iId.(int)
			iName, err := utils.ConvertFromInterfacePtr(got[0]["name"].Data)
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
