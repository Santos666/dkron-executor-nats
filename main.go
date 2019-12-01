package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/distribworks/dkron/v2/dkron"
	dkplugin "github.com/distribworks/dkron/v2/plugin"
	"github.com/hashicorp/go-plugin"
	stan "github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
)

var log = logrus.NewEntry(logrus.New())

func main() {
	nats := &NATS{}
	go nats.Connect()

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: dkplugin.Handshake,
		Plugins: map[string]plugin.Plugin{
			"executor": &dkplugin.ExecutorPlugin{Executor: nats},
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})

	nats.Close()
}

type NATS struct {
	conn  stan.Conn
	mutex sync.Mutex
}

func (e *NATS) Execute(args *dkron.ExecuteRequest) (*dkron.ExecuteResponse, error) {
	subject, ok := args.Config["subject"]
	if !ok || "" == subject {
		subject = "dkron"
	}

	message, ok := args.Config["message"]
	if !ok || "" == message {
		return &dkron.ExecuteResponse{Error: "nats: Invalid message"}, nil
	}

	if nil == e.conn {
		err := e.Connect()
		if nil != err {
			return nil, err
		}
	}

	nuid, err := e.Publish(subject, []byte(message))
	if nil != err {
		return &dkron.ExecuteResponse{Error: err.Error()}, nil
	}

	msg := fmt.Sprintf("Message NUID: %s", nuid)

	return &dkron.ExecuteResponse{Output: []byte(msg)}, nil
}

func (e *NATS) Publish(subject string, data []byte) (string, error) {
	nuid, err := e.conn.PublishAsync(subject, data, e.ackHandler)
	if nil != err {
		log.WithError(err).WithField("nuid", nuid).Error("nats: Failed to publish async")
		return "", err
	}

	return nuid, nil
}

func (e *NATS) Connect() error {
	clientID := os.Getenv("NATS_CLIENT_ID")
	if "" == clientID {
		clientID = "dkron-executor-nats"
	}

	clusterID := os.Getenv("NATS_CLUSTER_ID")
	if "" == clusterID {
		clusterID = "test-cluster"
	}

	uri := os.Getenv("NATS_URI")
	if "" == uri {
		uri = "nats://localhost:4222"
	}

	opts := []stan.Option{
		stan.NatsURL(uri),
		stan.Pings(3, 5),
		stan.SetConnectionLostHandler(e.connectionLostHandler),
	}

	conn, err := stan.Connect(clusterID, clientID, opts...)
	if nil != err {
		return err
	}

	e.mutex.Lock()
	e.conn = conn
	e.mutex.Unlock()

	return nil
}

func (e *NATS) Close() {
	if nil == e.conn {
		return
	}

	e.conn.Close()
}

func (e *NATS) connectionLostHandler(conn stan.Conn, err error) {
	log.WithError(err).Error("nats: Connection lost")

	for range time.Tick(3 * time.Second) {
		err := e.Connect()

		if nil != err {
			log.WithError(err).Error("nats: Failing to reconnect, retrying...")
			continue
		}

		log.Info("nats: Successfully reconnected")
		break
	}
}

func (e *NATS) ackHandler(nuid string, err error) {
	if nil != err {
		log.WithError(err).WithField("nuid", nuid).Error("nats: Failed to ack message")
	}
}
