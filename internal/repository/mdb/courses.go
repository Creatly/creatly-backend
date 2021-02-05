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

func (r *CoursesRepo) GetModules(ctx context.Context, courseId primitive.ObjectID) ([]domain.Module, error) {
	var modules []domain.Module
	cur, err := r.db.Collection(modulesCollection).Find(ctx, bson.M{"courseId": courseId})
	if err != nil {
		return nil, err
	}

	err = cur.All(ctx, &modules)
	return modules, err
}

func (r *CoursesRepo) GetModuleWithContent(ctx context.Context, moduleId primitive.ObjectID) (domain.Module, error) {
	var module domain.Module
	err := r.db.Collection(modulesCollection).FindOne(ctx, bson.M{"_id": moduleId, "published": true}).Decode(&module)
	if err != nil {
		return module, err
	}

	lessonIds := make([]primitive.ObjectID, len(module.Lessons))
	for _, lesson := range module.Lessons {
		lessonIds = append(lessonIds, lesson.ID)
	}

	var content []domain.LessonContent
	cur, err := r.db.Collection(contentCollection).Find(ctx, bson.M{"lessonId": bson.M{"$in": lessonIds}})
	if err != nil {
		return module, err
	}

	err = cur.All(ctx, &content)
	if err != nil {
		return module, err
	}

	for i := range module.Lessons {
		for _, lessonContent := range content {
			if module.Lessons[i].ID == lessonContent.LessonID {
				module.Lessons[i].Content = lessonContent.Content
			}
		}
	}

	return module, nil
}

func (r *CoursesRepo) GetModule(ctx context.Context, moduleId primitive.ObjectID) (domain.Module, error) {
	var module domain.Module
	err := r.db.Collection(modulesCollection).FindOne(ctx, bson.M{"_id": moduleId, "published": true}).Decode(&module)
	return module, err
}

func (r *CoursesRepo) GetPackagesModules(ctx context.Context, packageIds []primitive.ObjectID) ([]domain.Module, error) {
	var modules []domain.Module
	cur, err := r.db.Collection(modulesCollection).Find(ctx, bson.M{"packageId": bson.M{"$in": packageIds}})
	if err != nil {
		return nil, err
	}

	err = cur.All(ctx, &modules)
	return modules, err
}

func (r *CoursesRepo) Create(ctx context.Context, schoolId primitive.ObjectID, course domain.Course) (primitive.ObjectID, error) {
	course.ID = primitive.NewObjectID()
	_, err := r.db.Collection(schoolsCollection).UpdateOne(ctx, bson.M{"_id": schoolId}, bson.M{"$push": bson.M{"courses": course}})
	return course.ID, err
}

func (r *CoursesRepo) UpdateCourse(ctx context.Context, schoolId primitive.ObjectID, course domain.Course) error {
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

func (r *CoursesRepo) CreateModule(ctx context.Context, module domain.Module) (primitive.ObjectID, error) {
	res, err := r.db.Collection(modulesCollection).InsertOne(ctx, module)
	return res.InsertedID.(primitive.ObjectID), err
}
