package prepare

import (
	"errors"
	usecasehandler "github.com/AtakanPehlivanoglu/gymshark-shipment-calculator-api/internal/usecase/shipmentcalculator"
	"log"
	"sort"
)

func NewShipmentCalculatorHandler(handlerName string, logger *log.Logger, packSizes []int) (*usecasehandler.ShipmentCalculator, error) {
	if handlerName == "" {
		return nil, errors.New("handler name is missing")
	}
	if logger == nil {
		return nil, errors.New("logger is missing")
	}
	if len(packSizes) <= 0 {
		return nil, errors.New("pack sizes are missing")
	}

	descendingPackSizes := make([]int, len(packSizes))
	copy(descendingPackSizes, packSizes)
	sort.Sort(sort.Reverse(sort.IntSlice(descendingPackSizes)))
	// least common multiple
	lcm := lcmSlice(packSizes)

	return &usecasehandler.ShipmentCalculator{
		HandlerName:         handlerName,
		Logger:              logger,
		AscendingPackSizes:  packSizes,
		DescendingPackSizes: descendingPackSizes,
		LeastCommonMultiple: lcm,
	}, nil
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
