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

func (r *StudentsRepo) Create(ctx context.Context, student *domain.Student) error {
	res, err := r.db.InsertOne(ctx, student)
	if mongodb.IsDuplicate(err) {
		return domain.ErrUserAlreadyExists
	}

	student.ID = res.InsertedID.(primitive.ObjectID) //nolint:forcetypeassert

	return nil
}

func (r *StudentsRepo) Update(ctx context.Context, inp domain.UpdateStudentInput) error {
	updateQuery := bson.M{}

	if inp.Name != "" {
		updateQuery["name"] = inp.Name
	}

	if inp.Email != "" {
		updateQuery["email"] = inp.Email
	}

	if inp.Verified != nil {
		updateQuery["verification"] = domain.Verification{
			Verified: *inp.Verified,
		}
	}

	if inp.Blocked != nil {
		updateQuery["blocked"] = *inp.Blocked
	}

	_, err := r.db.UpdateOne(ctx,
		bson.M{"_id": inp.StudentID, "schoolId": inp.SchoolID}, bson.M{"$set": updateQuery})

	return err
}

func (r *StudentsRepo) Delete(ctx context.Context, schoolId, studentId primitive.ObjectID) error {
	_, err := r.db.DeleteOne(ctx, bson.M{"_id": studentId, "schoolId": schoolId})

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

	if err := r.db.FindOne(ctx, bson.M{"_id": id, "schoolId": schoolId}).Decode(&student); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Student{}, domain.ErrUserNotFound
		}

		return domain.Student{}, err
	}

	return student, nil
}

func (r *StudentsRepo) GetBySchool(ctx context.Context, schoolId primitive.ObjectID, query domain.GetStudentsQuery) ([]domain.Student, int64, error) {
	paginationOpts := getPaginationOpts(&query.PaginationQuery)
	paginationOpts.SetSort(bson.M{"registeredAt": -1})

	filter := bson.M{"$and": []bson.M{{"schoolId": schoolId}}}

	if query.Search != "" {
		expression := primitive.Regex{Pattern: query.Search}

		filter["$and"] = append(filter["$and"].([]bson.M), bson.M{
			"$or": []bson.M{
				{"name": expression},
				{"email": expression},
			},
		})
	}

	if query.Verified != nil {
		filter["$and"] = append(filter["$and"].([]bson.M), bson.M{
			"verification.verified": *query.Verified,
		})
	}

	if err := filterDateQueries(query.RegisterDateFrom, query.RegisterDateTo, "registeredAt", filter); err != nil {
		return nil, 0, err
	}

	if err := filterDateQueries(query.LastVisitDateFrom, query.LastVisitDateTo, "lastVisitAt", filter); err != nil {
		return nil, 0, err
	}

	cur, err := r.db.Find(ctx, filter, paginationOpts)
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

func (r *StudentsRepo) SetSession(ctx context.Context, studentID primitive.ObjectID, session domain.Session) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": studentID}, bson.M{"$set": bson.M{"session": session, "lastVisitAt": time.Now()}})

	return err
}

func (r *StudentsRepo) GiveAccessToModule(ctx context.Context, studentID, moduleID primitive.ObjectID) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": studentID}, bson.M{"$addToSet": bson.M{"availableModules": moduleID}})

	return err
}

func (r *StudentsRepo) AttachOffer(ctx context.Context, studentID, offerID primitive.ObjectID, moduleIds []primitive.ObjectID) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": studentID}, bson.M{"$addToSet": bson.M{
		"availableModules": bson.M{"$each": moduleIds},
		"availableOffers":  offerID,
	}})

	return err
}

func (r *StudentsRepo) DetachOffer(ctx context.Context, studentID, offerID primitive.ObjectID, moduleIds []primitive.ObjectID) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": studentID}, bson.M{"$pull": bson.M{
		"availableModules": bson.M{"$in": moduleIds},
		"availableOffers":  offerID,
	}})

	return err
}

func (r *StudentsRepo) Verify(ctx context.Context, code string) (domain.Student, error) {
	res := r.db.FindOneAndUpdate(ctx,
		bson.M{"verification.code": code},
		bson.M{"$set": bson.M{"verification.verified": true, "verification.code": ""}})
	if res.Err() != nil {
		return domain.Student{}, res.Err()
	}

	var student domain.Student
	err := res.Decode(&student)

	return student, err
}
