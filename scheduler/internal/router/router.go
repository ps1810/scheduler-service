package router

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	httpInterface "scheduler/internal/interface/http"
	httpError "scheduler/internal/interface/http/error"
)

func NewFiberRouter() *fiber.App {
	r := fiber.New(fiber.Config{
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
		ErrorHandler: httpError.ErrorHandler,
	})

	httpInterface.RegisterRoute(r)
	return r
}
