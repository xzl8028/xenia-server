// Copyright (c) 2016-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package utils

import (
	"net/http"
	"strings"

	"github.com/xzl8028/xenia-server/model"
)

func IsPasswordValidWithSettings(password string, settings *model.PasswordSettings) *model.AppError {
	id := "model.user.is_valid.pwd"
	isError := false

	if len(password) < *settings.MinimumLength || len(password) > model.PASSWORD_MAXIMUM_LENGTH {
		isError = true
	}

	if *settings.Lowercase {
		if !strings.ContainsAny(password, model.LOWERCASE_LETTERS) {
			isError = true
		}

		id = id + "_lowercase"
	}

	if *settings.Uppercase {
		if !strings.ContainsAny(password, model.UPPERCASE_LETTERS) {
			isError = true
		}

		id = id + "_uppercase"
	}

	if *settings.Number {
		if !strings.ContainsAny(password, model.NUMBERS) {
			isError = true
		}

		id = id + "_number"
	}

	if *settings.Symbol {
		if !strings.ContainsAny(password, model.SYMBOLS) {
			isError = true
		}

		id = id + "_symbol"
	}

	if isError {
		return model.NewAppError("User.IsValid", id+".app_error", map[string]interface{}{"Min": *settings.MinimumLength}, "", http.StatusBadRequest)
	}

	return nil
}
