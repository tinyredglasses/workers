package cloudflare

import (
	"context"
	"github.com/tinyredglasses/workers/internal/jsutil"
	"github.com/tinyredglasses/workers/internal/runtimecontext"
)

func CreateContext() context.Context {
	runtimeCtxObj := jsutil.RuntimeContext

	return runtimecontext.New(context.Background(), runtimeCtxObj)
}
