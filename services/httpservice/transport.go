// Copyright (c) 2017-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package httpservice

import (
	"net/http"
)

// XeniaTransport is an implementation of http.RoundTripper that ensures each request contains a custom user agent
// string to indicate that the request is coming from a Xenia instance.
type XeniaTransport struct {
	// Transport is the underlying http.RoundTripper that is actually used to make the request
	Transport http.RoundTripper
}

func (t *XeniaTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", defaultUserAgent)

	return t.Transport.RoundTrip(req)
}
