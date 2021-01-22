package mdb

import (
	"context"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type StudentsRepo struct {
	db *mongo.Collection
}

func NewStudentsRepo(db *mongo.Database) *StudentsRepo {
	return &StudentsRepo{
		db: db.Collection(studentsCollection),
	}
}

func (r *StudentsRepo) Create(ctx context.Context, student domain.Student) error {
	_, err := r.db.InsertOne(ctx, student)
	return err
}

func (r *StudentsRepo) GetByCredentials(ctx context.Context, email, password domain.Student) error {
	return nil
}

func (r *StudentsRepo) Verify(ctx context.Context, hash string) error {
	return nil
}
