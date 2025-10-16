package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"auth-service/internal/domain"
	"auth-service/internal/service/dto"
	"auth-service/internal/service/errlist"
)

type Service struct {
	repo           repo
	passwordHasher passwordHasher
	auth           auth
	emailSender    emailSender
}

func New(repo repo, passwordHasher passwordHasher, auth auth, emailSender emailSender) *Service {
	return &Service{
		repo:           repo,
		passwordHasher: passwordHasher,
		auth:           auth,
		emailSender:    emailSender,
	}
}

func (s *Service) GetUsers(ctx context.Context) ([]domain.User, error) {
	return s.repo.GetUsers(ctx)
}

func (s *Service) GetUser(ctx context.Context, req dto.GetUserRequest) (domain.User, error) {
	user, err := s.repo.GetUserByID(ctx, req.UserID)
	if err != nil {
		return domain.User{}, fmt.Errorf("GetUserByID(%d): %w", req.UserID, err)
	}

	if user == nil {
		return domain.User{}, errors.New("user not found")
	}

	return *user, nil
}

func (s *Service) Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return dto.LoginResponse{}, fmt.Errorf("GetUserByEmail(%q): %w", req.Email, err)
	}

	if user == nil {
		return dto.LoginResponse{}, errlist.ErrInvalidCredentials
	}

	if err := s.passwordHasher.VerifyHash(req.Password, user.PasswordHash); err != nil {
		return dto.LoginResponse{}, errlist.ErrInvalidCredentials
	}

	token, err := s.auth.GenerateToken(user.ID)
	if err != nil {
		return dto.LoginResponse{}, fmt.Errorf("GenerateToken(%d): %w", user.ID, err)
	}

	return dto.LoginResponse{Token: token}, nil
}

func (s *Service) Register(ctx context.Context, req dto.RegisterRequest) error {
	passwordHash, err := s.passwordHasher.Hash(req.Password)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	err = s.repo.InsertUser(ctx, req.ConvertToRepo(passwordHash))
	if err != nil {
		return fmt.Errorf("insert: %w", err)
	}

	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return fmt.Errorf("get user by email: %w", err)
	}

	err = s.SendVerification(ctx, dto.NewSendVerificationRequest(user.ID))
	if err != nil {
		return fmt.Errorf("send verification: %w", err)
	}

	return nil
}

func (s *Service) VerifyEmail(ctx context.Context, req dto.VerifyEmailRequest) error {
	userID, err := s.auth.ParseToken(req.Token)
	if err != nil {
		return fmt.Errorf("ParseToken(%s): %w", req.Token, err)
	}

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("GetUserByID(%d): %w", userID, err)
	}

	if user == nil {
		return errors.New("user not found")
	}

	if user.Verified {
		return errors.New("email already verified")
	}

	repoReq := dto.ConvertToRepoUserVerifiedRequest(user.ID, true)

	if err := s.repo.SetUserVerified(ctx, repoReq); err != nil {
		return fmt.Errorf("SetUserVerified(%d, %t): %w", userID, true, err)
	}

	return nil
}

func (s *Service) SendVerification(ctx context.Context, req dto.SendVerificationRequest) error {
	user, err := s.repo.GetUserByID(ctx, req.UserID)
	if err != nil {
		return fmt.Errorf("get user by id: %w", err)
	}

	token, err := s.auth.GenerateVerificationToken(user.ID)
	if err != nil {
		return fmt.Errorf("generate verification token: %w", err)
	}

	err = s.emailSender.SendVerificationEmail(user.Email, token)
	if err != nil {
		return fmt.Errorf("send verification email: %w", err)
	}

	return nil
}

func (s *Service) ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) error {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return fmt.Errorf("GetUserByEmail(%s): %w", req.Email, err)
	}

	verified, err := s.repo.GetUserVerified(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("GetUserVerified(%d): %w", user.ID, err)
	}

	if !verified {
		return errors.New("only verified users can reset password")
	}

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := make([]byte, 6)

	for i := range code {
		code[i] = byte('0' + random.Intn(10))
	}
	token := string(code)
	repoReq := dto.ConvertToRepoInsertPasswordResetTokenRequest(user.ID, token)

	err = s.repo.InsertPasswordResetToken(ctx, repoReq)
	if err != nil {
		return fmt.Errorf("InserResetToken(%d, %s): %w", user.ID, token, err)
	}

	err = s.emailSender.SendPasswordResetEmail(user.Email, token)
	if err != nil {
		return fmt.Errorf("SendPasswordResetEmail: %w", err)
	}

	return nil
}

func (s *Service) ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return fmt.Errorf("GetUserByEmail(%s): %w", req.Email, err)
	}

	repoGetPasswordResetTokenRequest := req.ConvertToRepo(user.ID)

	passwordResetToken, err := s.repo.GetPasswordResetToken(ctx, repoGetPasswordResetTokenRequest)
	if err != nil {
		return fmt.Errorf("GetPasswordResetToken(%d, %s): %w", user.ID, req.Token, err)
	}

	if passwordResetToken.Token.Used {
		return errors.New("token alredy used")
	}

	if !passwordResetToken.Token.ExpiredAt.After(time.Now()) {
		return errors.New("token expired")
	}

	passwordHash, err := s.passwordHasher.Hash(req.Password)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	repoSetPasswordHashRequest := dto.ConvertToRepoSetPasswordHashRequest(user.ID, passwordHash)

	err = s.repo.SetUserPasswordHash(ctx, repoSetPasswordHashRequest)
	if err != nil {
		return fmt.Errorf("set user password hash: %w", err)
	}

	err = s.repo.MarkPasswordResetTokenAsUsed(ctx, passwordResetToken.Token.ID)
	if err != nil {
		return fmt.Errorf("MarkPasswodResetTokenAsUsed(%d): %w", passwordResetToken.Token.ID, err)
	}

	return nil
}
