package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Survey struct {
	Title     string           `json:"title"`
	Questions []SurveyQuestion `json:"questions"`
	Required  bool             `json:"required"`
}

type SurveyQuestion struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Question      string             `json:"question"`
	AnswerType    string             `json:"answerType"`
	AnswerOptions []string           `json:"answerOptions"`
}

type SurveyResult struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	StudentID   primitive.ObjectID `json:"studentId" bson:"studentId"`
	ModuleID    primitive.ObjectID `json:"moduleId" bson:"moduleId"`
	SubmittedAt time.Time          `json:"submittedAt" bson:"submittedAt"`
	Answers     []SurveyAnswer     `json:"answers" bson:"answers"`
}

type SurveyAnswer struct {
	QuestionID primitive.ObjectID `json:"questionId" bson:"questionId"`
	Answer     string             `json:"answer"`
}
