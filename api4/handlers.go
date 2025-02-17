// Copyright (c) 2017-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package api4

import (
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/xzl8028/xenia-server/web"
)

type Context = web.Context

// ApiHandler provides a handler for API endpoints which do not require the user to be logged in order for access to be
// granted.
func (api *API) ApiHandler(h func(*Context, http.ResponseWriter, *http.Request)) http.Handler {
	handler := &web.Handler{
		GetGlobalAppOptions: api.GetGlobalAppOptions,
		HandleFunc:          h,
		RequireSession:      false,
		TrustRequester:      false,
		RequireMfa:          false,
		IsStatic:            false,
	}
	if *api.ConfigService.Config().ServiceSettings.WebserverMode == "gzip" {
		return gziphandler.GzipHandler(handler)
	}
	return handler
}

// ApiSessionRequired provides a handler for API endpoints which require the user to be logged in in order for access to
// be granted.
func (api *API) ApiSessionRequired(h func(*Context, http.ResponseWriter, *http.Request)) http.Handler {
	handler := &web.Handler{
		GetGlobalAppOptions: api.GetGlobalAppOptions,
		HandleFunc:          h,
		RequireSession:      true,
		TrustRequester:      false,
		RequireMfa:          true,
		IsStatic:            false,
	}
	if *api.ConfigService.Config().ServiceSettings.WebserverMode == "gzip" {
		return gziphandler.GzipHandler(handler)
	}
	return handler

}

// ApiSessionRequiredMfa provides a handler for API endpoints which require a logged-in user session  but when accessed,
// if MFA is enabled, the MFA process is not yet complete, and therefore the requirement to have completed the MFA
// authentication must be waived.
func (api *API) ApiSessionRequiredMfa(h func(*Context, http.ResponseWriter, *http.Request)) http.Handler {
	handler := &web.Handler{
		GetGlobalAppOptions: api.GetGlobalAppOptions,
		HandleFunc:          h,
		RequireSession:      true,
		TrustRequester:      false,
		RequireMfa:          false,
		IsStatic:            false,
	}
	if *api.ConfigService.Config().ServiceSettings.WebserverMode == "gzip" {
		return gziphandler.GzipHandler(handler)
	}
	return handler

}

// ApiHandlerTrustRequester provides a handler for API endpoints which do not require the user to be logged in and are
// allowed to be requested directly rather than via javascript/XMLHttpRequest, such as site branding images or the
// websocket.
func (api *API) ApiHandlerTrustRequester(h func(*Context, http.ResponseWriter, *http.Request)) http.Handler {
	handler := &web.Handler{
		GetGlobalAppOptions: api.GetGlobalAppOptions,
		HandleFunc:          h,
		RequireSession:      false,
		TrustRequester:      true,
		RequireMfa:          false,
		IsStatic:            false,
	}
	if *api.ConfigService.Config().ServiceSettings.WebserverMode == "gzip" {
		return gziphandler.GzipHandler(handler)
	}
	return handler

}

// ApiSessionRequiredTrustRequester provides a handler for API endpoints which do require the user to be logged in and
// are allowed to be requested directly rather than via javascript/XMLHttpRequest, such as emoji or file uploads.
func (api *API) ApiSessionRequiredTrustRequester(h func(*Context, http.ResponseWriter, *http.Request)) http.Handler {
	handler := &web.Handler{
		GetGlobalAppOptions: api.GetGlobalAppOptions,
		HandleFunc:          h,
		RequireSession:      true,
		TrustRequester:      true,
		RequireMfa:          true,
		IsStatic:            false,
	}
	if *api.ConfigService.Config().ServiceSettings.WebserverMode == "gzip" {
		return gziphandler.GzipHandler(handler)
	}
	return handler

}
