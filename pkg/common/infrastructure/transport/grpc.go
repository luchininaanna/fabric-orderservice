package transport

import (
	"context"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"time"
)

const waitTimeout = 5 * time.Second
const connectTimeout = 15 * time.Second

func DialGRPC(address string) *grpc.ClientConn {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	timer := time.After(connectTimeout)
	for conn.GetState() != connectivity.Ready {
		select {
		case <-timer:
			log.WithField("url", address).Fatal("GRPC connection timeout")
		default:
		}

		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(waitTimeout))
		conn.WaitForStateChange(ctx, conn.GetState())
		cancel()
	}

	return conn
}
