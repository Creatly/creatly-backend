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

func (s *OffersService) Create(ctx context.Context, inp CreateOfferInput) (primitive.ObjectID, error) {
	return s.repo.Create(ctx, domain.Offer{
		SchoolID:    inp.SchoolID,
		Name:        inp.Name,
		Description: inp.Description,
		Price:       inp.Price,
	})
}

func (s *OffersService) GetAll(ctx context.Context, schoolId primitive.ObjectID) ([]domain.Offer, error) {
	return s.repo.GetBySchool(ctx, schoolId)
}

func (s *OffersService) Update(ctx context.Context, inp UpdateOfferInput) error {
	id, err := primitive.ObjectIDFromHex(inp.ID)
	if err != nil {
		return err
	}

	updateInput := repository.UpdateOfferInput{
		ID:          id,
		Name:        inp.Name,
		Description: inp.Description,
		Price:       inp.Price,
	}

	if inp.Packages != nil {
		updateInput.Packages, err = stringArrayToObjectId(inp.Packages)
		if err != nil {
			return err
		}
	}

	return s.repo.Update(ctx, updateInput)
}

func (s *OffersService) Delete(ctx context.Context, id primitive.ObjectID) error {
	return s.repo.Delete(ctx, id)
}

func inArray(array []primitive.ObjectID, searchedItem primitive.ObjectID) bool {
	for i := range array {
		if array[i] == searchedItem {
			return true
		}
	}
	return false
}
