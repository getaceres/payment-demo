package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/getaceres/payment-demo/frontend"
	"github.com/getaceres/payment-demo/persistence/mongo"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	repository, err := mongo.NewMongoPaymentRepository("mongodb://localhost:27017", "payment-demo")
	if err != nil {
		fmt.Printf("Error initializing repository: %s", err.Error())
		os.Exit(-1)
	}
	frontend := frontend.FrontendV1{
		Router:            router,
		PaymentRepository: repository,
	}
	frontend.InitializeRoutes()
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Printf("Error initializing service: %s", err.Error())
		os.Exit(-1)
	}

}
