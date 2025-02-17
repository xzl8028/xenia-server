// Copyright (c) 2015-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package sqlstore

import (
	"testing"

	"github.com/xzl8028/xenia-server/model"
	"github.com/xzl8028/xenia-server/store"
	"github.com/xzl8028/xenia-server/store/storetest"

	"github.com/stretchr/testify/require"
)

func TestPreferenceStore(t *testing.T) {
	StoreTest(t, storetest.TestPreferenceStore)
}

func TestDeleteUnusedFeatures(t *testing.T) {
	StoreTest(t, func(t *testing.T, ss store.Store) {
		userId1 := model.NewId()
		userId2 := model.NewId()
		category := model.PREFERENCE_CATEGORY_ADVANCED_SETTINGS
		feature1 := "feature1"
		feature2 := "feature2"

		features := model.Preferences{
			{
				UserId:   userId1,
				Category: category,
				Name:     store.FEATURE_TOGGLE_PREFIX + feature1,
				Value:    "true",
			},
			{
				UserId:   userId2,
				Category: category,
				Name:     store.FEATURE_TOGGLE_PREFIX + feature1,
				Value:    "false",
			},
			{
				UserId:   userId1,
				Category: category,
				Name:     store.FEATURE_TOGGLE_PREFIX + feature2,
				Value:    "false",
			},
			{
				UserId:   userId2,
				Category: category,
				Name:     store.FEATURE_TOGGLE_PREFIX + feature2,
				Value:    "true",
			},
		}

		err := ss.Preference().Save(&features)
		require.Nil(t, err)

		ss.Preference().(*SqlPreferenceStore).DeleteUnusedFeatures()

		//make sure features with value "false" have actually been deleted from the database
		if val, err := ss.Preference().(*SqlPreferenceStore).GetReplica().SelectInt(`SELECT COUNT(*)
                            FROM Preferences
                    WHERE Category = :Category
                    AND Value = :Val
                    AND Name LIKE '`+store.FEATURE_TOGGLE_PREFIX+`%'`, map[string]interface{}{"Category": model.PREFERENCE_CATEGORY_ADVANCED_SETTINGS, "Val": "false"}); err != nil {
			t.Fatal(err)
		} else if val != 0 {
			t.Fatalf("Found %d features with value 'false', expected all to be deleted", val)
		}
		//
		// make sure features with value "true" remain saved
		if val, err := ss.Preference().(*SqlPreferenceStore).GetReplica().SelectInt(`SELECT COUNT(*)
                            FROM Preferences
                    WHERE Category = :Category
                    AND Value = :Val
                    AND Name LIKE '`+store.FEATURE_TOGGLE_PREFIX+`%'`, map[string]interface{}{"Category": model.PREFERENCE_CATEGORY_ADVANCED_SETTINGS, "Val": "true"}); err != nil {
			t.Fatal(err)
		} else if val == 0 {
			t.Fatalf("Found %d features with value 'true', expected to find at least %d features", val, 2)
		}
	})
}
