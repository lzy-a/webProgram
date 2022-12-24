package simpleweb

import (
	"fmt"
	"net/http"
	"strings"
)

func SayHello(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hello")
}

func SayHelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form)
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello Louis!")
}
