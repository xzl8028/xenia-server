// Copyright (c) 2015-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package einterfaces

import (
	"github.com/xzl8028/xenia-server/model"
)

type ClusterMessageHandler func(msg *model.ClusterMessage)

type ClusterInterface interface {
	StartInterNodeCommunication()
	StopInterNodeCommunication()
	RegisterClusterMessageHandler(event string, crm ClusterMessageHandler)
	GetClusterId() string
	IsLeader() bool
	GetMyClusterInfo() *model.ClusterInfo
	GetClusterInfos() []*model.ClusterInfo
	SendClusterMessage(cluster *model.ClusterMessage)
	NotifyMsg(buf []byte)
	GetClusterStats() ([]*model.ClusterStats, *model.AppError)
	GetLogs(page, perPage int) ([]string, *model.AppError)
	GetPluginStatuses() (model.PluginStatuses, *model.AppError)
	ConfigChanged(previousConfig *model.Config, newConfig *model.Config, sendToOtherServer bool) *model.AppError
}
