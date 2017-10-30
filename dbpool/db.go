package dbpool

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"sync"
)

type connectionPool struct {
	free  []*sql.DB
	mutex *sync.Mutex
}

var pool = &connectionPool{
	nil,
	&sync.Mutex{},
}

func InitDb() {
	con, err := GetConnection()

	if err != nil {
		log.Fatalln(err)
	}

	defer ReleaseConnection(con)

	rows, err := con.Query("SELECT id, token FROM users LIMIT 1")

	if err != nil {
		log.Fatalln(err)
	}

	for rows.Next() {
	}
}

func GetConnection() (*sql.DB, error) {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	if len(pool.free) > 0 {
		con := pool.free[0]
		pool.free = pool.free[1:]
		return con, nil

	} else {
		dbHost := os.Getenv("DB_HOST")

		if dbHost == "" {
			dbHost = "localhost"
		}

		connStr := fmt.Sprintf("host=%s port=5432 user=postgres dbname=postgres sslmode=disable", dbHost)
		con, err := sql.Open("postgres", connStr)

		if err != nil {
			return nil, err
		}

		return con, nil
	}
}

func ReleaseConnection(con *sql.DB) {
	pool.mutex.Lock()
	pool.free = append(pool.free, con)
	pool.mutex.Unlock()
}
