package infrastructure

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
)

func CloseRows(rows *sql.Rows) {
	err := rows.Close()
	if err != nil {
		log.Error(err)
	}
}
