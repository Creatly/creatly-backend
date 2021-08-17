package repository

import (
	"context"
	"time"

	"github.com/zhashkevych/creatly-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CoursesRepo struct {
	db *mongo.Collection
}

func NewCoursesRepo(db *mongo.Database) *CoursesRepo {
	return &CoursesRepo{db: db.Collection(schoolsCollection)}
}

func (r *CoursesRepo) Create(ctx context.Context, schoolId primitive.ObjectID, course domain.Course) (primitive.ObjectID, error) {
	course.ID = primitive.NewObjectID()
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": schoolId}, bson.M{"$push": bson.M{"courses": course}})

	return course.ID, err
}

func (r *CoursesRepo) Update(ctx context.Context, inp UpdateCourseInput) error {
	updateQuery := bson.M{}

	updateQuery["courses.$.updatedAt"] = time.Now()

	if inp.Name != nil {
		updateQuery["courses.$.name"] = *inp.Name
	}

	if inp.Description != nil {
		updateQuery["courses.$.description"] = *inp.Description
	}

	if inp.ImageURL != nil {
		updateQuery["courses.$.imageUrl"] = *inp.ImageURL
	}

	if inp.Color != nil {
		updateQuery["courses.$.color"] = *inp.Color
	}

	if inp.Published != nil {
		updateQuery["courses.$.published"] = *inp.Published
	}

	_, err := r.db.UpdateOne(ctx,
		bson.M{"_id": inp.SchoolID, "courses._id": inp.ID}, bson.M{"$set": updateQuery})

	return err
}

func (r *CoursesRepo) Delete(ctx context.Context, schoolId, courseId primitive.ObjectID) error {
	res, err := r.db.UpdateOne(ctx, bson.M{"_id": schoolId}, bson.M{"$pull": bson.M{"courses": bson.M{"_id": courseId}}})
	if err != nil {
		return err
	}

	if res.ModifiedCount == 0 {
		return domain.ErrCourseNotFound
	}

	return nil
}
