package main

import (
	"context"
	"fmt"
	"os/signal"
	"scheduler/cmd"
	"scheduler/config"
	"scheduler/internal/logger"
	"scheduler/internal/router"
	"syscall"
	"time"
)

func main() {
	// Setting up the configuration before starting the server
	cmd.SetUpConfig()
	cmd.SetUpLogger()
	cmd.SetUpDatabase()
	cmd.SetUpCron()
	cmd.LoadSchedules()

	r := router.NewFiberRouter()

	port := config.GetConfig().HttpServer.Port

	// checking for ctrl+c command
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Running the server at the specified port
	go func() {
		logger.Log.Info(fmt.Sprintf("Starting server on port %d", port))
		if err := r.Listen(fmt.Sprintf(":%d", port)); err != nil {
			logger.Log.Fatal(fmt.Sprintf("listen: %s\n", err))
		}
	}()

	<-ctx.Done()
	stop()

	// shutting down the server
	fmt.Println("\nShutting down gracefully, press Ctrl+C again to force")
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd.ShutDown()
	if err := r.ShutdownWithTimeout(5 * time.Second); err != nil {
		fmt.Println(err)
	}
}
