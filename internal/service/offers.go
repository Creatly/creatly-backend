package service

import (
	"context"

	"github.com/zhashkevych/creatly-backend/internal/domain"
	"github.com/zhashkevych/creatly-backend/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OffersService struct {
	repo            repository.Offers
	modulesService  Modules
	packagesService Packages
}

func NewOffersService(repo repository.Offers, modulesService Modules, packagesService Packages) *OffersService {
	return &OffersService{repo: repo, modulesService: modulesService, packagesService: packagesService}
}

func (s *OffersService) GetById(ctx context.Context, id primitive.ObjectID) (domain.Offer, error) {
	return s.repo.GetById(ctx, id)
}

func (s *OffersService) getByPackage(ctx context.Context, schoolId, packageId primitive.ObjectID) ([]domain.Offer, error) {
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

	return s.getByPackage(ctx, schoolId, module.PackageID)
}

func (s *OffersService) GetByCourse(ctx context.Context, courseId primitive.ObjectID) ([]domain.Offer, error) {
	packages, err := s.packagesService.GetByCourse(ctx, courseId)
	if err != nil {
		return nil, err
	}

	if len(packages) == 0 {
		return []domain.Offer{}, nil
	}

	packageIds := make([]primitive.ObjectID, len(packages))
	for i, pkg := range packages {
		packageIds[i] = pkg.ID
	}

	return s.repo.GetByPackages(ctx, packageIds)
}

func (s *OffersService) Create(ctx context.Context, inp CreateOfferInput) (primitive.ObjectID, error) {
	if inp.PaymentMethod.UsesProvider {
		if err := inp.PaymentMethod.Validate(); err != nil {
			return primitive.ObjectID{}, err
		}
	}

	var (
		packageIDs []primitive.ObjectID
		err        error
	)

	if inp.Packages != nil {
		packageIDs, err = stringArrayToObjectId(inp.Packages)
		if err != nil {
			return primitive.ObjectID{}, err
		}
	}

	return s.repo.Create(ctx, domain.Offer{
		SchoolID:      inp.SchoolID,
		Name:          inp.Name,
		Description:   inp.Description,
		Benefits:      inp.Benefits,
		Price:         inp.Price,
		PaymentMethod: inp.PaymentMethod,
		PackageIDs:    packageIDs,
	})
}

func (s *OffersService) GetAll(ctx context.Context, schoolId primitive.ObjectID) ([]domain.Offer, error) {
	return s.repo.GetBySchool(ctx, schoolId)
}

func (s *OffersService) Update(ctx context.Context, inp UpdateOfferInput) error {
	if err := inp.ValidatePayment(); err != nil {
		return err
	}

	id, err := primitive.ObjectIDFromHex(inp.ID)
	if err != nil {
		return err
	}

	schoolId, err := primitive.ObjectIDFromHex(inp.SchoolID)
	if err != nil {
		return err
	}

	updateInput := repository.UpdateOfferInput{
		ID:            id,
		SchoolID:      schoolId,
		Name:          inp.Name,
		Description:   inp.Description,
		Price:         inp.Price,
		Benefits:      inp.Benefits,
		PaymentMethod: inp.PaymentMethod,
	}

	if inp.Packages != nil {
		updateInput.Packages, err = stringArrayToObjectId(inp.Packages)
		if err != nil {
			return err
		}
	}

	return s.repo.Update(ctx, updateInput)
}

func (s *OffersService) Delete(ctx context.Context, schoolId, id primitive.ObjectID) error {
	return s.repo.Delete(ctx, schoolId, id)
}

func (s *OffersService) GetByIds(ctx context.Context, ids []primitive.ObjectID) ([]domain.Offer, error) {
	return s.repo.GetByIds(ctx, ids)
}

func inArray(array []primitive.ObjectID, searchedItem primitive.ObjectID) bool {
	for i := range array {
		if array[i] == searchedItem {
			return true
		}
	}

	return false
}
