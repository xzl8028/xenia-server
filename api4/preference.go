// // Copyright (c) 2017-present Xenia, Inc. All Rights Reserved.
// // See License.txt for license information.

package api4

import (
	"net/http"

	"github.com/xzl8028/xenia-server/model"
)

func (api *API) InitPreference() {
	api.BaseRoutes.Preferences.Handle("", api.ApiSessionRequired(getPreferences)).Methods("GET")
	api.BaseRoutes.Preferences.Handle("", api.ApiSessionRequired(updatePreferences)).Methods("PUT")
	api.BaseRoutes.Preferences.Handle("/delete", api.ApiSessionRequired(deletePreferences)).Methods("POST")
	api.BaseRoutes.Preferences.Handle("/{category:[A-Za-z0-9_]+}", api.ApiSessionRequired(getPreferencesByCategory)).Methods("GET")
	api.BaseRoutes.Preferences.Handle("/{category:[A-Za-z0-9_]+}/name/{preference_name:[A-Za-z0-9_]+}", api.ApiSessionRequired(getPreferenceByCategoryAndName)).Methods("GET")
}

func getPreferences(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireUserId()
	if c.Err != nil {
		return
	}

	if !c.App.SessionHasPermissionToUser(c.App.Session, c.Params.UserId) {
		c.SetPermissionError(model.PERMISSION_EDIT_OTHER_USERS)
		return
	}

	preferences, err := c.App.GetPreferencesForUser(c.Params.UserId)
	if err != nil {
		c.Err = err
		return
	}

	w.Write([]byte(preferences.ToJson()))
}

func getPreferencesByCategory(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireUserId().RequireCategory()
	if c.Err != nil {
		return
	}

	if !c.App.SessionHasPermissionToUser(c.App.Session, c.Params.UserId) {
		c.SetPermissionError(model.PERMISSION_EDIT_OTHER_USERS)
		return
	}

	preferences, err := c.App.GetPreferenceByCategoryForUser(c.Params.UserId, c.Params.Category)
	if err != nil {
		c.Err = err
		return
	}

	w.Write([]byte(preferences.ToJson()))
}

func getPreferenceByCategoryAndName(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireUserId().RequireCategory().RequirePreferenceName()
	if c.Err != nil {
		return
	}

	if !c.App.SessionHasPermissionToUser(c.App.Session, c.Params.UserId) {
		c.SetPermissionError(model.PERMISSION_EDIT_OTHER_USERS)
		return
	}

	preferences, err := c.App.GetPreferenceByCategoryAndNameForUser(c.Params.UserId, c.Params.Category, c.Params.PreferenceName)
	if err != nil {
		c.Err = err
		return
	}

	w.Write([]byte(preferences.ToJson()))
}

func updatePreferences(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireUserId()
	if c.Err != nil {
		return
	}

	if !c.App.SessionHasPermissionToUser(c.App.Session, c.Params.UserId) {
		c.SetPermissionError(model.PERMISSION_EDIT_OTHER_USERS)
		return
	}

	preferences, err := model.PreferencesFromJson(r.Body)
	if err != nil {
		c.SetInvalidParam("preferences")
		return
	}

	var sanitizedPreferences model.Preferences

	for _, pref := range preferences {
		if pref.Category == model.PREFERENCE_CATEGORY_FLAGGED_POST {
			post, err := c.App.GetSinglePost(pref.Name)
			if err != nil {
				c.SetInvalidParam("preference.name")
				return
			}

			if !c.App.SessionHasPermissionToChannel(c.App.Session, post.ChannelId, model.PERMISSION_READ_CHANNEL) {
				c.SetPermissionError(model.PERMISSION_READ_CHANNEL)
				return
			}
		}

		sanitizedPreferences = append(sanitizedPreferences, pref)
	}

	if err := c.App.UpdatePreferences(c.Params.UserId, sanitizedPreferences); err != nil {
		c.Err = err
		return
	}

	ReturnStatusOK(w)
}

func deletePreferences(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireUserId()
	if c.Err != nil {
		return
	}

	if !c.App.SessionHasPermissionToUser(c.App.Session, c.Params.UserId) {
		c.SetPermissionError(model.PERMISSION_EDIT_OTHER_USERS)
		return
	}

	preferences, err := model.PreferencesFromJson(r.Body)
	if err != nil {
		c.SetInvalidParam("preferences")
		return
	}

	if err := c.App.DeletePreferences(c.Params.UserId, preferences); err != nil {
		c.Err = err
		return
	}

	ReturnStatusOK(w)
}
