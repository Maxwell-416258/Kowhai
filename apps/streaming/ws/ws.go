package ws

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

// 维护 WebSocket 连接
var clients = make(map[*websocket.Conn]bool)
var lock = sync.Mutex{}

// WebSocket 升级器
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有跨域请求
	},
}

// WebSocket 处理函数
func HandleConnections(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("WebSocket 连接失败:", err)
		return
	}
	defer conn.Close()

	// 添加到连接列表
	lock.Lock()
	clients[conn] = true
	lock.Unlock()

	fmt.Println("WebSocket 连接成功")

	// 监听 WebSocket 关闭
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			lock.Lock()
			delete(clients, conn)
			lock.Unlock()
			fmt.Println("WebSocket 断开连接")
			break
		}
	}
}

// 发送消息给所有 WebSocket 连接
func BroadcastMessage(message string) {
	lock.Lock()
	defer lock.Unlock()

	for conn := range clients {
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			fmt.Println("发送 WebSocket 消息失败:", err)
			conn.Close()
			delete(clients, conn)
		}
	}
}
