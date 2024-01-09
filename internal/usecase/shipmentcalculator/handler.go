package shipmentcalculator

import (
	"context"
	"go.uber.org/zap"
	"math"
)

type ShipmentCalculator struct {
	HandlerName         string
	Logger              *zap.SugaredLogger
	AscendingPackSizes  []int
	DescendingPackSizes []int
	LeastCommonMultiple int
}

type Handler interface {
	Handle(ctx context.Context, itemCount int) (map[int]int, error)
}

func (h *ShipmentCalculator) Handle(ctx context.Context, itemCount int) (map[int]int, error) {
	h.Logger.Infof("%v started for item count %v ", h.HandlerName, itemCount)
	initialItemCount := itemCount
	defer func() {
		h.Logger.Infof("%v finished for item count %v ", h.HandlerName, initialItemCount)
	}()

	ascendingPackSizes := h.AscendingPackSizes
	lcm := h.LeastCommonMultiple
	neededPacksMap := make(map[int]int)
	biggestPackSize := ascendingPackSizes[len(ascendingPackSizes)-1]
	nearestCombination := make([]int, len(ascendingPackSizes))
	smallestDifference := math.MaxInt64
	currentCombination := make([]int, len(ascendingPackSizes))

	// update the target for larger values than the lcm for efficiency
	if lcm > biggestPackSize && itemCount > lcm {
		division := lcm / biggestPackSize
		mod := itemCount % lcm
		multiplier := (itemCount - mod) / lcm
		itemCount = mod
		multiplier = multiplier * division
		neededPacksMap[biggestPackSize] = multiplier
	}

	iterationRange := fixIterationRange(ascendingPackSizes, itemCount)
	h.findCombinationRecursively(itemCount, 0, iterationRange, &currentCombination, &nearestCombination, &smallestDifference)
	// update the neededPacksMap using the nearestCombination slice after recursion is done
	for i, size := range ascendingPackSizes {
		if size == biggestPackSize {
			neededPacksMap[biggestPackSize] += nearestCombination[len(nearestCombination)-1]
			break
		}
		neededPacksMap[size] = nearestCombination[i]
	}

	return neededPacksMap, nil
}

func (h *ShipmentCalculator) findCombinationRecursively(itemCount int, index, iterationRange int, currentCombination, nearestCombination *[]int, smallestDifference *int) {
	ascendingPackSizes := h.AscendingPackSizes
	if index == len(ascendingPackSizes) {
		// calculate the sum and difference
		currentSum := calculateSum(*currentCombination, ascendingPackSizes)
		difference := abs(currentSum - itemCount)

		// update the nearest combination if needed
		if difference < *smallestDifference && itemCount <= currentSum {
			*smallestDifference = difference
			copy(*nearestCombination, *currentCombination)

		}
		// update the nearest combination with the least packages possible
		if currentSum == calculateSum(*nearestCombination, ascendingPackSizes) {
			nearestCombinationPackSize := findTotalPackSize(*nearestCombination)
			currentCombinationPackSize := findTotalPackSize(*currentCombination)

			if currentCombinationPackSize < nearestCombinationPackSize {
				copy(*nearestCombination, *currentCombination)
			}
		}
		return
	}

	// recursive call for the current index
	for i := 0; i <= iterationRange/ascendingPackSizes[index]; i++ {
		(*currentCombination)[index] = i
		h.findCombinationRecursively(itemCount, index+1, iterationRange, currentCombination, nearestCombination, smallestDifference)
	}
}

// fixIterationRange updates iteration range for special cases to have all combinations properly
func fixIterationRange(numbers []int, target int) int {
	iterationRange := 0
	for i, number := range numbers {
		// target is smaller than the smallest size
		if i == 0 && target < number {
			iterationRange = numbers[i]
			return iterationRange
		}
		// round to nearest bigger size for in between case
		if target < numbers[i] && target > numbers[i-1] {
			iterationRange = numbers[i]
			return iterationRange
		}
	}
	return target
}

// abs function for calculating absolute value
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func calculateSum(combination []int, numbers []int) int {
	sum := 0
	for i := range combination {
		sum += combination[i] * numbers[i]
	}
	return sum
}

func findTotalPackSize(packs []int) int {
	sum := 0
	for _, pack := range packs {
		sum += pack
	}
	return sum
}
