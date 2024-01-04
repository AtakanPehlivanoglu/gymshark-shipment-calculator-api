package shipmentcalculator

import (
	"context"
	"go.uber.org/zap"
)

type ShipmentCalculator struct {
	HandlerName         string
	Logger              *zap.SugaredLogger
	AscendingPackSizes  []int
	DescendingPackSizes []int
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
	neededPacksMap := make(map[int]int)
	smallestPackSize := ascendingPackSizes[0]
	biggestPackSize := ascendingPackSizes[len(ascendingPackSizes)-1]

	// return if (size(n-1) + size(n-2)) <= itemCount < size(n) special case satisfies.
	for i := range ascendingPackSizes {
		if i >= 2 && itemCount <= (ascendingPackSizes[i-1]+ascendingPackSizes[i-2]) && itemCount > ascendingPackSizes[i] {
			consecutiveSum := ascendingPackSizes[i-1] + ascendingPackSizes[i-2]
			// first check the distance of item count to consecutive sum and current pack size + n.th size
			if isConsecutiveSumBigger(consecutiveSum, ascendingPackSizes[i], itemCount, ascendingPackSizes) {
				break
			}
			neededPacksMap[ascendingPackSizes[i-1]] = 1
			neededPacksMap[ascendingPackSizes[i-2]] = 1
			return neededPacksMap, nil
		}
	}
	// return if smaller than smallest
	if itemCount < smallestPackSize {
		neededPacksMap[smallestPackSize] += 1
		return neededPacksMap, nil
	}

	totalItems := h.calculateTotalItems(itemCount, smallestPackSize, biggestPackSize)

	h.calculatePackSizes(neededPacksMap, totalItems)

	return neededPacksMap, nil
}

// calculateTotalItems finds total number of items needed with respect to 'packSizes' for given 'itemCount'.
func (h *ShipmentCalculator) calculateTotalItems(itemCount, smallestPackSize, biggestPackSize int) int {
	var i, totalItems int
	descendingPackSizes := h.DescendingPackSizes

	// Calculate totalItems starting with the biggest case
	for {
		if itemCount > biggestPackSize {
			packCount := itemCount / biggestPackSize
			totalItems += packCount * biggestPackSize
			itemCount = itemCount % biggestPackSize
			break
		} else {
			break
		}
	}

	// Return totalItems in case of itemCount is multiple of biggest pack size.
	if itemCount == 0 {
		return totalItems
	}

	// Finalize totalItems calculation for rest of the scenarios
	for {
		size := descendingPackSizes[i]
		if itemCount == size { // equal sizes
			totalItems += size
			break
		} else if itemCount > size { // itemCount bigger than size
			count := itemCount / size
			totalItems += count * size
			mod := itemCount % size
			itemCount = mod
		} else if itemCount < smallestPackSize { // itemCount smaller than smallest
			totalItems += smallestPackSize
			break
		}

		if i < len(descendingPackSizes)-1 {
			i++
		}
	}

	return totalItems
}

// calculatePackSizes calculates the smallest pack sizes possible for 'totalItems'.
func (h *ShipmentCalculator) calculatePackSizes(neededPacksMap map[int]int, totalItems int) {
	descendingPackSizes := h.DescendingPackSizes

	for _, size := range descendingPackSizes {
		if totalItems == size {
			neededPacksMap[size] = 1
			break
		}
		if totalItems > size {
			count := totalItems / size
			mod := totalItems % size

			neededPacksMap[size] = count
			totalItems = mod
		}
	}
}

func isConsecutiveSumBigger(consecutiveSum, currentPackSize, itemCount int, packSizes []int) bool {
	for _, size := range packSizes {
		if currentPackSize+size < consecutiveSum && itemCount <= currentPackSize+size {
			diff1 := currentPackSize + size - itemCount
			diff2 := consecutiveSum - itemCount

			if diff1 > diff2 {
				return false
			} else {
				return true
			}
		}
	}
	return false
}
