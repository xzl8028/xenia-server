// Copyright (c) 2015-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package storetest

import (
	"testing"

	"github.com/xzl8028/xenia-server/model"
	"github.com/xzl8028/xenia-server/store"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionStore(t *testing.T, ss store.Store) {
	// Run serially to prevent interfering with other tests
	testSessionCleanup(t, ss)

	t.Run("Save", func(t *testing.T) { testSessionStoreSave(t, ss) })
	t.Run("SessionGet", func(t *testing.T) { testSessionGet(t, ss) })
	t.Run("SessionGetWithDeviceId", func(t *testing.T) { testSessionGetWithDeviceId(t, ss) })
	t.Run("SessionRemove", func(t *testing.T) { testSessionRemove(t, ss) })
	t.Run("SessionRemoveAll", func(t *testing.T) { testSessionRemoveAll(t, ss) })
	t.Run("SessionRemoveByUser", func(t *testing.T) { testSessionRemoveByUser(t, ss) })
	t.Run("SessionRemoveToken", func(t *testing.T) { testSessionRemoveToken(t, ss) })
	t.Run("SessionUpdateDeviceId", func(t *testing.T) { testSessionUpdateDeviceId(t, ss) })
	t.Run("SessionUpdateDeviceId2", func(t *testing.T) { testSessionUpdateDeviceId2(t, ss) })
	t.Run("UpdateLastActivityAt", func(t *testing.T) { testSessionStoreUpdateLastActivityAt(t, ss) })
	t.Run("SessionCount", func(t *testing.T) { testSessionCount(t, ss) })
}

func testSessionStoreSave(t *testing.T, ss store.Store) {
	s1 := &model.Session{}
	s1.UserId = model.NewId()

	if _, err := ss.Session().Save(s1); err != nil {
		t.Fatal(err)
	}
}

func testSessionGet(t *testing.T, ss store.Store) {
	s1 := &model.Session{}
	s1.UserId = model.NewId()

	s1, err := ss.Session().Save(s1)
	require.Nil(t, err)

	s2 := &model.Session{}
	s2.UserId = s1.UserId

	s2, err = ss.Session().Save(s2)
	require.Nil(t, err)

	s3 := &model.Session{}
	s3.UserId = s1.UserId
	s3.ExpiresAt = 1

	s3, err = ss.Session().Save(s3)
	require.Nil(t, err)

	if session, err := ss.Session().Get(s1.Id); err != nil {
		t.Fatal(err)
	} else {
		if session.Id != s1.Id {
			t.Fatal("should match")
		}
	}

	if session, err := ss.Session().GetSessions(s1.UserId); err != nil {
		t.Fatal(err)
	} else {
		if len(session) != 3 {
			t.Fatal("should match len")
		}
	}
}

func testSessionGetWithDeviceId(t *testing.T, ss store.Store) {
	s1 := &model.Session{}
	s1.UserId = model.NewId()
	s1.ExpiresAt = model.GetMillis() + 10000

	s1, err := ss.Session().Save(s1)
	require.Nil(t, err)

	s2 := &model.Session{}
	s2.UserId = s1.UserId
	s2.DeviceId = model.NewId()
	s2.ExpiresAt = model.GetMillis() + 10000

	s2, err = ss.Session().Save(s2)
	require.Nil(t, err)

	s3 := &model.Session{}
	s3.UserId = s1.UserId
	s3.ExpiresAt = 1
	s3.DeviceId = model.NewId()

	s3, err = ss.Session().Save(s3)
	require.Nil(t, err)

	if data, err := ss.Session().GetSessionsWithActiveDeviceIds(s1.UserId); err != nil {
		t.Fatal(err)
	} else {
		if len(data) != 1 {
			t.Fatal("should match len")
		}
	}
}

func testSessionRemove(t *testing.T, ss store.Store) {
	s1 := &model.Session{}
	s1.UserId = model.NewId()

	s1, err := ss.Session().Save(s1)
	require.Nil(t, err)

	if session, err := ss.Session().Get(s1.Id); err != nil {
		t.Fatal(err)
	} else {
		if session.Id != s1.Id {
			t.Fatal("should match")
		}
	}

	removeErr := ss.Session().Remove(s1.Id)
	require.Nil(t, removeErr)
	if _, err := ss.Session().Get(s1.Id); err == nil {
		t.Fatal("should have been removed")
	}
}

func testSessionRemoveAll(t *testing.T, ss store.Store) {
	s1 := &model.Session{}
	s1.UserId = model.NewId()

	s1, err := ss.Session().Save(s1)
	require.Nil(t, err)

	if session, err := ss.Session().Get(s1.Id); err != nil {
		t.Fatal(err)
	} else {
		if session.Id != s1.Id {
			t.Fatal("should match")
		}
	}

	removeErr := ss.Session().RemoveAllSessions()
	require.Nil(t, removeErr)

	if _, err := ss.Session().Get(s1.Id); err == nil {
		t.Fatal("should have been removed")
	}
}

func testSessionRemoveByUser(t *testing.T, ss store.Store) {
	s1 := &model.Session{}
	s1.UserId = model.NewId()

	s1, err := ss.Session().Save(s1)
	require.Nil(t, err)

	if session, err := ss.Session().Get(s1.Id); err != nil {
		t.Fatal(err)
	} else {
		if session.Id != s1.Id {
			t.Fatal("should match")
		}
	}

	deleteErr := ss.Session().PermanentDeleteSessionsByUser(s1.UserId)
	require.Nil(t, deleteErr)

	if _, err := ss.Session().Get(s1.Id); err == nil {
		t.Fatal("should have been removed")
	}
}

func testSessionRemoveToken(t *testing.T, ss store.Store) {
	s1 := &model.Session{}
	s1.UserId = model.NewId()

	s1, err := ss.Session().Save(s1)
	require.Nil(t, err)

	if session, err := ss.Session().Get(s1.Id); err != nil {
		t.Fatal(err)
	} else {
		if session.Id != s1.Id {
			t.Fatal("should match")
		}
	}

	removeErr := ss.Session().Remove(s1.Token)
	require.Nil(t, removeErr)

	if _, err := ss.Session().Get(s1.Id); err == nil {
		t.Fatal("should have been removed")
	}

	if session, err := ss.Session().GetSessions(s1.UserId); err != nil {
		t.Fatal(err)
	} else {
		if len(session) != 0 {
			t.Fatal("should match len")
		}
	}
}

func testSessionUpdateDeviceId(t *testing.T, ss store.Store) {
	s1 := &model.Session{}
	s1.UserId = model.NewId()

	s1, err := ss.Session().Save(s1)
	require.Nil(t, err)

	if _, err = ss.Session().UpdateDeviceId(s1.Id, model.PUSH_NOTIFY_APPLE+":1234567890", s1.ExpiresAt); err != nil {
		t.Fatal(err)
	}

	s2 := &model.Session{}
	s2.UserId = model.NewId()

	s2, err = ss.Session().Save(s2)
	require.Nil(t, err)

	if _, err := ss.Session().UpdateDeviceId(s2.Id, model.PUSH_NOTIFY_APPLE+":1234567890", s1.ExpiresAt); err != nil {
		t.Fatal(err)
	}
}

func testSessionUpdateDeviceId2(t *testing.T, ss store.Store) {
	s1 := &model.Session{}
	s1.UserId = model.NewId()

	s1, err := ss.Session().Save(s1)
	require.Nil(t, err)

	if _, err = ss.Session().UpdateDeviceId(s1.Id, model.PUSH_NOTIFY_APPLE_REACT_NATIVE+":1234567890", s1.ExpiresAt); err != nil {
		t.Fatal(err)
	}

	s2 := &model.Session{}
	s2.UserId = model.NewId()

	s2, err = ss.Session().Save(s2)
	require.Nil(t, err)

	if _, err := ss.Session().UpdateDeviceId(s2.Id, model.PUSH_NOTIFY_APPLE_REACT_NATIVE+":1234567890", s1.ExpiresAt); err != nil {
		t.Fatal(err)
	}
}

func testSessionStoreUpdateLastActivityAt(t *testing.T, ss store.Store) {
	s1 := &model.Session{}
	s1.UserId = model.NewId()

	s1, err := ss.Session().Save(s1)
	require.Nil(t, err)

	err = ss.Session().UpdateLastActivityAt(s1.Id, 1234567890)
	require.Nil(t, err)

	if session, err := ss.Session().Get(s1.Id); err != nil {
		t.Fatal(err)
	} else {
		if session.LastActivityAt != 1234567890 {
			t.Fatal("LastActivityAt not updated correctly")
		}
	}

}

func testSessionCount(t *testing.T, ss store.Store) {
	s1 := &model.Session{}
	s1.UserId = model.NewId()
	s1.ExpiresAt = model.GetMillis() + 100000

	s1, err := ss.Session().Save(s1)
	require.Nil(t, err)

	if count, err := ss.Session().AnalyticsSessionCount(); err != nil {
		t.Fatal(err)
	} else {
		if count == 0 {
			t.Fatal("should have at least 1 session")
		}
	}
}

func testSessionCleanup(t *testing.T, ss store.Store) {
	now := model.GetMillis()

	s1 := &model.Session{}
	s1.UserId = model.NewId()
	s1.ExpiresAt = 0 // never expires

	s1, err := ss.Session().Save(s1)
	require.Nil(t, err)

	s2 := &model.Session{}
	s2.UserId = s1.UserId
	s2.ExpiresAt = now + 1000000 // expires in the future

	s2, err = ss.Session().Save(s2)
	require.Nil(t, err)

	s3 := &model.Session{}
	s3.UserId = model.NewId()
	s3.ExpiresAt = 1 // expired

	s3, err = ss.Session().Save(s3)
	require.Nil(t, err)

	s4 := &model.Session{}
	s4.UserId = model.NewId()
	s4.ExpiresAt = 2 // expired

	s4, err = ss.Session().Save(s4)
	require.Nil(t, err)

	ss.Session().Cleanup(now, 1)

	_, err = ss.Session().Get(s1.Id)
	assert.Nil(t, err)

	_, err = ss.Session().Get(s2.Id)
	assert.Nil(t, err)

	_, err = ss.Session().Get(s3.Id)
	assert.NotNil(t, err)

	_, err = ss.Session().Get(s4.Id)
	assert.NotNil(t, err)

	removeErr := ss.Session().Remove(s1.Id)
	require.Nil(t, removeErr)

	removeErr = ss.Session().Remove(s2.Id)
	require.Nil(t, removeErr)
}
