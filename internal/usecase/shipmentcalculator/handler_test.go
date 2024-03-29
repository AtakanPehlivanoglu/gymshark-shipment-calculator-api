package shipmentcalculator

import (
	"bytes"
	"context"
	"log"
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
			expectedPackSizes:   map[int]int{250: 0, 500: 1, 1000: 0, 2000: 0, 5000: 0},
		},
		{
			name:                "happy path - 2",
			itemCount:           12001,
			ascendingPackSizes:  []int{250, 500, 1000, 2000, 5000},
			descendingPackSizes: []int{5000, 2000, 1000, 500, 250},
			expectedPackSizes:   map[int]int{250: 1, 500: 0, 1000: 0, 2000: 1, 5000: 2},
		},
		{
			name:                "happy path - different 'packSizes' configured",
			itemCount:           13501,
			ascendingPackSizes:  []int{500, 2500, 3000},
			descendingPackSizes: []int{3000, 2500, 500},
			expectedPackSizes:   map[int]int{500: 0, 2500: 2, 3000: 3},
		},
		{
			name:                "happy path - 'itemCount' smaller than smallest",
			itemCount:           1,
			ascendingPackSizes:  []int{500, 2500, 3000},
			descendingPackSizes: []int{3000, 2500, 500},
			expectedPackSizes:   map[int]int{500: 1, 2500: 0, 3000: 0},
		},
		{
			name:                "happy path - 'itemCount' is in between consecutive sum of last 2 sizes and current size",
			itemCount:           35,
			ascendingPackSizes:  []int{15, 21, 32},
			descendingPackSizes: []int{32, 21, 15},
			expectedPackSizes:   map[int]int{15: 1, 21: 1, 32: 0},
		},
		{
			name:                "happy path - 'itemCount' is too large",
			itemCount:           500000,
			ascendingPackSizes:  []int{23, 31, 53},
			descendingPackSizes: []int{53, 31, 23},
			expectedPackSizes:   map[int]int{23: 2, 31: 7, 53: 9429},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var buff bytes.Buffer
			lgr := log.New(&buff, "", log.LstdFlags)
			h := &ShipmentCalculator{
				HandlerName:         "handler_test",
				Logger:              lgr,
				AscendingPackSizes:  tc.ascendingPackSizes,
				DescendingPackSizes: tc.descendingPackSizes,
				LeastCommonMultiple: lcmSlice(tc.ascendingPackSizes),
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

// lcmSlice calculates the least common multiple for a slice of integers
func lcmSlice(numbers []int) int {
	if len(numbers) == 0 {
		return 0
	}

	result := numbers[0]
	for i := 1; i < len(numbers); i++ {
		result = calculateLcm(result, numbers[i])
	}

	return result
}

// calculateLcm calculates the least common multiple using the GCD
func calculateLcm(a, b int) int {
	return a * b / calculateGcd(a, b)
}

// calculateGcd calculates the greatest common divisor using the Euclidean algorithm
func calculateGcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
