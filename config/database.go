package config

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

var (
	err error
	DB  *sqlx.DB
)

func InitDB(dbc *DatabaseCfg) *sqlx.DB {
	conn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbc.Host, dbc.User, dbc.Pass, dbc.Name, dbc.Port)

	DB, err = sqlx.Open("postgres", conn)
	if err != nil {
		log.Panicf("Database connection failed. Error : %+v", err)
	}

	return DB
}

func GetDBInstance() *sqlx.DB {
	return DB
}
