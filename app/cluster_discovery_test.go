// Copyright (c) 2016-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package app

import (
	"testing"

	"time"

	"github.com/xzl8028/xenia-server/model"
)

func TestClusterDiscoveryService(t *testing.T) {
	th := Setup(t)
	defer th.TearDown()

	ds := th.App.NewClusterDiscoveryService()
	ds.Type = model.CDS_TYPE_APP
	ds.ClusterName = "ClusterA"
	ds.AutoFillHostname()

	ds.Start()
	time.Sleep(2 * time.Second)

	ds.Stop()
	time.Sleep(2 * time.Second)
}
