package cloudflare

import (
	"context"
	"github.com/tinyredglasses/workers/jsutil"
	"github.com/tinyredglasses/workers/runtimecontext"
)

func CreateContext() context.Context {
	runtimeCtxObj := jsutil.RuntimeContext
	return runtimecontext.New(context.Background(), runtimeCtxObj)
}
