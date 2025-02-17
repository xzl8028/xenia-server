// Copyright (c) 2015-present Xenia, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package main

import (
	"fmt"

	"github.com/xzl8028/xenia-server/model"
	"github.com/xzl8028/xenia-server/plugin"
)

type MyPlugin struct {
	plugin.XeniaPlugin
}

func (p *MyPlugin) MessageWillBePosted(c *plugin.Context, post *model.Post) (*model.Post, string) {

	teams, err := p.API.SearchTeams("{{.BasicTeam.Name}}")
	if err != nil {
		return nil, "search failed: " + err.Message
	}
	if len(teams) != 1 {
		return nil, fmt.Sprintf("search failed, wrong number of teams: %v", len(teams))
	}

	teams, err = p.API.SearchTeams("{{.BasicTeam.DisplayName}}")
	if err != nil {
		return nil, "search failed: " + err.Message
	}
	if len(teams) != 1 {
		return nil, fmt.Sprintf("search failed, wrong number of teams: %v", len(teams))
	}

	teams, err = p.API.SearchTeams("{{.BasicTeam.Name}}"[:3])

	if err != nil {
		return nil, "search failed: " + err.Message
	}
	if len(teams) != 1 {
		return nil, fmt.Sprintf("search failed, wrong number of teams: %v", len(teams))
	}

	teams, err = p.API.SearchTeams("not found")
	if err != nil {
		return nil, "search failed: " + err.Message
	}
	if len(teams) != 0 {
		return nil, fmt.Sprintf("search failed, wrong number of teams: %v", len(teams))
	}
	return nil, ""
}

func main() {
	plugin.ClientMain(&MyPlugin{})
}
