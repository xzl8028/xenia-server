// Copyright (c) 2015-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package storetest

import (
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/xzl8028/gorp"
	"github.com/xzl8028/xenia-server/model"
	"github.com/xzl8028/xenia-server/store"
)

type SqlSupplier interface {
	GetMaster() *gorp.DbMap
}

func cleanupChannels(t *testing.T, ss store.Store) {
	list, err := ss.Channel().GetAllChannels(0, 100000, store.ChannelSearchOpts{IncludeDeleted: true})
	if err != nil {
		t.Fatalf("error cleaning all channels: %v", err)
	}
	for _, channel := range *list {
		ss.Channel().PermanentDelete(channel.Id)
	}
}

func TestChannelStore(t *testing.T, ss store.Store, s SqlSupplier) {
	createDefaultRoles(t, ss)

	t.Run("Save", func(t *testing.T) { testChannelStoreSave(t, ss) })
	t.Run("SaveDirectChannel", func(t *testing.T) { testChannelStoreSaveDirectChannel(t, ss, s) })
	t.Run("CreateDirectChannel", func(t *testing.T) { testChannelStoreCreateDirectChannel(t, ss) })
	t.Run("Update", func(t *testing.T) { testChannelStoreUpdate(t, ss) })
	t.Run("GetChannelUnread", func(t *testing.T) { testGetChannelUnread(t, ss) })
	t.Run("Get", func(t *testing.T) { testChannelStoreGet(t, ss, s) })
	t.Run("GetChannelsByIds", func(t *testing.T) { testChannelStoreGetChannelsByIds(t, ss) })
	t.Run("GetForPost", func(t *testing.T) { testChannelStoreGetForPost(t, ss) })
	t.Run("Restore", func(t *testing.T) { testChannelStoreRestore(t, ss) })
	t.Run("Delete", func(t *testing.T) { testChannelStoreDelete(t, ss) })
	t.Run("GetByName", func(t *testing.T) { testChannelStoreGetByName(t, ss) })
	t.Run("GetByNames", func(t *testing.T) { testChannelStoreGetByNames(t, ss) })
	t.Run("GetDeletedByName", func(t *testing.T) { testChannelStoreGetDeletedByName(t, ss) })
	t.Run("GetDeleted", func(t *testing.T) { testChannelStoreGetDeleted(t, ss) })
	t.Run("ChannelMemberStore", func(t *testing.T) { testChannelMemberStore(t, ss) })
	t.Run("ChannelDeleteMemberStore", func(t *testing.T) { testChannelDeleteMemberStore(t, ss) })
	t.Run("GetChannels", func(t *testing.T) { testChannelStoreGetChannels(t, ss) })
	t.Run("GetAllChannels", func(t *testing.T) { testChannelStoreGetAllChannels(t, ss, s) })
	t.Run("GetMoreChannels", func(t *testing.T) { testChannelStoreGetMoreChannels(t, ss) })
	t.Run("GetPublicChannelsForTeam", func(t *testing.T) { testChannelStoreGetPublicChannelsForTeam(t, ss) })
	t.Run("GetPublicChannelsByIdsForTeam", func(t *testing.T) { testChannelStoreGetPublicChannelsByIdsForTeam(t, ss) })
	t.Run("GetChannelCounts", func(t *testing.T) { testChannelStoreGetChannelCounts(t, ss) })
	t.Run("GetMembersForUser", func(t *testing.T) { testChannelStoreGetMembersForUser(t, ss) })
	t.Run("GetMembersForUserWithPagination", func(t *testing.T) { testChannelStoreGetMembersForUserWithPagination(t, ss) })
	t.Run("UpdateLastViewedAt", func(t *testing.T) { testChannelStoreUpdateLastViewedAt(t, ss) })
	t.Run("IncrementMentionCount", func(t *testing.T) { testChannelStoreIncrementMentionCount(t, ss) })
	t.Run("UpdateChannelMember", func(t *testing.T) { testUpdateChannelMember(t, ss) })
	t.Run("GetMember", func(t *testing.T) { testGetMember(t, ss) })
	t.Run("GetMemberForPost", func(t *testing.T) { testChannelStoreGetMemberForPost(t, ss) })
	t.Run("GetMemberCount", func(t *testing.T) { testGetMemberCount(t, ss) })
	t.Run("SearchMore", func(t *testing.T) { testChannelStoreSearchMore(t, ss) })
	t.Run("SearchInTeam", func(t *testing.T) { testChannelStoreSearchInTeam(t, ss) })
	t.Run("SearchAllChannels", func(t *testing.T) { testChannelStoreSearchAllChannels(t, ss) })
	t.Run("AutocompleteInTeamForSearch", func(t *testing.T) { testChannelStoreAutocompleteInTeamForSearch(t, ss, s) })
	t.Run("GetMembersByIds", func(t *testing.T) { testChannelStoreGetMembersByIds(t, ss) })
	t.Run("SearchGroupChannels", func(t *testing.T) { testChannelStoreSearchGroupChannels(t, ss) })
	t.Run("AnalyticsDeletedTypeCount", func(t *testing.T) { testChannelStoreAnalyticsDeletedTypeCount(t, ss) })
	t.Run("GetPinnedPosts", func(t *testing.T) { testChannelStoreGetPinnedPosts(t, ss) })
	t.Run("MaxChannelsPerTeam", func(t *testing.T) { testChannelStoreMaxChannelsPerTeam(t, ss) })
	t.Run("GetChannelsByScheme", func(t *testing.T) { testChannelStoreGetChannelsByScheme(t, ss) })
	t.Run("MigrateChannelMembers", func(t *testing.T) { testChannelStoreMigrateChannelMembers(t, ss) })
	t.Run("ResetAllChannelSchemes", func(t *testing.T) { testResetAllChannelSchemes(t, ss) })
	t.Run("ClearAllCustomRoleAssignments", func(t *testing.T) { testChannelStoreClearAllCustomRoleAssignments(t, ss) })
	t.Run("MaterializedPublicChannels", func(t *testing.T) { testMaterializedPublicChannels(t, ss, s) })
	t.Run("GetAllChannelsForExportAfter", func(t *testing.T) { testChannelStoreGetAllChannelsForExportAfter(t, ss) })
	t.Run("GetChannelMembersForExport", func(t *testing.T) { testChannelStoreGetChannelMembersForExport(t, ss) })
	t.Run("RemoveAllDeactivatedMembers", func(t *testing.T) { testChannelStoreRemoveAllDeactivatedMembers(t, ss) })
	t.Run("ExportAllDirectChannels", func(t *testing.T) { testChannelStoreExportAllDirectChannels(t, ss, s) })
	t.Run("ExportAllDirectChannelsExcludePrivateAndPublic", func(t *testing.T) { testChannelStoreExportAllDirectChannelsExcludePrivateAndPublic(t, ss, s) })
	t.Run("ExportAllDirectChannelsDeletedChannel", func(t *testing.T) { testChannelStoreExportAllDirectChannelsDeletedChannel(t, ss, s) })
	t.Run("GetChannelsBatchForIndexing", func(t *testing.T) { testChannelStoreGetChannelsBatchForIndexing(t, ss) })
}

func testChannelStoreSave(t *testing.T, ss store.Store) {
	teamId := model.NewId()

	o1 := model.Channel{}
	o1.TeamId = teamId
	o1.DisplayName = "Name"
	o1.Name = "zz" + model.NewId() + "b"
	o1.Type = model.CHANNEL_OPEN

	if _, err := ss.Channel().Save(&o1, -1); err != nil {
		t.Fatal("couldn't save item", err)
	}

	if _, err := ss.Channel().Save(&o1, -1); err == nil {
		t.Fatal("shouldn't be able to update from save")
	}

	o1.Id = ""
	if _, err := ss.Channel().Save(&o1, -1); err == nil {
		t.Fatal("should be unique name")
	}

	o1.Id = ""
	o1.Name = "zz" + model.NewId() + "b"
	o1.Type = model.CHANNEL_DIRECT
	if _, err := ss.Channel().Save(&o1, -1); err == nil {
		t.Fatal("Should not be able to save direct channel")
	}
}

func testChannelStoreSaveDirectChannel(t *testing.T, ss store.Store, s SqlSupplier) {
	teamId := model.NewId()

	o1 := model.Channel{}
	o1.TeamId = teamId
	o1.DisplayName = "Name"
	o1.Name = "zz" + model.NewId() + "b"
	o1.Type = model.CHANNEL_DIRECT

	u1 := &model.User{}
	u1.Email = MakeEmail()
	u1.Nickname = model.NewId()
	store.Must(ss.User().Save(u1))
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: model.NewId(), UserId: u1.Id}, -1))

	u2 := &model.User{}
	u2.Email = MakeEmail()
	u2.Nickname = model.NewId()
	store.Must(ss.User().Save(u2))
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: model.NewId(), UserId: u2.Id}, -1))

	m1 := model.ChannelMember{}
	m1.ChannelId = o1.Id
	m1.UserId = u1.Id
	m1.NotifyProps = model.GetDefaultChannelNotifyProps()

	m2 := model.ChannelMember{}
	m2.ChannelId = o1.Id
	m2.UserId = u2.Id
	m2.NotifyProps = model.GetDefaultChannelNotifyProps()

	if _, err := ss.Channel().SaveDirectChannel(&o1, &m1, &m2); err != nil {
		t.Fatal("couldn't save direct channel", err)
	}

	members, err := ss.Channel().GetMembers(o1.Id, 0, 100)
	require.Nil(t, err)
	if len(*members) != 2 {
		t.Fatal("should have saved 2 members")
	}

	if _, err = ss.Channel().SaveDirectChannel(&o1, &m1, &m2); err == nil {
		t.Fatal("shouldn't be able to update from save")
	}

	// Attempt to save a direct channel that already exists
	o1a := model.Channel{
		TeamId:      o1.TeamId,
		DisplayName: o1.DisplayName,
		Name:        o1.Name,
		Type:        o1.Type,
	}

	returnedChannel, err := ss.Channel().SaveDirectChannel(&o1a, &m1, &m2)
	if err == nil {
		t.Fatal("should've failed to save a duplicate direct channel")
	} else if err.Id != store.CHANNEL_EXISTS_ERROR {
		t.Fatal("should've returned CHANNEL_EXISTS_ERROR")
	} else if returnedChannel.Id != o1.Id {
		t.Fatal("should've returned original channel when saving a duplicate direct channel")
	}

	// Attempt to save a non-direct channel
	o1.Id = ""
	o1.Name = "zz" + model.NewId() + "b"
	o1.Type = model.CHANNEL_OPEN
	if _, err = ss.Channel().SaveDirectChannel(&o1, &m1, &m2); err == nil {
		t.Fatal("Should not be able to save non-direct channel")
	}

	// Save yourself Direct Message
	o1.Id = ""
	o1.DisplayName = "Myself"
	o1.Name = "zz" + model.NewId() + "b"
	o1.Type = model.CHANNEL_DIRECT
	if _, err = ss.Channel().SaveDirectChannel(&o1, &m1, &m1); err != nil {
		t.Fatal("couldn't save direct channel", err)
	}

	members, err = ss.Channel().GetMembers(o1.Id, 0, 100)
	require.Nil(t, err)
	if len(*members) != 1 {
		t.Fatal("should have saved just 1 member")
	}

	// Manually truncate Channels table until testlib can handle cleanups
	s.GetMaster().Exec("TRUNCATE Channels")
}

func testChannelStoreCreateDirectChannel(t *testing.T, ss store.Store) {
	u1 := &model.User{}
	u1.Email = MakeEmail()
	u1.Nickname = model.NewId()
	store.Must(ss.User().Save(u1))
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: model.NewId(), UserId: u1.Id}, -1))

	u2 := &model.User{}
	u2.Email = MakeEmail()
	u2.Nickname = model.NewId()
	store.Must(ss.User().Save(u2))
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: model.NewId(), UserId: u2.Id}, -1))

	c1, err := ss.Channel().CreateDirectChannel(u1.Id, u2.Id)
	if err != nil {
		t.Fatal("couldn't create direct channel", err)
	}
	defer func() {
		ss.Channel().PermanentDeleteMembersByChannel(c1.Id)
		<-ss.Channel().PermanentDelete(c1.Id)
	}()

	members, err := ss.Channel().GetMembers(c1.Id, 0, 100)
	require.Nil(t, err)
	if len(*members) != 2 {
		t.Fatal("should have saved 2 members")
	}
}

func testChannelStoreUpdate(t *testing.T, ss store.Store) {
	o1 := model.Channel{}
	o1.TeamId = model.NewId()
	o1.DisplayName = "Name"
	o1.Name = "zz" + model.NewId() + "b"
	o1.Type = model.CHANNEL_OPEN

	_, err := ss.Channel().Save(&o1, -1)
	require.Nil(t, err)

	o2 := model.Channel{}
	o2.TeamId = o1.TeamId
	o2.DisplayName = "Name"
	o2.Name = "zz" + model.NewId() + "b"
	o2.Type = model.CHANNEL_OPEN

	_, err = ss.Channel().Save(&o2, -1)
	require.Nil(t, err)

	time.Sleep(100 * time.Millisecond)

	if _, err := ss.Channel().Update(&o1); err != nil {
		t.Fatal(err)
	}

	o1.DeleteAt = 100
	if _, err := ss.Channel().Update(&o1); err == nil {
		t.Fatal("Update should have failed because channel is archived")
	}

	o1.DeleteAt = 0
	o1.Id = "missing"
	if _, err := ss.Channel().Update(&o1); err == nil {
		t.Fatal("Update should have failed because of missing key")
	}

	o1.Id = model.NewId()
	if _, err := ss.Channel().Update(&o1); err == nil {
		t.Fatal("Update should have faile because id change")
	}

	o2.Name = o1.Name
	if _, err := ss.Channel().Update(&o2); err == nil {
		t.Fatal("Update should have failed because of existing name")
	}
}

func testGetChannelUnread(t *testing.T, ss store.Store) {
	teamId1 := model.NewId()
	teamId2 := model.NewId()

	uid := model.NewId()
	m1 := &model.TeamMember{TeamId: teamId1, UserId: uid}
	m2 := &model.TeamMember{TeamId: teamId2, UserId: uid}
	store.Must(ss.Team().SaveMember(m1, -1))
	store.Must(ss.Team().SaveMember(m2, -1))
	notifyPropsModel := model.GetDefaultChannelNotifyProps()

	// Setup Channel 1
	c1 := &model.Channel{TeamId: m1.TeamId, Name: model.NewId(), DisplayName: "Downtown", Type: model.CHANNEL_OPEN, TotalMsgCount: 100}
	_, err := ss.Channel().Save(c1, -1)
	require.Nil(t, err)

	cm1 := &model.ChannelMember{ChannelId: c1.Id, UserId: m1.UserId, NotifyProps: notifyPropsModel, MsgCount: 90}
	store.Must(ss.Channel().SaveMember(cm1))

	// Setup Channel 2
	c2 := &model.Channel{TeamId: m2.TeamId, Name: model.NewId(), DisplayName: "Cultural", Type: model.CHANNEL_OPEN, TotalMsgCount: 100}
	_, err = ss.Channel().Save(c2, -1)
	require.Nil(t, err)

	cm2 := &model.ChannelMember{ChannelId: c2.Id, UserId: m2.UserId, NotifyProps: notifyPropsModel, MsgCount: 90, MentionCount: 5}
	store.Must(ss.Channel().SaveMember(cm2))

	// Check for Channel 1
	if ch, err := ss.Channel().GetChannelUnread(c1.Id, uid); err != nil {
		t.Fatal(err)
	} else {
		if c1.Id != ch.ChannelId {
			t.Fatal("wrong channel id")
		}

		if teamId1 != ch.TeamId {
			t.Fatal("wrong team id for channel 1")
		}

		if ch.NotifyProps == nil {
			t.Fatal("wrong props for channel 1")
		}

		if ch.MentionCount != 0 {
			t.Fatal("wrong MentionCount for channel 1")
		}

		if ch.MsgCount != 10 {
			t.Fatal("wrong MsgCount for channel 1")
		}
	}

	// Check for Channel 2
	if ch2, err := ss.Channel().GetChannelUnread(c2.Id, uid); err != nil {
		t.Fatal(err)
	} else {
		if c2.Id != ch2.ChannelId {
			t.Fatal("wrong channel id")
		}

		if teamId2 != ch2.TeamId {
			t.Fatal("wrong team id")
		}

		if ch2.MentionCount != 5 {
			t.Fatal("wrong MentionCount for channel 2")
		}

		if ch2.MsgCount != 10 {
			t.Fatal("wrong MsgCount for channel 2")
		}
	}
}

func testChannelStoreGet(t *testing.T, ss store.Store, s SqlSupplier) {
	o1 := model.Channel{}
	o1.TeamId = model.NewId()
	o1.DisplayName = "Name"
	o1.Name = "zz" + model.NewId() + "b"
	o1.Type = model.CHANNEL_OPEN
	_, err := ss.Channel().Save(&o1, -1)
	require.Nil(t, err)

	c1 := &model.Channel{}
	if c1, err = ss.Channel().Get(o1.Id, false); err != nil {
		t.Fatal(err)
	} else {
		if c1.ToJson() != o1.ToJson() {
			t.Fatal("invalid returned channel")
		}
	}

	if _, err = ss.Channel().Get("", false); err == nil {
		t.Fatal("Missing id should have failed")
	}

	u1 := &model.User{}
	u1.Email = MakeEmail()
	u1.Nickname = model.NewId()
	store.Must(ss.User().Save(u1))
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: model.NewId(), UserId: u1.Id}, -1))

	u2 := model.User{}
	u2.Email = MakeEmail()
	u2.Nickname = model.NewId()
	store.Must(ss.User().Save(&u2))
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: model.NewId(), UserId: u2.Id}, -1))

	o2 := model.Channel{}
	o2.TeamId = model.NewId()
	o2.DisplayName = "Direct Name"
	o2.Name = "zz" + model.NewId() + "b"
	o2.Type = model.CHANNEL_DIRECT

	m1 := model.ChannelMember{}
	m1.ChannelId = o2.Id
	m1.UserId = u1.Id
	m1.NotifyProps = model.GetDefaultChannelNotifyProps()

	m2 := model.ChannelMember{}
	m2.ChannelId = o2.Id
	m2.UserId = u2.Id
	m2.NotifyProps = model.GetDefaultChannelNotifyProps()

	_, err = ss.Channel().SaveDirectChannel(&o2, &m1, &m2)
	require.Nil(t, err)

	if c2, err := ss.Channel().Get(o2.Id, false); err != nil {
		t.Fatal(err)
	} else {
		if c2.ToJson() != o2.ToJson() {
			t.Fatal("invalid returned channel")
		}
	}

	if c4, err := ss.Channel().Get(o2.Id, true); err != nil {
		t.Fatal(err)
	} else {
		if c4.ToJson() != o2.ToJson() {
			t.Fatal("invalid returned channel")
		}
	}

	if channels, chanErr := ss.Channel().GetAll(o1.TeamId); chanErr != nil {
		t.Fatal(chanErr)
	} else {
		if len(channels) == 0 {
			t.Fatal("too little")
		}
	}

	if channels, err := ss.Channel().GetTeamChannels(o1.TeamId); err != nil {
		t.Fatal(err)
	} else {
		if len(*channels) == 0 {
			t.Fatal("too little")
		}
	}
	// Manually truncate Channels table until testlib can handle cleanups
	s.GetMaster().Exec("TRUNCATE Channels")
}

func testChannelStoreGetChannelsByIds(t *testing.T, ss store.Store) {
	o1 := model.Channel{}
	o1.TeamId = model.NewId()
	o1.DisplayName = "Name"
	o1.Name = "aa" + model.NewId() + "b"
	o1.Type = model.CHANNEL_OPEN
	_, err := ss.Channel().Save(&o1, -1)
	require.Nil(t, err)

	u1 := &model.User{}
	u1.Email = MakeEmail()
	u1.Nickname = model.NewId()
	store.Must(ss.User().Save(u1))
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: model.NewId(), UserId: u1.Id}, -1))

	u2 := model.User{}
	u2.Email = MakeEmail()
	u2.Nickname = model.NewId()
	store.Must(ss.User().Save(&u2))
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: model.NewId(), UserId: u2.Id}, -1))

	o2 := model.Channel{}
	o2.TeamId = model.NewId()
	o2.DisplayName = "Direct Name"
	o2.Name = "bb" + model.NewId() + "b"
	o2.Type = model.CHANNEL_DIRECT

	m1 := model.ChannelMember{}
	m1.ChannelId = o2.Id
	m1.UserId = u1.Id
	m1.NotifyProps = model.GetDefaultChannelNotifyProps()

	m2 := model.ChannelMember{}
	m2.ChannelId = o2.Id
	m2.UserId = u2.Id
	m2.NotifyProps = model.GetDefaultChannelNotifyProps()

	_, err = ss.Channel().SaveDirectChannel(&o2, &m1, &m2)
	require.Nil(t, err)

	if r1, err := ss.Channel().GetChannelsByIds([]string{o1.Id, o2.Id}); err != nil {
		t.Fatal(err)
	} else {
		if len(r1) != 2 {
			t.Fatal("invalid returned channels, expected 2 and got " + strconv.Itoa(len(r1)))
		}
		if r1[0].ToJson() != o1.ToJson() {
			t.Fatal("invalid returned channel")
		}
		if r1[1].ToJson() != o2.ToJson() {
			t.Fatal("invalid returned channel")
		}
	}

	nonexistentId := "abcd1234"
	if r2, err := ss.Channel().GetChannelsByIds([]string{o1.Id, nonexistentId}); err != nil {
		t.Fatal(err)
	} else {
		if len(r2) != 1 {
			t.Fatal("invalid returned channels, expected 1 and got " + strconv.Itoa(len(r2)))
		}
		if r2[0].ToJson() != o1.ToJson() {
			t.Fatal("invalid returned channel")
		}
	}
}

func testChannelStoreGetForPost(t *testing.T, ss store.Store) {

	ch := &model.Channel{
		TeamId:      model.NewId(),
		DisplayName: "Name",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_OPEN,
	}
	o1, err := ss.Channel().Save(ch, -1)
	require.Nil(t, err)

	p1, err := ss.Post().Save(&model.Post{
		UserId:    model.NewId(),
		ChannelId: o1.Id,
		Message:   "test",
	})
	require.Nil(t, err)

	if channel, chanErr := ss.Channel().GetForPost(p1.Id); chanErr != nil {
		t.Fatal(chanErr)
	} else if channel.Id != o1.Id {
		t.Fatal("incorrect channel returned")
	}
}

func testChannelStoreRestore(t *testing.T, ss store.Store) {
	o1 := model.Channel{}
	o1.TeamId = model.NewId()
	o1.DisplayName = "Channel1"
	o1.Name = "zz" + model.NewId() + "b"
	o1.Type = model.CHANNEL_OPEN
	_, err := ss.Channel().Save(&o1, -1)
	require.Nil(t, err)

	if err := ss.Channel().Delete(o1.Id, model.GetMillis()); err != nil {
		t.Fatal(err)
	}

	if c, _ := ss.Channel().Get(o1.Id, false); c.DeleteAt == 0 {
		t.Fatal("should have been deleted")
	}

	if err := ss.Channel().Restore(o1.Id, model.GetMillis()); err != nil {
		t.Fatal(err)
	}

	if c, _ := ss.Channel().Get(o1.Id, false); c.DeleteAt != 0 {
		t.Fatal("should have been restored")
	}

}

func testChannelStoreDelete(t *testing.T, ss store.Store) {
	o1 := model.Channel{}
	o1.TeamId = model.NewId()
	o1.DisplayName = "Channel1"
	o1.Name = "zz" + model.NewId() + "b"
	o1.Type = model.CHANNEL_OPEN
	_, err := ss.Channel().Save(&o1, -1)
	require.Nil(t, err)

	o2 := model.Channel{}
	o2.TeamId = o1.TeamId
	o2.DisplayName = "Channel2"
	o2.Name = "zz" + model.NewId() + "b"
	o2.Type = model.CHANNEL_OPEN
	_, err = ss.Channel().Save(&o2, -1)
	require.Nil(t, err)

	o3 := model.Channel{}
	o3.TeamId = o1.TeamId
	o3.DisplayName = "Channel3"
	o3.Name = "zz" + model.NewId() + "b"
	o3.Type = model.CHANNEL_OPEN
	_, err = ss.Channel().Save(&o3, -1)
	require.Nil(t, err)

	o4 := model.Channel{}
	o4.TeamId = o1.TeamId
	o4.DisplayName = "Channel4"
	o4.Name = "zz" + model.NewId() + "b"
	o4.Type = model.CHANNEL_OPEN
	_, err = ss.Channel().Save(&o4, -1)
	require.Nil(t, err)

	m1 := model.ChannelMember{}
	m1.ChannelId = o1.Id
	m1.UserId = model.NewId()
	m1.NotifyProps = model.GetDefaultChannelNotifyProps()
	store.Must(ss.Channel().SaveMember(&m1))

	m2 := model.ChannelMember{}
	m2.ChannelId = o2.Id
	m2.UserId = m1.UserId
	m2.NotifyProps = model.GetDefaultChannelNotifyProps()
	store.Must(ss.Channel().SaveMember(&m2))

	if err = ss.Channel().Delete(o1.Id, model.GetMillis()); err != nil {
		t.Fatal(err)
	}

	if c, _ := ss.Channel().Get(o1.Id, false); c.DeleteAt == 0 {
		t.Fatal("should have been deleted")
	}

	if err = ss.Channel().Delete(o3.Id, model.GetMillis()); err != nil {
		t.Fatal(err)
	}

	list, err := ss.Channel().GetChannels(o1.TeamId, m1.UserId, false)
	require.Nil(t, err)

	if len(*list) != 1 {
		t.Fatal("invalid number of channels")
	}

	list, err = ss.Channel().GetMoreChannels(o1.TeamId, m1.UserId, 0, 100)
	require.Nil(t, err)

	if len(*list) != 1 {
		t.Fatal("invalid number of channels")
	}

	cresult := <-ss.Channel().PermanentDelete(o2.Id)
	require.Nil(t, cresult.Err)

	list, err = ss.Channel().GetChannels(o1.TeamId, m1.UserId, false)
	if assert.NotNil(t, err) {
		require.Equal(t, "store.sql_channel.get_channels.not_found.app_error", err.Id)
	} else {
		require.Equal(t, &model.ChannelList{}, list)
	}

	if r := <-ss.Channel().PermanentDeleteByTeam(o1.TeamId); r.Err != nil {
		t.Fatal(r.Err)
	}
}

func testChannelStoreGetByName(t *testing.T, ss store.Store) {
	o1 := model.Channel{}
	o1.TeamId = model.NewId()
	o1.DisplayName = "Name"
	o1.Name = "zz" + model.NewId() + "b"
	o1.Type = model.CHANNEL_OPEN
	_, err := ss.Channel().Save(&o1, -1)
	require.Nil(t, err)

	result, err := ss.Channel().GetByName(o1.TeamId, o1.Name, true)
	require.Nil(t, err)
	require.Equal(t, o1.ToJson(), result.ToJson(), "invalid returned channel")

	channelID := result.Id

	result, err = ss.Channel().GetByName(o1.TeamId, "", true)
	require.NotNil(t, err, "Missing id should have failed")

	result, err = ss.Channel().GetByName(o1.TeamId, o1.Name, false)
	require.Nil(t, err)
	require.Equal(t, o1.ToJson(), result.ToJson(), "invalid returned channel")

	result, err = ss.Channel().GetByName(o1.TeamId, "", false)
	require.NotNil(t, err, "Missing id should have failed")

	err = ss.Channel().Delete(channelID, model.GetMillis())
	require.Nil(t, err, "channel should have been deleted")

	result, err = ss.Channel().GetByName(o1.TeamId, o1.Name, false)
	require.NotNil(t, err, "Deleted channel should not be returned by GetByName()")
}

func testChannelStoreGetByNames(t *testing.T, ss store.Store) {
	o1 := model.Channel{
		TeamId:      model.NewId(),
		DisplayName: "Name",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_OPEN,
	}
	_, err := ss.Channel().Save(&o1, -1)
	require.Nil(t, err)

	o2 := model.Channel{
		TeamId:      o1.TeamId,
		DisplayName: "Name",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_OPEN,
	}
	_, err = ss.Channel().Save(&o2, -1)
	require.Nil(t, err)

	for index, tc := range []struct {
		TeamId      string
		Names       []string
		ExpectedIds []string
	}{
		{o1.TeamId, []string{o1.Name}, []string{o1.Id}},
		{o1.TeamId, []string{o1.Name, o2.Name}, []string{o1.Id, o2.Id}},
		{o1.TeamId, nil, nil},
		{o1.TeamId, []string{"foo"}, nil},
		{o1.TeamId, []string{o1.Name, "foo", o2.Name, o2.Name}, []string{o1.Id, o2.Id}},
		{"", []string{o1.Name, "foo", o2.Name, o2.Name}, []string{o1.Id, o2.Id}},
		{"asd", []string{o1.Name, "foo", o2.Name, o2.Name}, nil},
	} {
		var channels []*model.Channel
		channels, err = ss.Channel().GetByNames(tc.TeamId, tc.Names, true)
		require.Nil(t, err)
		var ids []string
		for _, channel := range channels {
			ids = append(ids, channel.Id)
		}
		sort.Strings(ids)
		sort.Strings(tc.ExpectedIds)
		assert.Equal(t, tc.ExpectedIds, ids, "tc %v", index)
	}

	err = ss.Channel().Delete(o1.Id, model.GetMillis())
	require.Nil(t, err, "channel should have been deleted")

	err = ss.Channel().Delete(o2.Id, model.GetMillis())
	require.Nil(t, err, "channel should have been deleted")

	channels, err := ss.Channel().GetByNames(o1.TeamId, []string{o1.Name}, false)
	require.Nil(t, err)
	assert.Len(t, channels, 0)
}

func testChannelStoreGetDeletedByName(t *testing.T, ss store.Store) {
	o1 := model.Channel{}
	o1.TeamId = model.NewId()
	o1.DisplayName = "Name"
	o1.Name = "zz" + model.NewId() + "b"
	o1.Type = model.CHANNEL_OPEN
	_, err := ss.Channel().Save(&o1, -1)
	require.Nil(t, err)

	now := model.GetMillis()
	err = ss.Channel().Delete(o1.Id, model.GetMillis())
	require.Nil(t, err, "channel should have been deleted")
	o1.DeleteAt = now
	o1.UpdateAt = now

	if r1, err := ss.Channel().GetDeletedByName(o1.TeamId, o1.Name); err != nil {
		t.Fatal(err)
	} else {
		if r1.ToJson() != o1.ToJson() {
			t.Fatal("invalid returned channel")
		}
	}

	if _, err := ss.Channel().GetDeletedByName(o1.TeamId, ""); err == nil {
		t.Fatal("Missing id should have failed")
	}
}

func testChannelStoreGetDeleted(t *testing.T, ss store.Store) {
	o1 := model.Channel{}
	o1.TeamId = model.NewId()
	o1.DisplayName = "Channel1"
	o1.Name = "zz" + model.NewId() + "b"
	o1.Type = model.CHANNEL_OPEN

	_, err := ss.Channel().Save(&o1, -1)
	require.Nil(t, err)

	err = ss.Channel().Delete(o1.Id, model.GetMillis())
	require.Nil(t, err, "channel should have been deleted")

	list, err := ss.Channel().GetDeleted(o1.TeamId, 0, 100)
	if err != nil {
		t.Fatal(err)
	}

	if len(*list) != 1 {
		t.Fatal("wrong list")
	}

	if (*list)[0].Name != o1.Name {
		t.Fatal("missing channel")
	}

	o2 := model.Channel{}
	o2.TeamId = o1.TeamId
	o2.DisplayName = "Channel2"
	o2.Name = "zz" + model.NewId() + "b"
	o2.Type = model.CHANNEL_OPEN
	_, err = ss.Channel().Save(&o2, -1)
	require.Nil(t, err)

	list, err = ss.Channel().GetDeleted(o1.TeamId, 0, 100)
	if err != nil {
		t.Fatal(err)
	}

	if len(*list) != 1 {
		t.Fatal("wrong list")
	}

	o3 := model.Channel{}
	o3.TeamId = o1.TeamId
	o3.DisplayName = "Channel3"
	o3.Name = "zz" + model.NewId() + "b"
	o3.Type = model.CHANNEL_OPEN

	_, err = ss.Channel().Save(&o3, -1)
	require.Nil(t, err)

	err = ss.Channel().Delete(o3.Id, model.GetMillis())
	require.Nil(t, err, "channel should have been deleted")

	list, err = ss.Channel().GetDeleted(o1.TeamId, 0, 100)
	if err != nil {
		t.Fatal(err)
	}

	if len(*list) != 2 {
		t.Fatal("wrong list length")
	}

	list, err = ss.Channel().GetDeleted(o1.TeamId, 0, 1)
	if err != nil {
		t.Fatal(err)
	}

	if len(*list) != 1 {
		t.Fatal("wrong list length")
	}

	list, err = ss.Channel().GetDeleted(o1.TeamId, 1, 1)
	if err != nil {
		t.Fatal(err)
	}

	if len(*list) != 1 {
		t.Fatal("wrong list length")
	}

}

func testChannelMemberStore(t *testing.T, ss store.Store) {
	c1 := &model.Channel{}
	c1.TeamId = model.NewId()
	c1.DisplayName = "NameName"
	c1.Name = "zz" + model.NewId() + "b"
	c1.Type = model.CHANNEL_OPEN
	c1, err := ss.Channel().Save(c1, -1)
	require.Nil(t, err)

	c1t1, _ := ss.Channel().Get(c1.Id, false)
	assert.EqualValues(t, 0, c1t1.ExtraUpdateAt, "ExtraUpdateAt should be 0")

	u1 := model.User{}
	u1.Email = MakeEmail()
	u1.Nickname = model.NewId()
	store.Must(ss.User().Save(&u1))
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: model.NewId(), UserId: u1.Id}, -1))

	u2 := model.User{}
	u2.Email = MakeEmail()
	u2.Nickname = model.NewId()
	store.Must(ss.User().Save(&u2))
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: model.NewId(), UserId: u2.Id}, -1))

	o1 := model.ChannelMember{}
	o1.ChannelId = c1.Id
	o1.UserId = u1.Id
	o1.NotifyProps = model.GetDefaultChannelNotifyProps()
	store.Must(ss.Channel().SaveMember(&o1))

	o2 := model.ChannelMember{}
	o2.ChannelId = c1.Id
	o2.UserId = u2.Id
	o2.NotifyProps = model.GetDefaultChannelNotifyProps()
	store.Must(ss.Channel().SaveMember(&o2))

	c1t2, _ := ss.Channel().Get(c1.Id, false)
	assert.EqualValues(t, 0, c1t2.ExtraUpdateAt, "ExtraUpdateAt should be 0")

	count, err := ss.Channel().GetMemberCount(o1.ChannelId, true)
	require.Nil(t, err)
	if count != 2 {
		t.Fatal("should have saved 2 members")
	}

	count, err = ss.Channel().GetMemberCount(o1.ChannelId, true)
	require.Nil(t, err)
	if count != 2 {
		t.Fatal("should have saved 2 members")
	}

	if ss.Channel().GetMemberCountFromCache(o1.ChannelId) != 2 {
		t.Fatal("should have saved 2 members")
	}

	if ss.Channel().GetMemberCountFromCache("junk") != 0 {
		t.Fatal("should have saved 0 members")
	}

	count, err = ss.Channel().GetMemberCount(o1.ChannelId, false)
	require.Nil(t, err)
	if count != 2 {
		t.Fatal("should have saved 2 members")
	}

	err = ss.Channel().RemoveMember(o2.ChannelId, o2.UserId)
	require.Nil(t, err)

	count, err = ss.Channel().GetMemberCount(o1.ChannelId, false)
	require.Nil(t, err)
	if count != 1 {
		t.Fatal("should have removed 1 member")
	}

	c1t3, _ := ss.Channel().Get(c1.Id, false)
	assert.EqualValues(t, 0, c1t3.ExtraUpdateAt, "ExtraUpdateAt should be 0")

	member, _ := ss.Channel().GetMember(o1.ChannelId, o1.UserId)
	if member.ChannelId != o1.ChannelId {
		t.Fatal("should have go member")
	}

	if err := (<-ss.Channel().SaveMember(&o1)).Err; err == nil {
		t.Fatal("Should have been a duplicate")
	}

	c1t4, _ := ss.Channel().Get(c1.Id, false)
	assert.EqualValues(t, 0, c1t4.ExtraUpdateAt, "ExtraUpdateAt should be 0")
}

func testChannelDeleteMemberStore(t *testing.T, ss store.Store) {
	c1 := &model.Channel{}
	c1.TeamId = model.NewId()
	c1.DisplayName = "NameName"
	c1.Name = "zz" + model.NewId() + "b"
	c1.Type = model.CHANNEL_OPEN
	c1, err := ss.Channel().Save(c1, -1)
	require.Nil(t, err)

	c1t1, _ := ss.Channel().Get(c1.Id, false)
	assert.EqualValues(t, 0, c1t1.ExtraUpdateAt, "ExtraUpdateAt should be 0")

	u1 := model.User{}
	u1.Email = MakeEmail()
	u1.Nickname = model.NewId()
	store.Must(ss.User().Save(&u1))
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: model.NewId(), UserId: u1.Id}, -1))

	u2 := model.User{}
	u2.Email = MakeEmail()
	u2.Nickname = model.NewId()
	store.Must(ss.User().Save(&u2))
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: model.NewId(), UserId: u2.Id}, -1))

	o1 := model.ChannelMember{}
	o1.ChannelId = c1.Id
	o1.UserId = u1.Id
	o1.NotifyProps = model.GetDefaultChannelNotifyProps()
	store.Must(ss.Channel().SaveMember(&o1))

	o2 := model.ChannelMember{}
	o2.ChannelId = c1.Id
	o2.UserId = u2.Id
	o2.NotifyProps = model.GetDefaultChannelNotifyProps()
	store.Must(ss.Channel().SaveMember(&o2))

	c1t2, _ := ss.Channel().Get(c1.Id, false)
	assert.EqualValues(t, 0, c1t2.ExtraUpdateAt, "ExtraUpdateAt should be 0")

	count, err := ss.Channel().GetMemberCount(o1.ChannelId, false)
	require.Nil(t, err)
	if count != 2 {
		t.Fatal("should have saved 2 members")
	}

	err = ss.Channel().PermanentDeleteMembersByUser(o2.UserId)
	require.Nil(t, err)

	count, err = ss.Channel().GetMemberCount(o1.ChannelId, false)
	require.Nil(t, err)
	if count != 1 {
		t.Fatal("should have removed 1 member")
	}

	if err = ss.Channel().PermanentDeleteMembersByChannel(o1.ChannelId); err != nil {
		t.Fatal(err)
	}

	count, err = ss.Channel().GetMemberCount(o1.ChannelId, false)
	require.Nil(t, err)
	if count != 0 {
		t.Fatal("should have removed all members")
	}
}

func testChannelStoreGetChannels(t *testing.T, ss store.Store) {
	o2 := model.Channel{}
	o2.TeamId = model.NewId()
	o2.DisplayName = "Channel2"
	o2.Name = "zz" + model.NewId() + "b"
	o2.Type = model.CHANNEL_OPEN
	_, err := ss.Channel().Save(&o2, -1)
	require.Nil(t, err)

	o1 := model.Channel{}
	o1.TeamId = model.NewId()
	o1.DisplayName = "Channel1"
	o1.Name = "zz" + model.NewId() + "b"
	o1.Type = model.CHANNEL_OPEN
	_, err = ss.Channel().Save(&o1, -1)
	require.Nil(t, err)

	m1 := model.ChannelMember{}
	m1.ChannelId = o1.Id
	m1.UserId = model.NewId()
	m1.NotifyProps = model.GetDefaultChannelNotifyProps()
	store.Must(ss.Channel().SaveMember(&m1))

	m2 := model.ChannelMember{}
	m2.ChannelId = o1.Id
	m2.UserId = model.NewId()
	m2.NotifyProps = model.GetDefaultChannelNotifyProps()
	store.Must(ss.Channel().SaveMember(&m2))

	m3 := model.ChannelMember{}
	m3.ChannelId = o2.Id
	m3.UserId = model.NewId()
	m3.NotifyProps = model.GetDefaultChannelNotifyProps()
	store.Must(ss.Channel().SaveMember(&m3))

	list, err := ss.Channel().GetChannels(o1.TeamId, m1.UserId, false)
	require.Nil(t, err)

	if (*list)[0].Id != o1.Id {
		t.Fatal("missing channel")
	}

	ids, _ := ss.Channel().GetAllChannelMembersForUser(m1.UserId, false, false)
	if _, ok := ids[o1.Id]; !ok {
		t.Fatal("missing channel")
	}

	ids2, _ := ss.Channel().GetAllChannelMembersForUser(m1.UserId, true, false)
	if _, ok := ids2[o1.Id]; !ok {
		t.Fatal("missing channel")
	}

	ids3, _ := ss.Channel().GetAllChannelMembersForUser(m1.UserId, true, false)
	if _, ok := ids3[o1.Id]; !ok {
		t.Fatal("missing channel")
	}

	if !ss.Channel().IsUserInChannelUseCache(m1.UserId, o1.Id) {
		t.Fatal("missing channel")
	}

	if ss.Channel().IsUserInChannelUseCache(m1.UserId, o2.Id) {
		t.Fatal("missing channel")
	}

	if ss.Channel().IsUserInChannelUseCache(m1.UserId, "blahblah") {
		t.Fatal("missing channel")
	}

	if ss.Channel().IsUserInChannelUseCache("blahblah", "blahblah") {
		t.Fatal("missing channel")
	}

	ss.Channel().InvalidateAllChannelMembersForUser(m1.UserId)
}

func testChannelStoreGetAllChannels(t *testing.T, ss store.Store, s SqlSupplier) {
	cleanupChannels(t, ss)

	t1 := model.Team{}
	t1.DisplayName = "Name"
	t1.Name = model.NewId()
	t1.Email = MakeEmail()
	t1.Type = model.TEAM_OPEN
	_, err := ss.Team().Save(&t1)
	require.Nil(t, err)

	t2 := model.Team{}
	t2.DisplayName = "Name2"
	t2.Name = model.NewId()
	t2.Email = MakeEmail()
	t2.Type = model.TEAM_OPEN
	_, err = ss.Team().Save(&t2)
	require.Nil(t, err)

	c1 := model.Channel{}
	c1.TeamId = t1.Id
	c1.DisplayName = "Channel1" + model.NewId()
	c1.Name = "zz" + model.NewId() + "b"
	c1.Type = model.CHANNEL_OPEN
	_, err = ss.Channel().Save(&c1, -1)
	require.Nil(t, err)

	group := &model.Group{
		Name:        model.NewId(),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewId(),
	}
	store.Must(ss.Group().Create(group))
	_, err = ss.Group().CreateGroupSyncable(model.NewGroupChannel(group.Id, c1.Id, true))
	require.Nil(t, err)

	c2 := model.Channel{}
	c2.TeamId = t1.Id
	c2.DisplayName = "Channel2" + model.NewId()
	c2.Name = "zz" + model.NewId() + "b"
	c2.Type = model.CHANNEL_OPEN
	_, err = ss.Channel().Save(&c2, -1)
	require.Nil(t, err)
	c2.DeleteAt = model.GetMillis()
	c2.UpdateAt = c2.DeleteAt
	err = ss.Channel().Delete(c2.Id, c2.DeleteAt)
	require.Nil(t, err, "channel should have been deleted")

	c3 := model.Channel{}
	c3.TeamId = t2.Id
	c3.DisplayName = "Channel3" + model.NewId()
	c3.Name = "zz" + model.NewId() + "b"
	c3.Type = model.CHANNEL_PRIVATE
	_, err = ss.Channel().Save(&c3, -1)
	require.Nil(t, err)

	_, err = ss.Channel().CreateDirectChannel(model.NewId(), model.NewId())
	require.Nil(t, err)

	userIds := []string{model.NewId(), model.NewId(), model.NewId()}

	c5 := model.Channel{}
	c5.Name = model.GetGroupNameFromUserIds(userIds)
	c5.DisplayName = "GroupChannel" + model.NewId()
	c5.Name = "zz" + model.NewId() + "b"
	c5.Type = model.CHANNEL_GROUP
	_, err = ss.Channel().Save(&c5, -1)
	require.Nil(t, err)

	list, err := ss.Channel().GetAllChannels(0, 10, store.ChannelSearchOpts{})
	require.Nil(t, err)
	assert.Len(t, *list, 2)
	assert.Equal(t, (*list)[0].Id, c1.Id)
	assert.Equal(t, (*list)[0].TeamDisplayName, "Name")
	assert.Equal(t, (*list)[1].Id, c3.Id)
	assert.Equal(t, (*list)[1].TeamDisplayName, "Name2")

	count1, err := ss.Channel().GetAllChannelsCount(store.ChannelSearchOpts{})
	require.Nil(t, err)

	list, err = ss.Channel().GetAllChannels(0, 10, store.ChannelSearchOpts{IncludeDeleted: true})
	require.Nil(t, err)
	assert.Len(t, *list, 3)
	assert.Equal(t, (*list)[0].Id, c1.Id)
	assert.Equal(t, (*list)[0].TeamDisplayName, "Name")
	assert.Equal(t, (*list)[1].Id, c2.Id)
	assert.Equal(t, (*list)[2].Id, c3.Id)

	count2, err := ss.Channel().GetAllChannelsCount(store.ChannelSearchOpts{IncludeDeleted: true})
	require.Nil(t, err)
	require.True(t, func() bool {
		return count2 > count1
	}())

	list, err = ss.Channel().GetAllChannels(0, 1, store.ChannelSearchOpts{IncludeDeleted: true})
	require.Nil(t, err)
	assert.Len(t, *list, 1)
	assert.Equal(t, (*list)[0].Id, c1.Id)
	assert.Equal(t, (*list)[0].TeamDisplayName, "Name")

	// Not associated to group
	list, err = ss.Channel().GetAllChannels(0, 10, store.ChannelSearchOpts{NotAssociatedToGroup: group.Id})
	require.Nil(t, err)
	assert.Len(t, *list, 1)

	// Exclude channel names
	list, err = ss.Channel().GetAllChannels(0, 10, store.ChannelSearchOpts{ExcludeChannelNames: []string{c1.Name}})
	require.Nil(t, err)
	assert.Len(t, *list, 1)

	// Manually truncate Channels table until testlib can handle cleanups
	s.GetMaster().Exec("TRUNCATE Channels")
}

func testChannelStoreGetMoreChannels(t *testing.T, ss store.Store) {
	teamId := model.NewId()
	otherTeamId := model.NewId()
	userId := model.NewId()
	otherUserId1 := model.NewId()
	otherUserId2 := model.NewId()

	// o1 is a channel on the team to which the user (and the other user 1) belongs
	o1 := model.Channel{
		TeamId:      teamId,
		DisplayName: "Channel1",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_OPEN,
	}
	_, err := ss.Channel().Save(&o1, -1)
	require.Nil(t, err)

	store.Must(ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   o1.Id,
		UserId:      userId,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}))

	store.Must(ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   o1.Id,
		UserId:      otherUserId1,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}))

	// o2 is a channel on the other team to which the user belongs
	o2 := model.Channel{
		TeamId:      otherTeamId,
		DisplayName: "Channel2",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_OPEN,
	}
	_, err = ss.Channel().Save(&o2, -1)
	require.Nil(t, err)

	store.Must(ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   o2.Id,
		UserId:      otherUserId2,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}))

	// o3 is a channel on the team to which the user does not belong, and thus should show up
	// in "more channels"
	o3 := model.Channel{
		TeamId:      teamId,
		DisplayName: "ChannelA",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_OPEN,
	}
	_, err = ss.Channel().Save(&o3, -1)
	require.Nil(t, err)

	// o4 is a private channel on the team to which the user does not belong
	o4 := model.Channel{
		TeamId:      teamId,
		DisplayName: "ChannelB",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_PRIVATE,
	}
	_, err = ss.Channel().Save(&o4, -1)
	require.Nil(t, err)

	// o5 is another private channel on the team to which the user does belong
	o5 := model.Channel{
		TeamId:      teamId,
		DisplayName: "ChannelC",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_PRIVATE,
	}
	_, err = ss.Channel().Save(&o5, -1)
	require.Nil(t, err)

	store.Must(ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   o5.Id,
		UserId:      userId,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}))

	t.Run("only o3 listed in more channels", func(t *testing.T) {
		list, channelErr := ss.Channel().GetMoreChannels(teamId, userId, 0, 100)
		require.Nil(t, channelErr)
		require.Equal(t, &model.ChannelList{&o3}, list)
	})

	// o6 is another channel on the team to which the user does not belong, and would thus
	// start showing up in "more channels".
	o6 := model.Channel{
		TeamId:      teamId,
		DisplayName: "ChannelD",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_OPEN,
	}
	_, err = ss.Channel().Save(&o6, -1)
	require.Nil(t, err)

	// o7 is another channel on the team to which the user does not belong, but is deleted,
	// and thus would not start showing up in "more channels"
	o7 := model.Channel{
		TeamId:      teamId,
		DisplayName: "ChannelD",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_OPEN,
	}
	_, err = ss.Channel().Save(&o7, -1)
	require.Nil(t, err)

	err = ss.Channel().Delete(o7.Id, model.GetMillis())
	require.Nil(t, err, "channel should have been deleted")

	t.Run("both o3 and o6 listed in more channels", func(t *testing.T) {
		list, err := ss.Channel().GetMoreChannels(teamId, userId, 0, 100)
		require.Nil(t, err)
		require.Equal(t, &model.ChannelList{&o3, &o6}, list)
	})

	t.Run("only o3 listed in more channels with offset 0, limit 1", func(t *testing.T) {
		list, err := ss.Channel().GetMoreChannels(teamId, userId, 0, 1)
		require.Nil(t, err)
		require.Equal(t, &model.ChannelList{&o3}, list)
	})

	t.Run("only o6 listed in more channels with offset 1, limit 1", func(t *testing.T) {
		list, err := ss.Channel().GetMoreChannels(teamId, userId, 1, 1)
		require.Nil(t, err)
		require.Equal(t, &model.ChannelList{&o6}, list)
	})

	t.Run("verify analytics for open channels", func(t *testing.T) {
		count, err := ss.Channel().AnalyticsTypeCount(teamId, model.CHANNEL_OPEN)
		require.Nil(t, err)
		require.EqualValues(t, 4, count)
	})

	t.Run("verify analytics for private channels", func(t *testing.T) {
		count, err := ss.Channel().AnalyticsTypeCount(teamId, model.CHANNEL_PRIVATE)
		require.Nil(t, err)
		require.EqualValues(t, 2, count)
	})
}

func testChannelStoreGetPublicChannelsForTeam(t *testing.T, ss store.Store) {
	teamId := model.NewId()

	// o1 is a public channel on the team
	o1 := model.Channel{
		TeamId:      teamId,
		DisplayName: "OpenChannel1Team1",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_OPEN,
	}
	_, err := ss.Channel().Save(&o1, -1)
	require.Nil(t, err)

	// o2 is a public channel on another team
	o2 := model.Channel{
		TeamId:      model.NewId(),
		DisplayName: "OpenChannel1Team2",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_OPEN,
	}
	_, err = ss.Channel().Save(&o2, -1)
	require.Nil(t, err)

	// o3 is a private channel on the team
	o3 := model.Channel{
		TeamId:      teamId,
		DisplayName: "PrivateChannel1Team1",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_PRIVATE,
	}
	_, err = ss.Channel().Save(&o3, -1)
	require.Nil(t, err)

	t.Run("only o1 initially listed in public channels", func(t *testing.T) {
		list, channelErr := ss.Channel().GetPublicChannelsForTeam(teamId, 0, 100)
		require.Nil(t, channelErr)
		require.Equal(t, &model.ChannelList{&o1}, list)
	})

	// o4 is another public channel on the team
	o4 := model.Channel{
		TeamId:      teamId,
		DisplayName: "OpenChannel2Team1",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_OPEN,
	}
	_, err = ss.Channel().Save(&o4, -1)
	require.Nil(t, err)

	// o5 is another public, but deleted channel on the team
	o5 := model.Channel{
		TeamId:      teamId,
		DisplayName: "OpenChannel3Team1",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_OPEN,
	}
	_, err = ss.Channel().Save(&o5, -1)
	require.Nil(t, err)
	err = ss.Channel().Delete(o5.Id, model.GetMillis())
	require.Nil(t, err, "channel should have been deleted")

	t.Run("both o1 and o4 listed in public channels", func(t *testing.T) {
		list, err := ss.Channel().GetPublicChannelsForTeam(teamId, 0, 100)
		require.Nil(t, err)
		require.Equal(t, &model.ChannelList{&o1, &o4}, list)
	})

	t.Run("only o1 listed in public channels with offset 0, limit 1", func(t *testing.T) {
		list, err := ss.Channel().GetPublicChannelsForTeam(teamId, 0, 1)
		require.Nil(t, err)
		require.Equal(t, &model.ChannelList{&o1}, list)
	})

	t.Run("only o4 listed in public channels with offset 1, limit 1", func(t *testing.T) {
		list, err := ss.Channel().GetPublicChannelsForTeam(teamId, 1, 1)
		require.Nil(t, err)
		require.Equal(t, &model.ChannelList{&o4}, list)
	})

	t.Run("verify analytics for open channels", func(t *testing.T) {
		count, err := ss.Channel().AnalyticsTypeCount(teamId, model.CHANNEL_OPEN)
		require.Nil(t, err)
		require.EqualValues(t, 3, count)
	})

	t.Run("verify analytics for private channels", func(t *testing.T) {
		count, err := ss.Channel().AnalyticsTypeCount(teamId, model.CHANNEL_PRIVATE)
		require.Nil(t, err)
		require.EqualValues(t, 1, count)
	})
}

func testChannelStoreGetPublicChannelsByIdsForTeam(t *testing.T, ss store.Store) {
	teamId := model.NewId()

	// oc1 is a public channel on the team
	oc1 := model.Channel{
		TeamId:      teamId,
		DisplayName: "OpenChannel1Team1",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_OPEN,
	}
	_, err := ss.Channel().Save(&oc1, -1)
	require.Nil(t, err)

	// oc2 is a public channel on another team
	oc2 := model.Channel{
		TeamId:      model.NewId(),
		DisplayName: "OpenChannel2TeamOther",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_OPEN,
	}
	_, err = ss.Channel().Save(&oc2, -1)
	require.Nil(t, err)

	// pc3 is a private channel on the team
	pc3 := model.Channel{
		TeamId:      teamId,
		DisplayName: "PrivateChannel3Team1",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_PRIVATE,
	}
	_, err = ss.Channel().Save(&pc3, -1)
	require.Nil(t, err)

	t.Run("oc1 by itself should be found as a public channel in the team", func(t *testing.T) {
		list, channelErr := ss.Channel().GetPublicChannelsByIdsForTeam(teamId, []string{oc1.Id})
		require.Nil(t, channelErr)
		require.Equal(t, &model.ChannelList{&oc1}, list)
	})

	t.Run("only oc1, among others, should be found as a public channel in the team", func(t *testing.T) {
		list, channelErr := ss.Channel().GetPublicChannelsByIdsForTeam(teamId, []string{oc1.Id, oc2.Id, model.NewId(), pc3.Id})
		require.Nil(t, channelErr)
		require.Equal(t, &model.ChannelList{&oc1}, list)
	})

	// oc4 is another public channel on the team
	oc4 := model.Channel{
		TeamId:      teamId,
		DisplayName: "OpenChannel4Team1",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_OPEN,
	}
	_, err = ss.Channel().Save(&oc4, -1)
	require.Nil(t, err)

	// oc4 is another public, but deleted channel on the team
	oc5 := model.Channel{
		TeamId:      teamId,
		DisplayName: "OpenChannel4Team1",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_OPEN,
	}
	_, err = ss.Channel().Save(&oc5, -1)
	require.Nil(t, err)

	err = ss.Channel().Delete(oc5.Id, model.GetMillis())
	require.Nil(t, err, "channel should have been deleted")

	t.Run("only oc1 and oc4, among others, should be found as a public channel in the team", func(t *testing.T) {
		list, err := ss.Channel().GetPublicChannelsByIdsForTeam(teamId, []string{oc1.Id, oc2.Id, model.NewId(), pc3.Id, oc4.Id})
		require.Nil(t, err)
		require.Equal(t, &model.ChannelList{&oc1, &oc4}, list)
	})

	t.Run("random channel id should not be found as a public channel in the team", func(t *testing.T) {
		_, err := ss.Channel().GetPublicChannelsByIdsForTeam(teamId, []string{model.NewId()})
		require.NotNil(t, err)
		require.Equal(t, err.Id, "store.sql_channel.get_channels_by_ids.not_found.app_error")
	})
}

func testChannelStoreGetChannelCounts(t *testing.T, ss store.Store) {
	o2 := model.Channel{}
	o2.TeamId = model.NewId()
	o2.DisplayName = "Channel2"
	o2.Name = "zz" + model.NewId() + "b"
	o2.Type = model.CHANNEL_OPEN
	_, err := ss.Channel().Save(&o2, -1)
	require.Nil(t, err)

	o1 := model.Channel{}
	o1.TeamId = model.NewId()
	o1.DisplayName = "Channel1"
	o1.Name = "zz" + model.NewId() + "b"
	o1.Type = model.CHANNEL_OPEN
	_, err = ss.Channel().Save(&o1, -1)
	require.Nil(t, err)

	m1 := model.ChannelMember{}
	m1.ChannelId = o1.Id
	m1.UserId = model.NewId()
	m1.NotifyProps = model.GetDefaultChannelNotifyProps()
	store.Must(ss.Channel().SaveMember(&m1))

	m2 := model.ChannelMember{}
	m2.ChannelId = o1.Id
	m2.UserId = model.NewId()
	m2.NotifyProps = model.GetDefaultChannelNotifyProps()
	store.Must(ss.Channel().SaveMember(&m2))

	m3 := model.ChannelMember{}
	m3.ChannelId = o2.Id
	m3.UserId = model.NewId()
	m3.NotifyProps = model.GetDefaultChannelNotifyProps()
	store.Must(ss.Channel().SaveMember(&m3))

	counts, _ := ss.Channel().GetChannelCounts(o1.TeamId, m1.UserId)

	if len(counts.Counts) != 1 {
		t.Fatal("wrong number of counts")
	}

	if len(counts.UpdateTimes) != 1 {
		t.Fatal("wrong number of update times")
	}
}

func testChannelStoreGetMembersForUser(t *testing.T, ss store.Store) {
	t1 := model.Team{}
	t1.DisplayName = "Name"
	t1.Name = model.NewId()
	t1.Email = MakeEmail()
	t1.Type = model.TEAM_OPEN
	_, err := ss.Team().Save(&t1)
	require.Nil(t, err)

	o1 := model.Channel{}
	o1.TeamId = t1.Id
	o1.DisplayName = "Channel1"
	o1.Name = "zz" + model.NewId() + "b"
	o1.Type = model.CHANNEL_OPEN
	_, err = ss.Channel().Save(&o1, -1)
	require.Nil(t, err)

	o2 := model.Channel{}
	o2.TeamId = o1.TeamId
	o2.DisplayName = "Channel2"
	o2.Name = "zz" + model.NewId() + "b"
	o2.Type = model.CHANNEL_OPEN
	_, err = ss.Channel().Save(&o2, -1)
	require.Nil(t, err)

	m1 := model.ChannelMember{}
	m1.ChannelId = o1.Id
	m1.UserId = model.NewId()
	m1.NotifyProps = model.GetDefaultChannelNotifyProps()
	store.Must(ss.Channel().SaveMember(&m1))

	m2 := model.ChannelMember{}
	m2.ChannelId = o2.Id
	m2.UserId = m1.UserId
	m2.NotifyProps = model.GetDefaultChannelNotifyProps()
	store.Must(ss.Channel().SaveMember(&m2))

	members, err := ss.Channel().GetMembersForUser(o1.TeamId, m1.UserId)
	require.Nil(t, err)

	// no unread messages
	if len(*members) != 2 {
		t.Fatal("wrong number of members")
	}
}

func testChannelStoreGetMembersForUserWithPagination(t *testing.T, ss store.Store) {
	t1 := model.Team{}
	t1.DisplayName = "Name"
	t1.Name = model.NewId()
	t1.Email = MakeEmail()
	t1.Type = model.TEAM_OPEN
	_, err := ss.Team().Save(&t1)
	require.Nil(t, err)

	o1 := model.Channel{}
	o1.TeamId = t1.Id
	o1.DisplayName = "Channel1"
	o1.Name = "zz" + model.NewId() + "b"
	o1.Type = model.CHANNEL_OPEN
	_, err = ss.Channel().Save(&o1, -1)
	require.Nil(t, err)

	o2 := model.Channel{}
	o2.TeamId = o1.TeamId
	o2.DisplayName = "Channel2"
	o2.Name = "zz" + model.NewId() + "b"
	o2.Type = model.CHANNEL_OPEN
	_, err = ss.Channel().Save(&o2, -1)
	require.Nil(t, err)

	m1 := model.ChannelMember{}
	m1.ChannelId = o1.Id
	m1.UserId = model.NewId()
	m1.NotifyProps = model.GetDefaultChannelNotifyProps()
	store.Must(ss.Channel().SaveMember(&m1))

	m2 := model.ChannelMember{}
	m2.ChannelId = o2.Id
	m2.UserId = m1.UserId
	m2.NotifyProps = model.GetDefaultChannelNotifyProps()
	store.Must(ss.Channel().SaveMember(&m2))

	members, err := ss.Channel().GetMembersForUserWithPagination(o1.TeamId, m1.UserId, 0, 1)
	require.Nil(t, err)
	assert.Len(t, *members, 1)

	members, err = ss.Channel().GetMembersForUserWithPagination(o1.TeamId, m1.UserId, 1, 1)
	require.Nil(t, err)
	assert.Len(t, *members, 1)
}

func testChannelStoreUpdateLastViewedAt(t *testing.T, ss store.Store) {
	o1 := model.Channel{}
	o1.TeamId = model.NewId()
	o1.DisplayName = "Channel1"
	o1.Name = "zz" + model.NewId() + "b"
	o1.Type = model.CHANNEL_OPEN
	o1.TotalMsgCount = 25
	o1.LastPostAt = 12345
	_, err := ss.Channel().Save(&o1, -1)
	require.Nil(t, err)

	m1 := model.ChannelMember{}
	m1.ChannelId = o1.Id
	m1.UserId = model.NewId()
	m1.NotifyProps = model.GetDefaultChannelNotifyProps()
	store.Must(ss.Channel().SaveMember(&m1))

	o2 := model.Channel{}
	o2.TeamId = model.NewId()
	o2.DisplayName = "Channel1"
	o2.Name = "zz" + model.NewId() + "c"
	o2.Type = model.CHANNEL_OPEN
	o2.TotalMsgCount = 26
	o2.LastPostAt = 123456
	_, err = ss.Channel().Save(&o2, -1)
	require.Nil(t, err)

	m2 := model.ChannelMember{}
	m2.ChannelId = o2.Id
	m2.UserId = m1.UserId
	m2.NotifyProps = model.GetDefaultChannelNotifyProps()
	store.Must(ss.Channel().SaveMember(&m2))

	var times map[string]int64
	if times, err = ss.Channel().UpdateLastViewedAt([]string{m1.ChannelId}, m1.UserId); err != nil {
		t.Fatal("failed to update", err)
	} else if times[o1.Id] != o1.LastPostAt {
		t.Fatal("last viewed at time incorrect")
	}

	if times, err = ss.Channel().UpdateLastViewedAt([]string{m1.ChannelId, m2.ChannelId}, m1.UserId); err != nil {
		t.Fatal("failed to update", err)
	} else if times[o2.Id] != o2.LastPostAt {
		t.Fatal("last viewed at time incorrect")
	}

	rm1, err := ss.Channel().GetMember(m1.ChannelId, m1.UserId)
	assert.Nil(t, err)
	assert.Equal(t, rm1.LastViewedAt, o1.LastPostAt)
	assert.Equal(t, rm1.LastUpdateAt, o1.LastPostAt)
	assert.Equal(t, rm1.MsgCount, o1.TotalMsgCount)

	rm2, err := ss.Channel().GetMember(m2.ChannelId, m2.UserId)
	assert.Nil(t, err)
	assert.Equal(t, rm2.LastViewedAt, o2.LastPostAt)
	assert.Equal(t, rm2.LastUpdateAt, o2.LastPostAt)
	assert.Equal(t, rm2.MsgCount, o2.TotalMsgCount)

	if _, err := ss.Channel().UpdateLastViewedAt([]string{m1.ChannelId}, "missing id"); err != nil {
		t.Fatal("failed to update")
	}
}

func testChannelStoreIncrementMentionCount(t *testing.T, ss store.Store) {
	o1 := model.Channel{}
	o1.TeamId = model.NewId()
	o1.DisplayName = "Channel1"
	o1.Name = "zz" + model.NewId() + "b"
	o1.Type = model.CHANNEL_OPEN
	o1.TotalMsgCount = 25
	_, err := ss.Channel().Save(&o1, -1)
	require.Nil(t, err)

	m1 := model.ChannelMember{}
	m1.ChannelId = o1.Id
	m1.UserId = model.NewId()
	m1.NotifyProps = model.GetDefaultChannelNotifyProps()
	store.Must(ss.Channel().SaveMember(&m1))

	err = ss.Channel().IncrementMentionCount(m1.ChannelId, m1.UserId)
	if err != nil {
		t.Fatal("failed to update")
	}

	err = ss.Channel().IncrementMentionCount(m1.ChannelId, "missing id")
	if err != nil {
		t.Fatal("failed to update")
	}

	err = ss.Channel().IncrementMentionCount("missing id", m1.UserId)
	if err != nil {
		t.Fatal("failed to update")
	}

	err = ss.Channel().IncrementMentionCount("missing id", "missing id")
	if err != nil {
		t.Fatal("failed to update")
	}
}

func testUpdateChannelMember(t *testing.T, ss store.Store) {
	userId := model.NewId()

	c1 := &model.Channel{
		TeamId:      model.NewId(),
		DisplayName: model.NewId(),
		Name:        model.NewId(),
		Type:        model.CHANNEL_OPEN,
	}
	_, err := ss.Channel().Save(c1, -1)
	require.Nil(t, err)

	m1 := &model.ChannelMember{
		ChannelId:   c1.Id,
		UserId:      userId,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}
	store.Must(ss.Channel().SaveMember(m1))

	m1.NotifyProps["test"] = "sometext"
	if _, err := ss.Channel().UpdateMember(m1); err != nil {
		t.Fatal(err)
	}

	m1.UserId = ""
	if _, err := ss.Channel().UpdateMember(m1); err == nil {
		t.Fatal("bad user id - should fail")
	}
}

func testGetMember(t *testing.T, ss store.Store) {
	userId := model.NewId()

	c1 := &model.Channel{
		TeamId:      model.NewId(),
		DisplayName: model.NewId(),
		Name:        model.NewId(),
		Type:        model.CHANNEL_OPEN,
	}
	_, err := ss.Channel().Save(c1, -1)
	require.Nil(t, err)

	c2 := &model.Channel{
		TeamId:      c1.TeamId,
		DisplayName: model.NewId(),
		Name:        model.NewId(),
		Type:        model.CHANNEL_OPEN,
	}
	_, err = ss.Channel().Save(c2, -1)
	require.Nil(t, err)

	m1 := &model.ChannelMember{
		ChannelId:   c1.Id,
		UserId:      userId,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}
	store.Must(ss.Channel().SaveMember(m1))

	m2 := &model.ChannelMember{
		ChannelId:   c2.Id,
		UserId:      userId,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}
	store.Must(ss.Channel().SaveMember(m2))

	if _, err := ss.Channel().GetMember(model.NewId(), userId); err == nil {
		t.Fatal("should've failed to get member for non-existent channel")
	}

	if _, err := ss.Channel().GetMember(c1.Id, model.NewId()); err == nil {
		t.Fatal("should've failed to get member for non-existent user")
	}

	if member, err := ss.Channel().GetMember(c1.Id, userId); err != nil {
		t.Fatal("shouldn't have errored when getting member", err)
	} else if member.ChannelId != c1.Id {
		t.Fatal("should've gotten member of channel 1")
	} else if member.UserId != userId {
		t.Fatal("should've gotten member for user")
	}

	if member, err := ss.Channel().GetMember(c2.Id, userId); err != nil {
		t.Fatal("shouldn't have errored when getting member", err)
	} else if member.ChannelId != c2.Id {
		t.Fatal("should've gotten member of channel 2")
	} else if member.UserId != userId {
		t.Fatal("should've gotten member for user")
	}

	if props, err := ss.Channel().GetAllChannelMembersNotifyPropsForChannel(c2.Id, false); err != nil {
		t.Fatal(err)
	} else if len(props) == 0 {
		t.Fatal("should not be empty")
	}

	if props, err := ss.Channel().GetAllChannelMembersNotifyPropsForChannel(c2.Id, true); err != nil {
		t.Fatal(err)
	} else if len(props) == 0 {
		t.Fatal("should not be empty")

	}

	ss.Channel().InvalidateCacheForChannelMembersNotifyProps(c2.Id)
}

func testChannelStoreGetMemberForPost(t *testing.T, ss store.Store) {
	ch := &model.Channel{
		TeamId:      model.NewId(),
		DisplayName: "Name",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_OPEN,
	}

	o1, err := ss.Channel().Save(ch, -1)
	require.Nil(t, err)

	m1 := store.Must(ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   o1.Id,
		UserId:      model.NewId(),
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	})).(*model.ChannelMember)

	p1, err := ss.Post().Save(&model.Post{
		UserId:    model.NewId(),
		ChannelId: o1.Id,
		Message:   "test",
	})
	require.Nil(t, err)

	if r1, err := ss.Channel().GetMemberForPost(p1.Id, m1.UserId); err != nil {
		t.Fatal(err)
	} else if r1.ToJson() != m1.ToJson() {
		t.Fatal("invalid returned channel member")
	}

	if _, err := ss.Channel().GetMemberForPost(p1.Id, model.NewId()); err == nil {
		t.Fatal("shouldn't have returned a member")
	}
}

func testGetMemberCount(t *testing.T, ss store.Store) {
	teamId := model.NewId()

	c1 := model.Channel{
		TeamId:      teamId,
		DisplayName: "Channel1",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_OPEN,
	}
	_, err := ss.Channel().Save(&c1, -1)
	require.Nil(t, err)

	c2 := model.Channel{
		TeamId:      teamId,
		DisplayName: "Channel2",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_OPEN,
	}
	_, err = ss.Channel().Save(&c2, -1)
	require.Nil(t, err)

	u1 := &model.User{
		Email:    MakeEmail(),
		DeleteAt: 0,
	}
	store.Must(ss.User().Save(u1))
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1))

	m1 := model.ChannelMember{
		ChannelId:   c1.Id,
		UserId:      u1.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}
	store.Must(ss.Channel().SaveMember(&m1))

	if count, err := ss.Channel().GetMemberCount(c1.Id, false); err != nil {
		t.Fatalf("failed to get member count: %v", err)
	} else if count != 1 {
		t.Fatalf("got incorrect member count %v", count)
	}

	u2 := model.User{
		Email:    MakeEmail(),
		DeleteAt: 0,
	}
	store.Must(ss.User().Save(&u2))
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u2.Id}, -1))

	m2 := model.ChannelMember{
		ChannelId:   c1.Id,
		UserId:      u2.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}
	store.Must(ss.Channel().SaveMember(&m2))

	if count, err := ss.Channel().GetMemberCount(c1.Id, false); err != nil {
		t.Fatalf("failed to get member count: %v", err)
	} else if count != 2 {
		t.Fatalf("got incorrect member count %v", count)
	}

	// make sure members of other channels aren't counted
	u3 := model.User{
		Email:    MakeEmail(),
		DeleteAt: 0,
	}
	store.Must(ss.User().Save(&u3))
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u3.Id}, -1))

	m3 := model.ChannelMember{
		ChannelId:   c2.Id,
		UserId:      u3.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}
	store.Must(ss.Channel().SaveMember(&m3))

	if count, err := ss.Channel().GetMemberCount(c1.Id, false); err != nil {
		t.Fatalf("failed to get member count: %v", err)
	} else if count != 2 {
		t.Fatalf("got incorrect member count %v", count)
	}

	// make sure inactive users aren't counted
	u4 := &model.User{
		Email:    MakeEmail(),
		DeleteAt: 10000,
	}
	store.Must(ss.User().Save(u4))
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u4.Id}, -1))

	m4 := model.ChannelMember{
		ChannelId:   c1.Id,
		UserId:      u4.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}
	store.Must(ss.Channel().SaveMember(&m4))

	if count, err := ss.Channel().GetMemberCount(c1.Id, false); err != nil {
		t.Fatalf("failed to get member count: %v", err)
	} else if count != 2 {
		t.Fatalf("got incorrect member count %v", count)
	}
}

func testChannelStoreSearchMore(t *testing.T, ss store.Store) {
	teamId := model.NewId()
	otherTeamId := model.NewId()

	o1 := model.Channel{
		TeamId:      teamId,
		DisplayName: "ChannelA",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_OPEN,
	}
	_, err := ss.Channel().Save(&o1, -1)
	require.Nil(t, err)

	m1 := model.ChannelMember{
		ChannelId:   o1.Id,
		UserId:      model.NewId(),
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}
	store.Must(ss.Channel().SaveMember(&m1))

	m2 := model.ChannelMember{
		ChannelId:   o1.Id,
		UserId:      model.NewId(),
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}
	store.Must(ss.Channel().SaveMember(&m2))

	o2 := model.Channel{
		TeamId:      otherTeamId,
		DisplayName: "Channel2",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_OPEN,
	}
	_, err = ss.Channel().Save(&o2, -1)
	require.Nil(t, err)

	m3 := model.ChannelMember{
		ChannelId:   o2.Id,
		UserId:      model.NewId(),
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}
	store.Must(ss.Channel().SaveMember(&m3))

	o3 := model.Channel{
		TeamId:      teamId,
		DisplayName: "ChannelA",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_OPEN,
	}
	_, err = ss.Channel().Save(&o3, -1)
	require.Nil(t, err)

	o4 := model.Channel{
		TeamId:      teamId,
		DisplayName: "ChannelB",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_PRIVATE,
	}
	_, err = ss.Channel().Save(&o4, -1)
	require.Nil(t, err)

	o5 := model.Channel{
		TeamId:      teamId,
		DisplayName: "ChannelC",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_PRIVATE,
	}
	_, err = ss.Channel().Save(&o5, -1)
	require.Nil(t, err)

	o6 := model.Channel{
		TeamId:      teamId,
		DisplayName: "Off-Topic",
		Name:        "off-topic",
		Type:        model.CHANNEL_OPEN,
	}
	_, err = ss.Channel().Save(&o6, -1)
	require.Nil(t, err)

	o7 := model.Channel{
		TeamId:      teamId,
		DisplayName: "Off-Set",
		Name:        "off-set",
		Type:        model.CHANNEL_OPEN,
	}
	_, err = ss.Channel().Save(&o7, -1)
	require.Nil(t, err)

	o8 := model.Channel{
		TeamId:      teamId,
		DisplayName: "Off-Limit",
		Name:        "off-limit",
		Type:        model.CHANNEL_PRIVATE,
	}
	_, err = ss.Channel().Save(&o8, -1)
	require.Nil(t, err)

	o9 := model.Channel{
		TeamId:      teamId,
		DisplayName: "Channel With Purpose",
		Purpose:     "This can now be searchable!",
		Name:        "with-purpose",
		Type:        model.CHANNEL_OPEN,
	}
	_, err = ss.Channel().Save(&o9, -1)
	require.Nil(t, err)

	o10 := model.Channel{
		TeamId:      teamId,
		DisplayName: "ChannelA",
		Name:        "channel-a-deleted",
		Type:        model.CHANNEL_OPEN,
	}
	_, err = ss.Channel().Save(&o10, -1)
	require.Nil(t, err)

	o10.DeleteAt = model.GetMillis()
	o10.UpdateAt = o10.DeleteAt
	err = ss.Channel().Delete(o10.Id, model.GetMillis())
	require.Nil(t, err, "channel should have been deleted")

	t.Run("three public channels matching 'ChannelA', but already a member of one and one deleted", func(t *testing.T) {
		channels, err := ss.Channel().SearchMore(m1.UserId, teamId, "ChannelA")
		require.Nil(t, err)
		require.Equal(t, &model.ChannelList{&o3}, channels)
	})

	t.Run("one public channels, but already a member", func(t *testing.T) {
		channels, err := ss.Channel().SearchMore(m1.UserId, teamId, o4.Name)
		require.Nil(t, err)
		require.Equal(t, &model.ChannelList{}, channels)
	})

	t.Run("three matching channels, but only two public", func(t *testing.T) {
		channels, err := ss.Channel().SearchMore(m1.UserId, teamId, "off-")
		require.Nil(t, err)
		require.Equal(t, &model.ChannelList{&o7, &o6}, channels)
	})

	t.Run("one channel matching 'off-topic'", func(t *testing.T) {
		channels, err := ss.Channel().SearchMore(m1.UserId, teamId, "off-topic")
		require.Nil(t, err)
		require.Equal(t, &model.ChannelList{&o6}, channels)
	})

	t.Run("search purpose", func(t *testing.T) {
		channels, err := ss.Channel().SearchMore(m1.UserId, teamId, "now searchable")
		require.Nil(t, err)
		require.Equal(t, &model.ChannelList{&o9}, channels)
	})
}

type ByChannelDisplayName model.ChannelList

func (s ByChannelDisplayName) Len() int { return len(s) }
func (s ByChannelDisplayName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByChannelDisplayName) Less(i, j int) bool {
	if s[i].DisplayName != s[j].DisplayName {
		return s[i].DisplayName < s[j].DisplayName
	}

	return s[i].Id < s[j].Id
}

func testChannelStoreSearchInTeam(t *testing.T, ss store.Store) {
	teamId := model.NewId()
	otherTeamId := model.NewId()

	o1 := model.Channel{
		TeamId:      teamId,
		DisplayName: "ChannelA",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_OPEN,
	}
	_, err := ss.Channel().Save(&o1, -1)
	require.Nil(t, err)

	o2 := model.Channel{
		TeamId:      otherTeamId,
		DisplayName: "ChannelA",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_OPEN,
	}
	_, err = ss.Channel().Save(&o2, -1)
	require.Nil(t, err)

	m1 := model.ChannelMember{
		ChannelId:   o1.Id,
		UserId:      model.NewId(),
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}
	store.Must(ss.Channel().SaveMember(&m1))

	m2 := model.ChannelMember{
		ChannelId:   o1.Id,
		UserId:      model.NewId(),
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}
	store.Must(ss.Channel().SaveMember(&m2))

	m3 := model.ChannelMember{
		ChannelId:   o2.Id,
		UserId:      model.NewId(),
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}
	store.Must(ss.Channel().SaveMember(&m3))

	o3 := model.Channel{
		TeamId:      teamId,
		DisplayName: "ChannelA (alternate)",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_OPEN,
	}
	_, err = ss.Channel().Save(&o3, -1)
	require.Nil(t, err)

	o4 := model.Channel{
		TeamId:      teamId,
		DisplayName: "Channel B",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_PRIVATE,
	}
	_, err = ss.Channel().Save(&o4, -1)
	require.Nil(t, err)

	o5 := model.Channel{
		TeamId:      teamId,
		DisplayName: "Channel C",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_PRIVATE,
	}
	_, err = ss.Channel().Save(&o5, -1)
	require.Nil(t, err)

	o6 := model.Channel{
		TeamId:      teamId,
		DisplayName: "Off-Topic",
		Name:        "off-topic",
		Type:        model.CHANNEL_OPEN,
	}
	_, err = ss.Channel().Save(&o6, -1)
	require.Nil(t, err)

	o7 := model.Channel{
		TeamId:      teamId,
		DisplayName: "Off-Set",
		Name:        "off-set",
		Type:        model.CHANNEL_OPEN,
	}
	_, err = ss.Channel().Save(&o7, -1)
	require.Nil(t, err)

	o8 := model.Channel{
		TeamId:      teamId,
		DisplayName: "Off-Limit",
		Name:        "off-limit",
		Type:        model.CHANNEL_PRIVATE,
	}
	_, err = ss.Channel().Save(&o8, -1)
	require.Nil(t, err)

	o9 := model.Channel{
		TeamId:      teamId,
		DisplayName: "Town Square",
		Name:        "town-square",
		Type:        model.CHANNEL_OPEN,
	}
	_, err = ss.Channel().Save(&o9, -1)
	require.Nil(t, err)

	o10 := model.Channel{
		TeamId:      teamId,
		DisplayName: "The",
		Name:        "the",
		Type:        model.CHANNEL_OPEN,
	}
	_, err = ss.Channel().Save(&o10, -1)
	require.Nil(t, err)

	o11 := model.Channel{
		TeamId:      teamId,
		DisplayName: "Native Mobile Apps",
		Name:        "native-mobile-apps",
		Type:        model.CHANNEL_OPEN,
	}
	_, err = ss.Channel().Save(&o11, -1)
	require.Nil(t, err)

	o12 := model.Channel{
		TeamId:      teamId,
		DisplayName: "ChannelZ",
		Purpose:     "This can now be searchable!",
		Name:        "with-purpose",
		Type:        model.CHANNEL_OPEN,
	}
	_, err = ss.Channel().Save(&o12, -1)
	require.Nil(t, err)

	o13 := model.Channel{
		TeamId:      teamId,
		DisplayName: "ChannelA (deleted)",
		Name:        model.NewId(),
		Type:        model.CHANNEL_OPEN,
	}
	_, err = ss.Channel().Save(&o13, -1)
	require.Nil(t, err)
	o13.DeleteAt = model.GetMillis()
	o13.UpdateAt = o13.DeleteAt
	err = ss.Channel().Delete(o13.Id, model.GetMillis())
	require.Nil(t, err, "channel should have been deleted")

	testCases := []struct {
		Description     string
		TeamId          string
		Term            string
		IncludeDeleted  bool
		ExpectedResults *model.ChannelList
	}{
		{"ChannelA", teamId, "ChannelA", false, &model.ChannelList{&o1, &o3}},
		{"ChannelA, include deleted", teamId, "ChannelA", true, &model.ChannelList{&o1, &o3, &o13}},
		{"ChannelA, other team", otherTeamId, "ChannelA", false, &model.ChannelList{&o2}},
		{"empty string", teamId, "", false, &model.ChannelList{&o1, &o3, &o12, &o11, &o7, &o6, &o10, &o9}},
		{"no matches", teamId, "blargh", false, &model.ChannelList{}},
		{"prefix", teamId, "off-", false, &model.ChannelList{&o7, &o6}},
		{"full match with dash", teamId, "off-topic", false, &model.ChannelList{&o6}},
		{"town square", teamId, "town square", false, &model.ChannelList{&o9}},
		{"the in name", teamId, "the", false, &model.ChannelList{&o10}},
		{"Mobile", teamId, "Mobile", false, &model.ChannelList{&o11}},
		{"search purpose", teamId, "now searchable", false, &model.ChannelList{&o12}},
		{"pipe ignored", teamId, "town square |", false, &model.ChannelList{&o9}},
	}

	for name, search := range map[string]func(teamId string, term string, includeDeleted bool) (*model.ChannelList, *model.AppError){
		"AutocompleteInTeam": ss.Channel().AutocompleteInTeam,
		"SearchInTeam":       ss.Channel().SearchInTeam,
	} {
		for _, testCase := range testCases {
			t.Run(testCase.Description, func(t *testing.T) {
				channels, err := search(testCase.TeamId, testCase.Term, testCase.IncludeDeleted)
				require.Nil(t, err)

				// AutoCompleteInTeam doesn't currently sort its output results.
				if name == "AutocompleteInTeam" {
					sort.Sort(ByChannelDisplayName(*channels))
				}

				require.Equal(t, testCase.ExpectedResults, channels)
			})
		}
	}
}

func testChannelStoreSearchAllChannels(t *testing.T, ss store.Store) {
	cleanupChannels(t, ss)

	t1 := model.Team{}
	t1.DisplayName = "Name"
	t1.Name = model.NewId()
	t1.Email = MakeEmail()
	t1.Type = model.TEAM_OPEN
	_, err := ss.Team().Save(&t1)
	require.Nil(t, err)

	t2 := model.Team{}
	t2.DisplayName = "Name2"
	t2.Name = model.NewId()
	t2.Email = MakeEmail()
	t2.Type = model.TEAM_OPEN
	_, err = ss.Team().Save(&t2)
	require.Nil(t, err)

	o1 := model.Channel{
		TeamId:      t1.Id,
		DisplayName: "ChannelA",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_OPEN,
	}
	_, err = ss.Channel().Save(&o1, -1)
	require.Nil(t, err)

	o2 := model.Channel{
		TeamId:      t2.Id,
		DisplayName: "ChannelA",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_OPEN,
	}
	_, err = ss.Channel().Save(&o2, -1)
	require.Nil(t, err)

	m1 := model.ChannelMember{
		ChannelId:   o1.Id,
		UserId:      model.NewId(),
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}
	store.Must(ss.Channel().SaveMember(&m1))

	m2 := model.ChannelMember{
		ChannelId:   o1.Id,
		UserId:      model.NewId(),
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}
	store.Must(ss.Channel().SaveMember(&m2))

	m3 := model.ChannelMember{
		ChannelId:   o2.Id,
		UserId:      model.NewId(),
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}
	store.Must(ss.Channel().SaveMember(&m3))

	o3 := model.Channel{
		TeamId:      t1.Id,
		DisplayName: "ChannelA (alternate)",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_OPEN,
	}
	_, err = ss.Channel().Save(&o3, -1)
	require.Nil(t, err)

	o4 := model.Channel{
		TeamId:      t1.Id,
		DisplayName: "ChannelB",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_PRIVATE,
	}
	_, err = ss.Channel().Save(&o4, -1)
	require.Nil(t, err)

	o5 := model.Channel{
		TeamId:      t1.Id,
		DisplayName: "ChannelC",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_PRIVATE,
	}
	_, err = ss.Channel().Save(&o5, -1)
	require.Nil(t, err)

	o6 := model.Channel{
		TeamId:      t1.Id,
		DisplayName: "Off-Topic",
		Name:        "off-topic",
		Type:        model.CHANNEL_OPEN,
	}
	_, err = ss.Channel().Save(&o6, -1)
	require.Nil(t, err)

	o7 := model.Channel{
		TeamId:      t1.Id,
		DisplayName: "Off-Set",
		Name:        "off-set",
		Type:        model.CHANNEL_OPEN,
	}
	_, err = ss.Channel().Save(&o7, -1)
	require.Nil(t, err)

	group := &model.Group{
		Name:        model.NewId(),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewId(),
	}
	store.Must(ss.Group().Create(group))
	_, err = ss.Group().CreateGroupSyncable(model.NewGroupChannel(group.Id, o7.Id, true))
	require.Nil(t, err)

	o8 := model.Channel{
		TeamId:      t1.Id,
		DisplayName: "Off-Limit",
		Name:        "off-limit",
		Type:        model.CHANNEL_PRIVATE,
	}
	_, err = ss.Channel().Save(&o8, -1)
	require.Nil(t, err)

	o9 := model.Channel{
		TeamId:      t1.Id,
		DisplayName: "Town Square",
		Name:        "town-square",
		Type:        model.CHANNEL_OPEN,
	}
	_, err = ss.Channel().Save(&o9, -1)
	require.Nil(t, err)

	o10 := model.Channel{
		TeamId:      t1.Id,
		DisplayName: "The",
		Name:        "the",
		Type:        model.CHANNEL_OPEN,
	}
	_, err = ss.Channel().Save(&o10, -1)
	require.Nil(t, err)

	o11 := model.Channel{
		TeamId:      t1.Id,
		DisplayName: "Native Mobile Apps",
		Name:        "native-mobile-apps",
		Type:        model.CHANNEL_OPEN,
	}
	_, err = ss.Channel().Save(&o11, -1)
	require.Nil(t, err)

	o12 := model.Channel{
		TeamId:      t1.Id,
		DisplayName: "ChannelZ",
		Purpose:     "This can now be searchable!",
		Name:        "with-purpose",
		Type:        model.CHANNEL_OPEN,
	}
	_, err = ss.Channel().Save(&o12, -1)
	require.Nil(t, err)

	o13 := model.Channel{
		TeamId:      t1.Id,
		DisplayName: "ChannelA (deleted)",
		Name:        model.NewId(),
		Type:        model.CHANNEL_OPEN,
	}
	_, err = ss.Channel().Save(&o13, -1)
	require.Nil(t, err)

	o13.DeleteAt = model.GetMillis()
	o13.UpdateAt = o13.DeleteAt
	err = ss.Channel().Delete(o13.Id, o13.DeleteAt)
	require.Nil(t, err, "channel should have been deleted")

	testCases := []struct {
		Description     string
		Term            string
		Opts            store.ChannelSearchOpts
		ExpectedResults *model.ChannelList
	}{
		{"ChannelA", "ChannelA", store.ChannelSearchOpts{IncludeDeleted: false}, &model.ChannelList{&o1, &o2, &o3}},
		{"ChannelA, include deleted", "ChannelA", store.ChannelSearchOpts{IncludeDeleted: true}, &model.ChannelList{&o1, &o2, &o3, &o13}},
		{"empty string", "", store.ChannelSearchOpts{IncludeDeleted: false}, &model.ChannelList{&o1, &o2, &o3, &o4, &o5, &o12, &o11, &o8, &o7, &o6, &o10, &o9}},
		{"no matches", "blargh", store.ChannelSearchOpts{IncludeDeleted: false}, &model.ChannelList{}},
		{"prefix", "off-", store.ChannelSearchOpts{IncludeDeleted: false}, &model.ChannelList{&o8, &o7, &o6}},
		{"full match with dash", "off-topic", store.ChannelSearchOpts{IncludeDeleted: false}, &model.ChannelList{&o6}},
		{"town square", "town square", store.ChannelSearchOpts{IncludeDeleted: false}, &model.ChannelList{&o9}},
		{"the in name", "the", store.ChannelSearchOpts{IncludeDeleted: false}, &model.ChannelList{&o10}},
		{"Mobile", "Mobile", store.ChannelSearchOpts{IncludeDeleted: false}, &model.ChannelList{&o11}},
		{"search purpose", "now searchable", store.ChannelSearchOpts{IncludeDeleted: false}, &model.ChannelList{&o12}},
		{"pipe ignored", "town square |", store.ChannelSearchOpts{IncludeDeleted: false}, &model.ChannelList{&o9}},
		{"exclude defaults search 'off'", "off-", store.ChannelSearchOpts{IncludeDeleted: false, ExcludeChannelNames: []string{"off-topic"}}, &model.ChannelList{&o8, &o7}},
		{"exclude defaults search 'town'", "town", store.ChannelSearchOpts{IncludeDeleted: false, ExcludeChannelNames: []string{"town-square"}}, &model.ChannelList{}},
		{"exclude by group association", "off", store.ChannelSearchOpts{IncludeDeleted: false, NotAssociatedToGroup: group.Id}, &model.ChannelList{&o8, &o6}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Description, func(t *testing.T) {
			channels, err := ss.Channel().SearchAllChannels(testCase.Term, testCase.Opts)
			require.Nil(t, err)
			require.Equal(t, len(*testCase.ExpectedResults), len(*channels))
			for i, expected := range *testCase.ExpectedResults {
				require.Equal(t, (*channels)[i].Id, expected.Id)
			}
		})
	}
}

func testChannelStoreAutocompleteInTeamForSearch(t *testing.T, ss store.Store, s SqlSupplier) {
	u1 := &model.User{}
	u1.Email = MakeEmail()
	u1.Username = "user1" + model.NewId()
	u1.Nickname = model.NewId()
	store.Must(ss.User().Save(u1))

	u2 := &model.User{}
	u2.Email = MakeEmail()
	u2.Username = "user2" + model.NewId()
	u2.Nickname = model.NewId()
	store.Must(ss.User().Save(u2))

	u3 := &model.User{}
	u3.Email = MakeEmail()
	u3.Username = "user3" + model.NewId()
	u3.Nickname = model.NewId()
	store.Must(ss.User().Save(u3))

	u4 := &model.User{}
	u4.Email = MakeEmail()
	u4.Username = "user4" + model.NewId()
	u4.Nickname = model.NewId()
	store.Must(ss.User().Save(u4))

	o1 := model.Channel{}
	o1.TeamId = model.NewId()
	o1.DisplayName = "ChannelA"
	o1.Name = "zz" + model.NewId() + "b"
	o1.Type = model.CHANNEL_OPEN
	_, err := ss.Channel().Save(&o1, -1)
	require.Nil(t, err)

	m1 := model.ChannelMember{}
	m1.ChannelId = o1.Id
	m1.UserId = u1.Id
	m1.NotifyProps = model.GetDefaultChannelNotifyProps()
	store.Must(ss.Channel().SaveMember(&m1))

	o2 := model.Channel{}
	o2.TeamId = model.NewId()
	o2.DisplayName = "Channel2"
	o2.Name = "zz" + model.NewId() + "b"
	o2.Type = model.CHANNEL_OPEN
	_, err = ss.Channel().Save(&o2, -1)
	require.Nil(t, err)

	m2 := model.ChannelMember{}
	m2.ChannelId = o2.Id
	m2.UserId = m1.UserId
	m2.NotifyProps = model.GetDefaultChannelNotifyProps()
	store.Must(ss.Channel().SaveMember(&m2))

	o3 := model.Channel{}
	o3.TeamId = o1.TeamId
	o3.DisplayName = "ChannelA"
	o3.Name = "zz" + model.NewId() + "b"
	o3.Type = model.CHANNEL_OPEN
	_, err = ss.Channel().Save(&o3, -1)
	require.Nil(t, err)

	m3 := model.ChannelMember{}
	m3.ChannelId = o3.Id
	m3.UserId = m1.UserId
	m3.NotifyProps = model.GetDefaultChannelNotifyProps()
	store.Must(ss.Channel().SaveMember(&m3))

	err = ss.Channel().SetDeleteAt(o3.Id, 100, 100)
	require.Nil(t, err, "channel should have been deleted")

	o4 := model.Channel{}
	o4.TeamId = o1.TeamId
	o4.DisplayName = "ChannelA"
	o4.Name = "zz" + model.NewId() + "b"
	o4.Type = model.CHANNEL_PRIVATE
	_, err = ss.Channel().Save(&o4, -1)
	require.Nil(t, err)

	m4 := model.ChannelMember{}
	m4.ChannelId = o4.Id
	m4.UserId = m1.UserId
	m4.NotifyProps = model.GetDefaultChannelNotifyProps()
	store.Must(ss.Channel().SaveMember(&m4))

	o5 := model.Channel{}
	o5.TeamId = o1.TeamId
	o5.DisplayName = "ChannelC"
	o5.Name = "zz" + model.NewId() + "b"
	o5.Type = model.CHANNEL_PRIVATE
	_, err = ss.Channel().Save(&o5, -1)
	require.Nil(t, err)

	_, err = ss.Channel().CreateDirectChannel(u1.Id, u2.Id)
	require.Nil(t, err)
	_, err = ss.Channel().CreateDirectChannel(u2.Id, u3.Id)
	require.Nil(t, err)

	tt := []struct {
		name            string
		term            string
		includeDeleted  bool
		expectedMatches int
	}{
		{"Empty search (list all)", "", false, 4},
		{"Narrow search", "ChannelA", false, 2},
		{"Wide search", "Cha", false, 3},
		{"Direct messages", "user", false, 1},
		{"Wide search with archived channels", "Cha", true, 4},
		{"Narrow with archived channels", "ChannelA", true, 3},
		{"Direct messages with archived channels", "user", true, 1},
		{"Search without results", "blarg", true, 0},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			channels, err := ss.Channel().AutocompleteInTeamForSearch(o1.TeamId, m1.UserId, "ChannelA", false)
			require.Nil(t, err)
			require.Len(t, *channels, 2)
		})
	}

	// Manually truncate Channels table until testlib can handle cleanups
	s.GetMaster().Exec("TRUNCATE Channels")
}

func testChannelStoreGetMembersByIds(t *testing.T, ss store.Store) {
	o1 := model.Channel{}
	o1.TeamId = model.NewId()
	o1.DisplayName = "ChannelA"
	o1.Name = "zz" + model.NewId() + "b"
	o1.Type = model.CHANNEL_OPEN
	_, err := ss.Channel().Save(&o1, -1)
	require.Nil(t, err)

	m1 := &model.ChannelMember{ChannelId: o1.Id, UserId: model.NewId(), NotifyProps: model.GetDefaultChannelNotifyProps()}
	store.Must(ss.Channel().SaveMember(m1))

	if members, err := ss.Channel().GetMembersByIds(m1.ChannelId, []string{m1.UserId}); err != nil {
		t.Fatal(err)
	} else {
		rm1 := (*members)[0]

		if rm1.ChannelId != m1.ChannelId {
			t.Fatal("bad team id")
		}

		if rm1.UserId != m1.UserId {
			t.Fatal("bad user id")
		}
	}

	m2 := &model.ChannelMember{ChannelId: o1.Id, UserId: model.NewId(), NotifyProps: model.GetDefaultChannelNotifyProps()}
	store.Must(ss.Channel().SaveMember(m2))

	if members, err := ss.Channel().GetMembersByIds(m1.ChannelId, []string{m1.UserId, m2.UserId, model.NewId()}); err != nil {
		t.Fatal(err)
	} else {
		if len(*members) != 2 {
			t.Fatal("return wrong number of results")
		}
	}

	if _, err := ss.Channel().GetMembersByIds(m1.ChannelId, []string{}); err == nil {
		t.Fatal("empty user ids - should have failed")
	}
}

func testChannelStoreSearchGroupChannels(t *testing.T, ss store.Store) {
	// Users
	u1 := &model.User{}
	u1.Username = "user.one"
	u1.Email = MakeEmail()
	u1.Nickname = model.NewId()
	store.Must(ss.User().Save(u1))

	u2 := &model.User{}
	u2.Username = "user.two"
	u2.Email = MakeEmail()
	u2.Nickname = model.NewId()
	store.Must(ss.User().Save(u2))

	u3 := &model.User{}
	u3.Username = "user.three"
	u3.Email = MakeEmail()
	u3.Nickname = model.NewId()
	store.Must(ss.User().Save(u3))

	u4 := &model.User{}
	u4.Username = "user.four"
	u4.Email = MakeEmail()
	u4.Nickname = model.NewId()
	store.Must(ss.User().Save(u4))

	// Group channels
	userIds := []string{u1.Id, u2.Id, u3.Id}
	gc1 := model.Channel{}
	gc1.Name = model.GetGroupNameFromUserIds(userIds)
	gc1.DisplayName = "GroupChannel" + model.NewId()
	gc1.Type = model.CHANNEL_GROUP
	_, err := ss.Channel().Save(&gc1, -1)
	require.Nil(t, err)

	for _, userId := range userIds {
		store.Must(ss.Channel().SaveMember(&model.ChannelMember{
			ChannelId:   gc1.Id,
			UserId:      userId,
			NotifyProps: model.GetDefaultChannelNotifyProps(),
		}))
	}

	userIds = []string{u1.Id, u4.Id}
	gc2 := model.Channel{}
	gc2.Name = model.GetGroupNameFromUserIds(userIds)
	gc2.DisplayName = "GroupChannel" + model.NewId()
	gc2.Type = model.CHANNEL_GROUP
	_, err = ss.Channel().Save(&gc2, -1)
	require.Nil(t, err)

	for _, userId := range userIds {
		store.Must(ss.Channel().SaveMember(&model.ChannelMember{
			ChannelId:   gc2.Id,
			UserId:      userId,
			NotifyProps: model.GetDefaultChannelNotifyProps(),
		}))
	}

	userIds = []string{u1.Id, u2.Id, u3.Id, u4.Id}
	gc3 := model.Channel{}
	gc3.Name = model.GetGroupNameFromUserIds(userIds)
	gc3.DisplayName = "GroupChannel" + model.NewId()
	gc3.Type = model.CHANNEL_GROUP
	_, err = ss.Channel().Save(&gc3, -1)
	require.Nil(t, err)

	for _, userId := range userIds {
		store.Must(ss.Channel().SaveMember(&model.ChannelMember{
			ChannelId:   gc3.Id,
			UserId:      userId,
			NotifyProps: model.GetDefaultChannelNotifyProps(),
		}))
	}

	defer func() {
		for _, gc := range []model.Channel{gc1, gc2, gc3} {
			ss.Channel().PermanentDeleteMembersByChannel(gc3.Id)
			<-ss.Channel().PermanentDelete(gc.Id)
		}
	}()

	testCases := []struct {
		Name           string
		UserId         string
		Term           string
		ExpectedResult []string
	}{
		{
			Name:           "Get all group channels for user1",
			UserId:         u1.Id,
			Term:           "",
			ExpectedResult: []string{gc1.Id, gc2.Id, gc3.Id},
		},
		{
			Name:           "Get group channels for user1 and term 'three'",
			UserId:         u1.Id,
			Term:           "three",
			ExpectedResult: []string{gc1.Id, gc3.Id},
		},
		{
			Name:           "Get group channels for user1 and term 'four two'",
			UserId:         u1.Id,
			Term:           "four two",
			ExpectedResult: []string{gc3.Id},
		},
		{
			Name:           "Get all group channels for user2",
			UserId:         u2.Id,
			Term:           "",
			ExpectedResult: []string{gc1.Id, gc3.Id},
		},
		{
			Name:           "Get group channels for user2 and term 'four'",
			UserId:         u2.Id,
			Term:           "four",
			ExpectedResult: []string{gc3.Id},
		},
		{
			Name:           "Get all group channels for user4",
			UserId:         u4.Id,
			Term:           "",
			ExpectedResult: []string{gc2.Id, gc3.Id},
		},
		{
			Name:           "Get group channels for user4 and term 'one five'",
			UserId:         u4.Id,
			Term:           "one five",
			ExpectedResult: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := ss.Channel().SearchGroupChannels(tc.UserId, tc.Term)
			require.Nil(t, err)

			resultIds := []string{}
			for _, gc := range *result {
				resultIds = append(resultIds, gc.Id)
			}

			require.ElementsMatch(t, tc.ExpectedResult, resultIds)
		})
	}
}

func testChannelStoreAnalyticsDeletedTypeCount(t *testing.T, ss store.Store) {
	o1 := model.Channel{}
	o1.TeamId = model.NewId()
	o1.DisplayName = "ChannelA"
	o1.Name = "zz" + model.NewId() + "b"
	o1.Type = model.CHANNEL_OPEN
	_, err := ss.Channel().Save(&o1, -1)
	require.Nil(t, err)

	o2 := model.Channel{}
	o2.TeamId = model.NewId()
	o2.DisplayName = "Channel2"
	o2.Name = "zz" + model.NewId() + "b"
	o2.Type = model.CHANNEL_OPEN
	_, err = ss.Channel().Save(&o2, -1)
	require.Nil(t, err)

	p3 := model.Channel{}
	p3.TeamId = model.NewId()
	p3.DisplayName = "Channel3"
	p3.Name = "zz" + model.NewId() + "b"
	p3.Type = model.CHANNEL_PRIVATE
	_, err = ss.Channel().Save(&p3, -1)
	require.Nil(t, err)

	u1 := &model.User{}
	u1.Email = MakeEmail()
	u1.Nickname = model.NewId()
	store.Must(ss.User().Save(u1))

	u2 := &model.User{}
	u2.Email = MakeEmail()
	u2.Nickname = model.NewId()
	store.Must(ss.User().Save(u2))

	d4, err := ss.Channel().CreateDirectChannel(u1.Id, u2.Id)
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer func() {
		ss.Channel().PermanentDeleteMembersByChannel(d4.Id)
		<-ss.Channel().PermanentDelete(d4.Id)
	}()

	var openStartCount int64
	if openStartCount, err = ss.Channel().AnalyticsDeletedTypeCount("", "O"); err != nil {
		t.Fatal(err)
	}

	var privateStartCount int64
	if privateStartCount, err = ss.Channel().AnalyticsDeletedTypeCount("", "P"); err != nil {
		t.Fatal(err)
	}

	var directStartCount int64
	if directStartCount, err = ss.Channel().AnalyticsDeletedTypeCount("", "D"); err != nil {
		t.Fatal(err)
	}

	err = ss.Channel().Delete(o1.Id, model.GetMillis())
	require.Nil(t, err, "channel should have been deleted")
	err = ss.Channel().Delete(o2.Id, model.GetMillis())
	require.Nil(t, err, "channel should have been deleted")
	err = ss.Channel().Delete(p3.Id, model.GetMillis())
	require.Nil(t, err, "channel should have been deleted")
	err = ss.Channel().Delete(d4.Id, model.GetMillis())
	require.Nil(t, err, "channel should have been deleted")

	var count int64

	if count, err = ss.Channel().AnalyticsDeletedTypeCount("", "O"); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, openStartCount+2, count, "Wrong open channel deleted count.")

	if count, err = ss.Channel().AnalyticsDeletedTypeCount("", "P"); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, privateStartCount+1, count, "Wrong private channel deleted count.")

	if count, err = ss.Channel().AnalyticsDeletedTypeCount("", "D"); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, directStartCount+1, count, "Wrong direct channel deleted count.")
}

func testChannelStoreGetPinnedPosts(t *testing.T, ss store.Store) {
	ch1 := &model.Channel{
		TeamId:      model.NewId(),
		DisplayName: "Name",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_OPEN,
	}

	o1, err := ss.Channel().Save(ch1, -1)
	require.Nil(t, err)

	p1, err := ss.Post().Save(&model.Post{
		UserId:    model.NewId(),
		ChannelId: o1.Id,
		Message:   "test",
		IsPinned:  true,
	})

	if pl, errGet := ss.Channel().GetPinnedPosts(o1.Id); errGet != nil {
		t.Fatal(errGet)
	} else if pl.Posts[p1.Id] == nil {
		t.Fatal("didn't return relevant pinned posts")
	}

	ch2 := &model.Channel{
		TeamId:      model.NewId(),
		DisplayName: "Name",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_OPEN,
	}

	o2, err := ss.Channel().Save(ch2, -1)
	require.Nil(t, err)

	_, err = ss.Post().Save(&model.Post{
		UserId:    model.NewId(),
		ChannelId: o2.Id,
		Message:   "test",
	})
	require.Nil(t, err)

	if pl, errGet := ss.Channel().GetPinnedPosts(o2.Id); errGet != nil {
		t.Fatal(errGet)
	} else if len(pl.Posts) != 0 {
		t.Fatal("wasn't supposed to return posts")
	}
}

func testChannelStoreMaxChannelsPerTeam(t *testing.T, ss store.Store) {
	channel := &model.Channel{
		TeamId:      model.NewId(),
		DisplayName: "Channel",
		Name:        model.NewId(),
		Type:        model.CHANNEL_OPEN,
	}
	_, err := ss.Channel().Save(channel, 0)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, err.Id, "store.sql_channel.save_channel.limit.app_error")

	channel.Id = ""
	_, err = ss.Channel().Save(channel, 1)
	assert.Nil(t, err)
}

func testChannelStoreGetChannelsByScheme(t *testing.T, ss store.Store) {
	// Create some schemes.
	s1 := &model.Scheme{
		DisplayName: model.NewId(),
		Name:        model.NewId(),
		Description: model.NewId(),
		Scope:       model.SCHEME_SCOPE_CHANNEL,
	}

	s2 := &model.Scheme{
		DisplayName: model.NewId(),
		Name:        model.NewId(),
		Description: model.NewId(),
		Scope:       model.SCHEME_SCOPE_CHANNEL,
	}

	result := <-ss.Scheme().Save(s1)
	s1 = result.Data.(*model.Scheme)
	s1 = (<-ss.Scheme().Save(s1)).Data.(*model.Scheme)
	s2 = (<-ss.Scheme().Save(s2)).Data.(*model.Scheme)

	// Create and save some teams.
	c1 := &model.Channel{
		TeamId:      model.NewId(),
		DisplayName: "Name",
		Name:        model.NewId(),
		Type:        model.CHANNEL_OPEN,
		SchemeId:    &s1.Id,
	}

	c2 := &model.Channel{
		TeamId:      model.NewId(),
		DisplayName: "Name",
		Name:        model.NewId(),
		Type:        model.CHANNEL_OPEN,
		SchemeId:    &s1.Id,
	}

	c3 := &model.Channel{
		TeamId:      model.NewId(),
		DisplayName: "Name",
		Name:        model.NewId(),
		Type:        model.CHANNEL_OPEN,
	}

	_, _ = ss.Channel().Save(c1, 100)
	_, _ = ss.Channel().Save(c2, 100)
	_, _ = ss.Channel().Save(c3, 100)

	// Get the channels by a valid Scheme ID.
	res1 := <-ss.Channel().GetChannelsByScheme(s1.Id, 0, 100)
	assert.Nil(t, res1.Err)
	d1 := res1.Data.(model.ChannelList)
	assert.Len(t, d1, 2)

	// Get the channels by a valid Scheme ID where there aren't any matching Channel.
	res2 := <-ss.Channel().GetChannelsByScheme(s2.Id, 0, 100)
	assert.Nil(t, res2.Err)
	d2 := res2.Data.(model.ChannelList)
	assert.Len(t, d2, 0)

	// Get the channels by an invalid Scheme ID.
	res3 := <-ss.Channel().GetChannelsByScheme(model.NewId(), 0, 100)
	assert.Nil(t, res3.Err)
	d3 := res3.Data.(model.ChannelList)
	assert.Len(t, d3, 0)
}

func testChannelStoreMigrateChannelMembers(t *testing.T, ss store.Store) {
	s1 := model.NewId()
	c1 := &model.Channel{
		TeamId:      model.NewId(),
		DisplayName: "Name",
		Name:        model.NewId(),
		Type:        model.CHANNEL_OPEN,
		SchemeId:    &s1,
	}
	c1, _ = ss.Channel().Save(c1, 100)

	cm1 := &model.ChannelMember{
		ChannelId:     c1.Id,
		UserId:        model.NewId(),
		ExplicitRoles: "channel_admin channel_user",
		NotifyProps:   model.GetDefaultChannelNotifyProps(),
	}
	cm2 := &model.ChannelMember{
		ChannelId:     c1.Id,
		UserId:        model.NewId(),
		ExplicitRoles: "channel_user",
		NotifyProps:   model.GetDefaultChannelNotifyProps(),
	}
	cm3 := &model.ChannelMember{
		ChannelId:     c1.Id,
		UserId:        model.NewId(),
		ExplicitRoles: "something_else",
		NotifyProps:   model.GetDefaultChannelNotifyProps(),
	}

	cm1 = (<-ss.Channel().SaveMember(cm1)).Data.(*model.ChannelMember)
	cm2 = (<-ss.Channel().SaveMember(cm2)).Data.(*model.ChannelMember)
	cm3 = (<-ss.Channel().SaveMember(cm3)).Data.(*model.ChannelMember)

	lastDoneChannelId := strings.Repeat("0", 26)
	lastDoneUserId := strings.Repeat("0", 26)

	for {
		data, err := ss.Channel().MigrateChannelMembers(lastDoneChannelId, lastDoneUserId)
		if assert.Nil(t, err) {
			if data == nil {
				break
			}
			lastDoneChannelId = data["ChannelId"]
			lastDoneUserId = data["UserId"]
		}
	}

	ss.Channel().ClearCaches()

	cm1b, err := ss.Channel().GetMember(cm1.ChannelId, cm1.UserId)
	assert.Nil(t, err)
	assert.Equal(t, "", cm1b.ExplicitRoles)
	assert.False(t, cm1b.SchemeGuest)
	assert.True(t, cm1b.SchemeUser)
	assert.True(t, cm1b.SchemeAdmin)

	cm2b, err := ss.Channel().GetMember(cm2.ChannelId, cm2.UserId)
	assert.Nil(t, err)
	assert.Equal(t, "", cm2b.ExplicitRoles)
	assert.False(t, cm1b.SchemeGuest)
	assert.True(t, cm2b.SchemeUser)
	assert.False(t, cm2b.SchemeAdmin)

	cm3b, err := ss.Channel().GetMember(cm3.ChannelId, cm3.UserId)
	assert.Nil(t, err)
	assert.Equal(t, "something_else", cm3b.ExplicitRoles)
	assert.False(t, cm1b.SchemeGuest)
	assert.False(t, cm3b.SchemeUser)
	assert.False(t, cm3b.SchemeAdmin)
}

func testResetAllChannelSchemes(t *testing.T, ss store.Store) {
	s1 := &model.Scheme{
		Name:        model.NewId(),
		DisplayName: model.NewId(),
		Description: model.NewId(),
		Scope:       model.SCHEME_SCOPE_CHANNEL,
	}
	s1 = (<-ss.Scheme().Save(s1)).Data.(*model.Scheme)

	c1 := &model.Channel{
		TeamId:      model.NewId(),
		DisplayName: "Name",
		Name:        model.NewId(),
		Type:        model.CHANNEL_OPEN,
		SchemeId:    &s1.Id,
	}

	c2 := &model.Channel{
		TeamId:      model.NewId(),
		DisplayName: "Name",
		Name:        model.NewId(),
		Type:        model.CHANNEL_OPEN,
		SchemeId:    &s1.Id,
	}

	c1, _ = ss.Channel().Save(c1, 100)
	c2, _ = ss.Channel().Save(c2, 100)

	assert.Equal(t, s1.Id, *c1.SchemeId)
	assert.Equal(t, s1.Id, *c2.SchemeId)

	err := ss.Channel().ResetAllChannelSchemes()
	assert.Nil(t, err)

	c1, _ = ss.Channel().Get(c1.Id, true)
	c2, _ = ss.Channel().Get(c2.Id, true)

	assert.Equal(t, "", *c1.SchemeId)
	assert.Equal(t, "", *c2.SchemeId)
}

func testChannelStoreClearAllCustomRoleAssignments(t *testing.T, ss store.Store) {
	c := &model.Channel{
		TeamId:      model.NewId(),
		DisplayName: "Name",
		Name:        model.NewId(),
		Type:        model.CHANNEL_OPEN,
	}

	c, _ = ss.Channel().Save(c, 100)

	m1 := &model.ChannelMember{
		ChannelId:     c.Id,
		UserId:        model.NewId(),
		NotifyProps:   model.GetDefaultChannelNotifyProps(),
		ExplicitRoles: "channel_user channel_admin system_user_access_token",
	}
	m2 := &model.ChannelMember{
		ChannelId:     c.Id,
		UserId:        model.NewId(),
		NotifyProps:   model.GetDefaultChannelNotifyProps(),
		ExplicitRoles: "channel_user custom_role channel_admin another_custom_role",
	}
	m3 := &model.ChannelMember{
		ChannelId:     c.Id,
		UserId:        model.NewId(),
		NotifyProps:   model.GetDefaultChannelNotifyProps(),
		ExplicitRoles: "channel_user",
	}
	m4 := &model.ChannelMember{
		ChannelId:     c.Id,
		UserId:        model.NewId(),
		NotifyProps:   model.GetDefaultChannelNotifyProps(),
		ExplicitRoles: "custom_only",
	}

	store.Must(ss.Channel().SaveMember(m1))
	store.Must(ss.Channel().SaveMember(m2))
	store.Must(ss.Channel().SaveMember(m3))
	store.Must(ss.Channel().SaveMember(m4))

	require.Nil(t, ss.Channel().ClearAllCustomRoleAssignments())

	member, err := ss.Channel().GetMember(m1.ChannelId, m1.UserId)
	require.Nil(t, err)
	assert.Equal(t, m1.ExplicitRoles, member.Roles)

	member, err = ss.Channel().GetMember(m2.ChannelId, m2.UserId)
	require.Nil(t, err)
	assert.Equal(t, "channel_user channel_admin", member.Roles)

	member, err = ss.Channel().GetMember(m3.ChannelId, m3.UserId)
	require.Nil(t, err)
	assert.Equal(t, m3.ExplicitRoles, member.Roles)

	member, err = ss.Channel().GetMember(m4.ChannelId, m4.UserId)
	require.Nil(t, err)
	assert.Equal(t, "", member.Roles)
}

// testMaterializedPublicChannels tests edge cases involving the triggers and stored procedures
// that materialize the PublicChannels table.
func testMaterializedPublicChannels(t *testing.T, ss store.Store, s SqlSupplier) {
	teamId := model.NewId()

	// o1 is a public channel on the team
	o1 := model.Channel{
		TeamId:      teamId,
		DisplayName: "Open Channel",
		Name:        model.NewId(),
		Type:        model.CHANNEL_OPEN,
	}
	_, err := ss.Channel().Save(&o1, -1)
	require.Nil(t, err)

	// o2 is another public channel on the team
	o2 := model.Channel{
		TeamId:      teamId,
		DisplayName: "Open Channel 2",
		Name:        model.NewId(),
		Type:        model.CHANNEL_OPEN,
	}
	_, err = ss.Channel().Save(&o2, -1)
	require.Nil(t, err)

	t.Run("o1 and o2 initially listed in public channels", func(t *testing.T) {
		channels, channelErr := ss.Channel().SearchInTeam(teamId, "", true)
		require.Nil(t, channelErr)
		require.Equal(t, &model.ChannelList{&o1, &o2}, channels)
	})

	o1.DeleteAt = model.GetMillis()
	o1.UpdateAt = model.GetMillis()

	e := ss.Channel().Delete(o1.Id, o1.DeleteAt)
	require.Nil(t, e, "channel should have been deleted")

	t.Run("o1 still listed in public channels when marked as deleted", func(t *testing.T) {
		channels, channelErr := ss.Channel().SearchInTeam(teamId, "", true)
		require.Nil(t, channelErr)
		require.Equal(t, &model.ChannelList{&o1, &o2}, channels)
	})

	<-ss.Channel().PermanentDelete(o1.Id)

	t.Run("o1 no longer listed in public channels when permanently deleted", func(t *testing.T) {
		channels, channelErr := ss.Channel().SearchInTeam(teamId, "", true)
		require.Nil(t, channelErr)
		require.Equal(t, &model.ChannelList{&o2}, channels)
	})

	o2.Type = model.CHANNEL_PRIVATE
	_, appErr := ss.Channel().Update(&o2)
	require.Nil(t, appErr)

	t.Run("o2 no longer listed since now private", func(t *testing.T) {
		channels, channelErr := ss.Channel().SearchInTeam(teamId, "", true)
		require.Nil(t, channelErr)
		require.Equal(t, &model.ChannelList{}, channels)
	})

	o2.Type = model.CHANNEL_OPEN
	_, appErr = ss.Channel().Update(&o2)
	require.Nil(t, appErr)

	t.Run("o2 listed once again since now public", func(t *testing.T) {
		channels, channelErr := ss.Channel().SearchInTeam(teamId, "", true)
		require.Nil(t, channelErr)
		require.Equal(t, &model.ChannelList{&o2}, channels)
	})

	// o3 is a public channel on the team that already existed in the PublicChannels table.
	o3 := model.Channel{
		Id:          model.NewId(),
		TeamId:      teamId,
		DisplayName: "Open Channel 3",
		Name:        model.NewId(),
		Type:        model.CHANNEL_OPEN,
	}

	_, execerr := s.GetMaster().ExecNoTimeout(`
		INSERT INTO
		    PublicChannels(Id, DeleteAt, TeamId, DisplayName, Name, Header, Purpose)
		VALUES
		    (:Id, :DeleteAt, :TeamId, :DisplayName, :Name, :Header, :Purpose);
	`, map[string]interface{}{
		"Id":          o3.Id,
		"DeleteAt":    o3.DeleteAt,
		"TeamId":      o3.TeamId,
		"DisplayName": o3.DisplayName,
		"Name":        o3.Name,
		"Header":      o3.Header,
		"Purpose":     o3.Purpose,
	})
	require.Nil(t, execerr)

	o3.DisplayName = "Open Channel 3 - Modified"

	_, execerr = s.GetMaster().ExecNoTimeout(`
		INSERT INTO
		    Channels(Id, CreateAt, UpdateAt, DeleteAt, TeamId, Type, DisplayName, Name, Header, Purpose, LastPostAt, TotalMsgCount, ExtraUpdateAt, CreatorId)
		VALUES
		    (:Id, :CreateAt, :UpdateAt, :DeleteAt, :TeamId, :Type, :DisplayName, :Name, :Header, :Purpose, :LastPostAt, :TotalMsgCount, :ExtraUpdateAt, :CreatorId);
	`, map[string]interface{}{
		"Id":            o3.Id,
		"CreateAt":      o3.CreateAt,
		"UpdateAt":      o3.UpdateAt,
		"DeleteAt":      o3.DeleteAt,
		"TeamId":        o3.TeamId,
		"Type":          o3.Type,
		"DisplayName":   o3.DisplayName,
		"Name":          o3.Name,
		"Header":        o3.Header,
		"Purpose":       o3.Purpose,
		"LastPostAt":    o3.LastPostAt,
		"TotalMsgCount": o3.TotalMsgCount,
		"ExtraUpdateAt": o3.ExtraUpdateAt,
		"CreatorId":     o3.CreatorId,
	})
	require.Nil(t, execerr)

	t.Run("verify o3 INSERT converted to UPDATE", func(t *testing.T) {
		channels, channelErr := ss.Channel().SearchInTeam(teamId, "", true)
		require.Nil(t, channelErr)
		require.Equal(t, &model.ChannelList{&o2, &o3}, channels)
	})

	// o4 is a public channel on the team that existed in the Channels table but was omitted from the PublicChannels table.
	o4 := model.Channel{
		TeamId:      teamId,
		DisplayName: "Open Channel 4",
		Name:        model.NewId(),
		Type:        model.CHANNEL_OPEN,
	}

	_, err = ss.Channel().Save(&o4, -1)
	require.Nil(t, err)

	_, execerr = s.GetMaster().ExecNoTimeout(`
		DELETE FROM
		    PublicChannels
		WHERE
		    Id = :Id
	`, map[string]interface{}{
		"Id": o4.Id,
	})
	require.Nil(t, execerr)

	o4.DisplayName += " - Modified"
	_, appErr = ss.Channel().Update(&o4)
	require.Nil(t, appErr)

	t.Run("verify o4 UPDATE converted to INSERT", func(t *testing.T) {
		channels, err := ss.Channel().SearchInTeam(teamId, "", true)
		require.Nil(t, err)
		require.Equal(t, &model.ChannelList{&o2, &o3, &o4}, channels)
	})
}

func testChannelStoreGetAllChannelsForExportAfter(t *testing.T, ss store.Store) {
	t1 := model.Team{}
	t1.DisplayName = "Name"
	t1.Name = model.NewId()
	t1.Email = MakeEmail()
	t1.Type = model.TEAM_OPEN
	_, err := ss.Team().Save(&t1)
	require.Nil(t, err)

	c1 := model.Channel{}
	c1.TeamId = t1.Id
	c1.DisplayName = "Channel1"
	c1.Name = "zz" + model.NewId() + "b"
	c1.Type = model.CHANNEL_OPEN
	_, err = ss.Channel().Save(&c1, -1)
	require.Nil(t, err)

	d1, err := ss.Channel().GetAllChannelsForExportAfter(10000, strings.Repeat("0", 26))
	assert.Nil(t, err)

	found := false
	for _, c := range d1 {
		if c.Id == c1.Id {
			found = true
			assert.Equal(t, t1.Id, c.TeamId)
			assert.Nil(t, c.SchemeId)
			assert.Equal(t, t1.Name, c.TeamName)
		}
	}
	assert.True(t, found)
}

func testChannelStoreGetChannelMembersForExport(t *testing.T, ss store.Store) {
	t1 := model.Team{}
	t1.DisplayName = "Name"
	t1.Name = model.NewId()
	t1.Email = MakeEmail()
	t1.Type = model.TEAM_OPEN
	_, err := ss.Team().Save(&t1)
	require.Nil(t, err)

	c1 := model.Channel{}
	c1.TeamId = t1.Id
	c1.DisplayName = "Channel1"
	c1.Name = "zz" + model.NewId() + "b"
	c1.Type = model.CHANNEL_OPEN
	_, err = ss.Channel().Save(&c1, -1)
	require.Nil(t, err)

	c2 := model.Channel{}
	c2.TeamId = model.NewId()
	c2.DisplayName = "Channel2"
	c2.Name = "zz" + model.NewId() + "b"
	c2.Type = model.CHANNEL_OPEN
	_, err = ss.Channel().Save(&c2, -1)
	require.Nil(t, err)

	u1 := model.User{}
	u1.Email = MakeEmail()
	u1.Nickname = model.NewId()
	store.Must(ss.User().Save(&u1))

	m1 := model.ChannelMember{}
	m1.ChannelId = c1.Id
	m1.UserId = u1.Id
	m1.NotifyProps = model.GetDefaultChannelNotifyProps()
	store.Must(ss.Channel().SaveMember(&m1))

	m2 := model.ChannelMember{}
	m2.ChannelId = c2.Id
	m2.UserId = u1.Id
	m2.NotifyProps = model.GetDefaultChannelNotifyProps()
	store.Must(ss.Channel().SaveMember(&m2))

	d1, err := ss.Channel().GetChannelMembersForExport(u1.Id, t1.Id)
	assert.Nil(t, err)

	assert.Len(t, d1, 1)

	cmfe1 := d1[0]
	assert.Equal(t, c1.Name, cmfe1.ChannelName)
	assert.Equal(t, c1.Id, cmfe1.ChannelId)
	assert.Equal(t, u1.Id, cmfe1.UserId)
}

func testChannelStoreRemoveAllDeactivatedMembers(t *testing.T, ss store.Store) {
	// Set up all the objects needed in the store.
	t1 := model.Team{}
	t1.DisplayName = "Name"
	t1.Name = model.NewId()
	t1.Email = MakeEmail()
	t1.Type = model.TEAM_OPEN
	_, err := ss.Team().Save(&t1)
	require.Nil(t, err)

	c1 := model.Channel{}
	c1.TeamId = t1.Id
	c1.DisplayName = "Channel1"
	c1.Name = "zz" + model.NewId() + "b"
	c1.Type = model.CHANNEL_OPEN
	_, err = ss.Channel().Save(&c1, -1)
	require.Nil(t, err)

	u1 := model.User{}
	u1.Email = MakeEmail()
	u1.Nickname = model.NewId()
	store.Must(ss.User().Save(&u1))

	u2 := model.User{}
	u2.Email = MakeEmail()
	u2.Nickname = model.NewId()
	store.Must(ss.User().Save(&u2))

	u3 := model.User{}
	u3.Email = MakeEmail()
	u3.Nickname = model.NewId()
	store.Must(ss.User().Save(&u3))

	m1 := model.ChannelMember{}
	m1.ChannelId = c1.Id
	m1.UserId = u1.Id
	m1.NotifyProps = model.GetDefaultChannelNotifyProps()
	store.Must(ss.Channel().SaveMember(&m1))

	m2 := model.ChannelMember{}
	m2.ChannelId = c1.Id
	m2.UserId = u2.Id
	m2.NotifyProps = model.GetDefaultChannelNotifyProps()
	store.Must(ss.Channel().SaveMember(&m2))

	m3 := model.ChannelMember{}
	m3.ChannelId = c1.Id
	m3.UserId = u3.Id
	m3.NotifyProps = model.GetDefaultChannelNotifyProps()
	store.Must(ss.Channel().SaveMember(&m3))

	// Get all the channel members. Check there are 3.
	d1, err := ss.Channel().GetMembers(c1.Id, 0, 1000)
	assert.Nil(t, err)
	assert.Len(t, *d1, 3)

	// Deactivate users 1 & 2.
	u1.DeleteAt = model.GetMillis()
	u2.DeleteAt = model.GetMillis()
	_, err = ss.User().Update(&u1, true)
	require.Nil(t, err)
	_, err = ss.User().Update(&u2, true)
	require.Nil(t, err)

	// Remove all deactivated users from the channel.
	assert.Nil(t, ss.Channel().RemoveAllDeactivatedMembers(c1.Id))

	// Get all the channel members. Check there is now only 1: m3.
	d2, err := ss.Channel().GetMembers(c1.Id, 0, 1000)
	assert.Nil(t, err)
	assert.Len(t, *d2, 1)
	assert.Equal(t, (*d2)[0].UserId, u3.Id)
}

func testChannelStoreExportAllDirectChannels(t *testing.T, ss store.Store, s SqlSupplier) {
	teamId := model.NewId()

	o1 := model.Channel{}
	o1.TeamId = teamId
	o1.DisplayName = "Name" + model.NewId()
	o1.Name = "zz" + model.NewId() + "b"
	o1.Type = model.CHANNEL_DIRECT

	userIds := []string{model.NewId(), model.NewId(), model.NewId()}

	o2 := model.Channel{}
	o2.Name = model.GetGroupNameFromUserIds(userIds)
	o2.DisplayName = "GroupChannel" + model.NewId()
	o2.Name = "zz" + model.NewId() + "b"
	o2.Type = model.CHANNEL_GROUP
	_, err := ss.Channel().Save(&o2, -1)
	require.Nil(t, err)

	u1 := &model.User{}
	u1.Email = MakeEmail()
	u1.Nickname = model.NewId()
	store.Must(ss.User().Save(u1))
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: model.NewId(), UserId: u1.Id}, -1))

	u2 := &model.User{}
	u2.Email = MakeEmail()
	u2.Nickname = model.NewId()
	store.Must(ss.User().Save(u2))
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: model.NewId(), UserId: u2.Id}, -1))

	m1 := model.ChannelMember{}
	m1.ChannelId = o1.Id
	m1.UserId = u1.Id
	m1.NotifyProps = model.GetDefaultChannelNotifyProps()

	m2 := model.ChannelMember{}
	m2.ChannelId = o1.Id
	m2.UserId = u2.Id
	m2.NotifyProps = model.GetDefaultChannelNotifyProps()

	ss.Channel().SaveDirectChannel(&o1, &m1, &m2)

	d1, err := ss.Channel().GetAllDirectChannelsForExportAfter(10000, strings.Repeat("0", 26))
	assert.Nil(t, err)

	assert.Equal(t, 2, len(d1))
	assert.ElementsMatch(t, []string{o1.DisplayName, o2.DisplayName}, []string{d1[0].DisplayName, d1[1].DisplayName})

	// Manually truncate Channels table until testlib can handle cleanups
	s.GetMaster().Exec("TRUNCATE Channels")
}

func testChannelStoreExportAllDirectChannelsExcludePrivateAndPublic(t *testing.T, ss store.Store, s SqlSupplier) {
	teamId := model.NewId()

	o1 := model.Channel{}
	o1.TeamId = teamId
	o1.DisplayName = "The Direct Channel" + model.NewId()
	o1.Name = "zz" + model.NewId() + "b"
	o1.Type = model.CHANNEL_DIRECT

	o2 := model.Channel{}
	o2.TeamId = teamId
	o2.DisplayName = "Channel2" + model.NewId()
	o2.Name = "zz" + model.NewId() + "b"
	o2.Type = model.CHANNEL_OPEN
	_, err := ss.Channel().Save(&o2, -1)
	require.Nil(t, err)

	o3 := model.Channel{}
	o3.TeamId = teamId
	o3.DisplayName = "Channel3" + model.NewId()
	o3.Name = "zz" + model.NewId() + "b"
	o3.Type = model.CHANNEL_PRIVATE
	_, err = ss.Channel().Save(&o3, -1)
	require.Nil(t, err)

	u1 := &model.User{}
	u1.Email = MakeEmail()
	u1.Nickname = model.NewId()
	store.Must(ss.User().Save(u1))
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: model.NewId(), UserId: u1.Id}, -1))

	u2 := &model.User{}
	u2.Email = MakeEmail()
	u2.Nickname = model.NewId()
	store.Must(ss.User().Save(u2))
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: model.NewId(), UserId: u2.Id}, -1))

	m1 := model.ChannelMember{}
	m1.ChannelId = o1.Id
	m1.UserId = u1.Id
	m1.NotifyProps = model.GetDefaultChannelNotifyProps()

	m2 := model.ChannelMember{}
	m2.ChannelId = o1.Id
	m2.UserId = u2.Id
	m2.NotifyProps = model.GetDefaultChannelNotifyProps()

	ss.Channel().SaveDirectChannel(&o1, &m1, &m2)

	d1, err := ss.Channel().GetAllDirectChannelsForExportAfter(10000, strings.Repeat("0", 26))
	assert.Nil(t, err)
	assert.Equal(t, 1, len(d1))
	assert.Equal(t, o1.DisplayName, d1[0].DisplayName)

	// Manually truncate Channels table until testlib can handle cleanups
	s.GetMaster().Exec("TRUNCATE Channels")
}

func testChannelStoreExportAllDirectChannelsDeletedChannel(t *testing.T, ss store.Store, s SqlSupplier) {
	teamId := model.NewId()

	o1 := model.Channel{}
	o1.TeamId = teamId
	o1.DisplayName = "Different Name" + model.NewId()
	o1.Name = "zz" + model.NewId() + "b"
	o1.Type = model.CHANNEL_DIRECT

	u1 := &model.User{}
	u1.Email = MakeEmail()
	u1.Nickname = model.NewId()
	store.Must(ss.User().Save(u1))
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: model.NewId(), UserId: u1.Id}, -1))

	u2 := &model.User{}
	u2.Email = MakeEmail()
	u2.Nickname = model.NewId()
	store.Must(ss.User().Save(u2))
	store.Must(ss.Team().SaveMember(&model.TeamMember{TeamId: model.NewId(), UserId: u2.Id}, -1))

	m1 := model.ChannelMember{}
	m1.ChannelId = o1.Id
	m1.UserId = u1.Id
	m1.NotifyProps = model.GetDefaultChannelNotifyProps()

	m2 := model.ChannelMember{}
	m2.ChannelId = o1.Id
	m2.UserId = u2.Id
	m2.NotifyProps = model.GetDefaultChannelNotifyProps()

	ss.Channel().SaveDirectChannel(&o1, &m1, &m2)

	o1.DeleteAt = 1
	err := ss.Channel().SetDeleteAt(o1.Id, 1, 1)
	require.Nil(t, err, "channel should have been deleted")

	d1, err := ss.Channel().GetAllDirectChannelsForExportAfter(10000, strings.Repeat("0", 26))
	assert.Nil(t, err)

	assert.Equal(t, 0, len(d1))

	// Manually truncate Channels table until testlib can handle cleanups
	s.GetMaster().Exec("TRUNCATE Channels")
}

func testChannelStoreGetChannelsBatchForIndexing(t *testing.T, ss store.Store) {
	// Set up all the objects needed
	c1 := &model.Channel{}
	c1.DisplayName = "Channel1"
	c1.Name = "zz" + model.NewId() + "b"
	c1.Type = model.CHANNEL_OPEN
	_, err := ss.Channel().Save(c1, -1)
	require.Nil(t, err)

	time.Sleep(10 * time.Millisecond)

	c2 := &model.Channel{}
	c2.DisplayName = "Channel2"
	c2.Name = "zz" + model.NewId() + "b"
	c2.Type = model.CHANNEL_OPEN
	_, err = ss.Channel().Save(c2, -1)
	require.Nil(t, err)

	time.Sleep(10 * time.Millisecond)
	startTime := c2.CreateAt

	c3 := &model.Channel{}
	c3.DisplayName = "Channel3"
	c3.Name = "zz" + model.NewId() + "b"
	c3.Type = model.CHANNEL_OPEN
	_, err = ss.Channel().Save(c3, -1)
	require.Nil(t, err)

	c4 := &model.Channel{}
	c4.DisplayName = "Channel4"
	c4.Name = "zz" + model.NewId() + "b"
	c4.Type = model.CHANNEL_PRIVATE
	_, err = ss.Channel().Save(c4, -1)
	require.Nil(t, err)

	c5 := &model.Channel{}
	c5.DisplayName = "Channel5"
	c5.Name = "zz" + model.NewId() + "b"
	c5.Type = model.CHANNEL_OPEN
	_, err = ss.Channel().Save(c5, -1)
	require.Nil(t, err)

	time.Sleep(10 * time.Millisecond)

	c6 := &model.Channel{}
	c6.DisplayName = "Channel6"
	c6.Name = "zz" + model.NewId() + "b"
	c6.Type = model.CHANNEL_OPEN
	_, err = ss.Channel().Save(c6, -1)
	require.Nil(t, err)

	endTime := c6.CreateAt

	// First and last channel should be outside the range
	channels, err := ss.Channel().GetChannelsBatchForIndexing(startTime, endTime, 1000)
	assert.Nil(t, err)
	assert.ElementsMatch(t, []*model.Channel{c2, c3, c5}, channels)

	// Update the endTime, last channel should be in
	endTime = model.GetMillis()
	channels, err = ss.Channel().GetChannelsBatchForIndexing(startTime, endTime, 1000)
	assert.Nil(t, err)
	assert.ElementsMatch(t, []*model.Channel{c2, c3, c5, c6}, channels)

	// Testing the limit
	channels, err = ss.Channel().GetChannelsBatchForIndexing(startTime, endTime, 2)
	assert.Nil(t, err)
	assert.ElementsMatch(t, []*model.Channel{c2, c3}, channels)
}
