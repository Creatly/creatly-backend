package mock_email

import (
	"github.com/stretchr/testify/mock"
	"github.com/zhashkevych/creatly-backend/pkg/email"
)

type EmailProvider struct {
	mock.Mock
}

func (m *EmailProvider) AddEmailToList(inp email.AddEmailInput) error {
	args := m.Called(inp)

	return args.Error(0)
}

type EmailSender struct {
	mock.Mock
}

func (m *EmailSender) Send(inp email.SendEmailInput) error {
	args := m.Called(inp)

	return args.Error(0)
}
