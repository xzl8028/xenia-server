// Copyright (c) 2017-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package commands

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/xzl8028/xenia-server/mlog"
	"github.com/xzl8028/viper"
	"github.com/spf13/cobra"
)

var JobserverCmd = &cobra.Command{
	Use:   "jobserver",
	Short: "Start the Xenia job server",
	RunE:  jobserverCmdF,
}

func init() {
	JobserverCmd.Flags().Bool("nojobs", false, "Do not run jobs on this jobserver.")
	JobserverCmd.Flags().Bool("noschedule", false, "Do not schedule jobs from this jobserver.")

	RootCmd.AddCommand(JobserverCmd)
}

func jobserverCmdF(command *cobra.Command, args []string) error {
	// Options
	noJobs, _ := command.Flags().GetBool("nojobs")
	noSchedule, _ := command.Flags().GetBool("noschedule")

	config := viper.GetString("config")

	// Initialize
	a, err := InitDBCommandContext(config)
	if err != nil {
		return err
	}
	defer a.Shutdown()

	a.LoadLicense()

	// Run jobs
	mlog.Info("Starting Xenia job server")
	defer mlog.Info("Stopped Xenia job server")

	if !noJobs {
		a.Srv.Jobs.StartWorkers()
		defer a.Srv.Jobs.StopWorkers()
	}
	if !noSchedule {
		a.Srv.Jobs.StartSchedulers()
		defer a.Srv.Jobs.StopSchedulers()
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	// Cleanup anything that isn't handled by a defer statement
	mlog.Info("Stopping Xenia job server")

	return nil
}
