package transport

import (
	"context"
	"database/sql"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
	orderservice "orderservice/api"
	"orderservice/pkg/common/cmd"
	"orderservice/pkg/common/infrastructure/transport"
	"orderservice/pkg/order/application/command"
	"orderservice/pkg/order/application/query"
	queryImpl "orderservice/pkg/order/infrastructure/query"
	"orderservice/pkg/order/infrastructure/repository"
)

type server struct {
	unitOfWork command.UnitOfWork
	oqs        query.OrderQueryService
}

func (s *server) CreateOrder(_ context.Context, request *orderservice.CreateOrderRequest) (*orderservice.CreateOrderResponse, error) {
	if len(request.Items) == 0 {
		return nil, OrderWithEmptyItemListError
	}

	orderItemList, err := convertOrderItem(request)
	if err != nil {
		log.Error("Can't parse json response with error")
		return nil, err
	}

	var h = command.NewAddOrderCommandHandler(s.unitOfWork)
	id, err := h.Handle(command.AddOrderCommand{
		Items:   orderItemList,
		Address: request.Address,
	})

	if err != nil {
		return nil, err
	}

	return &orderservice.CreateOrderResponse{Id: id.String()}, nil
}

func (s *server) CloseOrder(_ context.Context, request *orderservice.CloseOrderRequest) (*empty.Empty, error) {
	var h = command.NewCloseOrderCommandHandler(s.unitOfWork)

	orderUid, err := uuid.Parse(request.Id)
	if err != nil {
		return nil, err
	}

	err = h.Handle(command.CloseOrderCommand{ID: orderUid})
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *server) StartProcessingOrder(_ context.Context, request *orderservice.StartProcessingOrderRequest) (*empty.Empty, error) {
	var h = command.NewStartProcessingOrderCommandHandler(s.unitOfWork)

	orderUid, err := uuid.Parse(request.Id)
	if err != nil {
		return nil, err
	}

	err = h.Handle(command.StartProcessingOrderCommand{ID: orderUid})
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *server) GetOrderInfo(_ context.Context, request *orderservice.GetOrderInfoRequest) (*orderservice.OrderResponse, error) {
	order, err := s.oqs.GetOrder(request.Id)
	if err != nil {
		return nil, err
	}

	var orderItems []*orderservice.OrderResponse_OrderItems
	for _, orderItem := range order.OrderItems {
		orderItems = append(orderItems, &orderservice.OrderResponse_OrderItems{
			ItemId:   orderItem.ID.String(),
			Quantity: orderItem.Quantity,
		})
	}

	orderStatus, err := WrapOrderStatus(order.Status)
	if err != nil {
		return nil, err
	}

	response := orderservice.OrderResponse{
		OrderId:   order.ID,
		Items:     orderItems,
		Address:   order.Address,
		Cost:      order.Cost,
		Status:    orderStatus,
		CreatedAt: order.CreatedAt.String(),
	}

	return &response, nil
}

func (s *server) GetOrders(_ context.Context, empty *empty.Empty) (*orderservice.OrdersResponse, error) {
	orders, err := s.oqs.GetOrders()
	if err != nil {
		return nil, err
	}

	var ordersResponseList []*orderservice.OrderResponse
	for _, order := range orders {
		var orderItems []*orderservice.OrderResponse_OrderItems
		for _, orderItem := range order.OrderItems {
			orderItems = append(orderItems, &orderservice.OrderResponse_OrderItems{
				ItemId:   orderItem.ID.String(),
				Quantity: orderItem.Quantity,
			})
		}

		orderStatus, err := WrapOrderStatus(order.Status)
		if err != nil {
			return nil, err
		}

		ordersResponseList = append(ordersResponseList, &orderservice.OrderResponse{
			OrderId:   order.ID,
			Items:     orderItems,
			Address:   order.Address,
			Cost:      order.Cost,
			Status:    orderStatus,
			CreatedAt: order.CreatedAt.String(),
		})
	}

	response := orderservice.OrdersResponse{
		Orders: ordersResponseList,
	}

	return &response, nil
}

func Router(db *sql.DB) http.Handler {
	srv := &server{
		repository.NewUnitOfWork(db),
		queryImpl.NewOrderQueryService(db),
	}

	router := transport.NewServeMux()
	err := orderservice.RegisterOrderServiceHandlerServer(context.Background(), router, srv)
	if err != nil {
		log.Fatal(err)
	}

	return cmd.LogMiddleware(router)
}

func convertOrderItem(request *orderservice.CreateOrderRequest) ([]command.OrderItem, error) {
	orderItemList := make([]command.OrderItem, 0)
	for _, item := range request.Items {
		orderItemUid, err := uuid.Parse(item.Id)
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
