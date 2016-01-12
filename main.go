package main

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tmaiaroto/discfg/commands"
	"github.com/tmaiaroto/discfg/config"
	"time"
)

var _ time.Duration
var _ bytes.Buffer

var Options = config.Options{StorageInterfaceName: "dynamodb", Version: "0.2.0"}

var DiscfgCmd = &cobra.Command{
	Use:   "discfg",
	Short: "discfg is a distributed configuration service",
	Long:  `A distributed configuration service using Amazon Web Services.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "discfg version number",
	Long:  `Displays the version number for discfg`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("discfg v" + Options.Version)
	},
}

var useCmd = &cobra.Command{
	Use:   "use",
	Short: "use a specific discfg",
	Long:  `For the current path, always use a specific discfg`,
	Run: func(cmd *cobra.Command, args []string) {
		resp := commands.Use(Options)
		commands.Out(Options, resp)
	},
}
var whichCmd = &cobra.Command{
	Use:   "which",
	Short: "shows current discfg in use",
	Long:  `Shows which discfg is currently selected for use at the current path`,
	Run: func(cmd *cobra.Command, args []string) {
		resp := commands.Which(Options)
		commands.Out(Options, resp)
	},
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create config",
	Long:  `Creates a new discfg distributed configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		setOptsFromArgs(args)
		resp := commands.CreateCfg(Options)
		commands.Out(Options, resp)
	},
}
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "set key value",
	Long:  `Sets a key value for a given discfg`,
	Run: func(cmd *cobra.Command, args []string) {
		setOptsFromArgs(args)
		resp := commands.SetKey(Options)
		commands.Out(Options, resp)
	},
}
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get key value",
	Long:  `Gets a key value for a given discfg`,
	Run: func(cmd *cobra.Command, args []string) {
		setOptsFromArgs(args)
		resp := commands.GetKey(Options)
		commands.Out(Options, resp)
	},
}
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete key",
	Long:  `Deletes a key for a given discfg`,
	Run: func(cmd *cobra.Command, args []string) {
		setOptsFromArgs(args)
		resp := commands.DeleteKey(Options)
		commands.Out(Options, resp)
	},
}
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "config information",
	Long:  `Information about the config including version and modified time`,
	Run: func(cmd *cobra.Command, args []string) {
		setOptsFromArgs(args)
		resp := commands.Info(Options, args)
		commands.Out(Options, resp)
	},
}
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "export entire config",
	Long:  `Exports the entire discfg to a file in JSON format`,
	Run: func(cmd *cobra.Command, args []string) {
		commands.Export(Options, args)
	},
}

func main() {
	// Set up commands
	DiscfgCmd.AddCommand(versionCmd)
	DiscfgCmd.PersistentFlags().StringVarP(&Options.OutputFormat, "format", "f", "human", "Output format for responses (human|json|slient)")

	// AWS Options & Credentials
	DiscfgCmd.PersistentFlags().StringVarP(&Options.Storage.DynamoDB.Region, "region", "l", "us-east-1", "AWS Region to use")
	DiscfgCmd.PersistentFlags().StringVarP(&Options.Storage.DynamoDB.AccessKeyId, "keyId", "k", "", "AWS Access Key ID")
	DiscfgCmd.PersistentFlags().StringVarP(&Options.Storage.DynamoDB.SecretAccessKey, "secretKey", "s", "", "AWS Secret Access Key")
	DiscfgCmd.PersistentFlags().StringVarP(&Options.Storage.DynamoDB.CredProfile, "credProfile", "p", "", "AWS Credentials Profile to use")
	DiscfgCmd.PersistentFlags().Int64VarP(&Options.Storage.DynamoDB.ReadCapacityUnits, "readCapacity", "y", 1, "DynamoDB Table Read Capacity Units")
	DiscfgCmd.PersistentFlags().Int64VarP(&Options.Storage.DynamoDB.WriteCapacityUnits, "writeCapacity", "z", 2, "DynamoDB Table Write Capacity Units")

	// Additional options by some operations
	DiscfgCmd.PersistentFlags().StringVarP(&Options.ConditionalValue, "condition", "c", "", "Conditional operation value")
	// TODO: In a future version
	// DiscfgCmd.PersistentFlags().BoolVarP(&Options.Recursive, "recursive", "r", false, "Recursively return or delete child keys")
	DiscfgCmd.PersistentFlags().Int64VarP(&Options.TTL, "ttl", "t", 0, "Set a time to live for a key (0 is no TTL)")

	DiscfgCmd.AddCommand(useCmd)
	DiscfgCmd.AddCommand(whichCmd)
	DiscfgCmd.AddCommand(createCmd)
	DiscfgCmd.AddCommand(setCmd)
	DiscfgCmd.AddCommand(getCmd)
	DiscfgCmd.AddCommand(deleteCmd)
	DiscfgCmd.AddCommand(infoCmd)
	DiscfgCmd.Execute()
}

// Takes positional command arguments and sets options from them (because some may be optional)
func setOptsFromArgs(args []string) {
	// The user may have set a config name in a `.discfg` file, for convenience, to shorten the commands.
	name := commands.GetDiscfgNameFromFile()
	Options.CfgName = name

	switch len(args) {
	case 1:
		Options.Key = args[0]
		break
	case 2:
		Options.Key = args[0]
		Options.Value = []byte(args[1])
		break
	case 3:
		Options.CfgName = args[0]
		Options.Key = args[1]
		Options.Value = []byte(args[2])
		break
	}
}
