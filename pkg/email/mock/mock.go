package mock_email

import (
	"github.com/stretchr/testify/mock"
	"github.com/zhashkevych/courses-backend/pkg/email"
)

type EmailProvider struct {
	mock.Mock
}

func (m *EmailProvider) AddEmailToList(inp email.AddEmailInput) error {
	args := m.Called(inp)
	return args.Error(0)
}