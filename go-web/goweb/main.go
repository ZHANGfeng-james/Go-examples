package main

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/rpc/jsonrpc"
	"time"

	"opensource.com/goweb/session"
)

var globalSession *session.Manager

func init() {
	fmt.Println("main init()")
	// cookiename --> gosessionid
	globalSession, _ = session.NewManager("memory", "gosessionid", 3600)
	//globalSession.GC()
}

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

func main() {
	serverAddr := "127.0.0.1:1234"
	client, err := jsonrpc.Dial("tcp", serverAddr)
	checkErr(err)

	args := Args{17, 8}
	var reply int
	err = client.Call("Arith.Multiple", args, &reply)
	checkErr(err)
	fmt.Printf("Arith: %d*%d = %d\n", args.A, args.B, reply)

	var quot Quotient
	err = client.Call("Arith.Divide", args, &quot)
	checkErr(err)
	fmt.Printf("Arith: %d / %d = %d...%d\n", args.A, args.B, quot.Quo, quot.Rem)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal("arith error:", err)
	}
}

func module() {
	ch := make(chan struct{})

	sendMessage(ch)
	<-ch
}

func sendMessage(ch chan<- struct{}) {
	go func() {
		ch <- struct{}{}
	}()
}

func server() {
	http.HandleFunc("/", count)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	sess := globalSession.SessionStart(w, r)
	r.ParseForm()
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		// charset 编码问题，导致浏览器渲染出的中文字符乱码
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")

		fmt.Println("GET:", sess.Get("username"))

		t.Execute(w, sess.Get("username"))
	} else {
		fmt.Println("POST:", r.Form["username"])

		sess.Set("username", r.Form["username"])
		http.Redirect(w, r, "/", 302)
	}
}

func count(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fmt.Printf("path:%s.\n", r.URL.Path)

	if r.URL.Path == "/" {
		sess := globalSession.SessionStart(w, r)
		createTime := sess.Get("createtime")
		if createTime == nil {
			// Session 中 createtime 值单位为 秒
			sess.Set("createtime", time.Now().Unix())
		} else if (createTime.(int64) + 3600) < time.Now().Unix() {
			globalSession.SessionDestory(w, r)
			sess = globalSession.SessionStart(w, r)
		}
		ct := sess.Get("countnum")
		if ct == nil {
			sess.Set("countnum", 1)
		} else {
			sess.Set("countnum", (ct.(int) + 1))
		}

		saveToken := sess.Get("token")
		if r.Form.Get("token") == "" {
			fmt.Println("创建 token")
			h := md5.New()
			salt := "!..."
			io.WriteString(h, salt+time.Now().String())
			token := fmt.Sprintf("%x", h.Sum(nil))
			sess.Set("token", token)
		} else if r.Form.Get("token") != saveToken.(string) {
			// 提示登录
			fmt.Println("非法请求！")
			return
		}

		t, _ := template.ParseFiles("count.gtpl")
		w.Header().Set("Content-Type", "text/html")
		t.Execute(w, sess.Get("token"))
		t.Execute(w, sess.Get("countnum"))
	}
}
