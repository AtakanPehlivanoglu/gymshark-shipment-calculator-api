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

	var i, totalItems int
	neededPacksMap := make(map[int]int)

	// Calculate totalItems starting with smallest and biggest cases
	for {
		if itemCount < smallestPackSize { // Smaller than smallest
			neededPacksMap[smallestPackSize] += 1
			return neededPacksMap, nil
		} else if itemCount > biggestPackSize { // Bigger than biggest
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
			if itemCount == 0 {
				break
			}
			totalItems += smallestPackSize
			break
		}

		if i < len(packSizes)-1 {
			i++
		}
	}

	// Calculate the smallest pack sizes for totalItems
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

	return neededPacksMap, nil
}
