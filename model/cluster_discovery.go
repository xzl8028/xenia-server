// Copyright (c) 2017-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package model

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

const (
	CDS_OFFLINE_AFTER_MILLIS = 1000 * 60 * 30 // 30 minutes
	CDS_TYPE_APP             = "xenia_app"
)

type ClusterDiscovery struct {
	Id          string `json:"id"`
	Type        string `json:"type"`
	ClusterName string `json:"cluster_name"`
	Hostname    string `json:"hostname"`
	GossipPort  int32  `json:"gossip_port"`
	Port        int32  `json:"port"`
	CreateAt    int64  `json:"create_at"`
	LastPingAt  int64  `json:"last_ping_at"`
}

func (o *ClusterDiscovery) PreSave() {
	if o.Id == "" {
		o.Id = NewId()
	}

	if o.CreateAt == 0 {
		o.CreateAt = GetMillis()
		o.LastPingAt = o.CreateAt
	}
}

func (o *ClusterDiscovery) AutoFillHostname() {
	// attempt to set the hostname from the OS
	if len(o.Hostname) == 0 {
		if hn, err := os.Hostname(); err == nil {
			o.Hostname = hn
		}
	}
}

func (o *ClusterDiscovery) AutoFillIpAddress(iface string, ipAddress string) {
	// attempt to set the hostname to the first non-local IP address
	if len(o.Hostname) == 0 {
		if len(ipAddress) > 0 {
			o.Hostname = ipAddress
		} else {
			o.Hostname = GetServerIpAddress(iface)
		}
	}
}

func (o *ClusterDiscovery) IsEqual(in *ClusterDiscovery) bool {
	if in == nil {
		return false
	}

	if o.Type != in.Type {
		return false
	}

	if o.ClusterName != in.ClusterName {
		return false
	}

	if o.Hostname != in.Hostname {
		return false
	}

	return true
}

func FilterClusterDiscovery(vs []*ClusterDiscovery, f func(*ClusterDiscovery) bool) []*ClusterDiscovery {
	copy := make([]*ClusterDiscovery, 0)
	for _, v := range vs {
		if f(v) {
			copy = append(copy, v)
		}
	}

	return copy
}

func (o *ClusterDiscovery) IsValid() *AppError {
	if len(o.Id) != 26 {
		return NewAppError("ClusterDiscovery.IsValid", "model.cluster.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}

	if len(o.ClusterName) == 0 {
		return NewAppError("ClusterDiscovery.IsValid", "model.cluster.is_valid.name.app_error", nil, "", http.StatusBadRequest)
	}

	if len(o.Type) == 0 {
		return NewAppError("ClusterDiscovery.IsValid", "model.cluster.is_valid.type.app_error", nil, "", http.StatusBadRequest)
	}

	if len(o.Hostname) == 0 {
		return NewAppError("ClusterDiscovery.IsValid", "model.cluster.is_valid.hostname.app_error", nil, "", http.StatusBadRequest)
	}

	if o.CreateAt == 0 {
		return NewAppError("ClusterDiscovery.IsValid", "model.cluster.is_valid.create_at.app_error", nil, "", http.StatusBadRequest)
	}

	if o.LastPingAt == 0 {
		return NewAppError("ClusterDiscovery.IsValid", "model.cluster.is_valid.last_ping_at.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

func (o *ClusterDiscovery) ToJson() string {
	b, err := json.Marshal(o)
	if err != nil {
		return ""
	}

	return string(b)
}

func ClusterDiscoveryFromJson(data io.Reader) *ClusterDiscovery {
	decoder := json.NewDecoder(data)
	var me ClusterDiscovery
	err := decoder.Decode(&me)
	if err == nil {
		return &me
	}

	return nil
}
