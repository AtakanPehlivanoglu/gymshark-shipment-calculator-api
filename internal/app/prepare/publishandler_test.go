package prepare

import (
	"bytes"
	usecasehandler "github.com/AtakanPehlivanoglu/gymshark-shipment-calculator-api/internal/usecase/shipmentcalculator"
	"log"
	"reflect"
	"sort"
	"testing"
)

func TestNewShipmentCalculatorHandler(t *testing.T) {
	handlerName := "shipmentCalculatorHandler"
	var buff bytes.Buffer
	lgr := log.New(&buff, "", log.LstdFlags)
	packSizes := []int{250, 500, 1000, 2000, 5000}
	descendingPackSizes := make([]int, len(packSizes))
	copy(descendingPackSizes, packSizes)
	sort.Sort(sort.Reverse(sort.IntSlice(descendingPackSizes)))

	tt := []struct {
		name            string
		handlerName     string
		logger          *log.Logger
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
				LeastCommonMultiple: lcmSlice(packSizes),
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
