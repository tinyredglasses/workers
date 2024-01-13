package workers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"syscall/js"

	"github.com/tinyredglasses/workers/internal/jshttp"
	"github.com/tinyredglasses/workers/internal/jsutil"
	"github.com/tinyredglasses/workers/internal/runtimecontext"
)

var httpHandler http.Handler

func init() {
	var handleRequestCallback js.Func
	handleRequestCallback = js.FuncOf(func(this js.Value, args []js.Value) any {
		fmt.Println("inner1")
		if len(args) > 1 {
			panic(fmt.Errorf("too many args given to handleRequest: %d", len(args)))
		}
		reqObj := args[0]
		fmt.Println(reqObj)
		runtimeCtxObj := jsutil.RuntimeContext
		var cb js.Func
		cb = js.FuncOf(func(_ js.Value, pArgs []js.Value) any {
			fmt.Println("inner2")
			defer cb.Release()
			resolve := pArgs[0]
			go func() {
				fmt.Println("inner3")
				res, err := handleRequest(reqObj, runtimeCtxObj)
				if err != nil {
					panic(err)
				}
				resolve.Invoke(res)
			}()
			return js.Undefined()
		})
		return jsutil.NewPromise(cb)
	})
	jsutil.Binding.Set("handleRequest", handleRequestCallback)
}

// handleRequest accepts a Request object and returns Response object.
func handleRequest(reqObj js.Value, runtimeCtxObj js.Value) (js.Value, error) {
	if httpHandler == nil {
		return js.Value{}, fmt.Errorf("Serve must be called before handleRequest.")
	}
	req, err := jshttp.ToRequest(reqObj)
	if err != nil {
		panic(err)
	}
	ctx := runtimecontext.New(context.Background(), runtimeCtxObj)
	req = req.WithContext(ctx)
	reader, writer := io.Pipe()
	w := &jshttp.ResponseWriter{
		HeaderValue: http.Header{},
		StatusCode:  http.StatusOK,
		Reader:      reader,
		Writer:      writer,
		ReadyCh:     make(chan struct{}),
	}
	go func() {
		defer w.Ready()
		defer writer.Close()
		httpHandler.ServeHTTP(w, req)
	}()
	<-w.ReadyCh
	return w.ToJSResponse(), nil
}

//go:wasmimport workers ready
func ready()

// Server serves http.Handler on Cloudflare Workers.
// if the given handler is nil, http.DefaultServeMux will be used.
func Serve(handler http.Handler) {
	if handler == nil {
		handler = http.DefaultServeMux
	}
	httpHandler = handler
	ready()
	select {}
}

func Init() {

}
