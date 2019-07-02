// Copyright (c) 2015-present Xenia, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package plugin

import (
	"encoding/json"
	"time"

	"github.com/xzl8028/xenia-server/model"
	"github.com/pkg/errors"
)

func (p *HelpersImpl) EnsureBot(bot *model.Bot) (retBotId string, retErr error) {
	// Must provide a bot with a username
	if bot == nil || len(bot.Username) < 1 {
		return "", errors.New("passed a bad bot, nil or no username")
	}

	// If we fail for any reason, this could be a race between creation of bot and
	// retreval from anouther EnsureBot. Just try the basic retrieve existing again.
	defer func() {
		if retBotId == "" || retErr != nil {
			time.Sleep(time.Second)
			botIdBytes, err := p.API.KVGet(BOT_USER_KEY)
			if err == nil && botIdBytes != nil {
				retBotId = string(botIdBytes)
				retErr = nil
			}
		}
	}()

	botIdBytes, kvGetErr := p.API.KVGet(BOT_USER_KEY)
	if kvGetErr != nil {
		return "", errors.Wrap(kvGetErr, "failed to get bot")
	}

	// If the bot has already been created, there is nothing to do.
	if botIdBytes != nil {
		botId := string(botIdBytes)
		return botId, nil
	}

	// Check for an existing bot user with that username. If one exists, then use that.
	if user, userGetErr := p.API.GetUserByUsername(bot.Username); userGetErr == nil && user != nil {
		if user.IsBot {
			if kvSetErr := p.API.KVSet(BOT_USER_KEY, []byte(user.Id)); kvSetErr != nil {
				p.API.LogWarn("Failed to set claimed bot user id.", "userid", user.Id, "err", kvSetErr)
			}
		} else {
			p.API.LogError("Plugin attempted to use an account that already exists. Convert user to a bot account in the CLI by running 'xenia user convert <username> --bot'. If the user is an existing user account you want to preserve, change its username and restart the Xenia server, after which the plugin will create a bot account with that name. For more information about bot accounts, see https://xenia.com/pl/default-bot-accounts", "username", bot.Username, "user_id", user.Id)
		}
		return user.Id, nil
	}

	// Create a new bot user for the plugin
	createdBot, createBotErr := p.API.CreateBot(bot)
	if createBotErr != nil {
		return "", errors.Wrap(createBotErr, "failed to create bot")
	}

	if kvSetErr := p.API.KVSet(BOT_USER_KEY, []byte(createdBot.UserId)); kvSetErr != nil {
		p.API.LogWarn("Failed to set created bot user id.", "userid", createdBot.UserId, "err", kvSetErr)
	}

	return createdBot.UserId, nil
}

func (p *HelpersImpl) KVGetJSON(key string, value interface{}) error {
	data, err := p.API.KVGet(key)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, value)
}

func (p *HelpersImpl) KVSetJSON(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return p.API.KVSet(key, data)
}

func (p *HelpersImpl) KVCompareAndSetJSON(key string, oldValue interface{}, newValue interface{}) (bool, error) {
	oldData, err := json.Marshal(oldValue)
	if err != nil {
		return false, errors.Wrap(err, "unable to marshal old value")
	}

	newData, err := json.Marshal(newValue)
	if err != nil {
		return false, errors.Wrap(err, "unable to marshal new value")
	}

	return p.API.KVCompareAndSet(key, oldData, newData)
}

func (p *HelpersImpl) KVSetWithExpiryJSON(key string, value interface{}, expireInSeconds int64) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return p.API.KVSetWithExpiry(key, data, expireInSeconds)
}
