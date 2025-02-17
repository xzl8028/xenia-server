package plugin_test

import (
	"strings"
	"sync"

	"github.com/pkg/errors"

	"github.com/xzl8028/xenia-server/model"
	"github.com/xzl8028/xenia-server/plugin"
)

// configuration represents the configuration for this plugin as exposed via the Xenia
// server configuration.
type configuration struct {
	TeamName    string
	ChannelName string

	// channelId is resolved when the public configuration fields above change
	channelId string
}

type HelpPlugin struct {
	plugin.XeniaPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration
}

// getConfiguration retrieves the active configuration under lock, making it safe to use
// concurrently. The active configuration may change underneath the client of this method, but
// the struct returned by this API call is considered immutable.
func (p *HelpPlugin) getConfiguration() *configuration {
	p.configurationLock.RLock()
	defer p.configurationLock.RUnlock()

	if p.configuration == nil {
		return &configuration{}
	}

	return p.configuration
}

// setConfiguration replaces the active configuration under lock.
//
// Do not call setConfiguration while holding the configurationLock, as sync.Mutex is not
// reentrant.
func (p *HelpPlugin) setConfiguration(configuration *configuration) {
	// Replace the active configuration under lock.
	p.configurationLock.Lock()
	defer p.configurationLock.Unlock()
	p.configuration = configuration
}

// OnConfigurationChange updates the active configuration for this plugin under lock.
func (p *HelpPlugin) OnConfigurationChange() error {
	var configuration = new(configuration)

	// Load the public configuration fields from the Xenia server configuration.
	if err := p.API.LoadPluginConfiguration(configuration); err != nil {
		return errors.Wrap(err, "failed to load plugin configuration")
	}

	team, err := p.API.GetTeamByName(configuration.TeamName)
	if err != nil {
		return errors.Wrapf(err, "failed to find team %s", configuration.TeamName)
	}

	channel, err := p.API.GetChannelByName(configuration.ChannelName, team.Id, false)
	if err != nil {
		return errors.Wrapf(err, "failed to find channel %s", configuration.ChannelName)
	}

	configuration.channelId = channel.Id

	p.setConfiguration(configuration)

	return nil
}

func (p *HelpPlugin) MessageHasBeenPosted(c *plugin.Context, post *model.Post) {
	configuration := p.getConfiguration()

	// Ignore posts not in the configured channel
	if post.ChannelId != configuration.channelId {
		return
	}

	// Ignore posts this plugin made.
	if sentByPlugin, _ := post.Props["sent_by_plugin"].(bool); sentByPlugin {
		return
	}

	// Ignore posts without a plea for help.
	if !strings.Contains(post.Message, "help") {
		return
	}

	p.API.SendEphemeralPost(post.UserId, &model.Post{
		ChannelId: configuration.channelId,
		Message:   "You asked for help? Checkout https://about.xenia.com/help/",
		Props: map[string]interface{}{
			"sent_by_plugin": true,
		},
	})
}

func Example_helpPlugin() {
	plugin.ClientMain(&HelpPlugin{})
}
