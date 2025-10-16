package request

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Reader struct {
	validate *validator.Validate
}

func NewReader() Reader {
	return Reader{
		validate: validator.New(),
	}
}

func (r Reader) ReadAndValidateFiberBody(fiberCtx *fiber.Ctx, request any) error {
	if err := fiberCtx.BodyParser(request); err != nil {
		return fmt.Errorf("body parser: %w", err)
	}

	if err := r.validate.StructCtx(fiberCtx.Context(), request); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}
