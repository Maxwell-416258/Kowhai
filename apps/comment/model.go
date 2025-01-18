package comment

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"kowhai/apps/base"
)

type Comment struct {
	Id         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	VideoId    string             `json:"video_id" bson:"video_id"`
	ReviewerId string             `json:"reviewer_id" bson:"reviewer_id"`
	Content    string             `json:"content" bson:"content"`
	Likes      int                `json:"likes" bson:"likes"`
	Audit      base.Audit         `json:"audit" bson:"audit"`
}
