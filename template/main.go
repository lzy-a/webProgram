package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
)

type Person struct {
	UserName string
	Emails   []string
	Friends  []*Friend
}
type Friend struct {
	Fname string
}

func main() {
	//loop()
	// ifelse()
	// variable()
	// function()
	defineTemplate()
}

func loop() {
	f1 := Friend{Fname: "lzy"}
	f2 := Friend{Fname: "zyl"}
	t := template.New("example")
	t, _ = t.Parse(`hello {{.UserName}}!
	{{range .Emails}}
	an email {{.}}
	{{end}}
	{{with .Friends}}
	{{range .}}
	my friend name is {{.Fname}}
	{{end}}
	{{end}}
	`)
	p := Person{UserName: "Louis", Emails: []string{"lzy@qq.com", "lzy@gmail.com"}, Friends: []*Friend{&f1, &f2}}
	t.Execute(os.Stdout, p)
}

func ifelse() {
	tEmpty := template.New("template test")
	tEmpty = template.Must(tEmpty.Parse("空 pipeline if demo: {{if ``}} 不会输出. {{end}}\n"))
	tEmpty.Execute(os.Stdout, nil)

	tWithValue := template.New("template test")
	tWithValue = template.Must(tWithValue.Parse("不为空的 pipeline if demo: {{if `anything`}} 我有内容，我会输出. {{end}}\n"))
	tWithValue.Execute(os.Stdout, nil)

	tIfElse := template.New("template test")
	tIfElse = template.Must(tIfElse.Parse("if-else demo: {{if .Friends}} if部分 {{else}} else部分.{{end}}\n"))
	tIfElse.Execute(os.Stdout, Person{UserName: "louis"})
}

func variable() {
	f1 := Friend{Fname: "lzy"}
	f2 := Friend{Fname: "zyl"}
	t := template.New("example")
	t, _ = t.Parse(`{{with $x := "output" | printf "%q"}}{{$x}}{{end}}
{{with $x := "output"}}{{printf "%q" $x}}{{end}}
{{with $x := "output"}}{{$x | printf "%q"}}{{end}}`)
	p := Person{UserName: "Louis", Emails: []string{"lzy@qq.com", "lzy@gmail.com"}, Friends: []*Friend{&f1, &f2}}
	t.Execute(os.Stdout, p)
}

func function() {
	f1 := Friend{Fname: "lzy"}
	f2 := Friend{Fname: "zyl"}
	t := template.New("example")
	t = t.Funcs(template.FuncMap{"emailDeal": EmailDealWith})
	t, _ = t.Parse(`hello {{.UserName}}!
	{{range .Emails}}
	an email {{.|emailDeal}}
	{{end}}
	{{with .Friends}}
	{{range .}}
	my friend name is {{.Fname}}
	{{end}}
	{{end}}
	`)
	p := Person{UserName: "Louis", Emails: []string{"lzy@qq.com", "lzy@gmail.com"}, Friends: []*Friend{&f1, &f2}}
	t.Execute(os.Stdout, p)
}

func EmailDealWith(args ...interface{}) string {
	ok := false
	var s string
	if len(args) == 1 {
		s, ok = args[0].(string)
	}
	if !ok {
		s = fmt.Sprint(args...)
	}
	// find the @ symbol
	substrs := strings.Split(s, "@")
	if len(substrs) != 2 {
		return s
	}
	// replace the @ by " at "
	return (substrs[0] + " at " + substrs[1])

}

func defineTemplate() {
	s1, err := template.ParseFiles("./template/header.tmpl", "./template/content.tmpl", "./template/footer.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	s1.ExecuteTemplate(os.Stdout, "header", nil)
	fmt.Println()
	s1.ExecuteTemplate(os.Stdout, "content", nil)
	fmt.Println()
	s1.ExecuteTemplate(os.Stdout, "footer", nil)
	fmt.Println()
	s1.Execute(os.Stdout, nil)
}
