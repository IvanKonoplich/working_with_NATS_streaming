package storage

import (
	"errors"
	"fmt"
	"learning_NATS_streaming/internal/entities"
)

func (s *Storage) Save(inc entities.Order) error {
	query1 := "select exists(SELECT * FROM orders WHERE order_uid=$1)"
	var is bool
	s.db.Get(&is, query1, inc.OrderUid)
	if is {
		return errors.New(fmt.Sprintf("order: %s already exist", inc.OrderUid))
	}
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	query := "INSERT INTO delivery (name, phone, zip, city, address, region, email) values ($1, $2, $3, $4, $5, $6, $7) RETURNING delivery_id"
	var deliveryId int
	err = tx.QueryRow(query, inc.Delivery.Name, inc.Delivery.Phone, inc.Delivery.Zip, inc.Delivery.City, inc.Delivery.Address, inc.Delivery.Region, inc.Delivery.Email).Scan(&deliveryId)
	if err != nil {
		tx.Rollback()
		return err
	}

	var paymentId int
	query = "INSERT INTO payment (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING payment_id"
	err = tx.QueryRow(query, inc.Payment.Transaction, inc.Payment.RequestId, inc.Payment.Currency, inc.Payment.Provider, inc.Payment.Amount, inc.Payment.PaymentDt, inc.Payment.Bank, inc.Payment.DeliveryCost, inc.Payment.GoodsTotal, inc.Payment.CustomFee).Scan(&paymentId)
	if err != nil {
		tx.Rollback()
		return err
	}

	itemsIdList := make([]int, 0, len(inc.Items))
	for _, item := range inc.Items {

		var itemId int
		query = "INSERT INTO items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING item_id"
		err := tx.QueryRow(query, item.ChrtId, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmId, item.Brand, item.Status).Scan(&itemId)
		if err != nil {
			tx.Rollback()
			return err
		}
		itemsIdList = append(itemsIdList, itemId)
	}

	query = "INSERT INTO orders (order_uid, track_number, entry, delivery_id, payment_id, locate, internal_signature, customer_id, delivery_service, shard_key, sm_id, date_created, oof_shard) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)"
	_, err = tx.Exec(query, inc.OrderUid, inc.TrackNumber, inc.Entry, deliveryId, paymentId, inc.Locale, inc.InternalSignature, inc.CustomerId, inc.DeliveryService, inc.Shardkey, inc.SmId, inc.DateCreated, inc.OofShard)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, itemConvertedIdd := range itemsIdList {
		query = "INSERT INTO orders_to_items (order_id, item_id) values ($1, $2)"
		_, err := tx.Exec(query, inc.OrderUid, itemConvertedIdd)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}
func (s *Storage) Get(uid string) (entities.Order, error) {
	var result entities.Order
	query1 := "select exists(SELECT * FROM orders WHERE order_uid=$1)"
	var is bool
	s.db.Get(&is, query1, uid)
	if !is {
		return entities.Order{}, errors.New(fmt.Sprintf("order: %s does not exist", uid))
	}
	query2 := "SELECT order_uid, track_number, entry, locate, internal_signature, customer_id, delivery_service, shard_key, sm_id, date_created, oof_shard FROM orders WHERE order_uid=$1"
	row := s.db.QueryRowx(query2, uid)
	if err := row.StructScan(&result); err != nil {
		return entities.Order{}, err
	}

	delivery, err := s.getDelivery(uid)
	if err != nil {
		return entities.Order{}, err
	}
	result.Delivery = delivery

	payment, err := s.getPayment(uid)
	if err != nil {
		return entities.Order{}, err
	}
	result.Payment = payment

	items, err := s.getItems(uid)
	if err != nil {
		return entities.Order{}, err
	}
	result.Items = items

	return result, nil
}

func (s *Storage) getDelivery(uid string) (entities.Delivery, error) {
	var deliveryId int
	query := "SELECT delivery_id FROM orders WHERE order_uid=$1"
	row := s.db.QueryRow(query, uid)
	if err := row.Scan(&deliveryId); err != nil {
		return entities.Delivery{}, err
	}

	var delivery entities.Delivery
	query = "SELECT name, phone, zip, city, address, region, email FROM delivery WHERE delivery_id=$1"
	rows := s.db.QueryRowx(query, deliveryId)
	if err := rows.StructScan(&delivery); err != nil {
		return entities.Delivery{}, err
	}
	return delivery, nil
}

func (s *Storage) getPayment(uid string) (entities.Payment, error) {
	var paymentId int
	query := "SELECT payment_id FROM orders WHERE order_uid=$1"
	row := s.db.QueryRow(query, uid)
	if err := row.Scan(&paymentId); err != nil {
		return entities.Payment{}, err
	}

	var payment entities.Payment
	query = "SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee FROM payment WHERE payment_id=$1"
	rowx := s.db.QueryRowx(query, paymentId)
	if err := rowx.StructScan(&payment); err != nil {
		return entities.Payment{}, err
	}
	return payment, nil
}

func (s *Storage) getItems(uid string) ([]entities.Item, error) {

	var itemsId []int
	query := "SELECT item_id FROM orders_to_items WHERE order_id=$1"
	rows, err := s.db.Queryx(query, uid)
	for rows.Next() {
		var itemId int
		err = rows.Scan(&itemId)
		if err != nil {
			return nil, err
		}
		itemsId = append(itemsId, itemId)
	}
	var result []entities.Item
	for _, itemID := range itemsId {
		query = "SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM items WHERE item_id=$1"
		rowItem := s.db.QueryRowx(query, itemID)
		var item entities.Item
		err = rowItem.StructScan(&item)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, nil
}

func (s *Storage) GetOrdersForCache() ([]entities.Order, error) {
	var result = make([]entities.Order, 0, 10)
	query := "SELECT order_uid FROM orders LIMIT 10"
	rows, err := s.db.Queryx(query)
	for rows.Next() {
		var uid string
		err = rows.Scan(&uid)
		if err != nil {
			return nil, err
		}
		order, err := s.Get(uid)
		if err != nil {
			return nil, err
		}
		result = append(result, order)
	}
	return result, nil
}
