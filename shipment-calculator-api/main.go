package main

import (
	"fmt"
	"github.com/AtakanPehlivanoglu/gymshark-shipment-calculator-api/internal/app/config"
	apphandler "github.com/AtakanPehlivanoglu/gymshark-shipment-calculator-api/internal/app/handler"
	"github.com/AtakanPehlivanoglu/gymshark-shipment-calculator-api/internal/app/prepare"
	_ "github.com/AtakanPehlivanoglu/gymshark-shipment-calculator-api/shipment-calculator-api/docs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"log"
	"net/http"
)

//	@title			Gymshark Shipment Calculator API
//	@version		1.0
//	@description	Calculate number of packs that needed to be shipped.

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@host		localhost:3000
//	@BasePath	/

// Calculate godoc
//
//	@Summary	Calculate item packs
//	@Produce	json
//	@Param		itemCount	path	int	true	"Item Count"
//	@Router		/calculate/{itemCount} [get]
func main() {
	appConfig, err := config.GetConfig()
	if err != nil {
		log.Fatalf("error on GetConfig, %v", err)
	}

	logger := prepare.AppLogger(appConfig)

	// init chi router
	r := chi.NewRouter()

	// middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// business handlers
	handler, err := prepare.NewShipmentCalculatorHandler("shipmentCalculatorHandler", logger, appConfig.App.PackSizes)
	if err != nil {
		logger.Fatalf("error on prepare shipmentCalculatorHandler, %v", err)
	}

	// app handlers
	r.Get(apphandler.HealthEndpoint, apphandler.Health)

	shipmentCalculatorRoute := fmt.Sprintf("%v/{%v}",
		apphandler.ShipmentCalculatorEndpoint, apphandler.ShipmentCalculatorURLParam)

	r.Get(shipmentCalculatorRoute, apphandler.ShipmentCalculator(handler))

	// swagger
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(":3000/swagger/doc.json")), //The url pointing to API definition
	)

	logger.Fatal(http.ListenAndServe(":3000", r))
}
