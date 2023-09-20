package shipmentcalculator

import (
	"context"
	"go.uber.org/zap"
	"sort"
)

type ShipmentCalculator struct {
	HandlerName string
	Logger      *zap.SugaredLogger
	PackSizes   []int
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

	packSizes := h.PackSizes
	// packSizes in descending order
	sort.Sort(sort.Reverse(sort.IntSlice(packSizes)))
	smallestPackSize := packSizes[len(packSizes)-1]
	biggestPackSize := packSizes[0]

	neededPacksMap := make(map[int]int)

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
	packSizes := h.PackSizes

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

	// Finalize totalItems calculation for rest of the scenarios
	for {
		size := packSizes[i]
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

		if i < len(packSizes)-1 {
			i++
		}
	}

	return totalItems
}

// calculatePackSizes calculates the smallest pack sizes possible for 'totalItems'.
func (h *ShipmentCalculator) calculatePackSizes(neededPacksMap map[int]int, totalItems int) {
	packSizes := h.PackSizes

	for _, size := range packSizes {
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
