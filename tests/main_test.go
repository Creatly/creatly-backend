package tests

import (
	"context"
	"github.com/stretchr/testify/suite"
	"github.com/zhashkevych/courses-backend/internal/config"
	v1 "github.com/zhashkevych/courses-backend/internal/delivery/http/v1"
	"github.com/zhashkevych/courses-backend/internal/repository"
	"github.com/zhashkevych/courses-backend/internal/service"
	"github.com/zhashkevych/courses-backend/pkg/auth"
	"github.com/zhashkevych/courses-backend/pkg/cache"
	"github.com/zhashkevych/courses-backend/pkg/database/mongodb"
	emailmock "github.com/zhashkevych/courses-backend/pkg/email/mock"
	"github.com/zhashkevych/courses-backend/pkg/hash"
	"github.com/zhashkevych/courses-backend/pkg/otp"
	"github.com/zhashkevych/courses-backend/pkg/payment/fondy"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"testing"
	"time"
)

const (
	listId = "123456"
)

var (
	dbURI, dbName string
)

func init() {
	dbURI = os.Getenv("DB_URI")
	dbName = os.Getenv("DB_NAME")
}

type APITestSuite struct {
	suite.Suite

	db       *mongo.Database
	handler  *v1.Handler
	services *service.Services
	repos    *repository.Repositories

	hasher       hash.PasswordHasher
	tokenManager auth.TokenManager
	mocks        *mocks
}

type mocks struct {
	emailProvider *emailmock.EmailProvider
	emailSender   *emailmock.EmailSender
	otpGenerator  *otp.MockGenerator
}

func TestAPISuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	suite.Run(t, new(APITestSuite))
}

func (s *APITestSuite) SetupSuite() {
	if client, err := mongodb.NewClient(dbURI, "", ""); err != nil {
		s.FailNow("Failed to connect to mongo", err)
	} else {
		s.db = client.Database(dbName)
	}

	s.initMocks()
	s.initDeps()

	if err := s.populateDB(); err != nil {
		s.FailNow("Failed to populate DB", err)
	}
}

func (s *APITestSuite) TearDownSuite() {
	s.db.Client().Disconnect(context.Background())
}

func (s *APITestSuite) initDeps() {
	// Init domain deps
	repos := repository.NewRepositories(s.db)
	memCache := cache.NewMemoryCache()
	hasher := hash.NewSHA1Hasher("salt")
	tokenManager, err := auth.NewManager("signing_key")
	if err != nil {
		s.FailNow("Failed to initialize token manager", err)
	}
	paymentProvider := fondy.NewFondyClient("1396424", "test") // Fondy Testing Credentials

	services := service.NewServices(service.Deps{
		Repos:           repos,
		Cache:           memCache,
		Hasher:          hasher,
		TokenManager:    tokenManager,
		PaymentProvider: paymentProvider,
		EmailProvider:   s.mocks.emailProvider,
		EmailSender:     s.mocks.emailSender,
		EmailConfig: config.EmailConfig{
			SendPulse: config.SendPulseConfig{
				ListID: listId,
			},
			Templates: config.EmailTemplates{
				Verification:       "../templates/verification_email.html",
				PurchaseSuccessful: "../templates/purchase_successful.html",
			},
			Subjects: config.EmailSubjects{
				Verification: "Спасибо за регистрацию, %s!",
				PurchaseSuccessful: "Покупка прошла успешно!",
			},
		},
		AccessTokenTTL:         time.Minute * 15,
		RefreshTokenTTL:        time.Minute * 15,
		CacheTTL:               int64(time.Minute.Seconds()),
		OtpGenerator:           s.mocks.otpGenerator,
		VerificationCodeLength: 8,
		FrontendURL:            "http://localhost:1337",
	})

	s.repos = repos
	s.services = services
	s.handler = v1.NewHandler(services, tokenManager)
	s.hasher = hasher
	s.tokenManager = tokenManager
}

func (s *APITestSuite) initMocks() {
	s.mocks = &mocks{
		emailProvider: new(emailmock.EmailProvider),
		emailSender:   new(emailmock.EmailSender),
		otpGenerator:  new(otp.MockGenerator),
	}
}
func TestMain(m *testing.M) {
	rc := m.Run()
	os.Exit(rc)
}

func (s *APITestSuite) populateDB() error {
	_, err := s.db.Collection("schools").InsertOne(context.Background(), school)
	if err != nil {
		return err
	}

	_, err = s.db.Collection("packages").InsertMany(context.Background(), packages)
	if err != nil {
		return err
	}

	_, err = s.db.Collection("offers").InsertMany(context.Background(), offers)
	if err != nil {
		return err
	}

	_, err = s.db.Collection("modules").InsertMany(context.Background(), modules)
	if err != nil {
		return err
	}

	_, err = s.db.Collection("promocodes").InsertMany(context.Background(), promocodes)
	if err != nil {
		return err
	}

	return nil
}
