package service

import (
	"context"
	"encoding/json"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"github.com/zhashkevych/courses-backend/pkg/payment"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type PaymentsService struct {
	paymentProvider payment.FondyProvider
	ordersService   Orders
	offersService   Offers
	studentsService Students
}

func NewPaymentsService(paymentProvider payment.FondyProvider, ordersService Orders, offersService Offers, studentsService Students) *PaymentsService {
	return &PaymentsService{paymentProvider: paymentProvider, ordersService: ordersService, offersService: offersService, studentsService: studentsService}
}

// TODO callback data validation?
func (s *PaymentsService) ProcessTransaction(ctx context.Context, callbackData payment.Callback) error {
	orderId, err := primitive.ObjectIDFromHex(callbackData.OrderId)
	if err != nil {
		return err
	}

	transaction, err := createTransaction(callbackData)
	if err != nil {
		return err
	}

	order, err := s.ordersService.AddTransaction(ctx, orderId, transaction)
	if err != nil {
		return err
	}

	if transaction.Status != domain.OrderStatusPaid {
		return nil
	}

	offer, err := s.offersService.GetById(ctx, order.OfferID)
	if err != nil {
		return err
	}

	return s.studentsService.GiveAccessToPackages(ctx, order.StudentID, offer.PackageIDs)
}

func createTransaction(callbackData payment.Callback) (domain.Transaction, error) {
	var status string
	if !callbackData.Success() {
		status = domain.OrderStatusFailed
	}

	if callbackData.PaymentApproved() {
		status = domain.OrderStatusPaid
	} else {
		status = domain.OrderStatusOther
	}

	additionalInfo, err := json.Marshal(callbackData)
	if err != nil {
		return domain.Transaction{}, err
	}

	return domain.Transaction{
		Status:         status,
		CreatedAt:      time.Now(),
		AdditionalInfo: string(additionalInfo),
	}, nil
}
