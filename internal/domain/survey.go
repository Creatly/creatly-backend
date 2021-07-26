package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Survey struct {
	Title     string           `json:"title" bson:"title"`
	Questions []SurveyQuestion `json:"questions" bson:"questions"`
	Required  bool             `json:"required" bson:"required"`
}

type SurveyQuestion struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Question      string             `json:"question" bson:"question"`
	AnswerType    string             `json:"answerType" bson:"answerType"`
	AnswerOptions []string           `json:"answerOptions" bson:"answerOptions,omitempty"`
}

type SurveyResult struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Student     StudentInfoShort   `json:"student" bson:"student"`
	ModuleID    primitive.ObjectID `json:"moduleId" bson:"moduleId"`
	SubmittedAt time.Time          `json:"submittedAt" bson:"submittedAt"`
	Answers     []SurveyAnswer     `json:"answers" bson:"answers"`
}

type SurveyAnswer struct {
	QuestionID primitive.ObjectID `json:"questionId" bson:"questionId"`
	Answer     string             `json:"answer" bson:"answer"`
}
