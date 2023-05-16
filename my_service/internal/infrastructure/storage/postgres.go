package storage

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type ConfigDB struct {
	Host     string
	Port     string
	Username string
	DBName   string
	Password string
	SSLMode  string
}

const schema = `
CREATE TABLE IF NOT EXISTS orders(
	order_uid          varchar not null unique,    
	track_number       varchar not null,      
	entry             varchar not null,      
	delivery_id        varchar not null,  
	payment_id         varchar not null,     
	locate            varchar not null,     
	internal_signature varchar not null,      
	customer_id        varchar not null,      
	delivery_service   varchar not null,      
	shard_key          varchar not null,      
	sm_id              int not null,      
	date_created       timestamp not null,
	oof_shard          varchar not null   
                                 ) ;

CREATE TABLE IF NOT EXISTS delivery(
    delivery_id SERIAL PRIMARY KEY,
	name    varchar not null,  
	phone   varchar not null,  
	zip     varchar not null, 
	city    varchar not null,  
	address varchar not null,  
	region  varchar not null,  
	email   varchar not null 
                                   );

CREATE TABLE IF NOT EXISTS payment(
    payment_id  SERIAL PRIMARY KEY,
	transaction  varchar not null,   
	request_id    varchar not null,   
	currency     varchar not null,   
	provider     varchar not null,   
	amount       int not null,      
	payment_dt    int not null,      
	bank         varchar not null,   
	delivery_cost int not null,      
	goods_total   int not null,      
	custom_fee    int not null    
                                  );

CREATE TABLE IF NOT EXISTS orders_to_items(
    order_id varchar not null,
    item_id int not null
   );

CREATE TABLE IF NOT EXISTS items(
    item_id SERIAL PRIMARY KEY,
	chrt_id       int not null,    
	track_number varchar not null,  
	price        int not null,   
	rid         varchar not null,  
	name        varchar not null,  
	sale         int not null,    
	size        varchar not null,  
	total_price   int not null,    
	nm_id         int not null,    
	brand       varchar not null,  
	status       int not null    
   );
`

func OpenDBConnection(cfg ConfigDB) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}
	logrus.Info("postgres connection opened successfully")
	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}
	return db, nil
}
