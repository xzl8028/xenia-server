// Copyright (c) 2015-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package app

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
	"strconv"

	"github.com/xzl8028/xenia-server/mlog"
	"github.com/xzl8028/xenia-server/model"
	"github.com/xzl8028/xenia-server/services/mailservice"
	"github.com/xzl8028/xenia-server/utils"
)

const (
	SECURITY_URL           = "https://securityupdatecheck.xenia.com"
	SECURITY_UPDATE_PERIOD = 86400000 // 24 hours in milliseconds.

	PROP_SECURITY_ID                = "id"
	PROP_SECURITY_BUILD             = "b"
	PROP_SECURITY_ENTERPRISE_READY  = "be"
	PROP_SECURITY_DATABASE          = "db"
	PROP_SECURITY_OS                = "os"
	PROP_SECURITY_USER_COUNT        = "uc"
	PROP_SECURITY_TEAM_COUNT        = "tc"
	PROP_SECURITY_ACTIVE_USER_COUNT = "auc"
	PROP_SECURITY_UNIT_TESTS        = "ut"
)

func (s *Server) DoSecurityUpdateCheck() {
	if !*s.Config().ServiceSettings.EnableSecurityFixAlert {
		return
	}

	props, err := s.Store.System().Get()
	if err != nil {
		return
	}

	lastSecurityTime, _ := strconv.ParseInt(props[model.SYSTEM_LAST_SECURITY_TIME], 10, 0)
	currentTime := model.GetMillis()

	if (currentTime - lastSecurityTime) > SECURITY_UPDATE_PERIOD {
		mlog.Debug("Checking for security update from Xenia")

		v := url.Values{}

		v.Set(PROP_SECURITY_ID, s.diagnosticId)
		v.Set(PROP_SECURITY_BUILD, model.CurrentVersion+"."+model.BuildNumber)
		v.Set(PROP_SECURITY_ENTERPRISE_READY, model.BuildEnterpriseReady)
		v.Set(PROP_SECURITY_DATABASE, *s.Config().SqlSettings.DriverName)
		v.Set(PROP_SECURITY_OS, runtime.GOOS)

		if len(props[model.SYSTEM_RAN_UNIT_TESTS]) > 0 {
			v.Set(PROP_SECURITY_UNIT_TESTS, "1")
		} else {
			v.Set(PROP_SECURITY_UNIT_TESTS, "0")
		}

		systemSecurityLastTime := &model.System{Name: model.SYSTEM_LAST_SECURITY_TIME, Value: strconv.FormatInt(currentTime, 10)}
		if lastSecurityTime == 0 {
			s.Store.System().Save(systemSecurityLastTime)
		} else {
			s.Store.System().Update(systemSecurityLastTime)
		}

		if count, err := s.Store.User().Count(model.UserCountOptions{IncludeDeleted: true}); err == nil {
			v.Set(PROP_SECURITY_USER_COUNT, strconv.FormatInt(count, 10))
		}

		if ucr, err := s.Store.Status().GetTotalActiveUsersCount(); err == nil {
			v.Set(PROP_SECURITY_ACTIVE_USER_COUNT, strconv.FormatInt(ucr, 10))
		}

		if teamCount, err := s.Store.Team().AnalyticsTeamCount(); err == nil {
			v.Set(PROP_SECURITY_TEAM_COUNT, strconv.FormatInt(teamCount, 10))
		}

		res, err := http.Get(SECURITY_URL + "/security?" + v.Encode())
		if err != nil {
			mlog.Error("Failed to get security update information from Xenia.")
			return
		}

		defer res.Body.Close()

		bulletins := model.SecurityBulletinsFromJson(res.Body)

		for _, bulletin := range bulletins {
			if bulletin.AppliesToVersion == model.CurrentVersion {
				if props["SecurityBulletin_"+bulletin.Id] == "" {
					results := <-s.Store.User().GetSystemAdminProfiles()
					if results.Err != nil {
						mlog.Error("Failed to get system admins for security update information from Xenia.")
						return
					}
					users := results.Data.(map[string]*model.User)

					resBody, err := http.Get(SECURITY_URL + "/bulletins/" + bulletin.Id)
					if err != nil {
						mlog.Error("Failed to get security bulletin details")
						return
					}

					body, err := ioutil.ReadAll(resBody.Body)
					res.Body.Close()
					if err != nil || resBody.StatusCode != 200 {
						mlog.Error("Failed to read security bulletin details")
						return
					}

					for _, user := range users {
						mlog.Info(fmt.Sprintf("Sending security bulletin for %v to %v", bulletin.Id, user.Email))
						license := s.License()
						mailservice.SendMailUsingConfig(user.Email, utils.T("xenia.bulletin.subject"), string(body), s.Config(), license != nil && *license.Features.Compliance)
					}

					bulletinSeen := &model.System{Name: "SecurityBulletin_" + bulletin.Id, Value: bulletin.Id}
					s.Store.System().Save(bulletinSeen)
				}
			}
		}
	}
}
