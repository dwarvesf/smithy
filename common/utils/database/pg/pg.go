package pg

import (
	"context"
	"errors"
	"math/rand"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	backendConfig "github.com/dwarvesf/smithy/backend/config"
)

// DB info
const (
	DBHost     = "localhost"
	DBPort     = "5439"
	DBUser     = "postgres"
	DBPassword = "example"
	DBName     = "test"
)

// CreateDatabase create pg db for test
func CreateDatabase(t *testing.T, cfg *backendConfig.Config) func() {
	rand.Seed(time.Now().UnixNano())

	schemaName := "test" + strconv.FormatInt(rand.Int63(), 10)

	err := cfg.DB(DBName).Exec("CREATE SCHEMA " + schemaName).Error
	if err != nil {
		t.Fatalf("Fail to create schema. %s", err.Error())
	}

	// set schema for current db connection
	err = cfg.DB(DBName).Exec("SET search_path TO " + schemaName).Error
	if err != nil {
		t.Fatalf("Fail to set search_path to created schema. %s", err.Error())
	}

	// set schema name to config
	cfg.DBSchemaName = schemaName

	return func() {
		err := cfg.DB(DBName).Exec("DROP SCHEMA " + schemaName + " CASCADE").Error
		if err != nil {
			t.Fatalf("Fail to drop database. %s", err.Error())
		}
	}
}

// User struct for test
type User struct {
	ID   int `sql:"primary_key"`
	Name string
}

func getContainerByPort(port string) (string, error) {
	cmd := exec.Command("docker", "ps", "-a")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	outStr := string(out)
	outLine := strings.Split(outStr, "\n")
	dockerContainerName := ""
	for _, line := range outLine {
		c := strings.Contains(line, ":"+port)
		if c {
			s := strings.Split(line, " ")
			dockerContainerName = s[len(s)-1]
		}
	}

	if dockerContainerName == "" {
		return "", errors.New("Can't find docker container's name")
	}

	return dockerContainerName, nil
}

// CreatePGDatabase create a database with shell script
func CreatePGDatabase(port, databaseName string) error {
	// get container name
	dockerContainerName, err := getContainerByPort(port)
	if err != nil {
		return err
	}

	// create docker client
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	// find correct container
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return err
	}

	for _, container := range containers {
		if container.Names[0][1:] == dockerContainerName {
			execConfig := types.ExecConfig{
				Cmd:          []string{"psql", "-U", "postgres", "-c", "CREATE DATABASE " + databaseName},
				AttachStdin:  false,
				AttachStdout: true,
				AttachStderr: true,
				Tty:          false,
			}

			idResponse, err := cli.ContainerExecCreate(ctx, container.ID, execConfig)
			if err != nil {
				return err
			}

			if err := cli.ContainerExecStart(ctx, idResponse.ID, types.ExecStartCheck{Detach: false, Tty: false}); err != nil {
				return err
			}
		}
	}

	return nil
}
