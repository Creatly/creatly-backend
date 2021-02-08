package service

import (
	"context"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"github.com/zhashkevych/courses-backend/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PackagesService struct {
	repo repository.Packages
}

func NewPackagesService(repo repository.Packages) *PackagesService {
	return &PackagesService{repo: repo}
}

func (s *PackagesService) Create(ctx context.Context, inp CreatePackageInput) (primitive.ObjectID, error) {
	id, err := primitive.ObjectIDFromHex(inp.CourseID)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	return s.repo.Create(ctx, domain.Package{
		CourseID:    id,
		Name:        inp.Name,
		Description: inp.Description,
	})
}

func (s *PackagesService) GetByCourse(ctx context.Context, courseId primitive.ObjectID) ([]domain.Package, error) {
	return s.repo.GetByCourse(ctx, courseId)
}