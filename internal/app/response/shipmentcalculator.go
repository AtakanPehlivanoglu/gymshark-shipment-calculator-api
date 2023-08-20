package response

import (
	"github.com/go-chi/render"
	"net/http"
)

type ShipmentCalculatorResponse struct {
	HTTPStatusCode int `json:"-"` // http response status code

	Message   string      `json:"message"`   // user-level status message
	ItemPacks map[int]int `json:"itemPacks"` // calculated item packs for given number of ordered items
}

func (e *ShipmentCalculatorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func NewShipmentCalculatorResponse(statusCode int, message string, itemPacks map[int]int) *ShipmentCalculatorResponse {
	return &ShipmentCalculatorResponse{
		HTTPStatusCode: statusCode,
		Message:        message,
		ItemPacks:      itemPacks,
	}
}
