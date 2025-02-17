// Copyright (c) 2017-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package api4

import (
	"testing"
)

func TestGetSamlMetadata(t *testing.T) {
	th := Setup().InitBasic()
	defer th.TearDown()
	Client := th.Client

	_, resp := Client.GetSamlMetadata()
	CheckNotImplementedStatus(t, resp)

	// Rest is tested by enterprise tests
}
