package service

import (
	"fmt"

	"github.com/zhashkevych/creatly-backend/internal/config"
	emailProvider "github.com/zhashkevych/creatly-backend/pkg/email"
)

const (
	nameField            = "name"
	verificationLinkTmpl = "%s/verification?code=%s" // <frontend_url>/verification?code=<verification_code>
)

type (
	// Structures used for templates.
	verificationEmailInput struct {
		VerificationLink string
	}

	purchaseSuccessfulEmailInput struct {
		Name       string
		CourseName string
	}

	VerificationEmailInput struct {
		Email            string
		Name             string
		VerificationCode string
	}

	StudentPurchaseSuccessfulEmailInput struct {
		Email      string
		Name       string
		CourseName string
	}

	AddEmailInput struct {
		Email     string
		ListID    string
		Variables map[string]string
	}

	Provider interface {
		AddEmailToList(emailProvider.AddEmailInput) error
	}

	Sender interface {
		Send(input emailProvider.SendEmailInput) error
	}

	EmailService struct {
		provider    Provider
		sender      Sender
		config      config.EmailConfig
		frontendUrl string
	}
)

func NewEmailsService(provider emailProvider.Provider, sender emailProvider.Sender, config config.EmailConfig, frontendUrl string) *EmailService {
	return &EmailService{provider: provider, sender: sender, config: config, frontendUrl: frontendUrl}
}

func (s *EmailService) AddToList(name, email string) error {
	return s.provider.AddEmailToList(emailProvider.AddEmailInput{
		Email:  email,
		ListID: s.config.SendPulse.ListID,
		Variables: map[string]string{
			nameField: name,
		},
	})
}

func (s *EmailService) SendStudentVerificationEmail(input VerificationEmailInput) error {
	subject := fmt.Sprintf(s.config.Subjects.Verification, input.Name)

	templateInput := verificationEmailInput{s.createVerificationLink(input.VerificationCode)}
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

func (s *EmailService) createVerificationLink(code string) string {
	return fmt.Sprintf(verificationLinkTmpl, s.frontendUrl, code)
}
