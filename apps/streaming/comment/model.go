package comment

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"kowhai/apps/streaming/base"
)

type Comment struct {
	Id         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	VideoId    string             `json:"videoId" bson:"video_id"`
	ReviewerId string             `json:"reviewerId" bson:"reviewer_id"`
	UserName   string             `json:"userName" bson:"user_name"`
	Content    string             `json:"content" bson:"content"`
	Likes      int                `json:"likes" bson:"likes"`
	Audit      base.Audit         `json:"audit" bson:"audit"`
}
