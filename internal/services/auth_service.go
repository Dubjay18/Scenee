package services

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/Dubjay18/scenee/internal/domain"
	"github.com/Dubjay18/scenee/internal/models"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
)

type AuthService struct {
	usvc      *UserService
	jwtSecret string
	nsvc 	*NotificationService
}

func NewAuthService(usvc *UserService, jwtSecret string, ensendProjectId string, ensendProjectSecret string) *AuthService {
	return &AuthService{
		usvc:      usvc,
		jwtSecret: jwtSecret,
		nsvc: NewNotificationService(NotificationConfig{
			EnSendProjectID:     ensendProjectId,
			EnSendProjectSecret: ensendProjectSecret,
		}),
	}
}

func (s *AuthService) Register(ctx context.Context, email, username, password string) (*domain.User, error) {
	// Check if user exists
	if _, err := s.usvc.GetByEmail(ctx, email); err == nil {
		return nil, ErrUserExists
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:    email,
		Username: username,
		Password: string(hashed),
	}

	if err := s.usvc.Upsert(ctx, user); err != nil {
		return nil, err
	}

	// Send welcome email
	subject := "Welcome to Scenee!"
	body := "Hi " + username + ", Thank you for registering with Scenee. We're excited to have you on board!"

	err = s.nsvc.SendEmailNotification(email, subject, body)
	if err != nil {
		// Log the error, but don't fail the registration
		 log.Printf("Failed to send welcome email to %s: %v", email, err)
	}
	return domain.UserFromModel(user), nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, *domain.User, error) {
	user, err := s.usvc.GetByEmail(ctx, email)
	if err != nil {
		return "", nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, ErrInvalidCredentials
	}

	token, err := s.generateJWT(user.ID.String())
	if err != nil {
		return "", nil, err
	}

	return token, domain.UserFromModel(user), nil
}

func (s *AuthService) GetUser(ctx context.Context, id string) (*domain.User, error) {
	res, err := s.usvc.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return domain.UserFromModel(res), nil
}

func (s *AuthService) generateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
