// Copyright (c) 2017-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/blang/semver"
	"gopkg.in/yaml.v2"
)

type PluginOption struct {
	// The display name for the option.
	DisplayName string `json:"display_name" yaml:"display_name"`

	// The string value for the option.
	Value string `json:"value" yaml:"value"`
}

type PluginSetting struct {
	// The key that the setting will be assigned to in the configuration file.
	Key string `json:"key" yaml:"key"`

	// The display name for the setting.
	DisplayName string `json:"display_name" yaml:"display_name"`

	// The type of the setting.
	//
	// "bool" will result in a boolean true or false setting.
	//
	// "dropdown" will result in a string setting that allows the user to select from a list of
	// pre-defined options.
	//
	// "generated" will result in a string setting that is set to a random, cryptographically secure
	// string.
	//
	// "radio" will result in a string setting that allows the user to select from a short selection
	// of pre-defined options.
	//
	// "text" will result in a string setting that can be typed in manually.
	//
	// "longtext" will result in a multi line string that can be typed in manually.
	//
	// "username" will result in a text setting that will autocomplete to a username.
	Type string `json:"type" yaml:"type"`

	// The help text to display to the user. Supports Markdown formatting.
	HelpText string `json:"help_text" yaml:"help_text"`

	// The help text to display alongside the "Regenerate" button for settings of the "generated" type.
	RegenerateHelpText string `json:"regenerate_help_text,omitempty" yaml:"regenerate_help_text,omitempty"`

	// The placeholder to display for "text", "generated" and "username" types when blank.
	Placeholder string `json:"placeholder" yaml:"placeholder"`

	// The default value of the setting.
	Default interface{} `json:"default" yaml:"default"`

	// For "radio" or "dropdown" settings, this is the list of pre-defined options that the user can choose
	// from.
	Options []*PluginOption `json:"options,omitempty" yaml:"options,omitempty"`
}

type PluginSettingsSchema struct {
	// Optional text to display above the settings. Supports Markdown formatting.
	Header string `json:"header" yaml:"header"`

	// Optional text to display below the settings. Supports Markdown formatting.
	Footer string `json:"footer" yaml:"footer"`

	// A list of setting definitions.
	Settings []*PluginSetting `json:"settings" yaml:"settings"`
}

// The plugin manifest defines the metadata required to load and present your plugin. The manifest
// file should be named plugin.json or plugin.yaml and placed in the top of your
// plugin bundle.
//
// Example plugin.json:
//
//
//    {
//      "id": "com.mycompany.myplugin",
//      "name": "My Plugin",
//      "description": "This is my plugin",
//      "version": "0.1.0",
//      "min_server_version": "5.6.0",
//      "server": {
//        "executables": {
//          "linux-amd64": "server/dist/plugin-linux-amd64",
//          "darwin-amd64": "server/dist/plugin-darwin-amd64",
//          "windows-amd64": "server/dist/plugin-windows-amd64.exe"
//        }
//      },
//      "webapp": {
//          "bundle_path": "webapp/dist/main.js"
//      },
//      "settings_schema": {
//        "header": "Some header text",
//        "footer": "Some footer text",
//        "settings": [{
//          "key": "someKey",
//          "display_name": "Enable Extra Feature",
//          "type": "bool",
//          "help_text": "When true, an extra feature will be enabled!",
//          "default": "false"
//        }]
//      },
//      "props": {
//        "someKey": "someData"
//      }
//    }
type Manifest struct {
	// The id is a globally unique identifier that represents your plugin. Ids must be at least
	// 3 characters, at most 190 characters and must match ^[a-zA-Z0-9-_\.]+$.
	// Reverse-DNS notation using a name you control is a good option, e.g. "com.mycompany.myplugin".
	Id string `json:"id" yaml:"id"`

	// The name to be displayed for the plugin.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	// A description of what your plugin is and does.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// A version number for your plugin. Semantic versioning is recommended: http://semver.org
	Version string `json:"version" yaml:"version"`

	// The minimum Xenia server version required for your plugin.
	//
	// Minimum server version: 5.6
	MinServerVersion string `json:"min_server_version,omitempty" yaml:"min_server_version,omitempty"`

	// Server defines the server-side portion of your plugin.
	Server *ManifestServer `json:"server,omitempty" yaml:"server,omitempty"`

	// Backend is a deprecated flag for defining the server-side portion of your plugin. Going forward, use Server instead.
	Backend *ManifestServer `json:"backend,omitempty" yaml:"backend,omitempty"`

	// If your plugin extends the web app, you'll need to define webapp.
	Webapp *ManifestWebapp `json:"webapp,omitempty" yaml:"webapp,omitempty"`

	// To allow administrators to configure your plugin via the Xenia system console, you can
	// provide your settings schema.
	SettingsSchema *PluginSettingsSchema `json:"settings_schema,omitempty" yaml:"settings_schema,omitempty"`

	// Plugins can store any kind of data in Props to allow other plugins to use it.
	Props map[string]interface{} `json:"props,omitempty" yaml:"props,omitempty"`
}

type ManifestServer struct {
	// Executables are the paths to your executable binaries, specifying multiple entry points
	// for different platforms when bundled together in a single plugin.
	Executables *ManifestExecutables `json:"executables,omitempty" yaml:"executables,omitempty"`

	// Executable is the path to your executable binary. This should be relative to the root
	// of your bundle and the location of the manifest file.
	//
	// On Windows, this file must have a ".exe" extension.
	//
	// If your plugin is compiled for multiple platforms, consider bundling them together
	// and using the Executables field instead.
	Executable string `json:"executable" yaml:"executable"`
}

type ManifestExecutables struct {
	// LinuxAmd64 is the path to your executable binary for the corresponding platform
	LinuxAmd64 string `json:"linux-amd64,omitempty" yaml:"linux-amd64,omitempty"`
	// DarwinAmd64 is the path to your executable binary for the corresponding platform
	DarwinAmd64 string `json:"darwin-amd64,omitempty" yaml:"darwin-amd64,omitempty"`
	// WindowsAmd64 is the path to your executable binary for the corresponding platform
	// This file must have a ".exe" extension
	WindowsAmd64 string `json:"windows-amd64,omitempty" yaml:"windows-amd64,omitempty"`
}

type ManifestWebapp struct {
	// The path to your webapp bundle. This should be relative to the root of your bundle and the
	// location of the manifest file.
	BundlePath string `json:"bundle_path" yaml:"bundle_path"`

	// BundleHash is the 64-bit FNV-1a hash of the webapp bundle, computed when the plugin is loaded
	BundleHash []byte `json:"-"`
}

func (m *Manifest) ToJson() string {
	b, _ := json.Marshal(m)
	return string(b)
}

func ManifestListToJson(m []*Manifest) string {
	b, _ := json.Marshal(m)
	return string(b)
}

func ManifestFromJson(data io.Reader) *Manifest {
	var m *Manifest
	json.NewDecoder(data).Decode(&m)
	return m
}

func ManifestListFromJson(data io.Reader) []*Manifest {
	var manifests []*Manifest
	json.NewDecoder(data).Decode(&manifests)
	return manifests
}

func (m *Manifest) HasClient() bool {
	return m.Webapp != nil
}

func (m *Manifest) ClientManifest() *Manifest {
	cm := new(Manifest)
	*cm = *m
	cm.Name = ""
	cm.Description = ""
	cm.Server = nil
	if cm.Webapp != nil {
		cm.Webapp = new(ManifestWebapp)
		*cm.Webapp = *m.Webapp
		cm.Webapp.BundlePath = "/static/" + m.Id + "/" + fmt.Sprintf("%s_%x_bundle.js", m.Id, m.Webapp.BundleHash)
	}
	return cm
}

// GetExecutableForRuntime returns the path to the executable for the given runtime architecture.
//
// If the manifest defines multiple executables, but none match, or if only a single executable
// is defined, the Executable field will be returned. This method does not guarantee that the
// resulting binary can actually execute on the given platform.
func (m *Manifest) GetExecutableForRuntime(goOs, goArch string) string {
	server := m.Server

	// Support the deprecated backend parameter.
	if server == nil {
		server = m.Backend
	}

	if server == nil {
		return ""
	}

	var executable string
	if server.Executables != nil {
		if goOs == "linux" && goArch == "amd64" {
			executable = server.Executables.LinuxAmd64
		} else if goOs == "darwin" && goArch == "amd64" {
			executable = server.Executables.DarwinAmd64
		} else if goOs == "windows" && goArch == "amd64" {
			executable = server.Executables.WindowsAmd64
		}
	}

	if executable == "" {
		executable = server.Executable
	}

	return executable
}

func (m *Manifest) HasServer() bool {
	return m.Server != nil || m.Backend != nil
}

func (m *Manifest) HasWebapp() bool {
	return m.Webapp != nil
}

func (m *Manifest) MeetMinServerVersion(serverVersion string) (bool, error) {
	minServerVersion, err := semver.Parse(m.MinServerVersion)
	if err != nil {
		return false, errors.New("failed to parse MinServerVersion")
	}
	sv := semver.MustParse(serverVersion)
	if sv.LT(minServerVersion) {
		return false, nil
	}
	return true, nil
}

// FindManifest will find and parse the manifest in a given directory.
//
// In all cases other than a does-not-exist error, path is set to the path of the manifest file that was
// found.
//
// Manifests are JSON or YAML files named plugin.json, plugin.yaml, or plugin.yml.
func FindManifest(dir string) (manifest *Manifest, path string, err error) {
	for _, name := range []string{"plugin.yml", "plugin.yaml"} {
		path = filepath.Join(dir, name)
		f, ferr := os.Open(path)
		if ferr != nil {
			if !os.IsNotExist(ferr) {
				return nil, "", ferr
			}
			continue
		}
		b, ioerr := ioutil.ReadAll(f)
		f.Close()
		if ioerr != nil {
			return nil, path, ioerr
		}
		var parsed Manifest
		err = yaml.Unmarshal(b, &parsed)
		if err != nil {
			return nil, path, err
		}
		manifest = &parsed
		manifest.Id = strings.ToLower(manifest.Id)
		return manifest, path, nil
	}

	path = filepath.Join(dir, "plugin.json")
	f, ferr := os.Open(path)
	if ferr != nil {
		if os.IsNotExist(ferr) {
			path = ""
		}
		return nil, path, ferr
	}
	defer f.Close()
	var parsed Manifest
	err = json.NewDecoder(f).Decode(&parsed)
	if err != nil {
		return nil, path, err
	}
	manifest = &parsed
	manifest.Id = strings.ToLower(manifest.Id)
	return manifest, path, nil
}
