package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"time"
)

var counts = 0

func openDBPool() (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.New(context.Background(), "postgres://postgres:password@postgres:5432/lowserver?sslmode=disable")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		return nil, err
	}

	err = dbpool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return dbpool, nil
}

func ConnectToDB() *pgxpool.Pool {
	for {
		connection, err := openDBPool()
		if err != nil {
			log.Println("Could not connect to database, Postgres is not ready...")
			counts += 1
		} else {
			log.Println("Connected to database...")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Waiting for database to become ready...")
		time.Sleep(2 * time.Second)
		continue
	}
}
