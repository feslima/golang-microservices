package main

import (
	"errors"
	"net/http"

	"github.com/feslima/common"
	pb "github.com/feslima/common/api"
)

type handler struct {
	ordersClient pb.OrderServiceClient
}

func NewHandler(ordersClient pb.OrderServiceClient) *handler {
	return &handler{ordersClient}
}

func (h *handler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/customers/{customerID}/orders", h.HandleCreateOrder)
}

func (h *handler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customerID")
	var items []*pb.ItemsWithQuantity
	if err := common.ReadJSON(r, &items); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	err := validateItems(items)
	if err != nil {
		common.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	order, err := h.ordersClient.CreateOrder(r.Context(), &pb.CreateOrderRequest{
		CustomerID: customerID,
		Items:      items,
	})
	if err != nil {
		common.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	common.WriteJSON(w, http.StatusOK, order)
}

func validateItems(items []*pb.ItemsWithQuantity) error {
	if len(items) == 0 {
		return errors.New("items cannot be empty")
	}

	for _, item := range items {
		if item.ID == "" {
			return errors.New("item ID is required")
		}

		if item.Quantity <= 0 {
			return errors.New("item quantity must be positive")
		}
	}

	return nil
}
