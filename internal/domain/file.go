package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type FileStatus int
type FileType string

const (
	ClientUploadInProgress FileStatus = iota
	UploadedByClient
	ClientUploadError
	StorageUploadInProgress
	UploadedToStorage
	StorageUploadError
)

const (
	Image FileType = "image"
	Video FileType = "video"
	Other FileType = "other"
)

type File struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	SchoolID        primitive.ObjectID `json:"schoolId" bson:"schoolId"`
	Type            FileType           `json:"type" bson:"type"` // todo leave or remove ?
	ContentType     string             `json:"contentType" bson:"contentType"`
	Name            string             `json:"name" bson:"name"`
	Size            int64              `json:"size" bson:"size"`
	Status          FileStatus         `json:"status" bson:"status"`
	UploadStartedAt time.Time          `json:"uploadStartedAt" bson:"uploadStartedAt"`
	URL             string             `json:"url" bson:"url,omitempty"`
}
