package infrastructure

import (
	log "github.com/sirupsen/logrus"
	"io"
	"strings"
)

func Close(closer io.Closer, subject ...string) {
	subjectStr := strings.Join(subject, "")

	err := closer.Close()
	if err != nil {
		log.Errorf("Failed to close %v: %v", subjectStr, err)
	}
}
