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

var Config = config.Config{StorageInterfaceName: "dynamodb", Version: "0.1.0"}

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
		fmt.Println("discfg v" + Config.Version)
	},
}

var useCmd = &cobra.Command{
	Use:   "use",
	Short: "use a specific discfg",
	Long:  `For the current path, always use a specific discfg`,
	Run: func(cmd *cobra.Command, args []string) {
		resp := commands.Use(Config, args)
		commands.Out(Config, resp)
	},
}
var whichCmd = &cobra.Command{
	Use:   "which",
	Short: "shows current discfg in use",
	Long:  `Shows which discfg is currently selected for use at the current path`,
	Run: func(cmd *cobra.Command, args []string) {
		resp := commands.Which(Config, args)
		commands.Out(Config, resp)
	},
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create config",
	Long:  `Creates a new discfg distributed configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		resp := commands.CreateCfg(Config, args)
		commands.Out(Config, resp)
	},
}
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "set key value",
	Long:  `Sets a key value for a given discfg`,
	Run: func(cmd *cobra.Command, args []string) {
		resp := commands.SetKey(Config, args)
		commands.Out(Config, resp)
	},
}
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get key value",
	Long:  `Gets a key value for a given discfg`,
	Run: func(cmd *cobra.Command, args []string) {
		resp := commands.GetKey(Config, args)
		commands.Out(Config, resp)
	},
}
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete key",
	Long:  `Deletes a key for a given discfg`,
	Run: func(cmd *cobra.Command, args []string) {
		resp := commands.DeleteKey(Config, args)
		commands.Out(Config, resp)
	},
}
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "config information",
	Long:  `Information about the config including version and modified time`,
	Run: func(cmd *cobra.Command, args []string) {
		resp := commands.Info(Config, args)
		commands.Out(Config, resp)
	},
}
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "export entire config",
	Long:  `Exports the entire discfg to a file in JSON format`,
	Run: func(cmd *cobra.Command, args []string) {
		commands.Export(Config, args)
	},
}

func main() {
	// Set up commands
	DiscfgCmd.AddCommand(versionCmd)
	DiscfgCmd.PersistentFlags().StringVarP(&Config.OutputFormat, "format", "f", "human", "Output format for responses (human|json|slient)")

	// AWS Options & Credentials
	DiscfgCmd.PersistentFlags().StringVarP(&Config.Storage.DynamoDB.Region, "region", "l", "us-east-1", "AWS Region to use")
	DiscfgCmd.PersistentFlags().StringVarP(&Config.Storage.DynamoDB.AccessKeyId, "keyId", "k", "", "AWS Access Key ID")
	DiscfgCmd.PersistentFlags().StringVarP(&Config.Storage.DynamoDB.SecretAccessKey, "secretKey", "s", "", "AWS Secret Access Key")
	DiscfgCmd.PersistentFlags().StringVarP(&Config.Storage.DynamoDB.CredProfile, "credProfile", "p", "", "AWS Credentials Profile to use")

	// Additional options by some operations
	DiscfgCmd.PersistentFlags().StringVarP(&Config.ConditionalValue, "condition", "c", "", "Conditional operation value")
	DiscfgCmd.PersistentFlags().BoolVarP(&Config.Recursive, "recursive", "r", false, "Recursively return or delete child keys")
	DiscfgCmd.PersistentFlags().IntVarP(&Config.TTL, "ttl", "t", 0, "Set a time to live for a key (0 is no TTL)")

	DiscfgCmd.AddCommand(useCmd)
	DiscfgCmd.AddCommand(whichCmd)
	DiscfgCmd.AddCommand(createCmd)
	DiscfgCmd.AddCommand(setCmd)
	DiscfgCmd.AddCommand(getCmd)
	DiscfgCmd.AddCommand(deleteCmd)
	DiscfgCmd.AddCommand(infoCmd)
	DiscfgCmd.Execute()
}
