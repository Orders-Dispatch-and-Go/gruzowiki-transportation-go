package transport

import (
	"fmt"
	"time"

	"auth-service/internal/domain"
	servicedto "auth-service/internal/service/dto"
)

type errorResponse struct {
	Error string `json:"error"`
}

type getUsersResponse struct {
	Items []getUsersResponseItem `json:"items"`
}

func convertToGetUsersResponse(domainUsers []domain.User) (getUsersResponse, error) {
	resp := getUsersResponse{
		Items: make([]getUsersResponseItem, 0, len(domainUsers)),
	}

	for _, user := range domainUsers {
		userRole, err := convertFromDomainUserRole(user.Role)
		if err != nil {
			return resp, fmt.Errorf("convert user role: %w", err)
		}

		resp.Items = append(resp.Items, getUsersResponseItem{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			Role:      userRole,
			Bio:       user.Bio,
			AvatarURL: user.AvatarURL,
			CreatedAt: user.CreatedAt,
		})
	}

	return resp, nil
}

type getUsersResponseItem struct {
	ID        int64      `json:"id"`
	Email     string     `json:"email"`
	Name      string     `json:"name"`
	Role      userRole   `json:"role"`
	Bio       *string    `json:"bio"`
	AvatarURL *string    `json:"avatar_url"`
	CreatedAt *time.Time `json:"created_at"`
}

type getCurrentUserResponse struct {
	ID        int64      `json:"id"`
	Email     string     `json:"email"`
	Name      string     `json:"name"`
	Role      userRole   `json:"role"`
	Bio       *string    `json:"bio"`
	AvatarURL *string    `json:"avatar_url"`
	CreatedAt *time.Time `json:"created_at"`
}

func convertToGetCurrentUserResponse(user domain.User) (getCurrentUserResponse, error) {
	resp := getCurrentUserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Bio:       user.Bio,
		AvatarURL: user.AvatarURL,
		CreatedAt: user.CreatedAt,
	}

	var err error

	resp.Role, err = convertFromDomainUserRole(user.Role)
	if err != nil {
		return resp, fmt.Errorf("convert user role: %w", err)
	}

	return resp, nil
}

type getUserResponse struct {
	ID        int64      `json:"id"`
	Email     string     `json:"email"`
	Name      string     `json:"name"`
	Role      userRole   `json:"role"`
	Bio       *string    `json:"bio"`
	AvatarURL *string    `json:"avatar_url"`
	CreatedAt *time.Time `json:"created_at"`
}

func convertToGetUserResponse(user domain.User) (getUserResponse, error) {
	resp := getUserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Bio:       user.Bio,
		AvatarURL: user.AvatarURL,
		CreatedAt: user.CreatedAt,
	}

	var err error

	resp.Role, err = convertFromDomainUserRole(user.Role)
	if err != nil {
		return resp, fmt.Errorf("role: %w", err)
	}

	return resp, nil
}

type loginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (r loginRequest) convertToService() servicedto.LoginRequest {
	req := servicedto.LoginRequest{
		Email:    r.Email,
		Password: r.Password,
	}

	return req
}

type loginResponse struct {
	Token string `json:"token"`
}

func convertToLoginResponse(serviceResp servicedto.LoginResponse) loginResponse {
	return loginResponse{
		Token: serviceResp.Token,
	}
}

type registerRequest struct {
	Email     string   `json:"email" validate:"required"`
	Name      string   `json:"name" validate:"required"`
	Password  string   `json:"password" validate:"required"`
	Role      userRole `json:"role" validate:"required"`
	Bio       *string  `json:"bio"`
	AvatarURL *string  `json:"avatar_url"`
}

func (r registerRequest) convertToService() (req servicedto.RegisterRequest, err error) {
	req.Role, err = convertToDomainUserRole(r.Role)
	if err != nil {
		return req, fmt.Errorf("role: %w", err)
	}

	req.Email = r.Email
	req.Name = r.Name
	req.Password = r.Password
	req.Bio = r.Bio
	req.AvatarURL = r.AvatarURL

	return req, nil
}

type userRole string

const (
	userRoleReader = "reader"
	userRoleWriter = "writer"
)

func convertFromDomainUserRole(role domain.UserRole) (userRole, error) {
	switch role {
	case domain.UserRoleReader:
		return userRoleReader, nil
	case domain.UserRoleWriter:
		return userRoleWriter, nil
	default:
		return "", fmt.Errorf("unknown user role: %q", role)
	}
}

func convertToDomainUserRole(role userRole) (domain.UserRole, error) {
	switch role {
	case userRoleReader:
		return domain.UserRoleReader, nil
	case userRoleWriter:
		return domain.UserRoleWriter, nil
	default:
		return "", fmt.Errorf("unknown user role: %q", role)
	}
}

func convertToServiceVerifyEmailRequest(token string) servicedto.VerifyEmailRequest {
	return servicedto.VerifyEmailRequest{
		Token: token,
	}
}

func convertToServiceSendVerificationRequest(userID int64) servicedto.SendVerificationRequest {
	return servicedto.SendVerificationRequest{
		UserID: userID,
	}
}

type forgotPasswordRequest struct {
	Email string `json:"email" validate:"required"`
}

func (r forgotPasswordRequest) convertToService() servicedto.ForgotPasswordRequest {
	return servicedto.ForgotPasswordRequest{
		Email: r.Email,
	}
}

type resetPasswordRequest struct {
	Enail    string `json:"email" validate:"required"`
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (r resetPasswordRequest) convertToService() servicedto.ResetPasswordRequest {
	return servicedto.ResetPasswordRequest{
		Email:    r.Enail,
		Token:    r.Token,
		Password: r.Password,
	}
}
