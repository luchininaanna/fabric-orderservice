package transport

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"orderservice/pkg/common/errors"
	"orderservice/pkg/common/infrastructure"
	"time"
)

var grpcServeMuxOptions = &runtime.JSONPb{
	EmitDefaults: true,
	OrigName:     true,
}

func NewServeMux() *runtime.ServeMux {
	return runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, grpcServeMuxOptions))
}

func LogMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		h.ServeHTTP(w, r)
		log.WithFields(log.Fields{
			"method":     r.Method,
			"url":        r.URL,
			"remoteAddr": r.RemoteAddr,
			"userAgent":  r.UserAgent(),
			"duration":   time.Since(startTime).String(),
			"at":         startTime,
		}).Info("got request")
	})
}

func RenderJson(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Error(err)
		ProcessError(w, errors.InternalError)
		return
	}
}

func ProcessError(w http.ResponseWriter, e error) {
	if e == errors.InternalError {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	} else {
		http.Error(w, e.Error(), http.StatusBadRequest)
	}
}

func Parameter(r *http.Request, key string) (string, bool) {
	val, found := mux.Vars(r)[key]
	return val, found
}

func CloseService(closer io.Closer, subject ...string) {
	log.Infof("Close %v", subject)
	infrastructure.Close(closer, subject...)
}
