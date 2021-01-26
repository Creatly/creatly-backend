package service

import (
	"context"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"github.com/zhashkevych/courses-backend/internal/repository"
	"github.com/zhashkevych/courses-backend/pkg/auth"
	"github.com/zhashkevych/courses-backend/pkg/hash"
	"github.com/zhashkevych/courses-backend/pkg/logger"
	"github.com/zhashkevych/courses-backend/pkg/payment"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type StudentsService struct {
	repo            repository.Students
	coursesService  Courses
	hasher          hash.PasswordHasher
	tokenManager    auth.TokenManager
	emailService    Emails
	paymentProvider payment.Provider

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration

	callbackURL string
	responseURL string
}

func NewStudentsService(repo repository.Students, coursesService Courses, hasher hash.PasswordHasher, tokenManager auth.TokenManager,
	emailService Emails, paymentProvider payment.Provider, accessTTL, refreshTTL time.Duration, callbackURL, responseURL string) *StudentsService {
	return &StudentsService{
		repo:            repo,
		coursesService:  coursesService,
		hasher:          hasher,
		emailService:    emailService,
		tokenManager:    tokenManager,
		paymentProvider: paymentProvider,
		accessTokenTTL:  accessTTL,
		refreshTokenTTL: refreshTTL,
		callbackURL:     callbackURL,
		responseURL:     responseURL,
	}
}

func (s *StudentsService) SignUp(ctx context.Context, input StudentSignUpInput) error {
	verificationCode := primitive.NewObjectID()
	student := domain.Student{
		Name:           input.Name,
		Password:       s.hasher.Hash(input.Password),
		Email:          input.Email,
		RegisteredAt:   time.Now(),
		LastVisitAt:    time.Now(),
		SchoolID:       input.SchoolID,
		RegisterSource: input.RegisterSource,
		Verification: domain.Verification{
			Code: verificationCode,
		},
	}

	if err := s.repo.Create(ctx, student); err != nil {
		return err
	}

	// TODO: If it fails, what then?
	return s.emailService.AddToList(AddToListInput{
		Email:            student.Email,
		Name:             student.Name,
		RegisterSource:   student.RegisterSource,
		VerificationCode: verificationCode.Hex(),
	})
}

func (s *StudentsService) SignIn(ctx context.Context, input StudentSignInInput) (Tokens, error) {
	student, err := s.repo.GetByCredentials(ctx, input.SchoolID, input.Email, s.hasher.Hash(input.Password))
	if err != nil {
		return Tokens{}, err
	}

	return s.createSession(ctx, student.ID)
}

func (s *StudentsService) RefreshTokens(ctx context.Context, schoolId primitive.ObjectID, refreshToken string) (Tokens, error) {
	student, err := s.repo.GetByRefreshToken(ctx, schoolId, refreshToken)
	if err != nil {
		return Tokens{}, err
	}

	return s.createSession(ctx, student.ID)
}

func (s *StudentsService) Verify(ctx context.Context, hash string) error {
	return s.repo.Verify(ctx, hash)
}

func (s *StudentsService) GetStudentModuleWithLessons(ctx context.Context, schoolId, studentId, moduleId primitive.ObjectID) ([]domain.Lesson, error) {
	// Get module with lessons content, check if it is available for student
	module, err := s.coursesService.GetModuleWithContent(ctx, moduleId)
	if err != nil {
		return nil, err
	}

	student, err := s.repo.GetById(ctx, studentId)
	if err != nil {
		return nil, nil
	}

	if student.IsModuleAvailable(module) {
		return module.Lessons, nil
	}

	// Find module offers
	offers, err := s.coursesService.GetPackageOffers(ctx, schoolId, module.PackageID)
	if err != nil {
		return nil, err
	}

	if len(offers) != 0 {
		return nil, ErrModuleIsNotAvailable
	}

	// If module has no offers - it's free and available to everyone
	go func() {
		if err := s.repo.GiveModuleAccess(ctx, studentId, moduleId); err != nil {
			logger.Error(err)
		}
	}()

	return module.Lessons, nil
}

func (s *StudentsService) CreateOrder(ctx context.Context, studentId, offerId, promocodeId primitive.ObjectID) (string, error) {
	promocode, err := s.getOrderPromocode(ctx, promocodeId)
	if err != nil {
		return "", err
	}

	offer, err := s.coursesService.GetOfferById(ctx, offerId)
	if err != nil {
		return "", err
	}

	orderAmount := s.calculateOrderPrice(offer.Price.Value, promocode)

	id, err := s.repo.CreateOrder(ctx, studentId, domain.Order{
		OfferID: offerId,
		PromoID: promocodeId,
		Amount:  orderAmount,
		Status:  domain.OrderStatusCreated,
	})

	// TODO what if it fails?
	return s.paymentProvider.GeneratePaymentLink(payment.GeneratePaymentLinkInput{
		OrderId:     id.Hex(),
		Amount:      orderAmount,
		Currency:    offer.Price.Currency,
		OrderDesc:   offer.Description, // TODO proper order description
		CallbackURL: s.callbackURL,
		ResponseURL: s.responseURL,
	})
}

func (s *StudentsService) createSession(ctx context.Context, studentId primitive.ObjectID) (Tokens, error) {
	var (
		res Tokens
		err error
	)

	res.AccessToken, err = s.tokenManager.NewJWT(studentId.Hex(), s.accessTokenTTL)
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		return res, err
	}

	session := domain.Session{
		RefreshToken: res.RefreshToken,
		ExpiresAt:    time.Now().Add(s.refreshTokenTTL),
	}

	err = s.repo.SetSession(ctx, studentId, session)
	return res, err
}

func (s *StudentsService) getOrderPromocode(ctx context.Context, promocodeId primitive.ObjectID) (domain.Promocode, error) {
	var (
		promocode domain.Promocode
		err       error
	)

	if !promocodeId.IsZero() {
		promocode, err = s.coursesService.GetPromocodeById(ctx, promocodeId)
		if err != nil {
			return promocode, err
		}

		if promocode.ExpiresAt.Unix() < time.Now().Unix() {
			return promocode, ErrPromocodeExpired
		}
	}

	return promocode, nil
}

func (s *StudentsService) calculateOrderPrice(price int, promocode domain.Promocode) int {
	if promocode.ID.IsZero() {
		return price
	} else {
		return (price * (100 - promocode.DiscountPercentage)) / 100
	}
}
