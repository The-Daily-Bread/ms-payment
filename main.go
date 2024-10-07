package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tdb/ms-payment/src/cmd/utils"
	"github.com/tdb/ms-payment/src/pkg/server"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	config := read_configuration(read())

	server.Initialize(config)
}

func read_configuration(config utils.Configuration) utils.Configuration {

	mongoUri := os.Getenv("MONGODB_URL")
	port := os.Getenv("SERVER_PORT")
	dbName := os.Getenv("DB_NAME")
	appName := os.Getenv("APP_NAME")

	if mongoUri != "" || port != "" || dbName != "" || appName != "" {
		return utils.Configuration{
			App:      utils.Application{Name: appName},
			Database: utils.DatabaseSetting{Url: mongoUri, DbName: dbName},
			Server:   utils.ServerSettings{Port: port},
		}
	}

	return utils.Configuration{
		App:      utils.Application{Name: config.App.Name},
		Database: utils.DatabaseSetting{Url: config.Database.Url, DbName: config.Database.DbName},
		Server:   utils.ServerSettings{Port: config.Server.Port},
	}
}

func read() utils.Configuration {
	viper.SetConfigName("config")

	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	viper.SetConfigType("yml")
	var config utils.Configuration

	if err := viper.ReadInConfig(); err != nil {
		logrus.Errorf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		logrus.Errorf("Unable to decode into struct, %v", err)
	}

	logrus.Warnf("Config with variables %v", config)

	return config
}
