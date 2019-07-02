// Copyright (c) 2015-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package commands

import (
	"github.com/xzl8028/viper"
	"github.com/spf13/cobra"
)

type Command = cobra.Command

func Run(args []string) error {
	RootCmd.SetArgs(args)
	return RootCmd.Execute()
}

var RootCmd = &cobra.Command{
	Use:   "xenia",
	Short: "Open source, self-hosted Slack-alternative",
	Long:  `Xenia offers workplace messaging across web, PC and phones with archiving, search and integration with your existing systems. Documentation available at https://docs.xenia.com`,
}

func init() {
	RootCmd.PersistentFlags().StringP("config", "c", "config.json", "Configuration file to use.")
	RootCmd.PersistentFlags().Bool("disableconfigwatch", false, "When set config.json will not be loaded from disk when the file is changed.")
	RootCmd.PersistentFlags().Bool("platform", false, "This flag signifies that the user tried to start the command from the platform binary, so we can log a mssage")
	RootCmd.PersistentFlags().MarkHidden("platform")

	viper.SetEnvPrefix("mm")
	viper.BindEnv("config")
	viper.BindPFlag("config", RootCmd.PersistentFlags().Lookup("config"))
}
