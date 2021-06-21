package transport

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"orderservice/pkg/common/infrastructure"
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

type getOrderInfoRequest struct {
	ID string `json:"ID"`
}

type orderItemResponse struct {
	ItemId   string  `json:"item_id"`
	Quantity float32 `json:"quantity"`
}

func (s *server) getOrderInfo(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error("Can't read request body with error")
		return
	}

	defer infrastructure.LogError(r.Body.Close())

	var request getOrderInfoRequest
	err = json.Unmarshal(b, &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error("Can't parse json response with error")
		return
	}

	order, err := s.oqs.GetOrder(request.ID)
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

func (s *server) getOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := s.oqs.GetOrders()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var ordersResponse []orderResponse
	for _, order := range orders {
		var orderItems []orderItemResponse
		for _, orderItem := range order.OrderItems {
			orderItems = append(orderItems, orderItemResponse{ItemId: orderItem.ID.String(), Quantity: orderItem.Quantity})
		}

		orderStatus, err := WrapOrderStatus(order.Status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ordersResponse = append(ordersResponse, orderResponse{
			order.ID,
			orderItems,
			order.Address,
			order.Cost,
			orderStatus,
			order.CreatedAt,
		})
	}

	jsonOrders, err := json.Marshal(ordersResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if _, err = io.WriteString(w, string(jsonOrders)); err != nil {
		log.WithField("err", err).Error("write response error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
