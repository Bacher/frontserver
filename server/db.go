package main

import (
	"errors"
	"frontserver/dbpool"
)

var UserNotFound = errors.New("user not found")

func getUserId(token string) (uint64, error) {
	con, err := dbpool.GetConnection()

	if err != nil {
		return 0, err
	}

	defer dbpool.ReleaseConnection(con)

	rows, err := con.Query("SELECT id FROM users WHERE token = $1 LIMIT 1", token)

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
