package jobs

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"scheduler/internal/app/jobs"
	"scheduler/internal/interface/request"
	"scheduler/internal/interface/response"
	"scheduler/internal/logger"
	"scheduler/pkg/exception"
	"strconv"
)

type JobsHTTPHandler struct {
	app jobs.JobsApp
}

func NewJobHTTPHandler(app jobs.JobsApp) *JobsHTTPHandler {
	return &JobsHTTPHandler{app: app}
}

func (h *JobsHTTPHandler) GetAllJobs(c *fiber.Ctx) error {
	dtos, err := h.app.GetAllJobs(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(response.CommonResponse{
		ResponseCode:    200,
		ResponseMessage: "OK",
		Data:            dtos,
	})
}

func (h *JobsHTTPHandler) DeleteAJob(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}
	err = h.app.DeleteAJob(c.Context(), id)
	if err != nil {
		return err
	}
	logger.Log.Info("Schedule has been deleted", zap.String("job_id", idStr))
	return c.JSON(response.CommonResponse{
		ResponseCode:    200,
		ResponseMessage: "OK",
	})
}

func (h *JobsHTTPHandler) AddAJob(c *fiber.Ctx) error {
	var req request.SchedulerRequest

	if err := c.BodyParser(&req); err != nil {
		return exception.InvalidRequestBodyError
	}

	err := req.ValidateSchedulerRequest()
	if err != nil {
		return exception.ValidationFailedError
	}

	err = h.app.AddJob(c.Context(), req)
	if err != nil {
		return err
	}
	logger.Log.Info("New schedule is created", zap.String("job_name", req.Name))
	return c.JSON(response.CommonResponse{
		ResponseCode:    200,
		ResponseMessage: "OK",
	})
}
