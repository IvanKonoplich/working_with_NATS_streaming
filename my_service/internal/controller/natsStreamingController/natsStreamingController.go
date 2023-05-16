package natsStreamingController

import (
	"fmt"
	nats "github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
	"learning_NATS_streaming/internal/entities"
	"log"
	"os"
	"os/signal"
)

type UseCase interface {
	Save(inc entities.Order) error
}
type NatsStreamingController struct {
	uc UseCase
}

func New(uc UseCase) *NatsStreamingController {
	return &NatsStreamingController{uc: uc}
}
func (nsc *NatsStreamingController) InitNatsStreamingController(clusterID, clientID, URL string) {
	var (
		qgroup  string
		durable string
	)

	opts := []nats.Option{nats.Name("NATS Streaming Subscriber")}
	nc, err := nats.Connect(URL, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	sc, err := stan.Connect(clusterID, clientID, stan.NatsConn(nc),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)
		}))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, URL)
	}
	log.Printf("Connected to %s clusterID: [%s] clientID: [%s]\n", URL, clusterID, clientID)

	startOpt := stan.StartWithLastReceived()

	subj := "foo"
	mcb := func(msg *stan.Msg) {
		if err := nsc.HandlerNats(msg); err != nil {
			logrus.Errorf("error while reading from nats: %s", err.Error())
		}
	}
	fmt.Print(subj)
	sub, err := sc.QueueSubscribe(subj, qgroup, mcb, startOpt, stan.DurableName(durable))
	if err != nil {
		sc.Close()
		log.Fatal(err)
	}

	log.Printf("Listening on [%s], clientID=[%s], qgroup=[%s] durable=[%s]\n", subj, clientID, qgroup, durable)

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			fmt.Printf("\nReceived an interrupt, unsubscribing and closing connection...\n\n")

			sub.Unsubscribe()

			sc.Close()
			cleanupDone <- true
		}
	}()
	<-cleanupDone
	log.Fatal()
}
