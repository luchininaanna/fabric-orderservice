package infrastructure

import log "github.com/sirupsen/logrus"

func LogError(err error) {
	if err != nil {
		log.Error()
	}
}
