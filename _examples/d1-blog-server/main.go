package main

import (
	"net/http"

	"github.com/tinyredglasses/workers"
	"github.com/tinyredglasses/workers/_examples/d1-blog-server/app"
)

func main() {
	http.Handle("/articles", app.NewArticleHandler())
	workers.Serve(nil) // use http.DefaultServeMux
}
