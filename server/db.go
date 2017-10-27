package main

import (
	"database/sql"
	"github.com/go-errors/errors"
	_ "github.com/lib/pq"
	"log"
	"sync"
)

type connectionPool struct {
	free  []*sql.DB
	mutex *sync.Mutex
}

var UserNotFound = errors.New("user not found")

var pool = &connectionPool{
	nil,
	&sync.Mutex{},
}

func initDb() {
	con, err := getConnection()

	if err != nil {
		log.Fatalln(err)
	}

	defer releaseConnection(con)

	rows, err := con.Query("SELECT id, token FROM users LIMIT 1")

	if err != nil {
		log.Fatalln(err)
	}

	for rows.Next() {
	}
}

func getConnection() (*sql.DB, error) {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	if len(pool.free) > 0 {
		con := pool.free[0]
		pool.free = pool.free[1:]
		return con, nil

	} else {
		connStr := "user=postgres dbname=postgres sslmode=disable"
		con, err := sql.Open("postgres", connStr)

		if err != nil {
			return nil, err
		}

		return con, nil
	}
}

func releaseConnection(con *sql.DB) {
	pool.mutex.Lock()
	pool.free = append(pool.free, con)
	pool.mutex.Unlock()
}

func getUserId(token string) (uint64, error) {
	db, err := getConnection()

	if err != nil {
		return 0, err
	}

	rows, err := db.Query("SELECT id FROM users WHERE token = $1 LIMIT 1", token)

	if err != nil {
		return 0, err
	}

	if rows.Next() {
		var id uint64
		rows.Scan(&id)

		return id, nil
	} else {
		return 0, UserNotFound
	}
}
