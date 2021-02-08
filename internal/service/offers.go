package service

import (
	"context"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"github.com/zhashkevych/courses-backend/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OffersService struct {
	repo           repository.Offers
	modulesService Modules
}

func NewOffersService(repo repository.Offers, modulesService Modules) *OffersService {
	return &OffersService{repo: repo, modulesService: modulesService}
}

func (s *OffersService) GetById(ctx context.Context, id primitive.ObjectID) (domain.Offer, error) {
	return s.repo.GetById(ctx, id)
}

func (s *OffersService) GetByPackage(ctx context.Context, schoolId, packageId primitive.ObjectID) ([]domain.Offer, error) {
	offers, err := s.repo.GetBySchool(ctx, schoolId)
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

func (s *OffersService) GetByModule(ctx context.Context, schoolId, moduleId primitive.ObjectID) ([]domain.Offer, error) {
	module, err := s.modulesService.GetById(ctx, moduleId)
	if err != nil {
		return nil, err
	}

	return s.GetByPackage(ctx, schoolId, module.PackageID)
}

func inArray(array []primitive.ObjectID, searchedItem primitive.ObjectID) bool {
	for i := range array {
		if array[i] == searchedItem {
			return true
		}
	}
	return false
}
