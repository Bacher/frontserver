package main

import (
	"frontserver/dbpool"
)

const HttpPort = ":7000"

var tokensMap = make(map[string]uint64)

func main() {
	dbpool.InitDb()
	initApiServers()

	startRpc()
	startHttp(HttpPort)
}
