// Copyright (c) 2015-present Xenia, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/xzl8028/xenia-server/utils/fileutils"
)

func main() {
	// Print angry message to use xenia command directly
	fmt.Println(`
------------------------------------ ERROR ------------------------------------------------
The platform binary has been deprecated, please switch to using the new xenia binary.
The platform binary will be removed in a future version.
-------------------------------------------------------------------------------------------
	`)

	// Execve the real MM binary
	args := os.Args
	args[0] = "xenia"
	args = append(args, "--platform")

	realXenia := fileutils.FindFile("xenia")
	if realXenia == "" {
		realXenia = fileutils.FindFile("bin/xenia")
	}

	if realXenia == "" {
		fmt.Println("Could not start Xenia, use the xenia command directly: failed to find xenia")
	} else if err := syscall.Exec(realXenia, args, nil); err != nil {
		fmt.Printf("Could not start Xenia, use the xenia command directly: %s\n", err.Error())
	}
}
