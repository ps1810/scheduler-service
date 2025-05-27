package metadata

import (
	"github.com/gofiber/fiber/v2"
	"scheduler/internal/interface/response"
)

type MetadataHTTPHandler struct {
}

func NewMetadataHTTPHandler() *MetadataHTTPHandler {
	return &MetadataHTTPHandler{}
}

func (m *MetadataHTTPHandler) GetMetaData(c *fiber.Ctx) error {
	return c.JSON(response.CommonResponse{
		ResponseCode:    200,
		ResponseMessage: "OK",
		Data:            response.Metrics,
	})
}
