package comment

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"kowhai/global"
	"net/http"
	"strconv"
)

// 添加评论
func AddComment(c *gin.Context) {
	var client = global.Mongo
	var database = client.Database(global.Config.Mongo.Database)
	var comment_collection = database.Collection("comment")

	var comment Comment
	if err := c.BindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "添加失败", "err": err.Error()})
		return
	}

	insertResult, err := comment_collection.InsertOne(context.Background(), comment)
	if err != nil {
		global.Logger.Error("插入评论失败", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "添加失败", "err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "添加成功", "data": insertResult.InsertedID})

}

// 获取评论总数
func GetCommentTotal(c *gin.Context) {
	var client = global.Mongo
	var database = client.Database(global.Config.Mongo.Database)
	var comment_collection = database.Collection("comment")

	var video_id = c.Query("video_id")
	if video_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "video_id不能为空", "err": ""})
		return
	}

	var total int64
	total, err := comment_collection.CountDocuments(context.Background(), bson.M{"video_id": video_id})
	if err != nil {
		global.Logger.Error("获取评论总数失败", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "获取评论失败", "err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "获取评论成功", "data": total})
}

// 获取评论列表
func GetCommentList(c *gin.Context) {
	var client = global.Mongo
	var database = client.Database(global.Config.Mongo.Database)
	var comment_collection = database.Collection("comment")

	var CommentList []Comment

	var video_id = c.Query("video_id")
	if video_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "video_id不能为空", "err": ""})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "15"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 15
	}
	offset := int64((page - 1) * pageSize)
	limit := int64(pageSize)

	// 构建查询条件
	filter := bson.M{"video_id": video_id}

	// 查询mongodb
	cursor, err := comment_collection.Find(c, filter, &options.FindOptions{
		Skip:  &offset,
		Limit: &limit,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "无法获取评论列表", "err": ""})
		return
	}
	defer cursor.Close(c)

	// 解析查询结果
	if err := cursor.All(c, &CommentList); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "无法解析评论数据", "err": ""})
		return
	}
	// 返回结果
	data := map[string]interface{}{
		"page":     page,
		"pageSize": pageSize,
		"total":    len(CommentList),
		"comments": CommentList,
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取评论成功",
		"data": data,
	})
}
