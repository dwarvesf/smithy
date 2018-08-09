package main

import (
	"log"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/cobra"

	"github.com/dwarvesf/smithy/agent"
	"github.com/dwarvesf/smithy/agent/automigrate"
	agentConfig "github.com/dwarvesf/smithy/agent/config"
)

func main() {
	// TODO: remove static config file
	cfg, err := agent.NewConfig(agentConfig.NewYAMLConfigReader("example_agent_config.yaml"))
	if err != nil {
		panic(err)
	}

	var cmdAgentMigrate = &cobra.Command{
		Use:   "agent-migrate",
		Short: "Automigrate base on mode_list in agent config file",
		Long:  `agent-migrate migrate missing columns, tables described in config file`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			err = automigrate.AutoMigrate(cfg)
			if err != nil {
				log.Fatalln(err)
			}
		},
	}

	var rootCmd = &cobra.Command{Use: "smithy"}
	rootCmd.AddCommand(cmdAgentMigrate)
	rootCmd.Execute()
}
