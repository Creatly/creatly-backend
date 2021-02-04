package service

import (
	"context"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"github.com/zhashkevych/courses-backend/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type CoursesService struct {
	repo       repository.Courses
	offersRepo repository.Offers
	promoRepo  repository.Promocodes
}

func NewCoursesService(repo repository.Courses, offersRepo repository.Offers, promoRepo repository.Promocodes) *CoursesService {
	return &CoursesService{repo: repo, offersRepo: offersRepo, promoRepo: promoRepo}
}

func (s *CoursesService) GetCourseModules(ctx context.Context, courseId primitive.ObjectID) ([]domain.Module, error) {
	modules, err := s.repo.GetModules(ctx, courseId)
	if err != nil {
		return nil, err
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
	offers, err := s.offersRepo.GetBySchool(ctx, schoolId)
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

func (s *CoursesService) GetPackagesModules(ctx context.Context, packageIds []primitive.ObjectID) ([]domain.Module, error) {
	return s.repo.GetPackagesModules(ctx, packageIds)
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

func (s *CoursesService) GetPromocodeByCode(ctx context.Context, schoolId primitive.ObjectID, code string) (domain.Promocode, error) {
	return s.promoRepo.GetByCode(ctx, schoolId, code)
}

func (s *CoursesService) GetPromocodeById(ctx context.Context, id primitive.ObjectID) (domain.Promocode, error) {
	return s.promoRepo.GetById(ctx, id)
}

func (s *CoursesService) GetOfferById(ctx context.Context, id primitive.ObjectID) (domain.Offer, error) {
	return s.offersRepo.GetById(ctx, id)
}

func (s *CoursesService) Create(ctx context.Context, schoolId primitive.ObjectID, name string) (primitive.ObjectID, error) {
	return s.repo.Create(ctx, schoolId, domain.Course{
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
}
