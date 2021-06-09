package packages

import (
	"context"

	"github.com/zhashkevych/creatly-backend/internal/domain"
	"github.com/zhashkevych/creatly-backend/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	CreatePackageInput struct {
		CourseID    string
		SchoolID    string
		Name        string
		Description string
	}

	UpdatePackageInput struct {
		ID          string
		SchoolID    string
		Name        string
		Description string
		Modules     []string
	}

	Packages interface {
		Create(ctx context.Context, pkg domain.Package) (primitive.ObjectID, error)
		Update(ctx context.Context, inp repository.UpdatePackageInput) error
		Delete(ctx context.Context, schoolId, id primitive.ObjectID) error
		GetByCourse(ctx context.Context, courseId primitive.ObjectID) ([]domain.Package, error)
		GetById(ctx context.Context, id primitive.ObjectID) (domain.Package, error)
	}

	Modules interface {
		AttachPackage(ctx context.Context, schoolId, packageId primitive.ObjectID, modules []primitive.ObjectID) error
	}

	PackagesService struct {
		repo        Packages
		modulesRepo Modules
	}
)

func NewPackagesService(repo Packages, modulesRepo Modules) *PackagesService {
	return &PackagesService{repo: repo, modulesRepo: modulesRepo}
}

func (s *PackagesService) Create(ctx context.Context, inp CreatePackageInput) (primitive.ObjectID, error) {
	id, err := primitive.ObjectIDFromHex(inp.CourseID)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	schoolId, err := primitive.ObjectIDFromHex(inp.SchoolID)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	return s.repo.Create(ctx, domain.Package{
		CourseID:    id,
		SchoolID:    schoolId,
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

	schoolId, err := primitive.ObjectIDFromHex(inp.SchoolID)
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

		if err := s.modulesRepo.AttachPackage(ctx, schoolId, id, moduleIds); err != nil {
			return err
		}
	}

	return nil
}

func (s *PackagesService) Delete(ctx context.Context, schoolId, id primitive.ObjectID) error {
	return s.repo.Delete(ctx, schoolId, id)
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
