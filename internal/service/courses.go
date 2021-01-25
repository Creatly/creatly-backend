package service

import (
	"context"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"github.com/zhashkevych/courses-backend/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CoursesService struct {
	repo       repository.Courses
	offersRepo repository.Offers
}

func NewCoursesService(repo repository.Courses, offersRepo repository.Offers) *CoursesService {
	return &CoursesService{repo: repo, offersRepo: offersRepo}
}

func (s *CoursesService) GetCourseModules(ctx context.Context, courseId primitive.ObjectID) ([]domain.Module, error) {
	modules, err := s.repo.GetModules(ctx, courseId)
	if err != nil {
		return nil, err
	}

	if len(modules) == 0 {
		return nil, ErrCourseContentNotFound
	}

	return modules, nil
}

func (s *CoursesService) GetModule(ctx context.Context, moduleId primitive.ObjectID) (domain.Module, error) {
	return s.repo.GetModule(ctx, moduleId)
}

func (s *CoursesService) GetModuleWithContent(ctx context.Context, moduleId primitive.ObjectID) (domain.Module, error) {
	return s.repo.GetModuleWithContent(ctx, moduleId)
}

func (s *CoursesService) GetPackageOffers(ctx context.Context, schoolId, packageId primitive.ObjectID) ([]domain.Offer, error) {
	offers, err := s.offersRepo.GetSchoolOffers(ctx, schoolId)
	if err != nil {
		return nil, err
	}

	result := make([]domain.Offer, 0)
	for _, offer := range offers {
		if inArray(offer.PackageIDs, packageId) {
			result = append(result, offer)
		}
	}

	return result, nil
}

func inArray(array []primitive.ObjectID, searchedItem primitive.ObjectID) bool {
	for i := range array {
		if array[i] == searchedItem {
			return true
		}
	}
	return false
}

func (s *CoursesService) GetModuleOffers(ctx context.Context, schoolId, moduleId primitive.ObjectID) ([]domain.Offer, error) {
	module, err := s.GetModule(ctx, moduleId)
	if err != nil {
		return nil, err
	}

	return s.GetPackageOffers(ctx, schoolId, module.PackageID)
}
