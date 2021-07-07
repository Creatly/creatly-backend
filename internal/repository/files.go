package repository

import (
	"context"

	"github.com/zhashkevych/creatly-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FilesRepo struct {
	db *mongo.Collection
}

func NewFilesRepo(db *mongo.Database) *FilesRepo {
	return &FilesRepo{
		db: db.Collection(filesCollection),
	}
}

func (r *FilesRepo) Create(ctx context.Context, file domain.File) (primitive.ObjectID, error) {
	res, err := r.db.InsertOne(ctx, file)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	return res.InsertedID.(primitive.ObjectID), nil
}

func (r *FilesRepo) UpdateStatus(ctx context.Context, fileName string, status domain.FileStatus) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"name": fileName}, bson.M{"$set": bson.M{"status": status}})

	return err
}

func (r *FilesRepo) GetForUploading(ctx context.Context) (domain.File, error) {
	var file domain.File

	res := r.db.FindOneAndUpdate(ctx, bson.M{"status": domain.UploadedByClient}, bson.M{"$set": bson.M{"status": domain.StorageUploadInProgress}})
	err := res.Decode(&file)

	return file, err
}

func (r *FilesRepo) UpdateStatusAndSetURL(ctx context.Context, id primitive.ObjectID, url string) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"url": url, "status": domain.UploadedToStorage}})

	return err
}

func (r *FilesRepo) GetByID(ctx context.Context, id, schoolId primitive.ObjectID) (domain.File, error) {
	var file domain.File

	res := r.db.FindOne(ctx, bson.M{"_id": id, "schoolId": schoolId})
	err := res.Decode(&file)

	return file, err
}
