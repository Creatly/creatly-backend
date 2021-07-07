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
