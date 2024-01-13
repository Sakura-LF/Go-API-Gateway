package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

func WebSocketServer() {
	addr := "localhost:8002"
	http.HandleFunc("/wshandler", WebSocketUpgrade)
	log.Println("Starting websocket server at " + addr)

	go func() {
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("WebSocket 服务器正在运行。按Ctrl+C退出")
	select {}
}

func WebSocketUpgrade(resp http.ResponseWriter, req *http.Request) {
	// 初始化 Upgrader
	upgrader := websocket.Upgrader{} // 使用默认的选项
	// 第三个参数是响应头,默认会初始化
	conn, err := upgrader.Upgrade(resp, req, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// 主动向服务端推送消息
	go PushMessage(conn)

	// 读取客户端的发送额消息,并返回
	go ReadMessage(conn)
	select {}
}

func PushMessage(conn *websocket.Conn) {
	for {
		err := conn.WriteMessage(websocket.TextMessage, []byte("heart beat"))
		if err != nil {
			log.Println(err)
			return
		}
		time.Sleep(time.Second * 3)
	}
}

func ReadMessage(conn *websocket.Conn) {
	for {
		// 消息类型:文本消息和二进制消息
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("receive msg:", string(msg))

		err = conn.WriteMessage(messageType, msg)
		if err != nil {
			log.Println("write error:", err)
			return
		}

	}
}

//	var upgrader = websocket.Upgrader{
//		CheckOrigin: func(r *http.Request) bool {
//			return true
//		},
//	}
//func WebSocketUpgrade(w http.ResponseWriter, r *http.Request) {
//	upgrader := websocket.Upgrader{}
//	conn, err := upgrader.Upgrade(w, r, nil)
//	if err != nil {
//		logs.Println(err)
//		return
//	}
//	defer conn.Close()
//
//	for {
//		messageType, p, err := conn.ReadMessage()
//		if err != nil {
//			logs.Println(err)
//			return
//		}
//		if err := conn.WriteMessage(messageType, p); err != nil {
//			logs.Println(err)
//			return
//		}
//	}
//}

//func WebSocketServer() {
//	addr := "localhost:8002"
//	http.HandleFunc("/wshandler", WebSocketUpgrade)
//	logs.Println("正在监听地址 " + addr)
//
//	go func() {
//		err := http.ListenAndServe(addr, nil)
//		if err != nil {
//			logs.Fatal(err)
//		}
//	}()
//
//	logs.Println("WebSocket 服务器正在运行。按Ctrl+C退出。")
//	select {}
//}
