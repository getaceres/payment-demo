//go:generate swagger generate spec
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/getaceres/payment-demo/frontend"
	"github.com/getaceres/payment-demo/persistence/mongo"
	"github.com/gorilla/mux"

	"github.com/spf13/cobra"
)

func main() {

	var port int
	var connectionURL string

	var cmdServe = &cobra.Command{
		Use:   "serve",
		Short: "Start the API server",
		Long:  `This will start the server listening in the provided or the default port`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			startServer(port, connectionURL)
		},
	}

	cmdServe.Flags().IntVarP(&port, "port", "p", 8080, "Port to serve")
	cmdServe.Flags().StringVarP(&connectionURL, "mongourl", "m", "mongodb://localhost:27017", "Connection URL to a MongoDB database")

	var rootCmd = &cobra.Command{Use: "payment-demo"}
	rootCmd.AddCommand(cmdServe)

	rootCmd.Execute()
}

func startServer(port int, connectionURL string) {
	router := mux.NewRouter()
	repository, err := mongo.NewMongoPaymentRepository(connectionURL, "payment-demo")
	if err != nil {
		fmt.Printf("Error initializing MongoDB repository: %s", err.Error())
		os.Exit(-1)
	}
	frontend := frontend.FrontendV1{
		Router:            router,
		PaymentRepository: repository,
	}
	frontend.InitializeRoutes()
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	if err != nil {
		fmt.Printf("Error initializing service: %s", err.Error())
		os.Exit(-1)
	}
}
