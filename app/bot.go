// Copyright (c) 2016-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package app

import (
	"github.com/xzl8028/xenia-server/mlog"
	"github.com/xzl8028/xenia-server/model"
	"github.com/xzl8028/xenia-server/store"
	"github.com/xzl8028/xenia-server/utils"
)

// CreateBot creates the given bot and corresponding user.
func (a *App) CreateBot(bot *model.Bot) (*model.Bot, *model.AppError) {
	result := <-a.Srv.Store.User().Save(model.UserFromBot(bot))
	if result.Err != nil {
		return nil, result.Err
	}
	bot.UserId = result.Data.(*model.User).Id

	savedBot, err := a.Srv.Store.Bot().Save(bot)
	if err != nil {
		a.Srv.Store.User().PermanentDelete(bot.UserId)
		return nil, err
	}

	// Get the owner of the bot, if one exists. If not, don't send a message
	ownerUser, err := a.Srv.Store.User().Get(bot.OwnerId)
	if err != nil && err.Id != store.MISSING_ACCOUNT_ERROR {
		mlog.Error(err.Error())
		return nil, err
	} else if ownerUser != nil {
		// Send a message to the bot's creator to inform them that the bot needs to be added
		// to a team and channel after it's created
		channel, err := a.GetOrCreateDirectChannel(savedBot.UserId, bot.OwnerId)
		if err != nil {
			return nil, err
		}

		T := utils.GetUserTranslations(ownerUser.Locale)
		botAddPost := &model.Post{
			Type:      model.POST_ADD_BOT_TEAMS_CHANNELS,
			UserId:    savedBot.UserId,
			ChannelId: channel.Id,
			Message:   T("api.bot.teams_channels.add_message_mobile"),
		}

		if _, err := a.CreatePostAsUser(botAddPost, a.Session.Id); err != nil {
			return nil, err
		}
	}

	return savedBot, nil
}

// PatchBot applies the given patch to the bot and corresponding user.
func (a *App) PatchBot(botUserId string, botPatch *model.BotPatch) (*model.Bot, *model.AppError) {
	bot, err := a.GetBot(botUserId, true)
	if err != nil {
		return nil, err
	}

	bot.Patch(botPatch)

	user, err := a.Srv.Store.User().Get(botUserId)
	if err != nil {
		return nil, err
	}

	patchedUser := model.UserFromBot(bot)
	user.Id = patchedUser.Id
	user.Username = patchedUser.Username
	user.Email = patchedUser.Email
	user.FirstName = patchedUser.FirstName
	if _, err := a.Srv.Store.User().Update(user, true); err != nil {
		return nil, err
	}

	return a.Srv.Store.Bot().Update(bot)
}

// GetBot returns the given bot.
func (a *App) GetBot(botUserId string, includeDeleted bool) (*model.Bot, *model.AppError) {
	return a.Srv.Store.Bot().Get(botUserId, includeDeleted)
}

// GetBots returns the requested page of bots.
func (a *App) GetBots(options *model.BotGetOptions) (model.BotList, *model.AppError) {
	return a.Srv.Store.Bot().GetAll(options)
}

// UpdateBotActive marks a bot as active or inactive, along with its corresponding user.
func (a *App) UpdateBotActive(botUserId string, active bool) (*model.Bot, *model.AppError) {
	user, err := a.Srv.Store.User().Get(botUserId)
	if err != nil {
		return nil, err
	}

	if _, err = a.UpdateActive(user, active); err != nil {
		return nil, err
	}

	bot, err := a.Srv.Store.Bot().Get(botUserId, true)
	if err != nil {
		return nil, err
	}

	changed := true
	if active && bot.DeleteAt != 0 {
		bot.DeleteAt = 0
	} else if !active && bot.DeleteAt == 0 {
		bot.DeleteAt = model.GetMillis()
	} else {
		changed = false
	}

	if changed {
		bot, err = a.Srv.Store.Bot().Update(bot)
		if err != nil {
			return nil, err
		}
	}

	return bot, nil
}

// PermanentDeleteBot permanently deletes a bot and its corresponding user.
func (a *App) PermanentDeleteBot(botUserId string) *model.AppError {
	if err := a.Srv.Store.Bot().PermanentDelete(botUserId); err != nil {
		return err
	}

	if err := a.Srv.Store.User().PermanentDelete(botUserId); err != nil {
		return err
	}

	return nil
}

// UpdateBotOwner changes a bot's owner to the given value
func (a *App) UpdateBotOwner(botUserId, newOwnerId string) (*model.Bot, *model.AppError) {
	bot, err := a.Srv.Store.Bot().Get(botUserId, true)
	if err != nil {
		return nil, err
	}

	bot.OwnerId = newOwnerId

	bot, err = a.Srv.Store.Bot().Update(bot)
	if err != nil {
		return nil, err
	}

	return bot, nil
}

// disableUserBots disables all bots owned by the given user
func (a *App) disableUserBots(userId string) *model.AppError {
	perPage := 20
	for {
		options := &model.BotGetOptions{
			OwnerId:        userId,
			IncludeDeleted: false,
			OnlyOrphaned:   false,
			Page:           0,
			PerPage:        perPage,
		}
		userBots, err := a.GetBots(options)
		if err != nil {
			return err
		}

		for _, bot := range userBots {
			_, err := a.UpdateBotActive(bot.UserId, false)
			if err != nil {
				mlog.Error("Unable to deactivate bot.", mlog.String("bot_user_id", bot.UserId), mlog.Err(err))
			}
		}

		// Get next set of bots if we got the max number of bots
		if len(userBots) == perPage {
			options.Page += 1
			continue
		}
		break
	}

	return nil
}

// ConvertUserToBot converts a user to bot
func (a *App) ConvertUserToBot(user *model.User) (*model.Bot, *model.AppError) {
	return a.Srv.Store.Bot().Save(model.BotFromUser(user))
}
