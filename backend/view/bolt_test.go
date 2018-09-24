// +build integration

package view

import (
	"reflect"
	"testing"
	"time"

	utilDB "github.com/dwarvesf/smithy/common/utils/database/bolt"
)

func Test_boltImpl_Write(t *testing.T) {
	t.Parallel()

	persistenceFileName, clearDB := utilDB.CreateDatabase(t)
	defer clearDB()

	type args struct {
		sql *View
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "save",
			args: args{
				sql: &View{
					SQL:          "SELECT * FROM users",
					DatabaseName: "users",
					CreatedAt:    time.Now(),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBoltWriterReaderDeleter(persistenceFileName)
			if err := b.Write(tt.args.sql); (err != nil) != tt.wantErr {
				t.Errorf("boltImpl.Write() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_boltImpl_Read(t *testing.T) {
	t.Parallel()

	persistenceFileName, clearDB := utilDB.CreateDatabase(t)
	defer clearDB()

	sqlCMD := &View{
		SQL:          "SELECT * FROM users",
		DatabaseName: "users",
		CreatedAt:    time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
	}

	b := NewBoltWriterReaderDeleter(persistenceFileName)
	err := b.Write(sqlCMD)
	if err != nil {
		t.Fatalf("Fail to create sample data %v", err)
	}

	type args struct {
		sqlID int
	}
	tests := []struct {
		name    string
		args    args
		want    *View
		wantErr bool
	}{
		{
			name: "read exist id",
			args: args{
				sqlID: sqlCMD.ID,
			},
			want:    sqlCMD,
			wantErr: false,
		},
		{
			name: "read no-exist id",
			args: args{
				sqlID: 999,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := b.Read(tt.args.sqlID)
			if (err != nil) != tt.wantErr {
				t.Errorf("boltImpl.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("boltImpl.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_boltImpl_Delete(t *testing.T) {
	t.Parallel()

	persistenceFileName, clearDB := utilDB.CreateDatabase(t)
	defer clearDB()

	sqlCMD := &View{
		SQL:          "SELECT * FROM users",
		DatabaseName: "users",
		CreatedAt:    time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
	}

	b := NewBoltWriterReaderDeleter(persistenceFileName)
	err := b.Write(sqlCMD)
	if err != nil {
		t.Fatalf("Fail to create sample data %v", err)
	}

	type args struct {
		sqlID int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "delete exist id",
			args: args{
				sqlID: sqlCMD.ID,
			},
			wantErr: false,
		},
		{
			name: "delete no-exist id",
			args: args{
				sqlID: 999,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := b.Delete(tt.args.sqlID)
			if err != nil {
				t.Fatalf("boltImpl.Delete() error = %v", err)
			}

			sql, err := b.Read(tt.args.sqlID)
			if !(err != nil || sql == nil) {
				if tt.wantErr {
					t.Errorf("boltImpl.Delete() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func Test_boltImpl_ListCommands(t *testing.T) {
	t.Parallel()

	persistenceFileName, clearDB := utilDB.CreateDatabase(t)
	defer clearDB()

	listSQLCmd := []*View{
		{
			SQL:          "SELECT * FROM users",
			DatabaseName: "users",
			CreatedAt:    time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
		},
		{
			SQL:          "SELECT name FROM users WHERE id = 1",
			DatabaseName: "users",
			CreatedAt:    time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
		},
	}

	b := NewBoltWriterReaderDeleter(persistenceFileName)
	for _, sqlCMD := range listSQLCmd {
		err := b.Write(sqlCMD)
		if err != nil {
			t.Fatalf("Fail to create sample data %v", err)
		}
	}

	tests := []struct {
		name    string
		want    []*View
		wantErr bool
	}{
		{
			name: "list all view",
			want: listSQLCmd,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := b.ListCommands()
			if (err != nil) != tt.wantErr {
				t.Errorf("boltImpl.ListCommands() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("boltImpl.ListCommands() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_boltImpl_ListCommandsByDBName(t *testing.T) {
	t.Parallel()

	persistenceFileName, clearDB := utilDB.CreateDatabase(t)
	defer clearDB()

	listSQLCmd1 := []*View{
		{
			SQL:          "SELECT * FROM users",
			DatabaseName: "test1",
			CreatedAt:    time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
		},
		{
			SQL:          "SELECT name FROM users WHERE id = 1",
			DatabaseName: "test1",
			CreatedAt:    time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
		},
	}
	listSQLCmd2 := []*View{
		{
			SQL:          "SELECT * FROM books",
			DatabaseName: "test2",
			CreatedAt:    time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
		},
		{
			SQL:          "SELECT name FROM books WHERE id = 1",
			DatabaseName: "test2",
			CreatedAt:    time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
		},
	}

	b := NewBoltWriterReaderDeleter(persistenceFileName)
	for _, sqlCMD := range listSQLCmd1 {
		err := b.Write(sqlCMD)
		if err != nil {
			t.Fatalf("Fail to create sample data %v", err)
		}
	}
	for _, sqlCMD := range listSQLCmd2 {
		err := b.Write(sqlCMD)
		if err != nil {
			t.Fatalf("Fail to create sample data %v", err)
		}
	}

	type args struct {
		databaseName string
	}
	tests := []struct {
		name    string
		args    args
		want    []*View
		wantErr bool
	}{
		{
			name: "list all view in a db test1",
			args: args{
				databaseName: "test1",
			},
			want: listSQLCmd1,
		},
		{
			name: "list all view in a db test2",
			args: args{
				databaseName: "test2",
			},
			want: listSQLCmd2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := b.ListCommandsByDBName(tt.args.databaseName)
			if (err != nil) != tt.wantErr {
				t.Errorf("boltImpl.ListCommandsByDBName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("boltImpl.ListCommandsByDBName() = %v, want %v", got, tt.want)
			}
		})
	}
}
