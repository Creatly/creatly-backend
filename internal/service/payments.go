package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/zhashkevych/creatly-backend/internal/domain"
	"github.com/zhashkevych/creatly-backend/pkg/logger"
	"github.com/zhashkevych/creatly-backend/pkg/payment"
	"github.com/zhashkevych/creatly-backend/pkg/payment/fondy"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	redirectURLTmpl = "https://%s/" // TODO: generate link with URL params for popup on frontend ?
)

type PaymentsService struct {
	ordersService   Orders
	offersService   Offers
	studentsService Students
	emailService    Emails
	schoolsService  Schools

	fondyCallbackURL string
}

func NewPaymentsService(ordersService Orders, offersService Offers, studentsService Students,
	emailService Emails, schoolsService Schools, fondyCallbackURL string) *PaymentsService {
	return &PaymentsService{
		ordersService:    ordersService,
		offersService:    offersService,
		studentsService:  studentsService,
		emailService:     emailService,
		schoolsService:   schoolsService,
		fondyCallbackURL: fondyCallbackURL,
	}
}

func (s *PaymentsService) GeneratePaymentLink(ctx context.Context, orderId primitive.ObjectID) (string, error) {
	order, err := s.ordersService.GetById(ctx, orderId)
	if err != nil {
		return "", err
	}

	offer, err := s.offersService.GetById(ctx, order.Offer.ID)
	if err != nil {
		return "", err
	}

	if !offer.PaymentMethod.UsesProvider {
		return "", domain.ErrPaymentProviderNotUsed
	}

	paymentInput := payment.GeneratePaymentLinkInput{
		OrderId:   orderId.Hex(),
		Amount:    order.Amount,
		Currency:  offer.Price.Currency,
		OrderDesc: offer.Description, // TODO proper order description
	}

	switch offer.PaymentMethod.Provider {
	case domain.PaymentProviderFondy:
		return s.generateFondyPaymentLink(ctx, offer.SchoolID, paymentInput)
	default:
		return "", domain.ErrUnknownPaymentProvider
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
	orderID, err := primitive.ObjectIDFromHex(callback.OrderId)
	if err != nil {
		return err
	}

	order, err := s.ordersService.GetById(ctx, orderID)
	if err != nil {
		return err
	}

	school, err := s.schoolsService.GetById(ctx, order.SchoolID)
	if err != nil {
		return err
	}

	client, err := s.getFondyClient(school.Settings.Fondy)
	if err != nil {
		return err
	}

	if err := client.ValidateCallback(callback); err != nil {
		return domain.ErrTransactionInvalid
	}

	transaction, err := createTransaction(callback)
	if err != nil {
		return err
	}

	order, err = s.ordersService.AddTransaction(ctx, orderID, transaction)
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

func (s *PaymentsService) generateFondyPaymentLink(ctx context.Context, schoolId primitive.ObjectID,
	input payment.GeneratePaymentLinkInput) (string, error) {
	school, err := s.schoolsService.GetById(ctx, schoolId)
	if err != nil {
		return "", err
	}

	client, err := s.getFondyClient(school.Settings.Fondy)
	if err != nil {
		return "", err
	}

	input.CallbackURL = s.fondyCallbackURL
	input.RedirectURL = getRedirectURL(school.Settings.GetDomain())

	logger.Infof("%+v", input)

	return client.GeneratePaymentLink(input)
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

func (s *PaymentsService) getFondyClient(fondyConnectionInfo domain.Fondy) (*fondy.Client, error) {
	if !fondyConnectionInfo.Connected {
		return nil, domain.ErrFondyIsNotConnected
	}

	return fondy.NewFondyClient(fondyConnectionInfo.MerchantID, fondyConnectionInfo.MerchantPassword), nil
}

func getRedirectURL(domain string) string {
	return fmt.Sprintf(redirectURLTmpl, domain)
}
