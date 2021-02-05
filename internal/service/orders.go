package service

import (
	"context"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"github.com/zhashkevych/courses-backend/internal/repository"
	"github.com/zhashkevych/courses-backend/pkg/payment"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type OrdersService struct {
	coursesService  Courses
	repo            repository.Orders
	paymentProvider payment.FondyProvider

	callbackURL, responseURL string
}

func NewOrdersService(repo repository.Orders, coursesService Courses, paymentProvider payment.FondyProvider, callbackURL, responseURL string) *OrdersService {
	return &OrdersService{
		repo:            repo,
		coursesService:  coursesService,
		paymentProvider: paymentProvider,
		callbackURL:     callbackURL,
		responseURL:     responseURL,
	}
}

func (s *OrdersService) Create(ctx context.Context, studentId, offerId, promocodeId primitive.ObjectID) (string, error) {
	promocode, err := s.getOrderPromocode(ctx, promocodeId)
	if err != nil {
		return "", err
	}

	offer, err := s.coursesService.GetOfferById(ctx, offerId)
	if err != nil {
		return "", err
	}

	orderAmount := s.calculateOrderPrice(offer.Price.Value, promocode)

	id := primitive.NewObjectID()
	if err := s.repo.Create(ctx, domain.Order{
		ID:           id,
		StudentID:    studentId,
		OfferID:      offerId,
		PromoID:      promocodeId,
		Amount:       orderAmount,
		Status:       domain.OrderStatusCreated,
		Transactions: make([]domain.Transaction, 0),
	}); err != nil {
		return "", err
	}

	// TODO what if it fails?
	return s.paymentProvider.GeneratePaymentLink(payment.GeneratePaymentLinkInput{
		OrderId:     id.Hex(),
		Amount:      orderAmount,
		Currency:    offer.Price.Currency,
		OrderDesc:   offer.Description, // TODO proper order description
		CallbackURL: s.callbackURL,
		ResponseURL: s.responseURL,
	})
}

func (s *OrdersService) AddTransaction(ctx context.Context, id primitive.ObjectID, transaction domain.Transaction) (domain.Order, error) {
	return s.repo.AddTransaction(ctx, id, transaction)
}

func (s *OrdersService) getOrderPromocode(ctx context.Context, promocodeId primitive.ObjectID) (domain.PromoCode, error) {
	var (
		promocode domain.PromoCode
		err       error
	)

	if !promocodeId.IsZero() {
		promocode, err = s.coursesService.GetPromoById(ctx, promocodeId)
		if err != nil {
			return promocode, err
		}

		if promocode.ExpiresAt.Unix() < time.Now().Unix() {
			return promocode, ErrPromocodeExpired
		}
	}

	return promocode, nil
}

func (s *OrdersService) calculateOrderPrice(price int, promocode domain.PromoCode) int {
	if promocode.ID.IsZero() {
		return price
	} else {
		return (price * (100 - promocode.DiscountPercentage)) / 100
	}
}
