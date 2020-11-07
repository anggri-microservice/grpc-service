package db

import (
	dbpostgres "github.com/anggri-microservice/golang-service/internal/db/postgres"
	"github.com/jmoiron/sqlx"
)

// DBConn is PostgreSQL database connection instance shared for all packages & codes
var DBConn *sqlx.DB

//PostgreSQL call db connection
var PostgreSQL = dbpostgres.Get
