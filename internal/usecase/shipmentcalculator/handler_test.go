package shipmentcalculator

import (
	"context"
	"go.uber.org/zap"
	"reflect"
	"testing"
)

func TestShipmentCalculator_Handle(t *testing.T) {
	tt := []struct {
		name                string
		itemCount           int
		ascendingPackSizes  []int
		descendingPackSizes []int
		expectedPackSizes   map[int]int
	}{
		{
			name:                "happy path - 1",
			itemCount:           251,
			ascendingPackSizes:  []int{250, 500, 1000, 2000, 5000},
			descendingPackSizes: []int{5000, 2000, 1000, 500, 250},
			expectedPackSizes:   map[int]int{500: 1},
		},
		{
			name:                "happy path - 2",
			itemCount:           12001,
			ascendingPackSizes:  []int{250, 500, 1000, 2000, 5000},
			descendingPackSizes: []int{5000, 2000, 1000, 500, 250},
			expectedPackSizes:   map[int]int{250: 1, 2000: 1, 5000: 2},
		},
		{
			name:                "happy path - different 'packSizes' configured",
			itemCount:           13501,
			ascendingPackSizes:  []int{500, 2500, 3000},
			descendingPackSizes: []int{3000, 2500, 500},
			expectedPackSizes:   map[int]int{500: 4, 3000: 4},
		},
		{
			name:                "happy path - 'itemCount' smaller than smallest",
			itemCount:           1,
			ascendingPackSizes:  []int{500, 2500, 3000},
			descendingPackSizes: []int{3000, 2500, 500},
			expectedPackSizes:   map[int]int{500: 1},
		},
		{
			name:                "happy path - 'itemCount' is in between consecutive sum of last 2 sizes and current size",
			itemCount:           35,
			ascendingPackSizes:  []int{15, 21, 32},
			descendingPackSizes: []int{32, 21, 15},
			expectedPackSizes:   map[int]int{15: 1, 21: 1},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			lgr := zap.NewExample().Sugar()
			h := &ShipmentCalculator{
				HandlerName:         "handler_test",
				Logger:              lgr,
				AscendingPackSizes:  tc.ascendingPackSizes,
				DescendingPackSizes: tc.descendingPackSizes,
			}

			got, err := h.Handle(context.Background(), tc.itemCount)
			if err != nil {
				t.Errorf("unexpected Handle() error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, tc.expectedPackSizes) {
				t.Errorf("Handle() got = %v, want %v", got, tc.expectedPackSizes)
			}
		})
	}
}

func TestShipmentCalculator_calculateTotalItems(t *testing.T) {
	type args struct {
	}
	tt := []struct {
		name               string
		itemCount          int
		smallestPackSize   int
		biggestPackSize    int
		expectedTotalItems int
	}{
		{
			name:               "happy path",
			itemCount:          251,
			smallestPackSize:   250,
			biggestPackSize:    5000,
			expectedTotalItems: 500,
		},
		{
			name:               "happy path - bigger than biggest",
			itemCount:          5001,
			smallestPackSize:   250,
			biggestPackSize:    5000,
			expectedTotalItems: 5250,
		},
		{
			name:               "happy path - item count is multiple of biggest pack size",
			itemCount:          500000,
			smallestPackSize:   250,
			biggestPackSize:    5000,
			expectedTotalItems: 500000,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			lgr := zap.NewExample().Sugar()
			h := &ShipmentCalculator{
				HandlerName:         "handler_test",
				Logger:              lgr,
				DescendingPackSizes: []int{250, 500, 1000, 2000, 5000},
			}
			got := h.calculateTotalItems(tc.itemCount, tc.smallestPackSize, tc.biggestPackSize)
			if got != tc.expectedTotalItems {
				t.Errorf("calculateTotalItems() = %v, want %v", got, tc.expectedTotalItems)
			}
		})
	}
}
