package http_proxy

import (
	"fmt"
	"net/http"
	"time"
)

func RealHTTPServer() {
	// 初始化RealServer服务器结构体
	server := RealServer{Addr: "127.0.0.1:8001"}
	server.Run()
	select {}
}

func (r *RealServer) Run() {
	mux := http.NewServeMux()
	mux.HandleFunc("/sakura", r.HelloHandler)
	server := http.Server{
		Addr:         r.Addr,
		Handler:      mux,
		WriteTimeout: time.Second * 3,
	}
	go func() {
		fmt.Println("http服务器务已启动:", server.Addr)
		server.ListenAndServe()
	}()
}

type RealServer struct {
	Addr string
}

func (r *RealServer) HelloHandler(w http.ResponseWriter, req *http.Request) {
	// 拼接真实服务器地址
	URL := fmt.Sprintf("这是真实服务器,地址为: http://%s%s", req.RemoteAddr, req.URL.Path)
	w.Write([]byte(URL))
}
