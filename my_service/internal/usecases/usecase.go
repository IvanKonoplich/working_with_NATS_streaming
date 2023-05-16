package usecases

import (
	"learning_NATS_streaming/internal/entities"
)

type Storage interface {
	Save(inc entities.Order) error
	Get(uid string) (entities.Order, error)
	GetOrdersForCache() ([]entities.Order, error)
}

type MemoryCacheAdapter interface {
	Get(uid string) (entities.Order, bool, error)
	Set(order entities.Order) error
}

type UseCase struct {
	st  Storage
	mca MemoryCacheAdapter
}

func New(st Storage, mca MemoryCacheAdapter) *UseCase {
	return &UseCase{
		st:  st,
		mca: mca,
	}
}
