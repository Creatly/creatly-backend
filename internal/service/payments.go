package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/zhashkevych/creatly-backend/internal/domain"
	"github.com/zhashkevych/creatly-backend/pkg/logger"
	"github.com/zhashkevych/creatly-backend/pkg/payment"
	"github.com/zhashkevych/creatly-backend/pkg/payment/fondy"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentsService struct {
	paymentProvider payment.Provider
	ordersService   Orders
	offersService   Offers
	studentsService Students
	emailService    Emails
}

func NewPaymentsService(paymentProvider payment.Provider, ordersService Orders,
	offersService Offers, studentsService Students, emailService Emails) *PaymentsService {
	return &PaymentsService{
		paymentProvider: paymentProvider,
		ordersService:   ordersService,
		offersService:   offersService,
		studentsService: studentsService,
		emailService:    emailService,
	}
}

func (s *PaymentsService) ProcessTransaction(ctx context.Context, callback interface{}) error {
	switch callbackData := callback.(type) {
	case fondy.Callback:
		return s.processFondyCallback(ctx, callbackData)
	default:
		return domain.ErrUnknownCallbackType
	}
}

func (s *PaymentsService) processFondyCallback(ctx context.Context, callback fondy.Callback) error {
	if err := s.paymentProvider.ValidateCallback(callback); err != nil {
		return domain.ErrTransactionInvalid
	}

	orderId, err := primitive.ObjectIDFromHex(callback.OrderId)
	if err != nil {
		return err
	}

	transaction, err := createTransaction(callback)
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

	offer, err := s.offersService.GetById(ctx, order.Offer.ID)
	if err != nil {
		return err
	}

	if err := s.emailService.SendStudentPurchaseSuccessfulEmail(StudentPurchaseSuccessfulEmailInput{
		Name:       order.Student.Name,
		Email:      order.Student.Email,
		CourseName: order.Offer.Name,
	}); err != nil {
		logger.Errorf("failed to send email after purchase: %s", err.Error())
	}

	return s.studentsService.GiveAccessToOffer(ctx, order.Student.ID, offer)
}

func createTransaction(callbackData fondy.Callback) (domain.Transaction, error) {
	var status string
	if callbackData.PaymentApproved() {
		status = domain.OrderStatusPaid
	} else {
		status = domain.OrderStatusOther
	}

	if !callbackData.Success() {
		status = domain.OrderStatusFailed
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
