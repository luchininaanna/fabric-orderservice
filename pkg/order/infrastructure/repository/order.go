package repository

import (
	"database/sql"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"orderservice/pkg/common/infrastructure"
	"orderservice/pkg/order/model"
)

type orderRepository struct {
	tx *sql.Tx
}

func (or *orderRepository) Add(order model.Order) error {
	_, err := or.tx.Exec(
		"INSERT INTO `order`(id, cost, status, address, created_at) "+
			"VALUES (UUID_TO_BIN(?), ?, ?, ?, ?);", order.ID, order.Cost, model.OrderStatusOrderCreated, order.Address, order.CreatedAt)

	if err != nil {
		err = infrastructure.InternalError(err)
	}

	for _, orderItem := range order.Items {
		_, err = or.tx.Exec(""+
			"INSERT INTO `order_item`(order_id, menu_item_id, quantity) "+
			"VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?), ?);", order.ID, orderItem.ID, orderItem.Quantity)
		if err != nil {
			log.Error(or.tx.Rollback())
			return err
		}
	}

	return err
}

func (or *orderRepository) Delete(orderUuid uuid.UUID) error {
	_, err := or.tx.Exec(""+
		"UPDATE `order` "+
		"SET status = ?"+
		"WHERE id = UUID_TO_BIN(?);", model.OrderStatusOrderCanceled, orderUuid)
	if err != nil {
		return err
	}

	return nil
}
