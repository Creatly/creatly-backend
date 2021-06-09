package promocodes

import (
	"context"
	"errors"
	"time"

	"github.com/zhashkevych/creatly-backend/internal/domain"
	"github.com/zhashkevych/creatly-backend/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrUserNotFound            = errors.New("user doesn't exists")
	ErrOfferNotFound           = errors.New("offer doesn't exists")
	ErrPromoNotFound           = errors.New("promocode doesn't exists")
	ErrModuleIsNotAvailable    = errors.New("module's content is not available")
	ErrPromocodeExpired        = errors.New("promocode has expired")
	ErrTransactionInvalid      = errors.New("transaction is invalid")
	ErrUnknownCallbackType     = errors.New("unknown callback type")
	ErrVerificationCodeInvalid = errors.New("verification code is invalid")
	ErrUserAlreadyExists       = errors.New("user with such email already exists")
)

type (
	PromoCodes interface {
		Create(ctx context.Context, promocode domain.PromoCode) (primitive.ObjectID, error)
		Update(ctx context.Context, inp repository.UpdatePromoCodeInput) error
		Delete(ctx context.Context, schoolId, id primitive.ObjectID) error
		GetByCode(ctx context.Context, schoolId primitive.ObjectID, code string) (domain.PromoCode, error)
		GetById(ctx context.Context, schoolId, id primitive.ObjectID) (domain.PromoCode, error)
		GetBySchool(ctx context.Context, schoolId primitive.ObjectID) ([]domain.PromoCode, error)
	}

	PromoCodeService struct {
		repo PromoCodes
	}

	CreatePromoCodeInput struct {
		SchoolID           primitive.ObjectID
		Code               string
		DiscountPercentage int
		ExpiresAt          time.Time
		OfferIDs           []primitive.ObjectID
	}

	UpdatePromoCodeInput struct {
		ID                 primitive.ObjectID
		SchoolID           primitive.ObjectID
		Code               string
		DiscountPercentage int
		ExpiresAt          time.Time
		OfferIDs           []string
	}
)

func NewPromoCodeService(repo PromoCodes) *PromoCodeService {
	return &PromoCodeService{repo: repo}
}

func (s *PromoCodeService) Create(ctx context.Context, inp CreatePromoCodeInput) (primitive.ObjectID, error) {
	return s.repo.Create(ctx, domain.PromoCode{
		SchoolId:           inp.SchoolID,
		Code:               inp.Code,
		DiscountPercentage: inp.DiscountPercentage,
		ExpiresAt:          inp.ExpiresAt,
		OfferIDs:           inp.OfferIDs,
	})
}

func (s *PromoCodeService) Update(ctx context.Context, inp UpdatePromoCodeInput) error {
	updateInput := repository.UpdatePromoCodeInput{
		ID:                 inp.ID,
		SchoolID:           inp.SchoolID,
		Code:               inp.Code,
		DiscountPercentage: inp.DiscountPercentage,
		ExpiresAt:          inp.ExpiresAt,
	}

	if inp.OfferIDs != nil {
		var err error

		updateInput.OfferIDs, err = stringArrayToObjectId(inp.OfferIDs)
		if err != nil {
			return err
		}
	}

	return s.repo.Update(ctx, updateInput)
}

func (s *PromoCodeService) Delete(ctx context.Context, schoolId, id primitive.ObjectID) error {
	return s.repo.Delete(ctx, schoolId, id)
}

func (s *PromoCodeService) GetByCode(ctx context.Context, schoolId primitive.ObjectID, code string) (domain.PromoCode, error) {
	promo, err := s.repo.GetByCode(ctx, schoolId, code)
	if err != nil {
		if errors.Is(err, repository.ErrPromoNotFound) {
			return domain.PromoCode{}, ErrPromoNotFound
		}

		return domain.PromoCode{}, err
	}

	return promo, nil
}

func (s *PromoCodeService) GetById(ctx context.Context, schoolId, id primitive.ObjectID) (domain.PromoCode, error) {
	promo, err := s.repo.GetById(ctx, schoolId, id)
	if err != nil {
		if errors.Is(err, repository.ErrPromoNotFound) {
			return domain.PromoCode{}, ErrPromoNotFound
		}

		return domain.PromoCode{}, err
	}

	return promo, nil
}

func (s *PromoCodeService) GetBySchool(ctx context.Context, schoolId primitive.ObjectID) ([]domain.PromoCode, error) {
	return s.repo.GetBySchool(ctx, schoolId)
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
