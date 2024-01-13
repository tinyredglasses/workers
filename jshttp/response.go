package jshttp

import (
	jsutil2 "github.com/tinyredglasses/workers/jsutil"
	"io"
	"net/http"
	"strconv"
	"syscall/js"
)

// ToResponse converts JavaScript sides Response to *http.Response.
//   - Response: https://developer.mozilla.org/docs/Web/API/Response
func ToResponse(res js.Value) (*http.Response, error) {
	status := res.Get("status").Int()
	promise := res.Call("blob")
	blob, err := jsutil2.AwaitPromise(promise)
	if err != nil {
		return nil, err
	}
	header := ToHeader(res.Get("headers"))
	contentLength, _ := strconv.ParseInt(header.Get("Content-Length"), 10, 64)

	return &http.Response{
		Status:        strconv.Itoa(status) + " " + res.Get("statusText").String(),
		StatusCode:    status,
		Header:        header,
		Body:          io.NopCloser(jsutil2.ConvertStreamReaderToReader(blob.Call("stream").Call("getReader"))),
		ContentLength: contentLength,
	}, nil
}

// ToJSResponse converts *http.Response to JavaScript sides Response class object.
func ToJSResponse(res *http.Response) js.Value {
	return newJSResponse(res.StatusCode, res.Header, res.Body)
}

// newJSResponse creates JavaScript sides Response class object.
//   - Response: https://developer.mozilla.org/docs/Web/API/Response
func newJSResponse(statusCode int, headers http.Header, body io.ReadCloser) js.Value {
	status := statusCode
	if status == 0 {
		status = http.StatusOK
	}
	respInit := jsutil2.NewObject()
	respInit.Set("status", status)
	respInit.Set("statusText", http.StatusText(status))
	respInit.Set("headers", ToJSHeader(headers))
	if status == http.StatusSwitchingProtocols ||
		status == http.StatusNoContent ||
		status == http.StatusResetContent ||
		status == http.StatusNotModified {
		return jsutil2.ResponseClass.New(jsutil2.Null, respInit)
	}
	readableStream := jsutil2.ConvertReaderToReadableStream(body)
	return jsutil2.ResponseClass.New(readableStream, respInit)
}
