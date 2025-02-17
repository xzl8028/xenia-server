package api4

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/xzl8028/xenia-server/model"
	"github.com/stretchr/testify/assert"
)

const (
	acAllowOrigin      = "Access-Control-Allow-Origin"
	acExposeHeaders    = "Access-Control-Expose-Headers"
	acMaxAge           = "Access-Control-Max-Age"
	acAllowCredentials = "Access-Control-Allow-Credentials"
	acAllowMethods     = "Access-Control-Allow-Methods"
	acAllowHeaders     = "Access-Control-Allow-Headers"
)

func TestCORSRequestHandling(t *testing.T) {
	for name, testcase := range map[string]struct {
		AllowCorsFrom            string
		CorsExposedHeaders       string
		CorsAllowCredentials     bool
		ModifyRequest            func(req *http.Request)
		ExpectedAllowOrigin      string
		ExpectedExposeHeaders    string
		ExpectedAllowCredentials string
	}{
		"NoCORS": {
			"",
			"",
			false,
			func(req *http.Request) {
			},
			"",
			"",
			"",
		},
		"CORSEnabled": {
			"http://somewhere.com",
			"",
			false,
			func(req *http.Request) {
			},
			"",
			"",
			"",
		},
		"CORSEnabledStarOrigin": {
			"*",
			"",
			false,
			func(req *http.Request) {
				req.Header.Set("Origin", "http://pre-release.xenia.com")
			},
			"*",
			"",
			"",
		},
		"CORSEnabledStarNoOrigin": { // CORS spec requires this, not a bug.
			"*",
			"",
			false,
			func(req *http.Request) {
			},
			"",
			"",
			"",
		},
		"CORSEnabledMatching": {
			"http://xenia.com",
			"",
			false,
			func(req *http.Request) {
				req.Header.Set("Origin", "http://xenia.com")
			},
			"http://xenia.com",
			"",
			"",
		},
		"CORSEnabledMultiple": {
			"http://spinmint.com http://xenia.com",
			"",
			false,
			func(req *http.Request) {
				req.Header.Set("Origin", "http://xenia.com")
			},
			"http://xenia.com",
			"",
			"",
		},
		"CORSEnabledWithCredentials": {
			"http://xenia.com",
			"",
			true,
			func(req *http.Request) {
				req.Header.Set("Origin", "http://xenia.com")
			},
			"http://xenia.com",
			"",
			"true",
		},
		"CORSEnabledWithHeaders": {
			"http://xenia.com",
			"x-my-special-header x-blueberry",
			true,
			func(req *http.Request) {
				req.Header.Set("Origin", "http://xenia.com")
			},
			"http://xenia.com",
			"X-My-Special-Header, X-Blueberry",
			"true",
		},
	} {
		t.Run(name, func(t *testing.T) {
			th := SetupConfig(func(cfg *model.Config) {
				*cfg.ServiceSettings.AllowCorsFrom = testcase.AllowCorsFrom
				*cfg.ServiceSettings.CorsExposedHeaders = testcase.CorsExposedHeaders
				*cfg.ServiceSettings.CorsAllowCredentials = testcase.CorsAllowCredentials
			})
			defer th.TearDown()

			port := th.App.Srv.ListenAddr.Port
			host := fmt.Sprintf("http://localhost:%v", port)
			url := fmt.Sprintf("%v/api/v4/system/ping", host)

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				t.Fatal(err)
			}
			testcase.ModifyRequest(req)

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.Equal(t, testcase.ExpectedAllowOrigin, resp.Header.Get(acAllowOrigin))
			assert.Equal(t, testcase.ExpectedExposeHeaders, resp.Header.Get(acExposeHeaders))
			assert.Equal(t, "", resp.Header.Get(acMaxAge))
			assert.Equal(t, testcase.ExpectedAllowCredentials, resp.Header.Get(acAllowCredentials))
			assert.Equal(t, "", resp.Header.Get(acAllowMethods))
			assert.Equal(t, "", resp.Header.Get(acAllowHeaders))
		})
	}

}
