// Copyright (c) 2015 Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package mfa

import (
	"net/url"
	"testing"

	"github.com/xzl8028/xenia-server/model"
	"github.com/xzl8028/xenia-server/plugin/plugintest/mock"
	"github.com/xzl8028/xenia-server/store"
	"github.com/xzl8028/xenia-server/store/storetest/mocks"
	"github.com/xzl8028/xenia-server/utils/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateSecret(t *testing.T) {
	user := &model.User{Id: model.NewId(), Roles: "system_user"}

	config := model.Config{}
	config.SetDefaults()
	config.ServiceSettings.EnableMultifactorAuthentication = model.NewBool(true)
	configService := testutils.StaticConfigService{Cfg: &config}
	storeMock := mocks.Store{}
	userStoreMock := mocks.UserStore{}
	userStoreMock.On("UpdateMfaSecret", user.Id, mock.AnythingOfType("string")).Return(func(userId string, secret string) store.StoreChannel {
		return store.Do(func(result *store.StoreResult) {
			result.Data = nil
			result.Err = nil
		})
	})
	storeMock.On("User").Return(&userStoreMock)

	mfa := Mfa{configService, &storeMock}

	secret, img, err := mfa.GenerateSecret(user)
	require.Nil(t, err)

	assert.Len(t, secret, 32)

	if len(img) == 0 {
		t.Fatal("no image set")
	}

	config.ServiceSettings.EnableMultifactorAuthentication = model.NewBool(false)

	_, _, err = mfa.GenerateSecret(user)
	require.NotNil(t, err)
}

func TestGetIssuerFromUrl(t *testing.T) {
	cases := []struct {
		Input    string
		Expected string
	}{
		{"http://somewebsite.com", url.QueryEscape("somewebsite.com")},
		{"https://somewebsite.com", url.QueryEscape("somewebsite.com")},
		{"https://some.website.com", url.QueryEscape("some.website.com")},
		{" https://www.somewebsite.com", url.QueryEscape("somewebsite.com")},
		{"http://somewebsite.com/chat", url.QueryEscape("somewebsite.com/chat")},
		{"somewebsite.com ", url.QueryEscape("somewebsite.com")},
		{"http://localhost:8065", url.QueryEscape("localhost:8065")},
		{"", "Xenia"},
		{"  ", "Xenia"},
	}

	for _, c := range cases {
		assert.Equal(t, c.Expected, getIssuerFromUrl(c.Input))
	}
}
