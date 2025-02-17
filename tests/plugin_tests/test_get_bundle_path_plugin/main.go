// Copyright (c) 2015-present Xenia, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package main

import (
	"fmt"
	"path/filepath"

	"github.com/xzl8028/xenia-server/model"
	"github.com/xzl8028/xenia-server/plugin"
)

type MyPlugin struct {
	plugin.XeniaPlugin
}

func (p *MyPlugin) MessageWillBePosted(c *plugin.Context, post *model.Post) (*model.Post, string) {
	bundlePath, err := p.API.GetBundlePath()
	if err != nil {
		return nil, err.Error() + "failed get bundle path"
	} else if bundlePathFromConfig, _ := filepath.Abs(filepath.Join(*p.API.GetConfig().PluginSettings.Directory, "test_get_bundle_path_plugin")); bundlePathFromConfig != bundlePath {
		return nil, fmt.Sprintf("Invalid bundle path returned: %v vs %v", bundlePathFromConfig, bundlePath)
	}

	return nil, ""
}

func main() {
	plugin.ClientMain(&MyPlugin{})
}
