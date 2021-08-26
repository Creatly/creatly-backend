package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Course struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name,omitempty"`
	Code        string             `json:"code" bson:"code,omitempty"`
	Description string             `json:"description" bson:"description,omitempty"`
	Color       string             `json:"color" bson:"color,omitempty"`
	ImageURL    string             `json:"imageUrl" bson:"imageUrl,omitempty"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt,omitempty"`
	Published   bool               `json:"published" bson:"published,omitempty"`
}

type Module struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Position  uint               `json:"position" bson:"position"`
	Published bool               `json:"published"`
	CourseID  primitive.ObjectID `json:"courseId" bson:"courseId"`
	PackageID primitive.ObjectID `json:"packageId,omitempty" bson:"packageId,omitempty"`
	SchoolID  primitive.ObjectID `json:"schoolId" bson:"schoolId"`
	Lessons   []Lesson           `json:"lessons,omitempty" bson:"lessons,omitempty"`
	Survey    Survey             `json:"survey,omitempty" bson:"survey,omitempty"`
}

type Lesson struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Position  uint               `json:"position" bson:"position"`
	Published bool               `json:"published" bson:"published,omitempty"`
	Content   string             `json:"content,omitempty" bson:"content,omitempty"`
	SchoolID  primitive.ObjectID `json:"schoolId" bson:"schoolId"`
}

type LessonContent struct {
	LessonID primitive.ObjectID `json:"lessonId" bson:"lessonId"`
	SchoolID primitive.ObjectID `json:"schoolId" bson:"schoolId"`
	Content  string             `json:"content" bson:"content"`
}

type Package struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	CourseID primitive.ObjectID `json:"courseId" bson:"courseId"`
	SchoolID primitive.ObjectID `json:"schoolId" bson:"schoolId"`
	Modules  []Module           `json:"modules" bson:"-"`
}

type ModuleContent struct {
	Lessons []Lesson `json:"lessons" bson:"lessons"`
	Survey  Survey   `json:"survey" bson:"survey"`
}
