// Copyright (c) 2016-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package app

import (
	"github.com/xzl8028/xenia-server/model"
)

func (a *App) CreateTermsOfService(text, userId string) (*model.TermsOfService, *model.AppError) {
	termsOfService := &model.TermsOfService{
		Text:   text,
		UserId: userId,
	}

	if _, err := a.GetUser(userId); err != nil {
		return nil, err
	}

	result := <-a.Srv.Store.TermsOfService().Save(termsOfService)
	if result.Err != nil {
		return nil, result.Err
	}

	termsOfService = result.Data.(*model.TermsOfService)
	return termsOfService, nil
}

func (a *App) GetLatestTermsOfService() (*model.TermsOfService, *model.AppError) {
	if result := <-a.Srv.Store.TermsOfService().GetLatest(true); result.Err != nil {
		return nil, result.Err
	} else {
		termsOfService := result.Data.(*model.TermsOfService)
		return termsOfService, nil
	}
}

func (a *App) GetTermsOfService(id string) (*model.TermsOfService, *model.AppError) {
	if result := <-a.Srv.Store.TermsOfService().Get(id, true); result.Err != nil {
		return nil, result.Err
	} else {
		termsOfService := result.Data.(*model.TermsOfService)
		return termsOfService, nil
	}
}
