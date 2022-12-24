package simplemux

import (
	"net/http"
	"web/simpleweb"
)

type MyMux struct {
}

func (p *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		simpleweb.SayHelloName(w, r)
		return
	}
	http.NotFound(w, r)
}
