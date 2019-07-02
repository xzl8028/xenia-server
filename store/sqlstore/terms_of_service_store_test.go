package sqlstore

import (
	"testing"

	"github.com/xzl8028/xenia-server/store/storetest"
)

func TestTermsOfServiceStore(t *testing.T) {
	StoreTest(t, storetest.TestTermsOfServiceStore)
}
