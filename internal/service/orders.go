package service

import (
	"context"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"github.com/zhashkevych/courses-backend/internal/repository"
	"github.com/zhashkevych/courses-backend/pkg/logger"
	"github.com/zhashkevych/courses-backend/pkg/payment"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type OrdersService struct {
	offersService     Offers
	promoCodesService PromoCodes

	repo            repository.Orders
	paymentProvider payment.FondyProvider

	callbackURL, responseURL string
}

func NewOrdersService(repo repository.Orders, offersService Offers, promoCodesService PromoCodes, paymentProvider payment.FondyProvider, callbackURL, responseURL string) *OrdersService {
	return &OrdersService{
		repo:              repo,
		offersService:     offersService,
		promoCodesService: promoCodesService,
		paymentProvider:   paymentProvider,
		callbackURL:       callbackURL,
		responseURL:       responseURL,
	}
}

func (s *OrdersService) Create(ctx context.Context, studentId, offerId, promocodeId primitive.ObjectID) (string, error) {
	promocode, err := s.getOrderPromocode(ctx, promocodeId)
	if err != nil {
		return "", err
	}

	offer, err := s.offersService.GetById(ctx, offerId)
	if err != nil {
		return "", err
	}

	orderAmount := s.calculateOrderPrice(offer.Price.Value, promocode)

	id := primitive.NewObjectID()
	paymentLink, err := s.paymentProvider.GeneratePaymentLink(payment.GeneratePaymentLinkInput{
		OrderId:     id.Hex(),
		Amount:      orderAmount,
		Currency:    offer.Price.Currency,
		OrderDesc:   offer.Description, // TODO proper order description
		CallbackURL: s.callbackURL,
		ResponseURL: s.responseURL,
	})
	if err != nil {
		logger.Error("Failed to generate payment link: ", err.Error())
		return "", err
	}

	if err := s.repo.Create(ctx, domain.Order{
		ID:           id,
		StudentID:    studentId,
		OfferID:      offerId,
		PromoID:      promocodeId,
		Amount:       orderAmount,
		CreatedAt:    time.Now(),
		Status:       domain.OrderStatusCreated,
		Transactions: make([]domain.Transaction, 0),
	}); err != nil {
		return "", err
	}

	return paymentLink, nil
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
		promocode, err = s.promoCodesService.GetById(ctx, promocodeId)
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
