// Copyright (c) 2015-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package storetest

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/xzl8028/xenia-server/model"
	"github.com/xzl8028/xenia-server/store"
)

func makeBotWithUser(ss store.Store, bot *model.Bot) (*model.Bot, *model.User) {
	user := store.Must(ss.User().Save(model.UserFromBot(bot))).(*model.User)

	bot.UserId = user.Id
	bot, err := ss.Bot().Save(bot)
	if err != nil {
		time.Sleep(time.Second)
		panic(err)
	}

	return bot, user
}

func TestBotStore(t *testing.T, ss store.Store) {
	t.Run("Get", func(t *testing.T) { testBotStoreGet(t, ss) })
	t.Run("GetAll", func(t *testing.T) { testBotStoreGetAll(t, ss) })
	t.Run("Save", func(t *testing.T) { testBotStoreSave(t, ss) })
	t.Run("Update", func(t *testing.T) { testBotStoreUpdate(t, ss) })
	t.Run("PermanentDelete", func(t *testing.T) { testBotStorePermanentDelete(t, ss) })
}

func testBotStoreGet(t *testing.T, ss store.Store) {
	deletedBot, _ := makeBotWithUser(ss, &model.Bot{
		Username:    "deleted_bot",
		Description: "A deleted bot",
		OwnerId:     model.NewId(),
	})
	deletedBot.DeleteAt = 1
	deletedBot, err := ss.Bot().Update(deletedBot)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.Bot().PermanentDelete(deletedBot.UserId)) }()
	defer func() { require.Nil(t, ss.User().PermanentDelete(deletedBot.UserId)) }()

	permanentlyDeletedBot, _ := makeBotWithUser(ss, &model.Bot{
		Username:    "permanently_deleted_bot",
		Description: "A permanently deleted bot",
		OwnerId:     model.NewId(),
		DeleteAt:    0,
	})
	require.Nil(t, ss.Bot().PermanentDelete(permanentlyDeletedBot.UserId))
	defer func() { require.Nil(t, ss.User().PermanentDelete(permanentlyDeletedBot.UserId)) }()

	b1, _ := makeBotWithUser(ss, &model.Bot{
		Username:    "b1",
		Description: "The first bot",
		OwnerId:     model.NewId(),
	})
	defer func() { require.Nil(t, ss.Bot().PermanentDelete(b1.UserId)) }()
	defer func() { require.Nil(t, ss.User().PermanentDelete(b1.UserId)) }()

	b2, _ := makeBotWithUser(ss, &model.Bot{
		Username:    "b2",
		Description: "The second bot",
		OwnerId:     model.NewId(),
	})
	defer func() { require.Nil(t, ss.Bot().PermanentDelete(b2.UserId)) }()
	defer func() { require.Nil(t, ss.User().PermanentDelete(b2.UserId)) }()

	t.Run("get non-existent bot", func(t *testing.T) {
		_, err := ss.Bot().Get("unknown", false)
		require.NotNil(t, err)
		require.Equal(t, http.StatusNotFound, err.StatusCode)
	})

	t.Run("get deleted bot", func(t *testing.T) {
		_, err := ss.Bot().Get(deletedBot.UserId, false)
		require.NotNil(t, err)
		require.Equal(t, http.StatusNotFound, err.StatusCode)
	})

	t.Run("get deleted bot, include deleted", func(t *testing.T) {
		bot, err := ss.Bot().Get(deletedBot.UserId, true)
		require.Nil(t, err)
		require.Equal(t, deletedBot, bot)
	})

	t.Run("get permanently deleted bot", func(t *testing.T) {
		_, err := ss.Bot().Get(permanentlyDeletedBot.UserId, false)
		require.NotNil(t, err)
		require.Equal(t, http.StatusNotFound, err.StatusCode)
	})

	t.Run("get bot 1", func(t *testing.T) {
		bot, err := ss.Bot().Get(b1.UserId, false)
		require.Nil(t, err)
		require.Equal(t, b1, bot)
	})

	t.Run("get bot 2", func(t *testing.T) {
		bot, err := ss.Bot().Get(b2.UserId, false)
		require.Nil(t, err)
		require.Equal(t, b2, bot)
	})
}

func testBotStoreGetAll(t *testing.T, ss store.Store) {
	OwnerId1 := model.NewId()
	OwnerId2 := model.NewId()

	deletedBot, _ := makeBotWithUser(ss, &model.Bot{
		Username:    "deleted_bot",
		Description: "A deleted bot",
		OwnerId:     OwnerId1,
	})
	deletedBot.DeleteAt = 1
	deletedBot, err := ss.Bot().Update(deletedBot)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.Bot().PermanentDelete(deletedBot.UserId)) }()
	defer func() { require.Nil(t, ss.User().PermanentDelete(deletedBot.UserId)) }()

	permanentlyDeletedBot, _ := makeBotWithUser(ss, &model.Bot{
		Username:    "permanently_deleted_bot",
		Description: "A permanently deleted bot",
		OwnerId:     OwnerId1,
		DeleteAt:    0,
	})
	require.Nil(t, ss.Bot().PermanentDelete(permanentlyDeletedBot.UserId))
	defer func() { require.Nil(t, ss.User().PermanentDelete(permanentlyDeletedBot.UserId)) }()

	b1, _ := makeBotWithUser(ss, &model.Bot{
		Username:    "b1",
		Description: "The first bot",
		OwnerId:     OwnerId1,
	})
	defer func() { require.Nil(t, ss.Bot().PermanentDelete(b1.UserId)) }()
	defer func() { require.Nil(t, ss.User().PermanentDelete(b1.UserId)) }()

	b2, _ := makeBotWithUser(ss, &model.Bot{
		Username:    "b2",
		Description: "The second bot",
		OwnerId:     OwnerId1,
	})
	defer func() { require.Nil(t, ss.Bot().PermanentDelete(b2.UserId)) }()
	defer func() { require.Nil(t, ss.User().PermanentDelete(b2.UserId)) }()

	t.Run("get original bots", func(t *testing.T) {
		bot, err := ss.Bot().GetAll(&model.BotGetOptions{Page: 0, PerPage: 10})
		require.Nil(t, err)
		require.Equal(t, []*model.Bot{
			b1,
			b2,
		}, bot)
	})

	b3, _ := makeBotWithUser(ss, &model.Bot{
		Username:    "b3",
		Description: "The third bot",
		OwnerId:     OwnerId1,
	})
	defer func() { require.Nil(t, ss.Bot().PermanentDelete(b3.UserId)) }()
	defer func() { require.Nil(t, ss.User().PermanentDelete(b3.UserId)) }()

	b4, _ := makeBotWithUser(ss, &model.Bot{
		Username:    "b4",
		Description: "The fourth bot",
		OwnerId:     OwnerId2,
	})
	defer func() { require.Nil(t, ss.Bot().PermanentDelete(b4.UserId)) }()
	defer func() { require.Nil(t, ss.User().PermanentDelete(b4.UserId)) }()

	deletedUser := model.User{
		Email:    MakeEmail(),
		Username: model.NewId(),
	}
	if err := (<-ss.User().Save(&deletedUser)).Err; err != nil {
		t.Fatal("couldn't save user", err)
	}
	deletedUser.DeleteAt = model.GetMillis()
	if _, err := ss.User().Update(&deletedUser, true); err != nil {
		t.Fatal("couldn't delete user", err)
	}
	defer func() { require.Nil(t, ss.User().PermanentDelete(deletedUser.Id)) }()
	ob5, _ := makeBotWithUser(ss, &model.Bot{
		Username:    "ob5",
		Description: "Orphaned bot 5",
		OwnerId:     deletedUser.Id,
	})
	defer func() { require.Nil(t, ss.Bot().PermanentDelete(b4.UserId)) }()
	defer func() { require.Nil(t, ss.User().PermanentDelete(b4.UserId)) }()

	t.Run("get newly created bot stoo", func(t *testing.T) {
		bots, err := ss.Bot().GetAll(&model.BotGetOptions{Page: 0, PerPage: 10})
		require.Nil(t, err)
		require.Equal(t, []*model.Bot{
			b1,
			b2,
			b3,
			b4,
			ob5,
		}, bots)
	})

	t.Run("get orphaned", func(t *testing.T) {
		bots, err := ss.Bot().GetAll(&model.BotGetOptions{Page: 0, PerPage: 10, OnlyOrphaned: true})
		require.Nil(t, err)
		require.Equal(t, []*model.Bot{
			ob5,
		}, bots)
	})

	t.Run("get page=0, per_page=2", func(t *testing.T) {
		bots, err := ss.Bot().GetAll(&model.BotGetOptions{Page: 0, PerPage: 2})
		require.Nil(t, err)
		require.Equal(t, []*model.Bot{
			b1,
			b2,
		}, bots)
	})

	t.Run("get page=1, limit=2", func(t *testing.T) {
		bots, err := ss.Bot().GetAll(&model.BotGetOptions{Page: 1, PerPage: 2})
		require.Nil(t, err)
		require.Equal(t, []*model.Bot{
			b3,
			b4,
		}, bots)
	})

	t.Run("get page=5, perpage=1000", func(t *testing.T) {
		bots, err := ss.Bot().GetAll(&model.BotGetOptions{Page: 5, PerPage: 1000})
		require.Nil(t, err)
		require.Equal(t, []*model.Bot{}, bots)
	})

	t.Run("get offset=0, limit=2, include deleted", func(t *testing.T) {
		bots, err := ss.Bot().GetAll(&model.BotGetOptions{Page: 0, PerPage: 2, IncludeDeleted: true})
		require.Nil(t, err)
		require.Equal(t, []*model.Bot{
			deletedBot,
			b1,
		}, bots)
	})

	t.Run("get offset=2, limit=2, include deleted", func(t *testing.T) {
		bots, err := ss.Bot().GetAll(&model.BotGetOptions{Page: 1, PerPage: 2, IncludeDeleted: true})
		require.Nil(t, err)
		require.Equal(t, []*model.Bot{
			b2,
			b3,
		}, bots)
	})

	t.Run("get offset=0, limit=10, creator id 1", func(t *testing.T) {
		bots, err := ss.Bot().GetAll(&model.BotGetOptions{Page: 0, PerPage: 10, OwnerId: OwnerId1})
		require.Nil(t, err)
		require.Equal(t, []*model.Bot{
			b1,
			b2,
			b3,
		}, bots)
	})

	t.Run("get offset=0, limit=10, creator id 2", func(t *testing.T) {
		bots, err := ss.Bot().GetAll(&model.BotGetOptions{Page: 0, PerPage: 10, OwnerId: OwnerId2})
		require.Nil(t, err)
		require.Equal(t, []*model.Bot{
			b4,
		}, bots)
	})

	t.Run("get offset=0, limit=10, include deleted, creator id 1", func(t *testing.T) {
		bots, err := ss.Bot().GetAll(&model.BotGetOptions{Page: 0, PerPage: 10, IncludeDeleted: true, OwnerId: OwnerId1})
		require.Nil(t, err)
		require.Equal(t, []*model.Bot{
			deletedBot,
			b1,
			b2,
			b3,
		}, bots)
	})

	t.Run("get offset=0, limit=10, include deleted, creator id 2", func(t *testing.T) {
		bots, err := ss.Bot().GetAll(&model.BotGetOptions{Page: 0, PerPage: 10, IncludeDeleted: true, OwnerId: OwnerId2})
		require.Nil(t, err)
		require.Equal(t, []*model.Bot{
			b4,
		}, bots)
	})
}

func testBotStoreSave(t *testing.T, ss store.Store) {
	t.Run("invalid bot", func(t *testing.T) {
		bot := &model.Bot{
			UserId:      model.NewId(),
			Username:    "invalid bot",
			Description: "description",
		}

		_, err := ss.Bot().Save(bot)
		require.NotNil(t, err)
		require.Equal(t, "model.bot.is_valid.username.app_error", err.Id)
	})

	t.Run("normal bot", func(t *testing.T) {
		bot := &model.Bot{
			Username:    "normal_bot",
			Description: "description",
			OwnerId:     model.NewId(),
		}

		user := store.Must(ss.User().Save(model.UserFromBot(bot))).(*model.User)
		defer func() { require.Nil(t, ss.User().PermanentDelete(user.Id)) }()
		bot.UserId = user.Id

		returnedNewBot, err := ss.Bot().Save(bot)
		require.Nil(t, err)
		defer func() { require.Nil(t, ss.Bot().PermanentDelete(bot.UserId)) }()

		// Verify the returned bot matches the saved bot, modulo expected changes
		require.NotEqual(t, 0, returnedNewBot.CreateAt)
		require.NotEqual(t, 0, returnedNewBot.UpdateAt)
		require.Equal(t, returnedNewBot.CreateAt, returnedNewBot.UpdateAt)
		bot.UserId = returnedNewBot.UserId
		bot.CreateAt = returnedNewBot.CreateAt
		bot.UpdateAt = returnedNewBot.UpdateAt
		bot.DeleteAt = 0
		require.Equal(t, bot, returnedNewBot)

		// Verify the actual bot in the database matches the saved bot.
		actualNewBot, err := ss.Bot().Get(bot.UserId, false)
		require.Nil(t, err)
		require.Equal(t, bot, actualNewBot)
	})
}

func testBotStoreUpdate(t *testing.T, ss store.Store) {
	t.Run("invalid bot should fail to update", func(t *testing.T) {
		existingBot, _ := makeBotWithUser(ss, &model.Bot{
			Username: "existing_bot",
			OwnerId:  model.NewId(),
		})
		defer func() { require.Nil(t, ss.Bot().PermanentDelete(existingBot.UserId)) }()
		defer func() { require.Nil(t, ss.User().PermanentDelete(existingBot.UserId)) }()

		bot := existingBot.Clone()
		bot.Username = "invalid username"
		_, err := ss.Bot().Update(bot)
		require.NotNil(t, err)
		require.Equal(t, "model.bot.is_valid.username.app_error", err.Id)
	})

	t.Run("existing bot should update", func(t *testing.T) {
		existingBot, _ := makeBotWithUser(ss, &model.Bot{
			Username: "existing_bot",
			OwnerId:  model.NewId(),
		})
		defer func() { require.Nil(t, ss.Bot().PermanentDelete(existingBot.UserId)) }()
		defer func() { require.Nil(t, ss.User().PermanentDelete(existingBot.UserId)) }()

		bot := existingBot.Clone()
		bot.OwnerId = model.NewId()
		bot.Description = "updated description"
		bot.CreateAt = 999999 // Ignored
		bot.UpdateAt = 999999 // Ignored
		bot.DeleteAt = 100000 // Allowed

		returnedBot, err := ss.Bot().Update(bot)
		require.Nil(t, err)

		// Verify the returned bot matches the updated bot, modulo expected timestamp changes
		require.Equal(t, existingBot.CreateAt, returnedBot.CreateAt)
		require.NotEqual(t, bot.UpdateAt, returnedBot.UpdateAt, "update should have advanced UpdateAt")
		require.True(t, returnedBot.UpdateAt > bot.UpdateAt, "update should have advanced UpdateAt")
		require.NotEqual(t, 99999, returnedBot.UpdateAt, "should have ignored user-provided UpdateAt")
		bot.CreateAt = returnedBot.CreateAt
		bot.UpdateAt = returnedBot.UpdateAt

		// Verify the actual (now deleted) bot in the database
		actualBot, err := ss.Bot().Get(bot.UserId, true)
		require.Nil(t, err)
		require.Equal(t, bot, actualBot)
	})

	t.Run("deleted bot should update, restoring", func(t *testing.T) {
		existingBot, _ := makeBotWithUser(ss, &model.Bot{
			Username: "existing_bot",
			OwnerId:  model.NewId(),
		})
		defer func() { require.Nil(t, ss.Bot().PermanentDelete(existingBot.UserId)) }()
		defer func() { require.Nil(t, ss.User().PermanentDelete(existingBot.UserId)) }()

		existingBot.DeleteAt = 100000
		existingBot, err := ss.Bot().Update(existingBot)
		require.Nil(t, err)

		bot := existingBot.Clone()
		bot.DeleteAt = 0

		returnedBot, err := ss.Bot().Update(bot)
		require.Nil(t, err)

		// Verify the returned bot matches the updated bot, modulo expected timestamp changes
		require.EqualValues(t, 0, returnedBot.DeleteAt)
		bot.UpdateAt = returnedBot.UpdateAt

		// Verify the actual bot in the database
		actualBot, err := ss.Bot().Get(bot.UserId, false)
		require.Nil(t, err)
		require.Equal(t, bot, actualBot)
	})
}

func testBotStorePermanentDelete(t *testing.T, ss store.Store) {
	b1, _ := makeBotWithUser(ss, &model.Bot{
		Username: "b1",
		OwnerId:  model.NewId(),
	})
	defer func() { require.Nil(t, ss.Bot().PermanentDelete(b1.UserId)) }()
	defer func() { require.Nil(t, ss.User().PermanentDelete(b1.UserId)) }()

	b2, _ := makeBotWithUser(ss, &model.Bot{
		Username: "b2",
		OwnerId:  model.NewId(),
	})
	defer func() { require.Nil(t, ss.Bot().PermanentDelete(b2.UserId)) }()
	defer func() { require.Nil(t, ss.User().PermanentDelete(b2.UserId)) }()

	t.Run("permanently delete a non-existent bot", func(t *testing.T) {
		err := ss.Bot().PermanentDelete("unknown")
		require.Nil(t, err)
	})

	t.Run("permanently delete bot", func(t *testing.T) {
		err := ss.Bot().PermanentDelete(b1.UserId)
		require.Nil(t, err)

		_, err = ss.Bot().Get(b1.UserId, false)
		require.NotNil(t, err)
		require.Equal(t, http.StatusNotFound, err.StatusCode)
	})
}
