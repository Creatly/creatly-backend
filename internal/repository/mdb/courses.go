package mdb

import (
	"context"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

func (r *CoursesRepo) Update(ctx context.Context, schoolId primitive.ObjectID, course domain.Course) error {
	updateQuery := bson.M{}

	if course.Name != "" {
		updateQuery["courses.$.name"] = course.Name
	}

	if course.Description != "" {
		updateQuery["courses.$.description"] = course.Description
	}

	if course.Code != "" {
		updateQuery["courses.$.code"] = course.Code
	}

	if course.Published != nil {
		updateQuery["courses.$.published"] = *course.Published
	}

	_, err := r.db.Collection(schoolsCollection).UpdateOne(ctx,
		bson.M{"_id": schoolId, "courses._id": course.ID}, bson.M{"$set": updateQuery})

	return err
}
