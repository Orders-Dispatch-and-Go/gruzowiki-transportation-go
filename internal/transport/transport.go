package transport

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"

	servicedto "auth-service/internal/service/dto"
	serviceerrlist "auth-service/internal/service/errlist"
)

type Transport struct {
	requestReader requestReader
	service       service
}

func New(requestReader requestReader, service service) *Transport {
	return &Transport{
		requestReader: requestReader,
		service:       service,
	}
}

func (t *Transport) GetUsers(fiberCtx *fiber.Ctx) error {
	domainUsers, err := t.service.GetUsers(fiberCtx.Context())
	if err != nil {
		return fmt.Errorf("get users: %w", err)
	}

	resp, err := convertToGetUsersResponse(domainUsers)
	if err != nil {
		return fmt.Errorf("convert: %w", err)
	}

	return fiberCtx.Status(fiber.StatusOK).JSON(resp)
}

func (t *Transport) GetCurrentUser(fiberCtx *fiber.Ctx) error {
	userID, err := extractUserID(fiberCtx)
	if err != nil {
		return fmt.Errorf("extract user ID: %w", err)
	}

	user, err := t.service.GetUser(fiberCtx.Context(), servicedto.GetUserRequest{UserID: userID})
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	resp, err := convertToGetCurrentUserResponse(user)
	if err != nil {
		return fmt.Errorf("convert: %w", err)
	}

	return fiberCtx.Status(fiber.StatusOK).JSON(resp)
}

func (t *Transport) GetUser(fiberCtx *fiber.Ctx) error {
	userIDString := fiberCtx.Params("id")

	userID, err := strconv.ParseInt(userIDString, 10, 64)
	if err != nil {
		return fiberCtx.Status(fiber.StatusBadRequest).JSON(errorResponse{
			Error: err.Error(),
		})
	}

	user, err := t.service.GetUser(fiberCtx.Context(), servicedto.GetUserRequest{UserID: userID})
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	resp, err := convertToGetUserResponse(user)
	if err != nil {
		return fmt.Errorf("convert: %w", err)
	}

	return fiberCtx.Status(fiber.StatusOK).JSON(resp)
}

func (t *Transport) Login(fiberCtx *fiber.Ctx) error {
	var req loginRequest

	if err := t.requestReader.ReadAndValidateFiberBody(fiberCtx, &req); err != nil {
		return fiberCtx.Status(fiber.StatusBadRequest).JSON(errorResponse{
			Error: err.Error(),
		})
	}

	serviceResp, err := t.service.Login(fiberCtx.Context(), req.convertToService())
	if err != nil {
		if errors.Is(err, serviceerrlist.ErrInvalidCredentials) {
			return fiberCtx.Status(fiber.StatusBadRequest).JSON(errorResponse{
				Error: err.Error(),
			})
		}

		return fmt.Errorf("service: %w", err)
	}

	resp := convertToLoginResponse(serviceResp)

	return fiberCtx.Status(fiber.StatusOK).JSON(resp)
}

func (t *Transport) Register(fiberCtx *fiber.Ctx) error {
	var req registerRequest

	if err := t.requestReader.ReadAndValidateFiberBody(fiberCtx, &req); err != nil {
		return fiberCtx.Status(fiber.StatusBadRequest).JSON(errorResponse{
			Error: err.Error(),
		})
	}

	serviceReq, err := req.convertToService()
	if err != nil {
		return fiberCtx.Status(fiber.StatusBadRequest).JSON(errorResponse{
			Error: err.Error(),
		})
	}

	err = t.service.Register(fiberCtx.Context(), serviceReq)
	if err != nil {
		return fmt.Errorf("service: %w", err)
	}

	return fiberCtx.SendStatus(fiber.StatusCreated)
}

func (t *Transport) VerifyEmail(fiberCtx *fiber.Ctx) error {
	token := fiberCtx.Query("token")

	if token == "" {
		return fiberCtx.Status(fiber.StatusBadRequest).JSON(errorResponse{
			Error: "token parameter is required",
		})
	}

	serviceReq := convertToServiceVerifyEmailRequest(token)

	err := t.service.VerifyEmail(fiberCtx.Context(), serviceReq)
	if err != nil {
		return fmt.Errorf("service: %w", err)
	}

	return fiberCtx.SendStatus(fiber.StatusNoContent)
}

func (t *Transport) ResendVerification(fiberCtx *fiber.Ctx) error {
	userID, err := extractUserID(fiberCtx)
	if err != nil {
		return fmt.Errorf("extract user id: %w", err)
	}

	serviceReq := convertToServiceSendVerificationRequest(userID)

	err = t.service.SendVerification(fiberCtx.Context(), serviceReq)
	if err != nil {
		return fmt.Errorf("service: %w", err)
	}

	return fiberCtx.SendStatus(fiber.StatusNoContent)
}

func (t *Transport) ForgotPassword(fiberCtx *fiber.Ctx) error {
	var req forgotPasswordRequest
	if err := t.requestReader.ReadAndValidateFiberBody(fiberCtx, &req); err != nil {
		return fiberCtx.Status(fiber.StatusBadRequest).JSON(errorResponse{
			Error: err.Error(),
		})
	}

	serviceReq := req.convertToService()

	err := t.service.ForgotPassword(fiberCtx.Context(), serviceReq)
	if err != nil {
		return fmt.Errorf("service: %w", err)
	}

	return nil
}

func (t *Transport) ResetPassword(fiberCtx *fiber.Ctx) error {
	var req resetPasswordRequest
	if err := t.requestReader.ReadAndValidateFiberBody(fiberCtx, &req); err != nil {
		return fiberCtx.Status(fiber.StatusBadRequest).JSON(errorResponse{
			Error: err.Error(),
		})
	}

	serviceReq := req.convertToService()

	err := t.service.ResetPassword(fiberCtx.Context(), serviceReq)
	if err != nil {
		return fmt.Errorf("service: %w", err)
	}

	return nil
}
