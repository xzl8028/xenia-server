// Copyright (c) 2016 Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package mfa

import (
	b32 "encoding/base32"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/dgryski/dgoogauth"
	"github.com/xzl8028/xenia-server/model"
	"github.com/xzl8028/xenia-server/services/configservice"
	"github.com/xzl8028/xenia-server/store"
	"github.com/xzl8028/rsc/qr"
)

const (
	MFA_SECRET_SIZE = 20
)

type Mfa struct {
	ConfigService configservice.ConfigService
	Store         store.Store
}

func New(configService configservice.ConfigService, store store.Store) Mfa {
	return Mfa{configService, store}
}

func (m *Mfa) checkConfig() *model.AppError {
	if !*m.ConfigService.Config().ServiceSettings.EnableMultifactorAuthentication {
		return model.NewAppError("checkConfig", "mfa.mfa_disabled.app_error", nil, "", http.StatusNotImplemented)
	}

	return nil
}

func getIssuerFromUrl(uri string) string {
	issuer := "Xenia"
	siteUrl := strings.TrimSpace(uri)

	if len(siteUrl) > 0 {
		siteUrl = strings.TrimPrefix(siteUrl, "https://")
		siteUrl = strings.TrimPrefix(siteUrl, "http://")
		issuer = strings.TrimPrefix(siteUrl, "www.")
	}

	return url.QueryEscape(issuer)
}

func (m *Mfa) GenerateSecret(user *model.User) (string, []byte, *model.AppError) {
	if err := m.checkConfig(); err != nil {
		return "", nil, err
	}

	issuer := getIssuerFromUrl(*m.ConfigService.Config().ServiceSettings.SiteURL)

	secret := b32.StdEncoding.EncodeToString([]byte(model.NewRandomString(MFA_SECRET_SIZE)))

	authLink := fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s", issuer, user.Email, secret, issuer)

	code, err := qr.Encode(authLink, qr.H)

	if err != nil {
		return "", nil, model.NewAppError("GenerateQrCode", "mfa.generate_qr_code.create_code.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	img := code.PNG()

	if result := <-m.Store.User().UpdateMfaSecret(user.Id, secret); result.Err != nil {
		return "", nil, model.NewAppError("GenerateQrCode", "mfa.generate_qr_code.save_secret.app_error", nil, result.Err.Error(), http.StatusInternalServerError)
	}

	return secret, img, nil
}

func (m *Mfa) Activate(user *model.User, token string) *model.AppError {
	if err := m.checkConfig(); err != nil {
		return err
	}

	otpConfig := &dgoogauth.OTPConfig{
		Secret:      user.MfaSecret,
		WindowSize:  3,
		HotpCounter: 0,
	}

	trimmedToken := strings.TrimSpace(token)

	ok, err := otpConfig.Authenticate(trimmedToken)
	if err != nil {
		return model.NewAppError("Activate", "mfa.activate.authenticate.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	if !ok {
		return model.NewAppError("Activate", "mfa.activate.bad_token.app_error", nil, "", http.StatusUnauthorized)
	}

	if result := <-m.Store.User().UpdateMfaActive(user.Id, true); result.Err != nil {
		return model.NewAppError("Activate", "mfa.activate.save_active.app_error", nil, result.Err.Error(), http.StatusInternalServerError)
	}

	return nil
}

func (m *Mfa) Deactivate(userId string) *model.AppError {
	if err := m.checkConfig(); err != nil {
		return err
	}

	achan := m.Store.User().UpdateMfaActive(userId, false)
	schan := m.Store.User().UpdateMfaSecret(userId, "")

	if result := <-achan; result.Err != nil {
		return model.NewAppError("Deactivate", "mfa.deactivate.save_active.app_error", nil, result.Err.Error(), http.StatusInternalServerError)
	}

	if result := <-schan; result.Err != nil {
		return model.NewAppError("Deactivate", "mfa.deactivate.save_secret.app_error", nil, result.Err.Error(), http.StatusInternalServerError)
	}

	return nil
}

func (m *Mfa) ValidateToken(secret, token string) (bool, *model.AppError) {
	if err := m.checkConfig(); err != nil {
		return false, err
	}

	otpConfig := &dgoogauth.OTPConfig{
		Secret:      secret,
		WindowSize:  3,
		HotpCounter: 0,
	}

	trimmedToken := strings.TrimSpace(token)
	ok, err := otpConfig.Authenticate(trimmedToken)
	if err != nil {
		return false, model.NewAppError("ValidateToken", "mfa.validate_token.authenticate.app_error", nil, err.Error(), http.StatusBadRequest)
	}

	return ok, nil
}
