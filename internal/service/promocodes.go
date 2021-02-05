package service

import (
	"context"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"github.com/zhashkevych/courses-backend/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PromoCodeService struct {
	repo repository.Promocodes
}

func NewPromoCodeService(repo repository.Promocodes) *PromoCodeService {
	return &PromoCodeService{repo: repo}
}

func (s *PromoCodeService) GetByCode(ctx context.Context, schoolId primitive.ObjectID, code string) (domain.PromoCode, error) {
	return s.repo.GetByCode(ctx, schoolId, code)
}

func (s *PromoCodeService) GetById(ctx context.Context, id primitive.ObjectID) (domain.PromoCode, error) {
	return s.repo.GetById(ctx, id)
}
