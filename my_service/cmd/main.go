package main

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"learning_NATS_streaming/internal/app"
	"learning_NATS_streaming/internal/infrastructure/storage"
	"os"
)

//docker run --name tg_password_bot -p 5432:5432 -e POSTGRES_PASSWORD=qwerty -d postgres
func main() {

	postgresPassword := initGoDotEnv()
	if err := initConfig(); err != nil {
		logrus.Fatalf("error while reading config file:%s", err.Error())
	}
	postgresConfig := storage.ConfigDB{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.DBName"),
		Password: postgresPassword,
		SSLMode:  viper.GetString("db.SSLMode"),
	}
	app.RunApp(postgresConfig)
}

func initConfig() error {
	viper.SetDefault("port", "8000")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("config")
	return viper.ReadInConfig()
}

func initGoDotEnv() string {
	err := godotenv.Load()
	if err != nil {
		logrus.Error("Error loading .env file")
	}

	return os.Getenv("POSTGRES_PASSWORD")
}
