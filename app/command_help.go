// Copyright (c) 2017-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package app

import (
	goi18n "github.com/xzl8028/go-i18n/i18n"
	"github.com/xzl8028/xenia-server/model"
)

type HelpProvider struct {
}

const (
	CMD_HELP = "help"
)

func init() {
	RegisterCommandProvider(&HelpProvider{})
}

func (h *HelpProvider) GetTrigger() string {
	return CMD_HELP
}

func (h *HelpProvider) GetCommand(a *App, T goi18n.TranslateFunc) *model.Command {
	return &model.Command{
		Trigger:          CMD_HELP,
		AutoComplete:     true,
		AutoCompleteDesc: T("api.command_help.desc"),
		DisplayName:      T("api.command_help.name"),
	}
}

func (h *HelpProvider) DoCommand(a *App, args *model.CommandArgs, message string) *model.CommandResponse {
	helpLink := *a.Config().SupportSettings.HelpLink

	if helpLink == "" {
		helpLink = model.SUPPORT_SETTINGS_DEFAULT_HELP_LINK
	}

	return &model.CommandResponse{GotoLocation: helpLink}
}
