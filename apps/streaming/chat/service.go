package chat

import (
	"context"
	"github.com/gin-gonic/gin"
	"kowhai/global"
	"log"
	"net/http"
)

func GetHistoryMessages(c *gin.Context) {
	// **获取历史消息** ,默认返回30天的消息
	// Params: id1,id2

	senderID := c.Query("sender_id")
	receiverID := c.Query("receiver_id")
	if senderID == "" || receiverID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数 sender_id 和 receiver_id 不能为空"})
		return
	}

	// **执行 SQL 查询**
	query := `
		SELECT id, sender_id, receiver_id, content, created_at
		FROM messages
		WHERE ((sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1))
		AND created_at >= NOW() - INTERVAL '30 days'
		ORDER BY created_at ASC;
	`

	rows, err := global.DbPool.Query(context.Background(), query, senderID, receiverID)
	if err != nil {
		log.Println("数据库查询失败:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	defer rows.Close()

	// **解析数据**
	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.CreatedAt); err != nil {
			log.Println("数据解析失败:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "数据解析失败"})
			return
		}
		messages = append(messages, msg)
	}

	c.JSON(http.StatusOK, messages)
}
