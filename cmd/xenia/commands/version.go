// Copyright (c) 2016-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package commands

import (
	"github.com/xzl8028/xenia-server/app"
	"github.com/xzl8028/xenia-server/model"
	"github.com/xzl8028/xenia-server/store"
	"github.com/xzl8028/xenia-server/store/sqlstore"
	"github.com/spf13/cobra"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version information",
	RunE:  versionCmdF,
}

func init() {
	RootCmd.AddCommand(VersionCmd)
}

func versionCmdF(command *cobra.Command, args []string) error {
	a, err := InitDBCommandContextCobra(command)
	if err != nil {
		return err
	}

	printVersion(a)

	return nil
}

func printVersion(a *app.App) {
	CommandPrintln("Version: " + model.CurrentVersion)
	CommandPrintln("Build Number: " + model.BuildNumber)
	CommandPrintln("Build Date: " + model.BuildDate)
	CommandPrintln("Build Hash: " + model.BuildHash)
	CommandPrintln("Build Enterprise Ready: " + model.BuildEnterpriseReady)
	if supplier, ok := a.Srv.Store.(*store.LayeredStore).DatabaseLayer.(*sqlstore.SqlSupplier); ok {
		CommandPrintln("DB Version: " + supplier.GetCurrentSchemaVersion())
	}
}
