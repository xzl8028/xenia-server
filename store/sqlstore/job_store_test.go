// Copyright (c) 2015-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package sqlstore

import (
	"testing"

	"github.com/xzl8028/xenia-server/store/storetest"
)

func TestJobStore(t *testing.T) {
	StoreTest(t, storetest.TestJobStore)
}
