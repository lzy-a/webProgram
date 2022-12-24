package form

import (
	"log"
	"net/http"
	"web/simpleweb"
)

func Main() {
	http.HandleFunc("/upload", Upload)
	http.HandleFunc("/", simpleweb.SayHelloName)
	http.HandleFunc("/login", Login)
	// fmt.Println("111")
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("listen and serve:", err)
	}
}
