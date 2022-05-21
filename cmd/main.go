package main

import (
	"os"
	"time"

	"github.com/folins/biketrackserver"
	"github.com/folins/biketrackserver/package/handler"
	"github.com/folins/biketrackserver/package/repository"
	"github.com/folins/biketrackserver/package/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	logrus.SetLevel(logrus.DebugLevel)
	if err := initConfig(); err != nil {
		logrus.Fatalf("error init configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	logrus.Debugf("Server start time: %s", time.Now())

	// db, err := repository.NewPostgreDB(repository.Config{
	// 	Host:     viper.GetString("db.host"),
	// 	Port:     viper.GetString("db.port"),
	// 	Username: viper.GetString("db.username"),
	// 	DBNmae:   viper.GetString("db.dbname"),
	// 	SSLMode:  viper.GetString("db.sslmode"),
	// 	Password: os.Getenv("DB_PASSWORD"),
	// })

	db, err := repository.NewHerokuDB(os.Getenv("DATABASE_URL"))

	if err != nil {
		logrus.Fatalf("failed to init database: %s", err.Error())
	}

	smtp := service.NewSMTPService(
		viper.GetString("smtp.sender_email"),
		os.Getenv("SMTP_PASSWORD"),
		viper.GetString("smtp.host"),
		viper.GetInt("smtp.port"),
	)

	repos := repository.NewRepository(db)
	service := service.NewService(repos, smtp)
	handlers := handler.NewHandler(service)
	srv := new(biketrackserver.Server)

	if err := srv.Run(os.Getenv("PORT"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("Error occured while running server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
