package tests

import (
	"context"
	"github.com/stretchr/testify/suite"
	"github.com/zhashkevych/courses-backend/internal/delivery/http"
	"github.com/zhashkevych/courses-backend/internal/repository"
	"github.com/zhashkevych/courses-backend/internal/service"
	"github.com/zhashkevych/courses-backend/pkg/auth"
	"github.com/zhashkevych/courses-backend/pkg/cache"
	"github.com/zhashkevych/courses-backend/pkg/database/mongodb"
	emailmock "github.com/zhashkevych/courses-backend/pkg/email/mock"
	"github.com/zhashkevych/courses-backend/pkg/hash"
	"github.com/zhashkevych/courses-backend/pkg/payment"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"testing"
	"time"
)

var (
	dbURI, dbUsername, dbPassword string
	dbName                        = "coursesTesting"
)

func init() {
	dbURI = os.Getenv("DB_URI")
	dbUsername = os.Getenv("DB_USERNAME")
	dbPassword = os.Getenv("DB_PASSWORD")
}

type APITestSuite struct {
	suite.Suite

	db       *mongo.Database
	handler  *http.Handler
	services *service.Services
	repos    *repository.Repositories

	mocks *mocks
}

type mocks struct {
	emailProvider *emailmock.EmailProvider
}

func TestAPISuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	suite.Run(t, new(APITestSuite))
}

func (s *APITestSuite) SetupSuite() {
	if client, err := mongodb.NewClient(dbURI, dbUsername, dbPassword); err != nil {
		s.FailNow("Failed to connect to mongo", err)
	} else {
		s.db = client.Database(dbName)
	}

	s.initMocks()
	s.initDeps()
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
	paymentProvider := payment.NewFondyClient("1396424", "test") // Fondy Testing Credentials

	services := service.NewServices(service.Deps{
		Repos:                  repos,
		Cache:                  memCache,
		Hasher:                 hasher,
		TokenManager:           tokenManager,
		PaymentProvider:        paymentProvider,
		EmailProvider:          s.mocks.emailProvider,
		AccessTokenTTL:         time.Minute * 15,
		RefreshTokenTTL:        time.Minute * 15,
		CacheTTL:               int64(time.Minute.Seconds()),
		VerificationCodeLength: 8,
	})

	s.repos = repos
	s.services = services
	s.handler = http.NewHandler(services, tokenManager)
}

func (s *APITestSuite) initMocks() {
	s.mocks = &mocks{
		emailProvider: new(emailmock.EmailProvider),
	}
}
func TestMain(m *testing.M) {
	rc := m.Run()
	os.Exit(rc)
}
