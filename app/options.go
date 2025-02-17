// Copyright (c) 2017-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package app

import (
	"github.com/xzl8028/xenia-server/mlog"
	"github.com/pkg/errors"

	"github.com/xzl8028/xenia-server/config"
	"github.com/xzl8028/xenia-server/store"
)

type Option func(s *Server) error

// By default, the app will use the store specified by the configuration. This allows you to
// construct an app with a different store.
//
// The override parameter must be either a store.Store or func(App) store.Store.
func StoreOverride(override interface{}) Option {
	return func(s *Server) error {
		switch o := override.(type) {
		case store.Store:
			s.newStore = func() store.Store {
				return o
			}
			return nil

		case func(*Server) store.Store:
			s.newStore = func() store.Store {
				return o(s)
			}
			return nil

		default:
			return errors.New("invalid StoreOverride")
		}
	}
}

// Config applies the given config dsn, whether a path to config.json or a database connection string.
func Config(dsn string, watch bool) Option {
	return func(s *Server) error {
		configStore, err := config.NewStore(dsn, watch)
		if err != nil {
			return errors.Wrap(err, "failed to apply Config option")
		}

		s.configStore = configStore
		return nil
	}
}

// ConfigStore applies the given config store, typically to replace the traditional sources with a memory store for testing.
func ConfigStore(configStore config.Store) Option {
	return func(s *Server) error {
		s.configStore = configStore

		return nil
	}
}

func RunJobs(s *Server) error {
	s.runjobs = true

	return nil
}

func JoinCluster(s *Server) error {
	s.joinCluster = true

	return nil
}

func StartMetrics(s *Server) error {
	s.startMetrics = true

	return nil
}

func StartElasticsearch(s *Server) error {
	s.startElasticsearch = true

	return nil
}

func SetLogger(logger *mlog.Logger) Option {
	return func(s *Server) error {
		s.Log = logger
		return nil
	}
}

type AppOption func(a *App)
type AppOptionCreator func() []AppOption

func ServerConnector(s *Server) AppOption {
	return func(a *App) {
		a.Srv = s

		a.Log = s.Log
		a.NotificationsLog = s.NotificationsLog

		a.AccountMigration = s.AccountMigration
		a.Cluster = s.Cluster
		a.Compliance = s.Compliance
		a.DataRetention = s.DataRetention
		a.Elasticsearch = s.Elasticsearch
		a.Ldap = s.Ldap
		a.MessageExport = s.MessageExport
		a.Metrics = s.Metrics
		a.Saml = s.Saml

		a.HTTPService = s.HTTPService
		a.ImageProxy = s.ImageProxy
		a.Timezones = s.timezones
	}
}
