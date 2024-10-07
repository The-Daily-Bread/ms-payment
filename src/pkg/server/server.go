package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/tdb/ms-payment/src/cmd/handler"
	"github.com/tdb/ms-payment/src/cmd/services"
	"github.com/tdb/ms-payment/src/cmd/utils"
	"github.com/tdb/ms-payment/src/pkg/client/mongodb"
)

func Initialize(config utils.Configuration) {
	var log = logrus.New()

	log.WithFields(logrus.Fields{
		"mongo_url":   config.Database.Url,
		"server_port": config.Server.Port,
		"db_name":     config.Database.DbName,
	}).Info("\nConfiguration informations\n")

	logrus.Infof("Application Name %s is starting....", config.App.Name)

	client, err := mongodb.ConnectMongoDb(config.Database.Url)

	defer client.Disconnect(context.TODO())

	if err != nil {
		logrus.Fatal(err)
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Payment service is running"))
	})

	accountService := services.NewAccountService(client)
	accountHandler := handler.NewAccountHandler(accountService)

	paymentService := services.NewPaymentService(client)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /account", accountHandler.CreateAccount)
	mux.HandleFunc("POST /payment", paymentHandler.RegisterPayment)

	formattedUrl := fmt.Sprintf(":%s", config.Server.Port)

	http.ListenAndServe(formattedUrl, mux)
}
