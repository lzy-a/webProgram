# go web编程

[1.web基础](#1.web基础)

[2.表单](#2.表单)

[3.数据库](#3.数据库)

[4.字符串操作](#4.字符串操作)

## 1.web基础

- Go实现Web服务的工作模式的流程图:

![img](https://github.com/astaxie/build-web-application-with-golang/raw/master/zh/images/3.3.http.png?raw=true)

1. 创建Listen Socket, 监听指定的端口, 等待客户端请求到来。
2. Listen Socket接受客户端的请求, 得到Client Socket, 接下来通过Client Socket与客户端通信。
3. 处理客户端的请求, 首先从Client Socket读取HTTP请求的协议头, 如果是POST方法, 还可能要读取客户端提交的数据, 然后交给相应的handler处理请求, handler处理完毕准备好客户端需要的数据, 通过Client Socket写给客户端

```go
func ListenAndServe(addr string, handler Handler)
//初始化一个server对象，保存addr和handler信息,调用ListenAndServe方法

func (srv *Server) ListenAndServe()
//调用了net.Listen("tcp", addr)，也就是底层用TCP协议搭建了一个服务，最后调用srv.Serve监控我们设置的端口。

func (srv *Server) Serve(l net.Listener) 
//每当listener收到请求就先接收l.Accept()，再创建Conn，最后单独开了一个goroutine，把这个请求的数据当做参数扔给这个conn去服务：go c.serve(connCtx)

func (c *conn) serve(ctx context.Context)
//在conn中解析request，然后调用相应handler去处理

func (sh serverHandler) ServeHTTP(rw ResponseWriter, req *Request)
//若handler==nil，则调用DefaultServeMux。
```

整体流程：

![img](https://github.com/astaxie/build-web-application-with-golang/raw/master/zh/images/3.3.illustrator.png?raw=true)

==典型web编程代码==

```go
func sayhelloName(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello astaxie!") //这个写入到w的是输出到客户端的
}

func main() {
	http.HandleFunc("/", sayhelloName) //设置访问的路由
	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
```

通过对http包的分析之后，现在让我们来梳理一下整个的代码执行过程。

- 首先调用Http.HandleFunc

  按顺序做了几件事：

  1 调用了DefaultServeMux的HandleFunc

  2 调用了DefaultServeMux的Handle

  3 往DefaultServeMux的map[string]muxEntry中增加对应的handler和路由规则

- 其次调用http.ListenAndServe(":9090", nil)

  按顺序做了几件事情：

  1 实例化Server

  2 调用Server的ListenAndServe()

  3 调用net.Listen("tcp", addr)监听端口

  4 启动一个for循环，在循环体中Accept请求

  5 对**每个请求实例化一个Conn，并且开启一个goroutine为这个请求进行服务go c.serve()** （高并发）

  6 读取每个请求的内容w, err := c.readRequest()

  7 判断handler是否为空，如果没有设置handler（这个例子就没有设置handler），handler就设置为DefaultServeMux

  8 调用handler的ServeHttp

  9 在这个例子中，下面就进入到DefaultServeMux.ServeHttp

  10 根据request选择handler，并且进入到这个handler的ServeHTTP

  ```
    mux.handler(r).ServeHTTP(w, r)
  ```

  11 选择handler：

  A 判断是否有路由能满足这个request（循环遍历ServeMux的muxEntry）

  B 如果有路由满足，调用这个路由handler的ServeHTTP

  C 如果没有路由满足，调用NotFoundHandler的ServeHTTP

## 2.表单

==经典代码==

```go
func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		log.Println(t.Execute(w, nil))
	} else {
		//请求的是登录数据，那么执行登录的逻辑判断
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
	}
}
```

- 防止xss攻击
  - `func HTMLEscape(w io.Writer, b []byte)` //把b进行转义之后写到w
  - `func HTMLEscapeString(s string) string` //转义s之后返回结果字符串
  - `func HTMLEscaper(args ...interface{}) string `//支持多个参数一起转义，返回结果字符串

- 显示一个template

  ```go
  import "text/template"
  ...
  t, err := template.New("foo").Parse(`{{define "T"}}Hello, {{.}}!{{end}}`)
  err = t.ExecuteTemplate(out, "T", "<script>alert('you have been pwned')</script>")
  ```

- 处理文件上传

  1. 表单中增加enctype="multipart/form-data"

  2. 服务端调用`r.ParseMultipartForm`,把上传的文件存储在内存和临时文件中

  3. 使用`r.FormFile`获取文件句柄，然后对文件进行存储等处理。

     ```go
     r.ParseMultipartForm(32 << 20)
     file, handler, err := r.FormFile("uploadfile")
     //...
     f, err := os.OpenFile("./form/test"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
     ```

## 3.数据库

- mysql

  ```go
  db, err := sql.Open("mysql", "username:password@/dbname?charset=utf8")
  checkErr(err)
  
  //插入数据
  stmt, err := db.Prepare("INSERT INTO userinfo SET username=?,department=?,created=?")
  checkErr(err)
  res, err := stmt.Exec("astaxie", "研发部门", "2012-12-09")
  checkErr(err)
  id, err := res.LastInsertId()
  checkErr(err)
  ```

  - db.Prepare()函数用来返回准备要执行的sql操作，然后返回准备完毕的执行状态。

  - db.Query()函数用来直接执行Sql返回Rows结果。

  - stmt.Exec()函数用来执行stmt准备好的SQL语句

  - 我们可以看到我们传入的参数都是=?对应的数据，这样做的方式可以一定程度上防止SQL注入。

## 4.字符串操作

下面这些函数来自于strings包，这里介绍一些我平常经常用到的函数，更详细的请参考官方的文档。

- func Contains(s, substr string) bool

  字符串s中是否包含substr，返回bool值

```go
fmt.Println(strings.Contains("seafood", "foo"))
fmt.Println(strings.Contains("seafood", "bar"))
fmt.Println(strings.Contains("seafood", ""))
fmt.Println(strings.Contains("", ""))
//Output:
//true
//false
//true
//true
```

- func Join(a []string, sep string) string

  字符串链接，把slice a通过sep链接起来

```go
s := []string{"foo", "bar", "baz"}
fmt.Println(strings.Join(s, ", "))
//Output:foo, bar, baz		
```

- func Index(s, sep string) int

  在字符串s中查找sep所在的位置，返回位置值，找不到返回-1

```go
fmt.Println(strings.Index("chicken", "ken"))
fmt.Println(strings.Index("chicken", "dmr"))
//Output:4
//-1
```

- func Repeat(s string, count int) string

  重复s字符串count次，最后返回重复的字符串

```go
fmt.Println("ba" + strings.Repeat("na", 2))
//Output:banana
```

- func Replace(s, old, new string, n int) string

  在s字符串中，把old字符串替换为new字符串，n表示替换的次数，小于0表示全部替换

```go
fmt.Println(strings.Replace("oink oink oink", "k", "ky", 2))
fmt.Println(strings.Replace("oink oink oink", "oink", "moo", -1))
//Output:oinky oinky oink
//moo moo moo
```

- func Split(s, sep string) []string

  把s字符串按照sep分割，返回slice

```go
fmt.Printf("%q\n", strings.Split("a,b,c", ","))
fmt.Printf("%q\n", strings.Split("a man a plan a canal panama", "a "))
fmt.Printf("%q\n", strings.Split(" xyz ", ""))
fmt.Printf("%q\n", strings.Split("", "Bernardo O'Higgins"))
//Output:["a" "b" "c"]
//["" "man " "plan " "canal panama"]
//[" " "x" "y" "z" " "]
//[""]
```

- func Trim(s string, cutset string) string

  在s字符串的头部和尾部去除cutset指定的字符串

```go
fmt.Printf("[%q]", strings.Trim(" !!! Achtung !!! ", "! "))
//Output:["Achtung"]
```

- func Fields(s string) []string

  去除s字符串的空格符，并且按照空格分割返回slice

```go
fmt.Printf("Fields are: %q", strings.Fields("  foo bar  baz   "))
//Output:Fields are: ["foo" "bar" "baz"]
```