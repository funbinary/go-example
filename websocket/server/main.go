package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

func RunWsService() error {
	// Key is name of room, value is Room

	listen := "0.0.0.0:8888"
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			//解决跨域问题
			return true
		},
	}
	r := gin.Default()
	fmt.Println("---------")
	r.GET("/chat", func(c *gin.Context) {

		_, cancel := context.WithCancel(context.Background())
		defer cancel()
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println("Upgrade失败")
		}

		defer func() {
			fmt.Errorf("关闭连接%s\n", conn.RemoteAddr())
			conn.Close()
		}()
		//读

		ticker := time.NewTicker(10 * time.Second)
		defer func() {
			ticker.Stop()
		}()
		for {
			mt, message, err := conn.ReadMessage()
			if mt == websocket.PingMessage {
				continue
			}

			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					fmt.Errorf("receive Close from %s\n", conn.RemoteAddr())
				}
				fmt.Println(mt)
				fmt.Println(message)
				fmt.Println(err)
				break
			}
			if mt == websocket.TextMessage {
				fmt.Printf("receiver from %s:%s\n", conn.RemoteAddr(), string(message))
				conn.WriteMessage(websocket.TextMessage, message)
			}
			select {
			case <-ticker.C:
				conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
				fmt.Println("Write ping")
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					panic(err)
					return
				}
			}
		}

	})
	return r.Run(listen)
}

func main() {
	RunWsService()
}
