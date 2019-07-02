// Copyright (c) 2015-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package commands

import (
	"testing"
)

func TestVersion(t *testing.T) {
	th := Setup()
	defer th.TearDown()

	th.CheckCommand(t, "version")
}
