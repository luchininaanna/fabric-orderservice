package cmd

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

func CreateDBConnection(conf DatabaseConfig) *sql.DB {
	arguments := conf.DatabaseArguments
	if len(arguments) > 0 {
		arguments = "?" + arguments
	}

	dsn := fmt.Sprintf("%s:%s@%s/%s%s", conf.DatabaseUser, conf.DatabasePassword, conf.DatabaseAddress, conf.DatabaseName, arguments)
	db, err := sql.Open(conf.DatabaseDriver, dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Debugf("Connection to %s established", dsn)

	return db
}
