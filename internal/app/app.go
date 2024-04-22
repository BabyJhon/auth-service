package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/BabyJhon/auth-service/internal/handlers"
	"github.com/BabyJhon/auth-service/internal/repos"
	"github.com/BabyJhon/auth-service/internal/service"
	mongodbgo "github.com/BabyJhon/auth-service/pkg/db/mongodb.go"
	"github.com/BabyJhon/auth-service/pkg/httpserver"
	"github.com/sirupsen/logrus"
)

func Run() {
	//logrus
	logrus.SetFormatter(new(logrus.JSONFormatter))

	//.env
	// if err := godotenv.Load(); err != nil {
	// 	logrus.Fatalf("error loading env vars: %s", err.Error())
	// }

	//DB
	client, err := mongodbgo.NewMongoDB(context.Background())
	if err != nil {
		logrus.Fatalf("failed init db: %s", err.Error())
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			logrus.Fatalf("failed disconnect db: %s", err.Error())
		}
	}()

	// Repositories
	repos := repos.NewRepository(client)

	// Service
	services := service.NewService(repos)

	// Handlers
	handlers := handlers.NewHandler(services)

	//HTTP server
	srv := new(httpserver.Server)

	go func() {
		if err := srv.Run("8000", handlers.InitRoutes()); err != http.ErrServerClosed {
			logrus.Fatalf("error occured while running server: %s", err.Error())
		}
	}()

	logrus.Print("auth API started")

	//gracefull shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("shutting down")
	if err := srv.ShutDown(context.Background()); err != nil {
		logrus.Errorf("error while server shutting down: %s", err.Error())
	}

}
