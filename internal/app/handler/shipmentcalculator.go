package handler

import (
	"encoding/json"
	"fmt"
	"github.com/AtakanPehlivanoglu/gymshark-shipment-calculator-api/internal/app/response"
	"github.com/AtakanPehlivanoglu/gymshark-shipment-calculator-api/internal/usecase/shipmentcalculator"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

const ShipmentCalculatorURLParam = "itemCount"
const ShipmentCalculatorEndpoint = "/calculate"

func ShipmentCalculator(handler shipmentcalculator.Handler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		itemCountStr := chi.URLParam(r, ShipmentCalculatorURLParam)
		itemCount, err := strconv.Atoi(itemCountStr)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response.ErrInvalidRequest(err))
			return
		}

		if itemCount <= 0 {
			err = fmt.Errorf("item count should be greater than 0")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response.ErrInvalidRequest(err))
			return
		}

		packCount, err := handler.Handle(ctx, itemCount)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(response.ErrInternalServer(err))
			return
		}

		shipmentCalculatorResponse := response.NewShipmentCalculatorResponse(http.StatusOK, "Number of Item Packs", packCount)

		jsonResponse, err := json.Marshal(shipmentCalculatorResponse)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(response.ErrInternalServer(err))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}
