package transport

import (
	"database/sql"
	"encoding/json"
	"fmt"
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

func Router(db *sql.DB) http.Handler {
	srv := &server{
		repository.NewUnitOfWork(db),
		queryImpl.NewOrderQueryService(db),
	}
	r := mux.NewRouter()

	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/order", srv.addOrder).Methods(http.MethodPost)
	s.HandleFunc("/order/{ID:[0-9a-zA-Z-]+}", srv.closeOrder).Methods(http.MethodPost)
	s.HandleFunc("/order/{ID:[0-9a-zA-Z-]+}", srv.getOrderInfo).Methods(http.MethodGet)
	return cmd.LogMiddleware(r)
}

type addOrderResponse struct {
	Id string `json:"id"`
}

func (s *server) addOrder(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal("Can't read request body with error")
		return
	}

	defer infrastructure.LogError(r.Body.Close())

	var addOrderCommand command.AddOrderCommand
	err = json.Unmarshal(b, &addOrderCommand)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal("Can't parse json response with error")
		return
	}

	if len(addOrderCommand.Items) == 0 {
		http.Error(w, "Empty item list", http.StatusBadRequest)
		return
	}

	var h = command.NewAddOrderCommandHandler(s.unitOfWork)
	id, err := h.Handle(addOrderCommand)
	if err != nil {
		http.Error(w, WrapError(err).Error(), http.StatusBadRequest)
		return
	}

	RenderJson(w, &addOrderResponse{id.String()})
}

func (s *server) closeOrder(w http.ResponseWriter, r *http.Request) {
	id, found := mux.Vars(r)["ID"]
	if !found {
		w.WriteHeader(http.StatusBadRequest)
		_, err := fmt.Fprint(w, "Order not found")
		if err != nil {
			log.Error(err)
		}
		return
	}

	var h = command.NewCloseOrderCommandHandler(s.unitOfWork)
	err := h.Handle(command.CloseOrderCommand{ID: id})
	if err != nil {
		http.Error(w, WrapError(err).Error(), http.StatusBadRequest)
		return
	}
}
