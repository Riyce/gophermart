package db

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func GetConnection(dns string) (*sql.DB, error) {
	conn, err := sql.Open("pgx", dns)
	if err != nil {
		return nil, err
	}

	pingErr := conn.Ping()
	if pingErr != nil {
		return nil, pingErr
	}

	return conn, nil
}

func InitDB(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
    	id serial primary key ,
    	login varchar(32) not null unique,
    	password_hash varchar(128) not null,
    	api_key varchar(32) not null unique,
    	withdrawn double precision default 0,
    	current double precision default 0
	);	

	CREATE TABLE IF NOT EXISTS orders (
    	id serial primary key ,
    	order_number varchar(32) unique ,
    	user_id integer not null references users(id) on delete cascade ,
    	status varchar(16) not null ,
    	accrual double precision,
    	created_at timestamp not null default NOW()
	);

	CREATE TABLE IF NOT EXISTS withdraws (
    	id serial primary key ,
    	order_number varchar(32) unique ,
    	user_id integer not null references users(id) on delete cascade ,
    	sum double precision,
    	processed_at timestamp not null default NOW()
	);`
	_, err := db.Exec(query)
	return err
}
