package infrastructure

import (
	log "github.com/sirupsen/logrus"
	"orderservice/pkg/common/errors"
)

func InternalError(e error) error {
	log.Error(e)
	return errors.InternalError
}
