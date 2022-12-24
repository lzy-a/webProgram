package simpleweb

import (
	"log"
	"net/http"
)

func Main() {
	//使用内置Default Mux
	http.HandleFunc("/", SayHelloName)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
