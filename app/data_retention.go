// Copyright (c) 2017-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package app

import (
	"net/http"

	"github.com/xzl8028/xenia-server/model"
)

func (a *App) GetDataRetentionPolicy() (*model.DataRetentionPolicy, *model.AppError) {
	if a.DataRetention == nil {
		return nil, model.NewAppError("App.GetDataRetentionPolicy", "ent.data_retention.generic.license.error", nil, "", http.StatusNotImplemented)
	}

	return a.DataRetention.GetPolicy()
}
