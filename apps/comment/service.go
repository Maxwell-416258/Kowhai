package comment

import (
	"context"
	"github.com/gin-gonic/gin"
	"kowhai/global"
	"net/http"
)

func AddComment(c *gin.Context) {
	var client = global.Mongo
	var database = client.Database(global.Config.Mongo.Database)
	var comment_collection = database.Collection("comment")

	var comment Comment
	if err := c.BindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	insertResult, err := comment_collection.InsertOne(context.Background(), comment)
	if err != nil {
		global.Logger.Error("插入评论失败", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": insertResult.InsertedID})

}
