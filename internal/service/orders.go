package service

import (
	"context"
	"time"

	"github.com/zhashkevych/creatly-backend/internal/domain"
	"github.com/zhashkevych/creatly-backend/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrdersService struct {
	offersService     Offers
	promoCodesService PromoCodes
	studentsService   Students

	repo repository.Orders
}

func NewOrdersService(repo repository.Orders, offersService Offers, promoCodesService PromoCodes, studentsService Students) *OrdersService {
	return &OrdersService{
		repo:              repo,
		offersService:     offersService,
		promoCodesService: promoCodesService,
		studentsService:   studentsService,
	}
}

func (s *OrdersService) Create(ctx context.Context, studentId, offerId, promocodeId primitive.ObjectID) (primitive.ObjectID, error) { //nolint:funlen
	offer, err := s.offersService.GetById(ctx, offerId)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	promocode, err := s.getOrderPromocode(ctx, offer.SchoolID, promocodeId)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	student, err := s.studentsService.GetById(ctx, offer.SchoolID, studentId)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	orderAmount := s.calculateOrderPrice(offer.Price.Value, promocode)

	id := primitive.NewObjectID()

	order := domain.Order{
		ID:       id,
		SchoolID: offer.SchoolID,
		Student: domain.StudentInfoShort{
			ID:    student.ID,
			Name:  student.Name,
			Email: student.Email,
		},
		Offer: domain.OrderOfferInfo{
			ID:   offer.ID,
			Name: offer.Name,
		},
		Amount:       orderAmount,
		Currency:     offer.Price.Currency,
		CreatedAt:    time.Now(),
		Status:       domain.OrderStatusCreated,
		Transactions: make([]domain.Transaction, 0),
	}

	if !promocode.ID.IsZero() {
		order.Promo = domain.OrderPromoInfo{
			ID:   promocode.ID,
			Code: promocode.Code,
		}
	}

	err = s.repo.Create(ctx, order)

	return id, err
}

func (s *OrdersService) AddTransaction(ctx context.Context, id primitive.ObjectID, transaction domain.Transaction) (domain.Order, error) {
	return s.repo.AddTransaction(ctx, id, transaction)
}

func (s *OrdersService) GetBySchool(ctx context.Context, schoolId primitive.ObjectID, query domain.GetOrdersQuery) ([]domain.Order, int64, error) {
	return s.repo.GetBySchool(ctx, schoolId, query)
}

func (s *OrdersService) GetById(ctx context.Context, id primitive.ObjectID) (domain.Order, error) {
	return s.repo.GetById(ctx, id)
}

func (s *OrdersService) SetStatus(ctx context.Context, id primitive.ObjectID, status string) error {
	return s.repo.SetStatus(ctx, id, status)
}

func (s *OrdersService) getOrderPromocode(ctx context.Context, schoolId, promocodeId primitive.ObjectID) (domain.PromoCode, error) {
	var (
		promocode domain.PromoCode
		err       error
	)

	if !promocodeId.IsZero() {
		promocode, err = s.promoCodesService.GetById(ctx, schoolId, promocodeId)
		if err != nil {
			return promocode, err
		}

		if promocode.ExpiresAt.Unix() < time.Now().Unix() {
			return promocode, domain.ErrPromocodeExpired
		}
	}

	return promocode, nil
}

func (s *OrdersService) calculateOrderPrice(price uint, promocode domain.PromoCode) uint {
	if promocode.ID.IsZero() {
		return price
	} else {
		return (price * uint(100-promocode.DiscountPercentage)) / 100
	}
}
