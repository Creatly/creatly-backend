package service

import (
	"context"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"github.com/zhashkevych/courses-backend/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PromoCodeService struct {
	repo repository.PromoCodes
}

func NewPromoCodeService(repo repository.PromoCodes) *PromoCodeService {
	return &PromoCodeService{repo: repo}
}

func (s *PromoCodeService) GetByCode(ctx context.Context, schoolId primitive.ObjectID, code string) (domain.PromoCode, error) {
	promo, err := s.repo.GetByCode(ctx, schoolId, code)
	if err != nil {
		if err == repository.ErrPromoNotFound {
			return domain.PromoCode{}, ErrPromoNotFound
		}

		return domain.PromoCode{}, err
	}

	return promo, nil
}

func (s *PromoCodeService) GetById(ctx context.Context, schoolId, id primitive.ObjectID) (domain.PromoCode, error) {
	promo, err := s.repo.GetById(ctx, schoolId, id)
	if err != nil {
		if err == repository.ErrPromoNotFound {
			return domain.PromoCode{}, ErrPromoNotFound
		}

		return domain.PromoCode{}, err
	}

	return promo, nil
}
