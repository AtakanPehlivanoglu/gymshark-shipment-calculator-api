package prepare

import (
	usecasehandler "github.com/AtakanPehlivanoglu/gymshark-shipment-calculator-api/internal/usecase/shipmentcalculator"
	"go.uber.org/zap"
	"reflect"
	"sort"
	"testing"
)

func TestNewShipmentCalculatorHandler(t *testing.T) {
	handlerName := "shipmentCalculatorHandler"
	lgr := zap.NewExample().Sugar()
	packSizes := []int{250, 500, 1000, 2000, 5000}
	descendingPackSizes := make([]int, len(packSizes))
	copy(descendingPackSizes, packSizes)
	sort.Sort(sort.Reverse(sort.IntSlice(descendingPackSizes)))

	tt := []struct {
		name            string
		handlerName     string
		logger          *zap.SugaredLogger
		packSizes       []int
		expectedHandler *usecasehandler.ShipmentCalculator
		expectedError   bool
	}{
		{
			name:        "happy path",
			handlerName: handlerName,
			logger:      lgr,
			packSizes:   packSizes,
			expectedHandler: &usecasehandler.ShipmentCalculator{
				HandlerName:         handlerName,
				Logger:              lgr,
				AscendingPackSizes:  packSizes,
				DescendingPackSizes: descendingPackSizes,
			},
			expectedError: false,
		},
		{
			name:            "handler name missing",
			handlerName:     "",
			logger:          lgr,
			packSizes:       packSizes,
			expectedHandler: nil,
			expectedError:   true,
		},
		{
			name:            "logger missing",
			handlerName:     handlerName,
			logger:          nil,
			packSizes:       packSizes,
			expectedHandler: nil,
			expectedError:   true,
		},
		{
			name:            "pack sizes are missing",
			handlerName:     handlerName,
			logger:          lgr,
			packSizes:       []int{},
			expectedHandler: nil,
			expectedError:   true,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewShipmentCalculatorHandler(tc.handlerName, tc.logger, tc.packSizes)
			if (err != nil) != tc.expectedError {
				t.Errorf("NewShipmentCalculatorHandler() error = %v, wantErr %v", err, tc.expectedError)
				return
			}
			if !reflect.DeepEqual(got, tc.expectedHandler) {
				t.Errorf("NewShipmentCalculatorHandler() got = %v, want %v", got, tc.expectedHandler)
			}
		})
	}
}
