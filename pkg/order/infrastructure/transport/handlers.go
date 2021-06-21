package transport

import (
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"orderservice/pkg/common/cmd"
	"orderservice/pkg/common/infrastructure"
	"orderservice/pkg/order/application/command"
	"orderservice/pkg/order/application/query"
	queryImpl "orderservice/pkg/order/infrastructure/query"
	"orderservice/pkg/order/infrastructure/repository"
)

type server struct {
	unitOfWork command.UnitOfWork
	oqs        query.OrderQueryService
}

type orderItem struct {
	ID       string  `json:"ID"`
	Quantity float32 `json:"quantity"`
}

type addOrderRequest struct {
	Items   []orderItem `json:"items"`
	Address string      `json:"address"`
}

type closeOrderRequest struct {
	ID string `json:"ID"`
}

type startProcessingOrderRequest struct {
	ID string `json:"ID"`
}

type addOrderResponse struct {
	Id string `json:"id"`
}

func Router(db *sql.DB) http.Handler {
	srv := &server{
		repository.NewUnitOfWork(db),
		queryImpl.NewOrderQueryService(db),
	}
	r := mux.NewRouter()

	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/order", srv.addOrder).Methods(http.MethodPost)
	s.HandleFunc("/close/order", srv.closeOrder).Methods(http.MethodPost)
	s.HandleFunc("/process/order", srv.startProcessingOrder).Methods(http.MethodPost)
	s.HandleFunc("/order", srv.getOrderInfo).Methods(http.MethodGet)
	s.HandleFunc("/orders", srv.getOrders).Methods(http.MethodGet)
	return cmd.LogMiddleware(r)
}

func (s *server) addOrder(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal("Can't read request body with error")
		return
	}

	defer infrastructure.LogError(r.Body.Close())

	var request addOrderRequest
	err = json.Unmarshal(b, &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error("Can't parse json response with error")
		return
	}

	if len(request.Items) == 0 {
		http.Error(w, "Empty item list", http.StatusBadRequest)
		return
	}

	orderItemList, err := convertOrderItem(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error("Can't parse json response with error")
		return
	}

	var h = command.NewAddOrderCommandHandler(s.unitOfWork)
	id, err := h.Handle(command.AddOrderCommand{
		Items:   orderItemList,
		Address: request.Address,
	})

	if err != nil {
		http.Error(w, WrapError(err).Error(), http.StatusBadRequest)
		return
	}

	RenderJson(w, &addOrderResponse{id.String()})
}

func convertOrderItem(request addOrderRequest) ([]command.OrderItem, error) {
	orderItemList := make([]command.OrderItem, 0)
	for _, item := range request.Items {
		orderItemUid, err := uuid.Parse(item.ID)
		if err != nil {
			log.Error("Can't parse order item uid")
			return orderItemList, err
		}

		orderItem := command.OrderItem{
			ID:       orderItemUid,
			Quantity: item.Quantity,
		}

		orderItemList = append(orderItemList, orderItem)
	}

	return orderItemList, nil
}

func (s *server) closeOrder(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error("Can't read request body with error")
		return
	}

	defer infrastructure.LogError(r.Body.Close())

	var request closeOrderRequest
	err = json.Unmarshal(b, &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error("Can't parse json response with error")
		return
	}

	var h = command.NewCloseOrderCommandHandler(s.unitOfWork)

	orderUid, err := uuid.Parse(request.ID)
	if err != nil {
		log.Error("Can't parse order uid")
		return
	}

	err = h.Handle(command.CloseOrderCommand{ID: orderUid})
	if err != nil {
		http.Error(w, WrapError(err).Error(), http.StatusBadRequest)
		return
	}
}

func (s *server) startProcessingOrder(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error("Can't read request body with error")
		return
	}

	defer infrastructure.LogError(r.Body.Close())

	var request startProcessingOrderRequest
	err = json.Unmarshal(b, &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error("Can't parse json response with error")
		return
	}

	var h = command.NewStartProcessingOrderCommandHandler(s.unitOfWork)

	orderUid, err := uuid.Parse(request.ID)
	if err != nil {
		log.Error("Can't parse order uid")
		return
	}

	err = h.Handle(command.StartProcessingOrderCommand{ID: orderUid})
	if err != nil {
		http.Error(w, WrapError(err).Error(), http.StatusBadRequest)
		return
	}
}
