package http

import (
	"github.com/gofiber/fiber/v2"
	"scheduler/internal/app/jobs"
	httpJobs "scheduler/internal/interface/http/jobs"
	"scheduler/internal/interface/http/metadata"
	"scheduler/internal/repository"
)

var repo *repository.Repository

func RegisterRoute(r *fiber.App) {

	repo = repository.NewRepository()
	api := r.Group("/api")
	v1 := api.Group("/v1")

	// Job API
	jobAPI := v1.Group("/cron")
	jobApp := jobs.NewJobApp(repo)
	jobHandler := httpJobs.NewJobHTTPHandler(jobApp)
	jobAPI.Post("/add", jobHandler.AddAJob)
	jobAPI.Get("/jobs", jobHandler.GetAllJobs)
	jobAPI.Delete("/job/:id", jobHandler.DeleteAJob)

	// metadata API
	metadataAPI := v1.Group("/metadata")
	metadataHandler := metadata.NewMetadataHTTPHandler()
	metadataAPI.Get("/", metadataHandler.GetMetaData)

}
