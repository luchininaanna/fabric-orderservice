package transport

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

type orderResponse struct {
	OrderId    string              `json:"order_id"`
	OrderItems []orderItemResponse `json:"orderItems"`
	Address    string              `json:"address"`
	Cost       float32             `json:"cost"`
	Status     string              `json:"status"`
	CreatedAt  time.Time           `json:"created_at"`
}

type orderItemResponse struct {
	ItemId   string  `json:"item_id"`
	Quantity float32 `json:"quantity"`
}

func (s *server) getOrderInfo(w http.ResponseWriter, r *http.Request) {
	id, found := mux.Vars(r)["ID"]
	if !found {
		w.WriteHeader(http.StatusBadRequest)
		_, err := fmt.Fprint(w, "Order not found")
		if err != nil {
			log.Error(err)
		}
		return
	}

	order, err := s.oqs.GetOrder(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var orderItems []orderItemResponse
	for _, orderItem := range order.OrderItems {
		orderItems = append(orderItems, orderItemResponse{ItemId: orderItem.ID.String(), Quantity: orderItem.Quantity})
	}

	orderStatus, err := WrapOrderStatus(order.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonOrder, err := json.Marshal(orderResponse{
		order.ID,
		orderItems,
		order.Address,
		order.Cost,
		orderStatus,
		order.CreatedAt,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if _, err = io.WriteString(w, string(jsonOrder)); err != nil {
		log.WithField("err", err).Error("write response error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
