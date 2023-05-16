package httpServerController

import "learning_NATS_streaming/internal/entities"

type useCase interface {
	Get(id string) (entities.Order, error)
}

type HttpServerController struct {
	uc useCase
}

func New(uc useCase) *HttpServerController {
	return &HttpServerController{
		uc: uc,
	}
}
