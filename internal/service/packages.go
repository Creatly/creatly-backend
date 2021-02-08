package service

import (
	"context"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"github.com/zhashkevych/courses-backend/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PackagesService struct {
	repo        repository.Packages
	modulesRepo repository.Modules
}

func NewPackagesService(repo repository.Packages, modulesRepo repository.Modules) *PackagesService {
	return &PackagesService{repo: repo, modulesRepo: modulesRepo}
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

func (s *PackagesService) GetById(ctx context.Context, id primitive.ObjectID) (domain.Package, error) {
	return s.repo.GetById(ctx, id)
}

func (s *PackagesService) Update(ctx context.Context, inp UpdatePackageInput) error {
	id, err := primitive.ObjectIDFromHex(inp.ID)
	if err != nil {
		return err
	}

	if inp.Name != "" || inp.Description != "" {
		if err := s.repo.Update(ctx, repository.UpdatePackageInput{
			ID:          id,
			Name:        inp.Name,
			Description: inp.Description,
		}); err != nil {
			return err
		}
	}

	if inp.Modules != nil {
		moduleIds, err := stringArrayToObjectId(inp.Modules)
		if err != nil {
			return err
		}

		if err := s.modulesRepo.AttachPackage(ctx, moduleIds, id); err != nil {
			return err
		}
	}

	return nil
}

func (s *PackagesService) Delete(ctx context.Context, id primitive.ObjectID) error {
	return s.repo.Delete(ctx, id)
}

func stringArrayToObjectId(stringIds []string) ([]primitive.ObjectID, error) {
	var err error
	ids := make([]primitive.ObjectID, len(stringIds))

	for i, id := range stringIds {
		ids[i], err = primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
	}

	return ids, nil
}
