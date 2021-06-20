package query

import (
	"database/sql"
	"github.com/google/uuid"
	"orderservice/pkg/common/infrastructure"
	"orderservice/pkg/order/application/errors"
	"orderservice/pkg/order/application/query"
	"orderservice/pkg/order/application/query/data"
	"strconv"
	"strings"
	"time"
)

func NewOrderQueryService(db *sql.DB) query.OrderQueryService {
	return &orderQueryService{db: db}
}

type orderQueryService struct {
	db *sql.DB
}

func (qs *orderQueryService) GetOrder(id string) (*data.OrderData, error) {
	rows, err := qs.db.Query(""+
		getSelectOrderSQL()+
		`WHERE o.id = UUID_TO_BIN(?) 
		GROUP BY o.id`, id)

	if err != nil {
		return nil, infrastructure.InternalError(err)
	}
	defer infrastructure.CloseRows(rows)

	if rows.Next() {
		order, err := parseOrder(rows)
		if err != nil {
			return nil, infrastructure.InternalError(err)
		}

		return order, nil
	}

	return nil, errors.OrderNotExistError
}

func getSelectOrderSQL() string {
	return `SELECT
		BIN_TO_UUID(o.id) AS id,
		GROUP_CONCAT(CONCAT(BIN_TO_UUID(oi.fabric_id), \"=\", oi.quantity)) AS menuItems,
		o.created_at AS time,
		o.cost AS cost,
		o.status AS status,
		o.address AS address` +
		"FROM `order` o" +
		`LEFT JOIN order_item oi ON o.id = oi.order_id`
}

func parseOrder(r *sql.Rows) (*data.OrderData, error) {
	var orderId string
	var orderItems string
	var t time.Time
	var cost float32
	var status int
	var address string

	err := r.Scan(&orderId, &orderItems, &t, &cost, &address)
	if err != nil {
		return nil, err
	}

	orderItemsArray := strings.Split(orderItems, ",")

	var modelOrderItems []data.OrderItemData
	for _, orderItem := range orderItemsArray {
		s := strings.Split(orderItem, "=")
		itemUuid, err := uuid.Parse(s[0])
		if err != nil {
			return nil, err
		}
		quantity, err := strconv.ParseFloat(s[1], 32)
		if err != nil {
			return nil, err
		}
		modelOrderItems = append(modelOrderItems, data.OrderItemData{ID: itemUuid, Quantity: float32(quantity)})
	}

	return &data.OrderData{
		ID:         orderId,
		OrderItems: modelOrderItems,
		CreatedAt:  t,
		Cost:       cost,
		Status:     status,
		Address:    address,
	}, nil
}
