package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type CourseEntity struct {
	Name      string `json:"name" bson:"name"`
	Position  int    `json:"position" bson:"position"`
	Published bool   `json:"published"`
}

type Course struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Code        string             `json:"code" bson:"code"`
	Description string             `json:"description" bson:"description"`
	ImageURL    string             `json:"imageUrl" bson:"imageUrl"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
	Published   bool               `json:"published" bson:"published"`
}

type Module struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	Position  int                `json:"position" bson:"position"`
	Published bool               `json:"published"`
	CourseID  primitive.ObjectID `json:"courseId" bson:"courseId"`
	PackageID primitive.ObjectID `json:"packageId" bson:"packageId"`
	Lessons   []Lesson           `json:"lessons" bson:"lessons"`
}

type Lesson struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	Position  int                `json:"position" bson:"position"`
	Published bool               `json:"published"`
}

type CourseContent struct {
	LessonID primitive.ObjectID `json:"lessonId" bson:"lessonId"`
	Content  string             `json:"content" bson:"content"`
}

type CoursePackages struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	CourseID    primitive.ObjectID `json:"courseId" bson:"courseId"`
}
