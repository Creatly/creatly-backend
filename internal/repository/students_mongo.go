package repository

import (
	"context"
	"errors"
	"time"

	"github.com/zhashkevych/creatly-backend/internal/domain"
	"github.com/zhashkevych/creatly-backend/pkg/database/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	if mongodb.IsDuplicate(err) {
		return domain.ErrUserAlreadyExists
	}

	return err
}

func (r *StudentsRepo) GetByCredentials(ctx context.Context, schoolId primitive.ObjectID, email, password string) (domain.Student, error) {
	var student domain.Student
	if err := r.db.FindOne(ctx, bson.M{"email": email, "password": password, "schoolId": schoolId, "verification.verified": true}).
		Decode(&student); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Student{}, domain.ErrUserNotFound
		}

		return domain.Student{}, err
	}

	return student, nil
}

func (r *StudentsRepo) GetByRefreshToken(ctx context.Context, schoolId primitive.ObjectID, refreshToken string) (domain.Student, error) {
	var student domain.Student
	if err := r.db.FindOne(ctx, bson.M{
		"session.refreshToken": refreshToken, "schoolId": schoolId,
		"session.expiresAt": bson.M{"$gt": time.Now()},
	}).Decode(&student); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Student{}, domain.ErrUserNotFound
		}

		return domain.Student{}, err
	}

	return student, nil
}

func (r *StudentsRepo) GetById(ctx context.Context, schoolId, id primitive.ObjectID) (domain.Student, error) {
	var student domain.Student

	if err := r.db.FindOne(ctx, bson.M{"_id": id, "schoolId": schoolId, "verification.verified": true}).Decode(&student); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Student{}, domain.ErrUserNotFound
		}

		return domain.Student{}, err
	}

	return student, nil
}

func (r *StudentsRepo) GetBySchool(ctx context.Context, schoolId primitive.ObjectID, pagination *domain.PaginationQuery) ([]domain.Student, int64, error) {
	opts := getPaginationOpts(pagination)
	opts.SetSort(bson.M{"registeredAt": -1})

	filter := bson.M{"schoolId": schoolId}

	cur, err := r.db.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}

	var students []domain.Student
	if err := cur.All(ctx, &students); err != nil {
		return nil, 0, err
	}

	count, err := r.db.CountDocuments(ctx, filter)

	return students, count, err
}

func (r *StudentsRepo) SetSession(ctx context.Context, studentId primitive.ObjectID, session domain.Session) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": studentId}, bson.M{"$set": bson.M{"session": session, "lastVisitAt": time.Now()}})

	return err
}

func (r *StudentsRepo) GiveAccessToModule(ctx context.Context, studentId, moduleId primitive.ObjectID) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": studentId}, bson.M{"$addToSet": bson.M{"availableModules": moduleId}})

	return err
}

func (r *StudentsRepo) AttachOffer(ctx context.Context, studentId, offerId primitive.ObjectID, moduleIds []primitive.ObjectID) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": studentId}, bson.M{"$addToSet": bson.M{
		"availableModules": bson.M{"$each": moduleIds},
		"availableOffers":  offerId,
	}})

	return err
}

func (r *StudentsRepo) DetachOffer(ctx context.Context, studentId, offerId primitive.ObjectID, moduleIds []primitive.ObjectID) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": studentId}, bson.M{"$pull": bson.M{
		"availableModules": bson.M{"$in": moduleIds},
		"availableOffers":  offerId,
	}})

	return err
}

func (r *StudentsRepo) Verify(ctx context.Context, code string) error {
	res, err := r.db.UpdateOne(ctx,
		bson.M{"verification.code": code},
		bson.M{"$set": bson.M{"verification.verified": true, "verification.code": ""}})
	if err != nil {
		return err
	}

	if res.ModifiedCount == 0 {
		return domain.ErrVerificationCodeInvalid
	}

	return nil
}
