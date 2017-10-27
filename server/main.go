package main

const HttpPort = ":7000"

var tokensMap = make(map[string]uint64)

func main() {
	initDb()

	startRpc()
	startHttp(HttpPort)
}
