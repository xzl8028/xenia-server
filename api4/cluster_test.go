// Copyright (c) 2017-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package api4

import (
	"testing"

	"github.com/xzl8028/xenia-server/model"
)

func TestGetClusterStatus(t *testing.T) {
	th := Setup().InitBasic()
	defer th.TearDown()

	t.Run("as system user", func(t *testing.T) {
		_, resp := th.Client.GetClusterStatus()
		CheckForbiddenStatus(t, resp)
	})

	t.Run("as system admin", func(t *testing.T) {
		infos, resp := th.SystemAdminClient.GetClusterStatus()
		CheckNoError(t, resp)

		if infos == nil {
			t.Fatal("should not be nil")
		}
	})

	t.Run("as restricted system admin", func(t *testing.T) {
		th.App.UpdateConfig(func(cfg *model.Config) { *cfg.ExperimentalSettings.RestrictSystemAdmin = true })

		_, resp := th.SystemAdminClient.GetClusterStatus()
		CheckForbiddenStatus(t, resp)
	})
}
