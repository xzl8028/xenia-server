// Copyright (c) 2015-present Xenia, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package main

import (
	"github.com/xzl8028/xenia-server/model"
	"github.com/xzl8028/xenia-server/plugin"
)

type MyPlugin struct {
	plugin.XeniaPlugin
}

func (p *MyPlugin) MessageWillBePosted(c *plugin.Context, post *model.Post) (*model.Post, string) {

	channels, err := p.API.GetChannelsForTeamForUser("{{.BasicTeam.Id}}", "{{.BasicUser.Id}}", false)
	if err != nil {
		return nil, err.Error()
	}
	if len(channels) != 3 {
		return nil, "Returned invalid number of channels"
	}
	channels, err = p.API.GetChannelsForTeamForUser("invalidid", "{{.BasicUser.Id}}", false)
	if err == nil {
		return nil, "Expected to get an error while retrieving channels for invalid id"
	}
	if len(channels) != 0 {
		return nil, "Returned invalid number of channels"
	}
	return nil, ""
}

func main() {
	plugin.ClientMain(&MyPlugin{})
}
