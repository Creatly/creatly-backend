package service

import (
	"context"
	"encoding/json"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"github.com/zhashkevych/courses-backend/pkg/logger"
	"github.com/zhashkevych/courses-backend/pkg/payment"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type PaymentsService struct {
	paymentProvider payment.FondyProvider
	ordersService   Orders
}

func NewPaymentsService(paymentProvider payment.FondyProvider, ordersService Orders) *PaymentsService {
	return &PaymentsService{paymentProvider: paymentProvider, ordersService: ordersService}
}

// TODO groom code, add modules, decide with callback validation
func (s *PaymentsService) ProcessTransaction(ctx context.Context, callbackData payment.Callback) error {
	orderId, err := primitive.ObjectIDFromHex(callbackData.OrderId)
	if err != nil {
		return err
	}

	transaction, err := createTransaction(callbackData)
	if err != nil {
		return err
	}

	if err := s.ordersService.AddTransaction(ctx, orderId, transaction); err != nil {
		return err
	}

	if transaction.Status == domain.OrderStatusPaid {
		logger.Info("TOOD: Add modules")
	}

	return nil
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
