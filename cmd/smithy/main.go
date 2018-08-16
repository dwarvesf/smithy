package main

import (
	"fmt"
	"log"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/cobra"

	"github.com/dwarvesf/smithy/agent"
	agentConfig "github.com/dwarvesf/smithy/agent/config"

	"crypto/rand"
	"encoding/base64"
)

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

func main() {
	cfg, err := agent.NewConfig(agentConfig.ReadYAML("example_agent_config.yaml"))
	if err != nil {
		panic(err)
	}
	var configFile string

	var cmdAgentMigrate = &cobra.Command{
		Use:   "agent-migrate",
		Short: "Automigrate base on mode_list in agent config file",
		Long:  `agent-migrate migrate missing columns, tables described in config file`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			err = agent.AutoMigrate(cfg)
			if err != nil {
				log.Fatalln(err)
			}
		},
	}

	var cmdGenerate = &cobra.Command{
		Use:   "generate",
		Short: "Generate",
		Long:  `generate use to generate things, such as PSK use to authenticate with app`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Generate: ")
		},
	}

	var cmdPSK = &cobra.Command{
		Use:              "psk",
		TraverseChildren: true,
		Short:            "Generate PSK for authenticate with app",
		Long:             `generate use to generate PSK use to authenticate with app`,
		Args:             cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			token, err := GenerateRandomString(128)
			if err != nil {
				return
			}

			cfg, err := agent.NewConfig(agentConfig.ReadYAML(configFile))
			if err != nil {
				cfg = &agentConfig.Config{}
			}

			cfg.SerectKey = token

			wr := agentConfig.WriteYAML(configFile)
			wr.Write(cfg)

			if configFile == "" {
				fmt.Println(token)
			}
		},
	}

	var rootCmd = &cobra.Command{Use: "smithy"}
	rootCmd.AddCommand(cmdAgentMigrate, cmdGenerate)
	cmdGenerate.AddCommand(cmdPSK)
	cmdPSK.Flags().StringVarP(&configFile, "config-file", "f", "", "put your name of config file here, with extension")
	rootCmd.Execute()
}
