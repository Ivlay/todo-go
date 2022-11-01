package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Ivlay/todo-go"
	"github.com/Ivlay/todo-go/pkg/handler"
	"github.com/Ivlay/todo-go/pkg/repository"
	"github.com/Ivlay/todo-go/pkg/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// @title           Todo App Api
// @version         1.0
// @description Api server TodoList Application

// @host      localhost:8080
// @BasePath /

// @securityDefinitions.apiKey  ApiKeyAuth
// @in Header
// @name Authorization

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("error init config: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("errpr locading env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host: viper.GetString("db.host"),
		Port: viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName: viper.GetString("db.dbname"),
		SSLMode: viper.GetString("db.sslmode"),
	})

	if err != nil {
		logrus.Fatalf("failed to init db: %s", err.Error())
	}

	repos := repository.NewRspository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(todo.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error: %s", err.Error())
		}
	}()

	logrus.Print("Server started at port %s", viper.GetString("port"))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<- quit

	logrus.Print("Server Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error while server shutting down: %s", err.Error())	
	}
	
	if err := db.Close(); err != nil {
		logrus.Errorf("error on db connection close: %s", err.Error())	
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
