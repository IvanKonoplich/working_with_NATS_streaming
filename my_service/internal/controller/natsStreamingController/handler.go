package natsStreamingController

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
	"learning_NATS_streaming/internal/entities"
)

func (nsc *NatsStreamingController) HandlerNats(m *stan.Msg) error {
	var order entities.Order
	var rawData []byte
	rawData = m.MsgProto.Data
	err := json.Unmarshal(rawData, &order)
	if err != nil {
		return errors.New(fmt.Sprintf("incorrect message: %s - error: %s", rawData, err.Error()))
	}
	logrus.Infof("getting new order from nats streaming: %s", fmt.Sprint(order))
	return nsc.uc.Save(order)
}
