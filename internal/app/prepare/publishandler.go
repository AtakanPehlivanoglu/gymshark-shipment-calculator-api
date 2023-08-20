package prepare

import (
	"errors"
	usecasehandler "github.com/AtakanPehlivanoglu/gymshark-shipment-calculator-api/internal/usecase/shipmentcalculator"
	"go.uber.org/zap"
)

func NewShipmentCalculatorHandler(handlerName string, logger *zap.SugaredLogger, packSizes []int) (*usecasehandler.ShipmentCalculator, error) {
	if handlerName == "" {
		return nil, errors.New("handler name is missing")
	}
	if logger == nil {
		return nil, errors.New("logger is missing")
	}
	if len(packSizes) <= 0 {
		return nil, errors.New("pack sizes are missing")
	}
	return &usecasehandler.ShipmentCalculator{
		HandlerName: handlerName,
		Logger:      logger,
		PackSizes:   packSizes,
	}, nil
}
