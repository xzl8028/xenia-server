// Copyright (c) 2017-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package imageproxy

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/xzl8028/xenia-server/model"
	"github.com/xzl8028/xenia-server/services/httpservice"
	"github.com/xzl8028/xenia-server/utils/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeTestLocalProxy() *ImageProxy {
	configService := &testutils.StaticConfigService{
		Cfg: &model.Config{
			ServiceSettings: model.ServiceSettings{
				SiteURL:                             model.NewString("https://xenia.example.com"),
				AllowedUntrustedInternalConnections: model.NewString("127.0.0.1"),
			},
			ImageProxySettings: model.ImageProxySettings{
				Enable:         model.NewBool(true),
				ImageProxyType: model.NewString(model.IMAGE_PROXY_TYPE_LOCAL),
			},
		},
	}

	return MakeImageProxy(configService, httpservice.MakeHTTPService(configService), nil)
}

func TestLocalBackend_GetImage(t *testing.T) {
	t.Run("image", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "max-age=2592000, private")
			w.Header().Set("Content-Type", "image/png")
			w.Header().Set("Content-Length", "10")

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("1111111111"))
		})

		mock := httptest.NewServer(handler)
		defer mock.Close()

		proxy := makeTestLocalProxy()

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodGet, "", nil)
		proxy.GetImage(recorder, request, mock.URL+"/image.png")
		resp := recorder.Result()

		require.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "max-age=2592000, private", resp.Header.Get("Cache-Control"))
		assert.Equal(t, "10", resp.Header.Get("Content-Length"))

		respBody, _ := ioutil.ReadAll(resp.Body)
		assert.Equal(t, []byte("1111111111"), respBody)
	})

	t.Run("not an image", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotAcceptable)
		})

		mock := httptest.NewServer(handler)
		defer mock.Close()

		proxy := makeTestLocalProxy()

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodGet, "", nil)
		proxy.GetImage(recorder, request, mock.URL+"/file.pdf")
		resp := recorder.Result()

		require.Equal(t, http.StatusNotAcceptable, resp.StatusCode)
	})

	t.Run("not an image, but remote server ignores accept header", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "max-age=2592000, private")
			w.Header().Set("Content-Type", "application/pdf")
			w.Header().Set("Content-Length", "10")

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("1111111111"))
		})

		mock := httptest.NewServer(handler)
		defer mock.Close()

		proxy := makeTestLocalProxy()

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodGet, "", nil)
		proxy.GetImage(recorder, request, mock.URL+"/file.pdf")
		resp := recorder.Result()

		require.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("not found", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})

		mock := httptest.NewServer(handler)
		defer mock.Close()

		proxy := makeTestLocalProxy()

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodGet, "", nil)
		proxy.GetImage(recorder, request, mock.URL+"/image.png")
		resp := recorder.Result()

		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("other server error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})

		mock := httptest.NewServer(handler)
		defer mock.Close()

		proxy := makeTestLocalProxy()

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodGet, "", nil)
		proxy.GetImage(recorder, request, mock.URL+"/image.png")
		resp := recorder.Result()

		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("timeout", func(t *testing.T) {
		wait := make(chan bool, 1)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			<-wait
		})

		mock := httptest.NewServer(handler)
		defer mock.Close()

		proxy := makeTestLocalProxy()

		// Modify the timeout to be much shorter than the default 30 seconds
		proxy.backend.(*LocalBackend).impl.Timeout = time.Millisecond

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodGet, "", nil)
		proxy.GetImage(recorder, request, mock.URL+"/image.png")
		resp := recorder.Result()

		require.Equal(t, http.StatusGatewayTimeout, resp.StatusCode)

		wait <- true
	})
}

func TestLocalBackend_GetImageDirect(t *testing.T) {
	t.Run("image", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "max-age=2592000, private")
			w.Header().Set("Content-Type", "image/png")
			w.Header().Set("Content-Length", "10")

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("1111111111"))
		})

		mock := httptest.NewServer(handler)
		defer mock.Close()

		proxy := makeTestLocalProxy()

		body, contentType, err := proxy.GetImageDirect(mock.URL + "/image.png")

		assert.Nil(t, err)
		assert.Equal(t, "image/png", contentType)

		respBody, _ := ioutil.ReadAll(body)
		assert.Equal(t, []byte("1111111111"), respBody)
	})

	t.Run("not an image", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotAcceptable)
		})

		mock := httptest.NewServer(handler)
		defer mock.Close()

		proxy := makeTestLocalProxy()

		body, contentType, err := proxy.GetImageDirect(mock.URL + "/file.pdf")

		assert.NotNil(t, err)
		assert.Equal(t, "", contentType)
		assert.Equal(t, ErrLocalRequestFailed, err)
		assert.Nil(t, body)
	})

	t.Run("not an image, but remote server ignores accept header", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "max-age=2592000, private")
			w.Header().Set("Content-Type", "application/pdf")
			w.Header().Set("Content-Length", "10")

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("1111111111"))
		})

		mock := httptest.NewServer(handler)
		defer mock.Close()

		proxy := makeTestLocalProxy()

		body, contentType, err := proxy.GetImageDirect(mock.URL + "/file.pdf")

		assert.NotNil(t, err)
		assert.Equal(t, "", contentType)
		assert.Equal(t, ErrLocalRequestFailed, err)
		assert.Nil(t, body)
	})

	t.Run("not found", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})

		mock := httptest.NewServer(handler)
		defer mock.Close()

		proxy := makeTestLocalProxy()

		body, contentType, err := proxy.GetImageDirect(mock.URL + "/image.png")

		assert.NotNil(t, err)
		assert.Equal(t, "", contentType)
		assert.Equal(t, ErrLocalRequestFailed, err)
		assert.Nil(t, body)
	})

	t.Run("other server error", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})

		mock := httptest.NewServer(handler)
		defer mock.Close()

		proxy := makeTestLocalProxy()

		body, contentType, err := proxy.GetImageDirect(mock.URL + "/image.png")

		assert.NotNil(t, err)
		assert.Equal(t, "", contentType)
		assert.Equal(t, ErrLocalRequestFailed, err)
		assert.Nil(t, body)
	})

	t.Run("timeout", func(t *testing.T) {
		wait := make(chan bool, 1)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			<-wait
		})

		mock := httptest.NewServer(handler)
		defer mock.Close()

		proxy := makeTestLocalProxy()

		// Modify the timeout to be much shorter than the default 30 seconds
		proxy.backend.(*LocalBackend).impl.Timeout = time.Millisecond

		body, contentType, err := proxy.GetImageDirect(mock.URL + "/image.png")

		assert.NotNil(t, err)
		assert.Equal(t, "", contentType)
		assert.Equal(t, ErrLocalRequestFailed, err)
		assert.Nil(t, body)

		wait <- true
	})
}
