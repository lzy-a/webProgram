package simplemux

import (
	"net/http"
)

func Main() {
	//使用内置Default Mux
	// http.HandleFunc("/", simpleweb.SayHelloName)
	// err := http.ListenAndServe(":9090", nil)
	// if err != nil {
	// 	log.Fatal("ListenAndServe: ", err)
	// }

	//使用MyMux
	mux := &MyMux{}
	http.ListenAndServe(":9090", mux)
}
