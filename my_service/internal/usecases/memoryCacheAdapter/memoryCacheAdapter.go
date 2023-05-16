package memoryCacheAdapter

import (
	"fmt"
	"learning_NATS_streaming/internal/entities"
	"time"
)

type MemoryCache interface {
	Set(key string, value interface{}, duration time.Duration)
	Get(key string) (interface{}, bool)
}

type MemoryCacheAdapter struct {
	mc MemoryCache
}

func New(mc MemoryCache) *MemoryCacheAdapter {
	return &MemoryCacheAdapter{
		mc: mc,
	}

}

func (mca *MemoryCacheAdapter) Get(uid string) (entities.Order, bool, error) {
	order, ok := mca.mc.Get(fmt.Sprint(uid))
	if !ok {
		return entities.Order{}, ok, nil
	}
	orderConverted := order.(entities.Order)
	return orderConverted, ok, nil
}
func (mca *MemoryCacheAdapter) Set(order entities.Order) error {
	mca.mc.Set(order.OrderUid, order, 0)
	return nil
}
