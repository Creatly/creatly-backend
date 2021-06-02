package repository

import (
	"context"
	"time"

	"github.com/zhashkevych/creatly-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AdminsRepo struct {
	db *mongo.Collection
}

func NewAdminsRepo(db *mongo.Database) *AdminsRepo {
	return &AdminsRepo{db: db.Collection(adminsCollection)}
}

func (r *AdminsRepo) GetByCredentials(ctx context.Context, schoolId primitive.ObjectID, email, password string) (domain.Admin, error) {
	var admin domain.Admin
	err := r.db.FindOne(ctx, bson.M{"schoolId": schoolId, "email": email, "password": password}).Decode(&admin)

	return admin, err
}

func (r *AdminsRepo) GetByRefreshToken(ctx context.Context, schoolId primitive.ObjectID, refreshToken string) (domain.Admin, error) {
	var admin domain.Admin
	err := r.db.FindOne(ctx, bson.M{
		"session.refreshToken": refreshToken, "schoolId": schoolId,
		"session.expiresAt": bson.M{"$gt": time.Now()},
	}).Decode(&admin)

	return admin, err
}

func (r *AdminsRepo) SetSession(ctx context.Context, id primitive.ObjectID, session domain.Session) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"session": session}})

	return err
}

func (r *AdminsRepo) GetById(ctx context.Context, id primitive.ObjectID) (domain.Admin, error) {
	var admin domain.Admin

	err := r.db.FindOne(ctx, bson.M{"_id": id}).Decode(&admin)

	return admin, err
}
