package tests

import (
	"github.com/zhashkevych/courses-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	school = domain.School{
		ID: primitive.NewObjectID(),
		Courses: []domain.Course{
			{
				ID:        primitive.NewObjectID(),
				Name:      "Course #1",
				Published: true,
			},
			{
				ID:   primitive.NewObjectID(),
				Name: "Course #2", // Unpublished course, shouldn't be available to student
			},
		},
	}

	packages = []interface{}{
		domain.Package{
			ID:       primitive.NewObjectID(),
			Name:     "Package #1",
			CourseID: school.Courses[0].ID,
		},
	}

	offers = []interface{}{
		domain.Offer{
			ID:          primitive.NewObjectID(),
			Name:        "Offer #1",
			Description: "Offer #1 Description",
			SchoolID:    school.ID,
			PackageIDs:  []primitive.ObjectID{packages[0].(domain.Package).ID},
			Price:       domain.Price{Value: 6900, Currency: "USD"},
		},
	}

	modules = []interface{}{
		domain.Module{
			ID:        primitive.NewObjectID(),
			Name:      "Module #1", // Free Module, should be available to anyone
			CourseID:  school.Courses[0].ID,
			Published: true,
			Lessons: []domain.Lesson{
				{
					ID: primitive.NewObjectID(),
					Name: "Lesson #1",
					Published: true,
				},
			},
		},
		domain.Module{
			ID:        primitive.NewObjectID(),
			Name:      "Module #2", // Part of paid package, should be available only after purchase
			CourseID:  school.Courses[0].ID,
			Published: true,
			PackageID: packages[0].(domain.Package).ID,
			Lessons: []domain.Lesson{
				{
					ID: primitive.NewObjectID(),
					Name: "Lesson #1",
					Published: true,
				},
				{
					ID: primitive.NewObjectID(),
					Name: "Lesson #2",
					Published: true,
				},
			},
		},
	}
)
