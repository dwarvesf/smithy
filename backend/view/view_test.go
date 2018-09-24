// +build integration

package view

import (
	"testing"
	"time"

	"github.com/dwarvesf/smithy/backend/sqlmapper"
	"github.com/dwarvesf/smithy/backend/sqlmapper/drivers"
	utilTest "github.com/dwarvesf/smithy/common/utils/database/pg/test/set1"
)

func TestView_Validate(t *testing.T) {
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

	s := drivers.NewPGStore(cfg.DBs(), cfg.ModelMap)

	type fields struct {
		ID           int
		SQL          string
		DatabaseName string
		CreatedAt    time.Time
	}
	type args struct {
		m sqlmapper.Mapper
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid view",
			fields: fields{
				DatabaseName: "test1",
				SQL:          "SELECT * FROM users",
			},
			args:    args{m: s},
			wantErr: false,
		},
		{
			name: "view with no-exist table",
			fields: fields{
				DatabaseName: "test1",
				SQL:          "SELECT * FROM userss",
			},
			args:    args{m: s},
			wantErr: true,
		},
		{
			name: "view with no-exist column",
			fields: fields{
				DatabaseName: "test1",
				SQL:          "SELECT namee FROM users",
			},
			args:    args{m: s},
			wantErr: true,
		},
		{
			name: "view with INSERT",
			fields: fields{
				DatabaseName: "test1",
				SQL:          "INSERT INTO users VALUES(11, 'Hieu Dep Trai')",
			},
			args:    args{m: s},
			wantErr: true,
		},
		{
			name: "view with UPDATE",
			fields: fields{
				DatabaseName: "test1",
				SQL:          "UPDATE users WHERE id=1",
			},
			args:    args{m: s},
			wantErr: true,
		},
		{
			name: "view with DELETE",
			fields: fields{
				DatabaseName: "test1",
				SQL:          "DELETE FROM users WHERE id = 1",
			},
			args:    args{m: s},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &View{
				ID:           tt.fields.ID,
				SQL:          tt.fields.SQL,
				DatabaseName: tt.fields.DatabaseName,
				CreatedAt:    tt.fields.CreatedAt,
			}
			_, err := s.Validate(tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("View.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
