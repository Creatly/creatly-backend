package service

import (
	"context"
	"errors"

	"github.com/zhashkevych/creatly-backend/internal/domain"
	"github.com/zhashkevych/creatly-backend/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PromoCodeService struct {
	repo repository.PromoCodes
}

func NewPromoCodeService(repo repository.PromoCodes) *PromoCodeService {
	return &PromoCodeService{repo: repo}
}

func (s *PromoCodeService) Create(ctx context.Context, inp CreatePromoCodeInput) (primitive.ObjectID, error) {
	return s.repo.Create(ctx, domain.PromoCode{
		SchoolID:           inp.SchoolID,
		Code:               inp.Code,
		DiscountPercentage: inp.DiscountPercentage,
		ExpiresAt:          inp.ExpiresAt,
		OfferIDs:           inp.OfferIDs,
	})
}

func (s *PromoCodeService) Update(ctx context.Context, inp domain.UpdatePromoCodeInput) error {
	return s.repo.Update(ctx, inp)
}

func (s *PromoCodeService) Delete(ctx context.Context, schoolId, id primitive.ObjectID) error {
	return s.repo.Delete(ctx, schoolId, id)
}

func (s *PromoCodeService) GetByCode(ctx context.Context, schoolId primitive.ObjectID, code string) (domain.PromoCode, error) {
	promo, err := s.repo.GetByCode(ctx, schoolId, code)
	if err != nil {
		if errors.Is(err, domain.ErrPromoNotFound) {
			return domain.PromoCode{}, err
		}

		return domain.PromoCode{}, err
	}

	return promo, nil
}

func (s *PromoCodeService) GetById(ctx context.Context, schoolId, id primitive.ObjectID) (domain.PromoCode, error) {
	promo, err := s.repo.GetById(ctx, schoolId, id)
	if err != nil {
		if errors.Is(err, domain.ErrPromoNotFound) {
			return domain.PromoCode{}, err
		}

		return domain.PromoCode{}, err
	}

	return promo, nil
}

func (s *PromoCodeService) GetBySchool(ctx context.Context, schoolId primitive.ObjectID) ([]domain.PromoCode, error) {
	return s.repo.GetBySchool(ctx, schoolId)
}
