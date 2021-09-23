package service

import (
	"context"
	"fmt"

	"github.com/zhashkevych/creatly-backend/internal/config"
	"github.com/zhashkevych/creatly-backend/internal/domain"
	"github.com/zhashkevych/creatly-backend/pkg/cache"
	emailProvider "github.com/zhashkevych/creatly-backend/pkg/email"
	"github.com/zhashkevych/creatly-backend/pkg/email/sendpulse"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	verificationLinkTmpl = "https://%s/verification?code=%s" // https://<school host>/verification?code=<verification_code>
)

type EmailService struct {
	sender  emailProvider.Sender
	config  config.EmailConfig
	schools SchoolsService

	cache cache.Cache

	sendpulseClients map[primitive.ObjectID]*sendpulse.Client
}

// Structures used for templates.
type verificationEmailInput struct {
	VerificationLink string
}

type purchaseSuccessfulEmailInput struct {
	Name       string
	CourseName string
}

func NewEmailsService(sender emailProvider.Sender, config config.EmailConfig, schools SchoolsService, cache cache.Cache) *EmailService {
	return &EmailService{
		sender:           sender,
		config:           config,
		schools:          schools,
		cache:            cache,
		sendpulseClients: make(map[primitive.ObjectID]*sendpulse.Client),
	}
}

func (s *EmailService) SendStudentVerificationEmail(input VerificationEmailInput) error {
	subject := fmt.Sprintf(s.config.Subjects.Verification, input.Name)

	templateInput := verificationEmailInput{s.createVerificationLink(input.Domain, input.VerificationCode)}
	sendInput := emailProvider.SendEmailInput{Subject: subject, To: input.Email}

	if err := sendInput.GenerateBodyFromHTML(s.config.Templates.Verification, templateInput); err != nil {
		return err
	}

	return s.sender.Send(sendInput)
}

func (s *EmailService) SendStudentPurchaseSuccessfulEmail(input StudentPurchaseSuccessfulEmailInput) error {
	templateInput := purchaseSuccessfulEmailInput{Name: input.Name, CourseName: input.CourseName}
	sendInput := emailProvider.SendEmailInput{Subject: s.config.Subjects.PurchaseSuccessful, To: input.Email}

	if err := sendInput.GenerateBodyFromHTML(s.config.Templates.PurchaseSuccessful, templateInput); err != nil {
		return err
	}

	return s.sender.Send(sendInput)
}

func (s *EmailService) SendUserVerificationEmail(input VerificationEmailInput) error {
	// todo implement
	return nil
}

func (s *EmailService) createVerificationLink(domain, code string) string {
	return fmt.Sprintf(verificationLinkTmpl, domain, code)
}

func (s *EmailService) AddStudentToList(ctx context.Context, email, name string, schoolID primitive.ObjectID) error {
	// TODO refactor
	school, err := s.schools.GetById(ctx, schoolID)
	if err != nil {
		return err
	}

	if !school.Settings.SendPulse.Connected {
		return domain.ErrSendPulseIsNotConnected
	}

	client, ex := s.sendpulseClients[schoolID]
	if !ex {
		client = sendpulse.NewClient(school.Settings.SendPulse.ID, school.Settings.SendPulse.Secret, s.cache)
		s.sendpulseClients[schoolID] = client
	}

	return client.AddEmailToList(emailProvider.AddEmailInput{
		Email:  email,
		ListID: school.Settings.SendPulse.ListID,
		Variables: map[string]string{
			"Name":   name,
			"source": "registration",
		},
	})
}
