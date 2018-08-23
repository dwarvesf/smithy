package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/cobra"

	"github.com/dwarvesf/smithy/agent"
	agentConfig "github.com/dwarvesf/smithy/agent/config"
)

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func generateRandomString(s int) (string, error) {
	b, err := generateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

func main() {
	var (
		configFile     string
		configFilePath string
		forceCreate    bool
	)

	var cmdAgentMigrate = &cobra.Command{
		Use:   "agent-migrate",
		Short: "Automigrate base on mode_list in agent config file",
		Long:  `agent-migrate migrate missing columns, tables described in config file`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {

			cfg, err := agent.NewConfig(agentConfig.ReadYAML(configFile))
			if err != nil {
				log.Fatalln(err)
			}
			err = agent.AutoMigrate(cfg)
			if err != nil {
				log.Fatalln(err)
			}

			fmt.Println("Finish auto-migrate")
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

			// If command doesn't have flag, print PSK on CLI
			token, err := generateRandomString(128)
			if err != nil {
				log.Fatalln(err)
				return
			}
			if configFilePath == "" {
				fmt.Println(token)
				return
			}

			// If file doesn't exist, create a new file and write PSK into 'secrect_key'
			cfg, err := agent.NewConfig(agentConfig.ReadYAML(configFilePath))
			if err != nil {
				cfg = &agentConfig.Config{}
			}

			// If file already existed, update 'secrect_key'
			cfg.SerectKey = token
			wr := agentConfig.WriteYAML(configFilePath)
			if err := wr.Write(cfg); err != nil {
				log.Fatalln(err)
			}
		},
	}

	var cmdGenerateUser = &cobra.Command{
		Use:              "user",
		TraverseChildren: true,
		Short:            "generate user with ACL describe in model list",
		Long:             `generate user with ACL describe in model list`,
		Args:             cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := agent.NewConfig(agentConfig.ReadYAML(configFile))
			if err != nil {
				cfg = &agentConfig.Config{}
			}

			user, err := agent.CreateUserWithACL(cfg, forceCreate)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Printf("UserName: %s\n", user.Username)
			fmt.Printf("Password: %s\n", user.Password)
		},
	}

	var rootCmd = &cobra.Command{Use: "smithy"}
	rootCmd.AddCommand(cmdAgentMigrate, cmdGenerate)
	cmdGenerate.AddCommand(cmdPSK)
	cmdGenerate.AddCommand(cmdGenerateUser)

	// Set flags
	cmdAgentMigrate.Flags().StringVarP(&configFile, "config-file", "c", "example_agent_config.yaml", "put your name of config file here, with extension")
	cmdGenerateUser.Flags().StringVarP(&configFile, "config-file", "c", "example_agent_config.yaml", "put your name of config file here, with extension")
	cmdGenerateUser.Flags().BoolVarP(&forceCreate, "force-create", "f", false, "put your name of config file here, with extension")
	cmdPSK.Flags().StringVarP(&configFilePath, "config-file", "c", "", "put your name of config file here, with extension")

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
