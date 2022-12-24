package form

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"
)

func Login(w http.ResponseWriter, r *http.Request) {

	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		timestamp := strconv.Itoa(time.Now().Nanosecond())
		hashWr := md5.New()
		hashWr.Write([]byte(timestamp))
		token := fmt.Sprintf("%x", hashWr.Sum(nil))
		fmt.Println("token: ", hashWr.Sum(nil))
		t, _ := template.ParseFiles("./form/login.gtpl")
		log.Println(t.Execute(w, token))
	} else {
		r.ParseForm()
		token := r.Form.Get("token")
		if token == "" {
			fmt.Println("w.o. token")
			return
		}
		fmt.Println("username:", template.HTMLEscapeString(r.Form.Get("username")))
		fmt.Println("password:", template.HTMLEscapeString(r.Form.Get("password")))
		template.HTMLEscape(w, []byte(r.Form.Get("username")))
		// err := PrintAlert(w, r.Form.Get("username"))
		// if err != nil {
		// 	log.Fatal("print alert:", err)
		// }
	}
}

func Upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		curtime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(curtime, 10))
		toke := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("./form/upload.gtpl")
		t.Execute(w, toke)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("./form/test"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		io.Copy(f, file)
	}
}

func TimeandDate() {
	t := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	fmt.Printf("Go launched at %s\n", t.Local())
}

func PrintAlert(w http.ResponseWriter, str string) error {
	t, _ := template.New("foo").Parse(`{{define "T"}}hello: {{.}}!{{end}}`)
	err := t.ExecuteTemplate(w, "T", str)
	return err
}
