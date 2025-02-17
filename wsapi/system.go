// Copyright (c) 2017-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package wsapi

import (
	"github.com/xzl8028/xenia-server/model"
)

func (api *API) InitSystem() {
	api.Router.Handle("ping", api.ApiWebSocketHandler(ping))
}

func ping(req *model.WebSocketRequest) (map[string]interface{}, *model.AppError) {
	data := map[string]interface{}{}
	data["text"] = "pong"
	data["version"] = model.CurrentVersion
	data["server_time"] = model.GetMillis()
	data["node_id"] = ""

	return data, nil
}
