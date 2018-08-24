// Package main defines a command line interface for the sqlboiler package
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/boilingcore"
	"github.com/volatiletech/sqlboiler/drivers"
	"github.com/volatiletech/sqlboiler/importers"
)

//go:generate go-bindata -nometadata -pkg templatebin -o templatebin/bindata.go templates templates/singleton templates_test templates_test/singleton

const sqlBoilerVersion = "3.0.0"

var (
	flagConfigFile string
	cmdState       *boilingcore.State
	cmdConfig      *boilingcore.Config
)

func initConfig() {
	if len(flagConfigFile) != 0 {
		viper.SetConfigFile(flagConfigFile)
		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("Can't read config:", err)
			os.Exit(1)
		}
		return
	}

	var err error
	viper.SetConfigName("sqlboiler")

	configHome := os.Getenv("XDG_CONFIG_HOME")
	homePath := os.Getenv("HOME")
	wd, err := os.Getwd()
	if err != nil {
		wd = "."
	}

	configPaths := []string{wd}
	if len(configHome) > 0 {
		configPaths = append(configPaths, filepath.Join(configHome, "sqlboiler"))
	} else {
		configPaths = append(configPaths, filepath.Join(homePath, ".config/sqlboiler"))
	}

	for _, p := range configPaths {
		viper.AddConfigPath(p)
	}

	// Ignore errors here, fallback to other validation methods.
	// Users can use environment variables if a config is not found.
	_ = viper.ReadInConfig()
}

func main() {
	// Too much happens between here and cobra's argument handling, for
	// something so simple just do it immediately.
	for _, arg := range os.Args {
		if arg == "--version" {
			fmt.Println("SQLBoiler v" + sqlBoilerVersion)
			return
		}
	}

	// Set up the cobra root command
	var rootCmd = &cobra.Command{
		Use:   "sqlboiler [flags] <driver>",
		Short: "SQL Boiler generates an ORM tailored to your database schema.",
		Long: "SQL Boiler generates a Go ORM from template files, tailored to your database schema.\n" +
			`Complete documentation is available at http://github.com/volatiletech/sqlboiler`,
		Example:       `sqlboiler psql`,
		PreRunE:       preRun,
		RunE:          run,
		PostRunE:      postRun,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cobra.OnInitialize(initConfig)

	// Set up the cobra root command flags
	rootCmd.PersistentFlags().StringVarP(&flagConfigFile, "config", "c", "", "Filename of config file to override default lookup")
	rootCmd.PersistentFlags().StringP("output", "o", "models", "The name of the folder to output to")
	rootCmd.PersistentFlags().StringP("pkgname", "p", "models", "The name you wish to assign to your generated package")
	rootCmd.PersistentFlags().StringSliceP("templates", "", nil, "A templates directory, overrides the bindata'd template folders in sqlboiler")
	rootCmd.PersistentFlags().StringSliceP("tag", "t", nil, "Struct tags to be included on your models in addition to json, yaml, toml")
	rootCmd.PersistentFlags().StringSliceP("replace", "", nil, "Replace templates by directory: relpath/to_file.tpl:relpath/to_replacement.tpl")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Debug mode prints stack traces on error")
	rootCmd.PersistentFlags().BoolP("no-context", "", false, "Disable context.Context usage in the generated code")
	rootCmd.PersistentFlags().BoolP("no-tests", "", false, "Disable generated go test files")
	rootCmd.PersistentFlags().BoolP("no-hooks", "", false, "Disable hooks feature for your models")
	rootCmd.PersistentFlags().BoolP("no-rows-affected", "", false, "Disable rows affected in the generated API")
	rootCmd.PersistentFlags().BoolP("no-auto-timestamps", "", false, "Disable automatic timestamps for created_at/updated_at")
	rootCmd.PersistentFlags().BoolP("add-global-variants", "", false, "Enable generation for global variants")
	rootCmd.PersistentFlags().BoolP("add-panic-variants", "", false, "Enable generation for panic variants")
	rootCmd.PersistentFlags().BoolP("version", "", false, "Print the version")
	rootCmd.PersistentFlags().BoolP("wipe", "", false, "Delete the output folder (rm -rf) before generation to ensure sanity")
	rootCmd.PersistentFlags().StringP("struct-tag-casing", "", "snake", "Decides the casing for go structure tag names. camel or snake (default snake)")

	// hide flags not recommended for use
	rootCmd.PersistentFlags().MarkHidden("replace")

	viper.BindPFlags(rootCmd.PersistentFlags())
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := rootCmd.Execute(); err != nil {
		if e, ok := err.(commandFailure); ok {
			fmt.Printf("Error: %v\n\n", string(e))
			rootCmd.Help()
		} else if !viper.GetBool("debug") {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("Error: %+v\n", err)
		}

		os.Exit(1)
	}
}

type commandFailure string

func (c commandFailure) Error() string {
	return string(c)
}

func preRun(cmd *cobra.Command, args []string) error {
	var err error

	if len(args) == 0 {
		return commandFailure("must provide a driver name")
	}

	driverName := args[0]
	driverPath := args[0]

	if strings.ContainsRune(driverName, os.PathSeparator) {
		driverName = strings.Replace(filepath.Base(driverName), "sqlboiler-", "", 1)
		driverName = strings.Replace(driverName, ".exe", "", 1)
	} else {
		driverPath = "sqlboiler-" + driverPath
		if p, err := exec.LookPath(driverPath); err == nil {
			driverPath = p
		}
	}

	driverPath, err = filepath.Abs(driverPath)
	if err != nil {
		return errors.Wrap(err, "could not find absolute path to driver")
	}
	drivers.RegisterBinary(driverName, driverPath)

	cmdConfig = &boilingcore.Config{
		DriverName:       driverName,
		OutFolder:        viper.GetString("output"),
		PkgName:          viper.GetString("pkgname"),
		Debug:            viper.GetBool("debug"),
		AddGlobal:        viper.GetBool("add-global-variants"),
		AddPanic:         viper.GetBool("add-panic-variants"),
		NoContext:        viper.GetBool("no-context"),
		NoTests:          viper.GetBool("no-tests"),
		NoHooks:          viper.GetBool("no-hooks"),
		NoRowsAffected:   viper.GetBool("no-rows-affected"),
		NoAutoTimestamps: viper.GetBool("no-auto-timestamps"),
		Wipe:             viper.GetBool("wipe"),
		StructTagCasing:  strings.ToLower(viper.GetString("struct-tag-casing")), // camel | snake
		TemplateDirs:     viper.GetStringSlice("templates"),
		Tags:             viper.GetStringSlice("tag"),
		Replacements:     viper.GetStringSlice("replace"),
		Aliases:          boilingcore.ConvertAliases(viper.Get("aliases")),
		TypeReplaces:     boilingcore.ConvertTypeReplace(viper.Get("types")),
	}

	if cmdConfig.Debug {
		fmt.Fprintln(os.Stderr, "using driver:", driverPath)
	}

	// Configure the driver
	cmdConfig.DriverConfig = map[string]interface{}{
		"whitelist": viper.GetStringSlice(driverName + ".whitelist"),
		"blacklist": viper.GetStringSlice(driverName + ".blacklist"),
	}

	keys := allKeys(driverName)
	for _, key := range keys {
		prefixedKey := fmt.Sprintf("%s.%s", driverName, key)
		cmdConfig.DriverConfig[key] = viper.Get(prefixedKey)
	}

	cmdConfig.Imports = configureImports()

	cmdState, err = boilingcore.New(cmdConfig)
	return err
}

func configureImports() importers.Collection {
	imports := importers.NewDefaultImports()

	mustMap := func(m importers.Map, err error) importers.Map {
		if err != nil {
			panic("failed to change viper interface into importers.Map: " + err.Error())
		}

		return m
	}

	if viper.IsSet("imports.all.standard") {
		imports.All.Standard = viper.GetStringSlice("imports.all.standard")
	}
	if viper.IsSet("imports.all.third_party") {
		imports.All.ThirdParty = viper.GetStringSlice("imports.all.third_party")
	}
	if viper.IsSet("imports.test.standard") {
		imports.Test.Standard = viper.GetStringSlice("imports.test.standard")
	}
	if viper.IsSet("imports.test.third_party") {
		imports.Test.ThirdParty = viper.GetStringSlice("imports.test.third_party")
	}
	if viper.IsSet("imports.singleton") {
		imports.Singleton = mustMap(importers.MapFromInterface(viper.Get("imports.singleton")))
	}
	if viper.IsSet("imports.test_singleton") {
		imports.TestSingleton = mustMap(importers.MapFromInterface(viper.Get("imports.test_singleton")))
	}
	if viper.IsSet("imports.based_on_type") {
		imports.BasedOnType = mustMap(importers.MapFromInterface(viper.Get("imports.based_on_type")))
	}

	return imports
}

func run(cmd *cobra.Command, args []string) error {
	return cmdState.Run()
}

func postRun(cmd *cobra.Command, args []string) error {
	return cmdState.Cleanup()
}

func allKeys(prefix string) []string {
	keys := make(map[string]bool)

	prefix = prefix + "."

	for _, e := range os.Environ() {
		splits := strings.SplitN(e, "=", 2)
		key := strings.Replace(strings.ToLower(splits[0]), "_", ".", -1)

		if strings.HasPrefix(key, prefix) {
			keys[strings.Replace(key, prefix, "", -1)] = true
		}
	}

	for _, key := range viper.AllKeys() {
		if strings.HasPrefix(key, prefix) {
			keys[strings.Replace(key, prefix, "", -1)] = true
		}
	}

	keySlice := make([]string, 0, len(keys))
	for k := range keys {
		keySlice = append(keySlice, k)
	}
	return keySlice
}
