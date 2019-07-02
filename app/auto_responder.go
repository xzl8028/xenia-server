// Copyright (c) 2015-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package app

import (
	"github.com/xzl8028/xenia-server/mlog"
	"github.com/xzl8028/xenia-server/model"
)

func (a *App) SendAutoResponse(channel *model.Channel, receiver *model.User) {
	if receiver == nil || receiver.NotifyProps == nil {
		return
	}

	active := receiver.NotifyProps[model.AUTO_RESPONDER_ACTIVE_NOTIFY_PROP] == "true"
	message := receiver.NotifyProps[model.AUTO_RESPONDER_MESSAGE_NOTIFY_PROP]

	if active && message != "" {
		autoResponderPost := &model.Post{
			ChannelId: channel.Id,
			Message:   message,
			RootId:    "",
			ParentId:  "",
			Type:      model.POST_AUTO_RESPONDER,
			UserId:    receiver.Id,
		}

		if _, err := a.CreatePost(autoResponderPost, channel, false); err != nil {
			mlog.Error(err.Error())
		}
	}
}

func (a *App) SetAutoResponderStatus(user *model.User, oldNotifyProps model.StringMap) {
	active := user.NotifyProps[model.AUTO_RESPONDER_ACTIVE_NOTIFY_PROP] == "true"
	oldActive := oldNotifyProps[model.AUTO_RESPONDER_ACTIVE_NOTIFY_PROP] == "true"

	autoResponderEnabled := !oldActive && active
	autoResponderDisabled := oldActive && !active

	if autoResponderEnabled {
		a.SetStatusOutOfOffice(user.Id)
	} else if autoResponderDisabled {
		a.SetStatusOnline(user.Id, true)
	}
}

func (a *App) DisableAutoResponder(userId string, asAdmin bool) *model.AppError {
	user, err := a.GetUser(userId)
	if err != nil {
		return err
	}

	active := user.NotifyProps[model.AUTO_RESPONDER_ACTIVE_NOTIFY_PROP] == "true"

	if active {
		patch := &model.UserPatch{}
		patch.NotifyProps = user.NotifyProps
		patch.NotifyProps[model.AUTO_RESPONDER_ACTIVE_NOTIFY_PROP] = "false"

		_, err := a.PatchUser(userId, patch, asAdmin)
		if err != nil {
			return err
		}
	}

	return nil
}
