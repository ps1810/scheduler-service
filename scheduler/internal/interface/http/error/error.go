package http

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"scheduler/internal/interface/response"
	"scheduler/pkg/exception"
)

func ErrorHandler(c *fiber.Ctx, err error) error {

	responseCode := fiber.StatusInternalServerError
	responseMessage := err.Error()
	var cErrs *exception.ExceptionErrors

	//cErrs, ok := err.(*exception.ExceptionErrors)
	if errors.As(err, &cErrs) {
		responseCode = cErrs.HttpStatusCode
	}

	return c.Status(responseCode).JSON(
		&response.CommonResponse{
			ResponseCode:    responseCode,
			ResponseMessage: responseMessage,
			Errors:          cErrs,
		},
	)
}
