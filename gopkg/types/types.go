package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Collection struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Album struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Media struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name            string             `json:"name" bson:"name"`
	Original        string             `json:"original" bson:"original"`
	View            string             `json:"view" bson:"view"`
	ViewWidth       int32              `json:"viewWidth" bson:"viewWidth"`
	ViewHeight      int32              `json:"viewHeight" bson:"viewHeight"`
	VideoDuration   int32              `json:"videoDuration" bson:"videoDuration"`
	Thumbnail       string             `json:"thumbnail" bson:"thumbnail"`
	ThumbnailWidth  int32              `json:"thumbnailWidth" bson:"thumbnailWidth"`
	ThumbnailHeight int32              `json:"thumbnailHeight" bson:"thumbnailHeight"`
	Preview         string             `json:"preview" bson:"preview"`
	Type            string             `json:"type" bson:"type"`
	ProcessStatus   string             `json:"processStatus" bson:"processStatus"`
}

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Session struct {
	ID string `json:"id"`
	// expiresAt: Date;

	// createdAt: Date;
	// updatedAt: Date;
}
