package repository

import (
	"context"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type CoursesRepo struct {
	db *mongo.Database
}

func NewCoursesRepo(db *mongo.Database) *CoursesRepo {
	return &CoursesRepo{db: db}
}

func (r *CoursesRepo) Create(ctx context.Context, schoolId primitive.ObjectID, course domain.Course) (primitive.ObjectID, error) {
	course.ID = primitive.NewObjectID()
	_, err := r.db.Collection(schoolsCollection).UpdateOne(ctx, bson.M{"_id": schoolId}, bson.M{"$push": bson.M{"courses": course}})
	return course.ID, err
}

func (r *CoursesRepo) Update(ctx context.Context, schoolId primitive.ObjectID, inp UpdateCourseInput) error {
	updateQuery := bson.M{}

	updateQuery["courses.$.updatedAt"] = time.Now()

	if inp.Name != "" {
		updateQuery["courses.$.name"] = inp.Name
	}

	if inp.Description != "" {
		updateQuery["courses.$.description"] = inp.Description
	}

	if inp.Code != "" {
		updateQuery["courses.$.code"] = inp.Code
	}

	if inp.Published != nil {
		updateQuery["courses.$.published"] = *inp.Published
	}

	_, err := r.db.Collection(schoolsCollection).UpdateOne(ctx,
		bson.M{"_id": schoolId, "courses._id": inp.ID}, bson.M{"$set": updateQuery})

	return err
}
