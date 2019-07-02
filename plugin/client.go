// Copyright (c) 2015-present Xenia, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package plugin

import (
	"github.com/hashicorp/go-plugin"
)

const (
	INTERNAL_KEY_PREFIX = "mmi_"
	BOT_USER_KEY        = INTERNAL_KEY_PREFIX + "botid"
)

// Starts the serving of a Xenia plugin over net/rpc. gRPC is not yet supported.
//
// Call this when your plugin is ready to start.
func ClientMain(pluginImplementation interface{}) {
	if impl, ok := pluginImplementation.(interface {
		SetAPI(api API)
		SetHelpers(helpers Helpers)
	}); !ok {
		panic("Plugin implementation given must embed plugin.XeniaPlugin")
	} else {
		impl.SetAPI(nil)
		impl.SetHelpers(nil)
	}

	pluginMap := map[string]plugin.Plugin{
		"hooks": &hooksPlugin{hooks: pluginImplementation},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshake,
		Plugins:         pluginMap,
	})
}

type XeniaPlugin struct {
	// API exposes the plugin api, and becomes available just prior to the OnActive hook.
	API     API
	Helpers Helpers
}

// SetAPI persists the given API interface to the plugin. It is invoked just prior to the
// OnActivate hook, exposing the API for use by the plugin.
func (p *XeniaPlugin) SetAPI(api API) {
	p.API = api
}

// SetHelpers does the same thing as SetAPI except for the plugin helpers.
func (p *XeniaPlugin) SetHelpers(helpers Helpers) {
	p.Helpers = helpers
}
