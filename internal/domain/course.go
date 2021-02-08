package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Course struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name,omitempty"`
	Code        string             `json:"code" bson:"code,omitempty"`
	Description string             `json:"description" bson:"description,omitempty"`
	ImageURL    string             `json:"imageUrl" bson:"imageUrl,omitempty"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt,omitempty"`
	Published   bool               `json:"published" bson:"published,omitempty"`
}

type Module struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Position  int                `json:"position" bson:"position"`
	Published bool               `json:"published"`
	CourseID  primitive.ObjectID `json:"courseId" bson:"courseId"`
	PackageID primitive.ObjectID `json:"packageId,omitempty" bson:"packageId,omitempty"`
	Lessons   []Lesson           `json:"lessons,omitempty" bson:"lessons,omitempty"`
}

type Lesson struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Position  int                `json:"position" bson:"position"`
	Published bool               `json:"published" bson:"published,omitempty"`
	Content   string             `json:"content,omitempty" bson:"content,omitempty"`
}

type LessonContent struct {
	LessonID primitive.ObjectID `json:"lessonId" bson:"lessonId"`
	Content  string             `json:"content" bson:"content"`
}

type CoursePackages struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	CourseID    primitive.ObjectID `json:"courseId" bson:"courseId"`
}
