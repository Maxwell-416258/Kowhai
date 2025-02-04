package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	pb "kowhai/api/pb"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketServer struct {
	connections map[int64]*websocket.Conn
	mu          sync.Mutex
	grpcClient  pb.ChatServiceClient
}

type Message struct {
	SenderID   int64  `json:"sender_id"`
	ReceiverID int64  `json:"receiver_id"`
	Content    string `json:"content"`
}

func (s *WebSocketServer) handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket 连接失败:", err)
		return
	}
	defer conn.Close() // WebSocket 断开后，确保连接被关闭

	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		log.Println("无效的 user_id:", userIDStr)
		return
	}

	s.mu.Lock()
	s.connections[userID] = conn
	s.mu.Unlock()

	log.Printf("用户 %d 连接 WebSocket", userID)

	// **增加一个 context 控制 gRPC 消息监听**
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // **WebSocket 断开时，自动取消 gRPC 监听**

	go s.receiveMessagesFromGrpc(ctx, userID) // **传递 ctx 以便监听断开**

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket 读取错误:", err)
			break
		}

		var message Message
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Println("JSON 解析失败:", err)
			continue
		}

		s.sendMessageToGrpc(message)
	}

	// **确保删除用户连接信息**
	s.mu.Lock()
	delete(s.connections, userID)
	s.mu.Unlock()

	log.Printf("用户 %d WebSocket 断开", userID)
}

func (s *WebSocketServer) sendMessageToGrpc(msg Message) {
	_, err := s.grpcClient.SendMessage(context.Background(), &pb.SendMessageRequest{
		SenderId:   msg.SenderID,
		ReceiverId: msg.ReceiverID,
		Message:    msg.Content,
	})
	if err != nil {
		log.Println("gRPC 发送消息失败:", err)
	}
}

func (s *WebSocketServer) receiveMessagesFromGrpc(ctx context.Context, userID int64) {
	stream, err := s.grpcClient.ReceiveMessages(ctx, &pb.ListenMessagesRequest{UserId: userID})
	if err != nil {
		log.Println("gRPC 监听消息失败:", err)
		return
	}

	for {
		select {
		case <-ctx.Done(): // **监听 WebSocket 是否关闭**
			log.Printf("检测到用户 %d WebSocket 断开，关闭 gRPC 消息监听", userID)
			return
		default:
			msg, err := stream.Recv()
			if err != nil {
				log.Println("gRPC 消息流断开:", err)
				return
			}

			s.mu.Lock()
			conn, exists := s.connections[msg.ReceiverId]
			s.mu.Unlock()

			if exists {
				if err := conn.WriteJSON(msg); err != nil {
					log.Println("WebSocket 发送消息失败:", err)
				}
			}
		}
	}
}

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal("无法连接 gRPC 服务器:", err)
	}
	defer conn.Close()

	server := &WebSocketServer{
		connections: make(map[int64]*websocket.Conn),
		grpcClient:  pb.NewChatServiceClient(conn),
	}

	http.HandleFunc("/ws", server.handleConnection)
	log.Println("WebSocket 服务器启动，监听 8082 端口")
	http.ListenAndServe(":8082", nil)
}
