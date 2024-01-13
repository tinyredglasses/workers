package fetch

import (
	"errors"
	jshttp2 "github.com/tinyredglasses/workers/jshttp"
	"github.com/tinyredglasses/workers/jsutil"
	"net/http"
	"syscall/js"
)

// fetch is a function that reproduces cloudflare fetch.
// Docs: https://developers.cloudflare.com/workers/runtime-apis/fetch/
func fetch(namespace js.Value, req *http.Request, init *RequestInit) (*http.Response, error) {
	if namespace.IsUndefined() {
		return nil, errors.New("fetch function not found")
	}
	promise := namespace.Call("fetch",
		// The Request object to fetch.
		// Docs: https://developers.cloudflare.com/workers/runtime-apis/request
		jshttp2.ToJSRequest(req),
		// The content of the request.
		// Docs: https://developers.cloudflare.com/workers/runtime-apis/request#requestinit
		init.ToJS(),
	)
	jsRes, err := jsutil.AwaitPromise(promise)
	if err != nil {
		return nil, err
	}
	return jshttp2.ToResponse(jsRes)
}
