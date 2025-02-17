// Copyright (c) 2017-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package plugin

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/xzl8028/xenia-server/mlog"
	"github.com/xzl8028/xenia-server/model"
	"github.com/xzl8028/xenia-server/utils"
	"github.com/stretchr/testify/require"
)

func TestPluginHealthCheck(t *testing.T) {
	for name, f := range map[string]func(*testing.T){
		"PluginHealthCheck_Success": testPluginHealthCheck_Success,
		"PluginHealthCheck_Panic":   testPluginHealthCheck_Panic,
	} {
		t.Run(name, f)
	}
}

func testPluginHealthCheck_Success(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	backend := filepath.Join(dir, "backend.exe")
	utils.CompileGo(t, `
		package main

		import (
			"github.com/xzl8028/xenia-server/plugin"
		)

		type MyPlugin struct {
			plugin.XeniaPlugin
		}

		func main() {
			plugin.ClientMain(&MyPlugin{})
		}
	`, backend)

	err = ioutil.WriteFile(filepath.Join(dir, "plugin.json"), []byte(`{"id": "foo", "backend": {"executable": "backend.exe"}}`), 0600)
	require.NoError(t, err)

	bundle := model.BundleInfoForPath(dir)
	log := mlog.NewLogger(&mlog.LoggerConfiguration{
		EnableConsole: true,
		ConsoleJson:   true,
		ConsoleLevel:  "error",
		EnableFile:    false,
	})

	supervisor, err := newSupervisor(bundle, log, nil)
	require.Nil(t, err)
	require.NotNil(t, supervisor)

	err = supervisor.PerformHealthCheck()
	require.Nil(t, err)
}

func testPluginHealthCheck_Panic(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	backend := filepath.Join(dir, "backend.exe")
	utils.CompileGo(t, `
		package main

		import (
			"github.com/xzl8028/xenia-server/model"
			"github.com/xzl8028/xenia-server/plugin"
		)

		type MyPlugin struct {
			plugin.XeniaPlugin
		}

		func (p *MyPlugin) MessageWillBePosted(c *plugin.Context, post *model.Post) (*model.Post, string) {
			panic("Uncaught error")
		}

		func main() {
			plugin.ClientMain(&MyPlugin{})
		}
	`, backend)

	err = ioutil.WriteFile(filepath.Join(dir, "plugin.json"), []byte(`{"id": "foo", "backend": {"executable": "backend.exe"}}`), 0600)
	require.NoError(t, err)

	bundle := model.BundleInfoForPath(dir)
	log := mlog.NewLogger(&mlog.LoggerConfiguration{
		EnableConsole: true,
		ConsoleJson:   true,
		ConsoleLevel:  "error",
		EnableFile:    false,
	})

	supervisor, err := newSupervisor(bundle, log, nil)
	require.Nil(t, err)
	require.NotNil(t, supervisor)

	err = supervisor.PerformHealthCheck()
	require.Nil(t, err)

	supervisor.hooks.MessageWillBePosted(&Context{}, &model.Post{})

	err = supervisor.PerformHealthCheck()
	require.NotNil(t, err)
}

func TestShouldDeactivatePlugin(t *testing.T) {
	bundle := &model.BundleInfo{}
	rp := newRegisteredPlugin(bundle)
	require.NotNil(t, rp)

	// No failures, don't restart
	result := shouldDeactivatePlugin(rp)
	require.Equal(t, false, result)

	now := time.Now()

	// Failures are recent enough to restart
	rp = newRegisteredPlugin(bundle)
	rp.failTimeStamps = append(rp.failTimeStamps, now.Add(-HEALTH_CHECK_DISABLE_DURATION*0.2*time.Minute))
	rp.failTimeStamps = append(rp.failTimeStamps, now.Add(-HEALTH_CHECK_DISABLE_DURATION*0.1*time.Minute))
	rp.failTimeStamps = append(rp.failTimeStamps, now)

	result = shouldDeactivatePlugin(rp)
	require.Equal(t, true, result)

	// Failures are too spaced out to warrant a restart
	rp = newRegisteredPlugin(bundle)
	rp.failTimeStamps = append(rp.failTimeStamps, now.Add(-HEALTH_CHECK_DISABLE_DURATION*2*time.Minute))
	rp.failTimeStamps = append(rp.failTimeStamps, now.Add(-HEALTH_CHECK_DISABLE_DURATION*1*time.Minute))
	rp.failTimeStamps = append(rp.failTimeStamps, now)

	result = shouldDeactivatePlugin(rp)
	require.Equal(t, false, result)

	// Not enough failures are present to warrant a restart
	rp = newRegisteredPlugin(bundle)
	rp.failTimeStamps = append(rp.failTimeStamps, now.Add(-HEALTH_CHECK_DISABLE_DURATION*0.1*time.Minute))
	rp.failTimeStamps = append(rp.failTimeStamps, now)

	result = shouldDeactivatePlugin(rp)
	require.Equal(t, false, result)
}
