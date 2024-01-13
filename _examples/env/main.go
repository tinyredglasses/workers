package main

import (
	"fmt"
	"net/http"

	"github.com/tinyredglasses/workers"
	"github.com/tinyredglasses/workers/cloudflare"
)

func main() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "MY_ENV: %s", cloudflare.Getenv(req.Context(), "MY_ENV"))
	})
	workers.Serve(handler)
}
