package app

import (
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"learning_NATS_streaming/internal/controller/httpServerController"
	"learning_NATS_streaming/internal/controller/natsStreamingController"
	"learning_NATS_streaming/internal/infrastructure/storage"
	"learning_NATS_streaming/internal/usecases"
	"learning_NATS_streaming/internal/usecases/memoryCacheAdapter"
	"learning_NATS_streaming/pkg/memorycache"
)

func RunApp(postgresConfig storage.ConfigDB) {
	db, err := storage.OpenDBConnection(postgresConfig)
	if err != nil {
		logrus.Fatalf("error opening postgres connection:%s", err.Error())
	}
	store := storage.New(db)
	mc := memorycache.New(0, 0)
	mca := memoryCacheAdapter.New(mc)
	uc := usecases.New(store, mca)
	if err := uc.FillCash(); err != nil {
		logrus.Fatalf("cant fill cash: %s", err.Error())
	}
	hsc := httpServerController.New(uc)
	nsc := natsStreamingController.New(uc)
	router := hsc.InitRouter()
	server := new(httpServerController.Server)
	go nsc.InitNatsStreamingController("test-cluster", "myID", nats.DefaultURL)
	logrus.Info("starting nats streaming controller...")
	if err := server.RunServer(viper.GetString("port"), router); err != nil {
		logrus.Fatalf("error while starting server:%s", err.Error())
	}
	logrus.Info("starting server...")

}
