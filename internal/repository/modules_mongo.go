package repository

import (
	"context"

	"github.com/zhashkevych/creatly-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ModulesRepo struct {
	db *mongo.Collection
}

func NewModulesRepo(db *mongo.Database) *ModulesRepo {
	return &ModulesRepo{db: db.Collection(modulesCollection)}
}

func (r *ModulesRepo) Create(ctx context.Context, module domain.Module) (primitive.ObjectID, error) {
	res, err := r.db.InsertOne(ctx, module)

	return res.InsertedID.(primitive.ObjectID), err
}

func (r *ModulesRepo) GetPublishedByCourseId(ctx context.Context, courseId primitive.ObjectID) ([]domain.Module, error) {
	var modules []domain.Module

	opts := options.Find()
	opts.SetSort(bson.M{"position": 1})

	cur, err := r.db.Find(ctx, bson.M{"courseId": courseId, "published": true}, opts)
	if err != nil {
		return nil, err
	}

	err = cur.All(ctx, &modules)

	return modules, err
}

func (r *ModulesRepo) GetByCourseId(ctx context.Context, courseId primitive.ObjectID) ([]domain.Module, error) {
	var modules []domain.Module

	opts := options.Find()
	opts.SetSort(bson.M{"position": 1})

	cur, err := r.db.Find(ctx, bson.M{"courseId": courseId}, opts)
	if err != nil {
		return nil, err
	}

	err = cur.All(ctx, &modules)

	return modules, err
}

func (r *ModulesRepo) GetPublishedById(ctx context.Context, moduleId primitive.ObjectID) (domain.Module, error) {
	var module domain.Module

	err := r.db.FindOne(ctx, bson.M{"_id": moduleId, "published": true}).Decode(&module)

	return module, err
}

func (r *ModulesRepo) GetById(ctx context.Context, moduleId primitive.ObjectID) (domain.Module, error) {
	var module domain.Module

	err := r.db.FindOne(ctx, bson.M{"_id": moduleId}).Decode(&module)

	return module, err
}

func (r *ModulesRepo) GetByPackages(ctx context.Context, packageIds []primitive.ObjectID) ([]domain.Module, error) {
	var modules []domain.Module

	opts := options.Find()
	opts.SetSort(bson.M{"position": 1})

	cur, err := r.db.Find(ctx, bson.M{"packageId": bson.M{"$in": packageIds}}, opts)
	if err != nil {
		return nil, err
	}

	err = cur.All(ctx, &modules)

	return modules, err
}

func (r *ModulesRepo) Update(ctx context.Context, inp UpdateModuleInput) error {
	updateQuery := bson.M{}

	if inp.Name != "" {
		updateQuery["name"] = inp.Name
	}

	if inp.Position != nil {
		updateQuery["position"] = *inp.Position
	}

	if inp.Published != nil {
		updateQuery["published"] = *inp.Published
	}

	_, err := r.db.UpdateOne(ctx,
		bson.M{"_id": inp.ID, "schoolId": inp.SchoolID}, bson.M{"$set": updateQuery})

	return err
}

func (r *ModulesRepo) Delete(ctx context.Context, schoolId, id primitive.ObjectID) error {
	_, err := r.db.DeleteOne(ctx, bson.M{"_id": id, "schoolId": schoolId})

	return err
}

func (r *ModulesRepo) AddLesson(ctx context.Context, schoolId, id primitive.ObjectID, lesson domain.Lesson) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": id, "schoolId": schoolId}, bson.M{"$push": bson.M{"lessons": lesson}})

	return err
}

func (r *ModulesRepo) GetByLesson(ctx context.Context, lessonId primitive.ObjectID) (domain.Module, error) {
	var module domain.Module
	err := r.db.FindOne(ctx, bson.M{"lessons._id": lessonId}).Decode(&module)

	return module, err
}

func (r *ModulesRepo) UpdateLesson(ctx context.Context, inp UpdateLessonInput) error {
	updateQuery := bson.M{}

	if inp.Name != "" {
		updateQuery["lessons.$.name"] = inp.Name
	}

	if inp.Position != nil {
		updateQuery["lessons.$.position"] = *inp.Position
	}

	if inp.Published != nil {
		updateQuery["lessons.$.published"] = *inp.Published
	}

	_, err := r.db.UpdateOne(ctx,
		bson.M{"lessons._id": inp.ID, "schoolId": inp.SchoolID}, bson.M{"$set": updateQuery})

	return err
}

func (r *ModulesRepo) DeleteLesson(ctx context.Context, schoolId, id primitive.ObjectID) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"lessons._id": id, "schoolId": schoolId}, bson.M{"$pull": bson.M{"lessons": bson.M{"_id": id}}})

	return err
}

func (r *ModulesRepo) DetachPackageFromAll(ctx context.Context, schoolId, packageId primitive.ObjectID) error {
	_, err := r.db.UpdateMany(ctx, bson.M{"schoolId": schoolId, "packageId": packageId}, bson.M{"$unset": bson.M{"packageId": ""}})

	return err
}

func (r *ModulesRepo) AttachPackage(ctx context.Context, schoolId, packageId primitive.ObjectID, modules []primitive.ObjectID) error {
	_, err := r.db.UpdateMany(ctx, bson.M{"_id": bson.M{"$in": modules}, "schoolId": schoolId}, bson.M{"$set": bson.M{"packageId": packageId}})

	return err
}

func (r *ModulesRepo) DeleteByCourse(ctx context.Context, schoolId, courseId primitive.ObjectID) error {
	_, err := r.db.DeleteMany(ctx, bson.M{"courseId": courseId, "schoolId": schoolId})

	return err
}

func (r *ModulesRepo) AttachSurvey(ctx context.Context, schoolId, id primitive.ObjectID, survey domain.Survey) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": id, "schoolId": schoolId}, bson.M{"$set": bson.M{"survey": survey}})

	return err
}

func (r *ModulesRepo) DetachSurvey(ctx context.Context, schoolId, id primitive.ObjectID) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": id, "schoolId": schoolId}, bson.M{"$unset": bson.M{"survey": ""}})

	return err
}
