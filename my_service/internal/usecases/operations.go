package usecases

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"learning_NATS_streaming/internal/entities"
)

func (uc *UseCase) Save(inc entities.Order) error {
	err := uc.st.Save(inc)
	if err != nil {
		logrus.Errorf("error while saving order: %s in db. error: %s", fmt.Sprint(inc), err.Error())
		return err
	}
	logrus.Infof("saving order: %s in db", fmt.Sprint(inc))

	err = uc.mca.Set(inc)
	if err != nil {
		logrus.Errorf("error while saving order: %s in cache. error: %s", fmt.Sprint(inc), err.Error())
		return err
	}
	logrus.Infof("saving order: %s in cache", fmt.Sprint(inc))
	return nil
}
func (uc *UseCase) Get(id string) (entities.Order, error) {
	order, ok, err := uc.mca.Get(id)
	if err != nil {
		logrus.Errorf("error while getting order: %s from cache. error: %s", fmt.Sprint(id), err.Error())
		return entities.Order{}, err
	}
	if ok {
		logrus.Infof("getting order: %s from cache", fmt.Sprint(id))
		return order, nil
	}
	logrus.Infof("order: %s not cached", fmt.Sprint(id))
	order, err = uc.st.Get(id)
	if err != nil {
		logrus.Errorf("error while getting order: %s from db. error: %s", fmt.Sprint(id), err.Error())
		return entities.Order{}, err
	}
	logrus.Infof("getting order: %s from db", fmt.Sprint(id))
	return order, nil
}

func (uc *UseCase) FillCash() error {
	orders, err := uc.st.GetOrdersForCache()
	if err != nil {
		logrus.Info("error while filling cache: %s", err.Error())
		return err
	}
	logrus.Info("filling cache")
	for _, j := range orders {
		err = uc.mca.Set(j)
		if err != nil {
			logrus.Info("error while filling cache: %s", err.Error())
			return err
		}
	}
	logrus.Info("cache successfully filled")
	return nil
}
