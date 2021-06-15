package repository

import (
	"database/sql"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"orderservice/pkg/common/infrastructure"
	"orderservice/pkg/order/model"
	"strconv"
	"strings"
	"time"
)

type orderRepository struct {
	tx *sql.Tx
}

func (or *orderRepository) Store(order model.Order) error {
	_, err := or.tx.Exec(
		"INSERT INTO `order`(id, status, cost, address, created_at, closed_at) "+
			"VALUES (UUID_TO_BIN(?), ?, ?, ?, ?, ?);", order.ID, order.Status, order.Cost, order.Address, order.CreatedAt, nil)

	if err != nil {
		err = infrastructure.InternalError(err)
	}

	for _, orderItem := range order.Items {
		_, err = or.tx.Exec(""+
			"INSERT INTO `order_item`(order_id, fabric_id, quantity) "+
			"VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?), ?);", order.ID, orderItem.ID, orderItem.Quantity)
		if err != nil {
			log.Error(or.tx.Rollback())
			return err
		}
	}

	return err
}

func (or *orderRepository) Get(orderUuid uuid.UUID) (*model.Order, error) {
	orderIdBin, err := orderUuid.MarshalBinary()
	if err != nil {
		return nil, err
	}

	rows, err := or.tx.Query(""+
		`SELECT
		   BIN_TO_UUID(o.id) AS id,
		   GROUP_CONCAT(CONCAT(BIN_TO_UUID(oi.menu_item_id), \"=\", oi.quantity)) AS menuItems,
		   o.created_at AS created_at,
		   o.cost AS cost,
		   o.status AS status,
		   o.address AS address `+
		"FROM `order` o "+
		`LEFT JOIN order_item oi ON o.id = oi.order_id 
		WHERE o.id = ? 
		GROUP BY o.id`, orderIdBin)

	if err != nil {
		return nil, err
	}

	if rows.Next() {
		order, err := parseOrder(rows)
		if err != nil {
			return nil, err
		}

		return order, nil
	}

	return nil, model.OrderNotExistError
}

func parseOrder(r *sql.Rows) (*model.Order, error) {
	var orderId string
	var menuItems string
	var createdAt time.Time
	var cost int
	var status int
	var address string

	err := r.Scan(&orderId, &menuItems, &createdAt, &cost, &status, &address)
	if err != nil {
		return nil, err
	}

	orderUuid, err := uuid.Parse(orderId)
	if err != nil {
		return nil, err
	}

	menuItemsArray := strings.Split(menuItems, ",")

	var orderItems []model.OrderItemDto
	for _, menuItem := range menuItemsArray {
		s := strings.Split(menuItem, "=")
		itemUuid, err := uuid.Parse(s[0])
		if err != nil {
			return nil, err
		}
		quantity, err := strconv.Atoi(s[1])
		if err != nil {
			return nil, err
		}

		orderItem := model.OrderItemDto{
			ID:       itemUuid,
			Quantity: quantity,
		}

		orderItems = append(orderItems, orderItem)
	}

	order, err := model.NewOrder(orderUuid, orderItems, createdAt, cost, status, address)
	if err != nil {
		return nil, err
	}

	return &order, nil
}
