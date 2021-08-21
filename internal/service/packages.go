package service

import (
	"context"

	"github.com/zhashkevych/creatly-backend/internal/domain"
	"github.com/zhashkevych/creatly-backend/internal/repository"
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

func (s *PackagesService) GetByCourse(ctx context.Context, courseID primitive.ObjectID) ([]domain.Package, error) {
	pkgs, err := s.repo.GetByCourse(ctx, courseID)
	if err != nil {
		return nil, err
	}

	for i := range pkgs {
		modules, err := s.modulesRepo.GetByPackages(ctx, []primitive.ObjectID{pkgs[i].ID})
		if err != nil {
			return nil, err
		}

		pkgs[i].Modules = modules
	}

	return pkgs, nil
}

func (s *PackagesService) GetById(ctx context.Context, id primitive.ObjectID) (domain.Package, error) {
	pkg, err := s.repo.GetById(ctx, id)
	if err != nil {
		return pkg, err
	}

	modules, err := s.modulesRepo.GetByPackages(ctx, []primitive.ObjectID{pkg.ID})
	if err != nil {
		return pkg, err
	}

	pkg.Modules = modules

	return pkg, nil
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
