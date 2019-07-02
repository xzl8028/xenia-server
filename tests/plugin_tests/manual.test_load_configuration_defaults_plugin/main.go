// Copyright (c) 2015-present Xenia, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package main

import (
	"github.com/xzl8028/xenia-server/model"
	"github.com/xzl8028/xenia-server/plugin"
)

type configuration struct {
	MyStringSetting string
	MyIntSetting    int
	MyBoolSetting   bool
}

type MyPlugin struct {
	plugin.XeniaPlugin

	configuration configuration
}

func (p *MyPlugin) OnConfigurationChange() error {
	if err := p.API.LoadPluginConfiguration(&p.configuration); err != nil {
		return err
	}

	return nil
}

func (p *MyPlugin) MessageWillBePosted(c *plugin.Context, post *model.Post) (*model.Post, string) {
	if p.configuration.MyStringSetting != "override" {
		return nil, "MyStringSetting has invalid value"
	}
	if p.configuration.MyIntSetting != 35 {
		return nil, "MyIntSetting has invalid value"
	}
	if p.configuration.MyBoolSetting != true {
		return nil, "MyBoolSetting has invalid value"
	}
	return nil, ""
}

func main() {
	plugin.ClientMain(&MyPlugin{})
}
