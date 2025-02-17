// Copyright (c) 2017-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package imageproxy

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/xzl8028/xenia-server/mlog"
	"github.com/xzl8028/xenia-server/model"
	"github.com/xzl8028/xenia-server/services/configservice"
	"github.com/xzl8028/xenia-server/services/httpservice"
)

var ErrNotEnabled = Error{errors.New("imageproxy.ImageProxy: image proxy not enabled")}

// An ImageProxy is the public interface for Xenia's image proxy. An instance of ImageProxy should be created
// using MakeImageProxy which requires a configService and an HTTPService provided by the server.
type ImageProxy struct {
	ConfigService    configservice.ConfigService
	configListenerId string

	HTTPService httpservice.HTTPService

	Logger *mlog.Logger

	lock    sync.RWMutex
	backend ImageProxyBackend
}

// An ImageProxyBackend provides the functionality for different types of image proxies. An ImageProxy will construct
// the required backend depending on the ImageProxySettings provided by the ConfigService.
type ImageProxyBackend interface {
	// GetImage provides a proxied image in response to an HTTP request.
	GetImage(w http.ResponseWriter, r *http.Request, imageURL string)

	// GetImageDirect returns a proxied image along with its content type.
	GetImageDirect(imageURL string) (io.ReadCloser, string, error)
}

func MakeImageProxy(configService configservice.ConfigService, httpService httpservice.HTTPService, logger *mlog.Logger) *ImageProxy {
	proxy := &ImageProxy{
		ConfigService: configService,
		HTTPService:   httpService,
		Logger:        logger,
	}

	proxy.configListenerId = proxy.ConfigService.AddConfigListener(proxy.OnConfigChange)

	config := proxy.ConfigService.Config()
	proxy.backend = proxy.makeBackend(*config.ImageProxySettings.Enable, *config.ImageProxySettings.ImageProxyType)

	return proxy
}

func (proxy *ImageProxy) makeBackend(enable bool, proxyType string) ImageProxyBackend {
	if !enable {
		return nil
	}

	switch proxyType {
	case model.IMAGE_PROXY_TYPE_LOCAL:
		return makeLocalBackend(proxy)
	case model.IMAGE_PROXY_TYPE_ATMOS_CAMO:
		return makeAtmosCamoBackend(proxy)
	default:
		return nil
	}
}

func (proxy *ImageProxy) Close() {
	proxy.lock.Lock()
	defer proxy.lock.Unlock()

	proxy.ConfigService.RemoveConfigListener(proxy.configListenerId)
}

func (proxy *ImageProxy) OnConfigChange(oldConfig, newConfig *model.Config) {
	if *oldConfig.ImageProxySettings.Enable != *newConfig.ImageProxySettings.Enable ||
		*oldConfig.ImageProxySettings.ImageProxyType != *newConfig.ImageProxySettings.ImageProxyType {
		proxy.lock.Lock()
		defer proxy.lock.Unlock()

		proxy.backend = proxy.makeBackend(*newConfig.ImageProxySettings.Enable, *newConfig.ImageProxySettings.ImageProxyType)
	}
}

// GetImage takes an HTTP request for an image and requests that image using the image proxy.
func (proxy *ImageProxy) GetImage(w http.ResponseWriter, r *http.Request, imageURL string) {
	proxy.lock.RLock()
	defer proxy.lock.RUnlock()

	if proxy.backend == nil {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	proxy.backend.GetImage(w, r, imageURL)
}

// GetImageDirect takes the URL of an image and returns the image along with its content type.
func (proxy *ImageProxy) GetImageDirect(imageURL string) (io.ReadCloser, string, error) {
	proxy.lock.RLock()
	defer proxy.lock.RUnlock()

	if proxy.backend == nil {
		return nil, "", ErrNotEnabled
	}

	return proxy.backend.GetImageDirect(imageURL)
}

// GetProxiedImageURL takes the URL of an image and returns a URL that can be used to view that image through the
// image proxy.
func (proxy *ImageProxy) GetProxiedImageURL(imageURL string) string {
	return getProxiedImageURL(imageURL, *proxy.ConfigService.Config().ServiceSettings.SiteURL)
}

func getProxiedImageURL(imageURL, siteURL string) string {
	if imageURL == "" || imageURL[0] == '/' || strings.HasPrefix(imageURL, siteURL) {
		return imageURL
	}

	return siteURL + "/api/v4/image?url=" + url.QueryEscape(imageURL)
}

// GetUnproxiedImageURL takes the URL of an image on the image proxy and returns the original URL of the image.
func (proxy *ImageProxy) GetUnproxiedImageURL(proxiedURL string) string {
	return getUnproxiedImageURL(proxiedURL, *proxy.ConfigService.Config().ServiceSettings.SiteURL)
}

func getUnproxiedImageURL(proxiedURL, siteURL string) string {
	if !strings.HasPrefix(proxiedURL, siteURL+"/api/v4/image?url=") {
		return proxiedURL
	}

	parsed, err := url.Parse(proxiedURL)
	if err != nil {
		return proxiedURL
	}

	u := parsed.Query()["url"]
	if len(u) == 0 {
		return proxiedURL
	}

	return u[0]
}
