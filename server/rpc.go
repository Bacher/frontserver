package main

import (
	"gorpc/rpc"
	"log"
)

func startRpc() {
	server := rpc.NewServer(func(apiName string, body []byte) ([]byte, error) {
		if apiName == "disconnect" {

		}
		log.Printf("Unknown apiMethod for parse %s\n", apiName)
		return nil, rpc.ApiNotFound
	})

	server.SetHandlers(
		func(c *rpc.Connection) {
			addApiServer(c)
			log.Printf("Connected    | count: %d\n", len(apiServers))
		},
		func(c *rpc.Connection) {
			removeApiServer(c)
			log.Printf("Disconnected | count: %d\n", len(apiServers))
		})

	err := server.Listen()

	if err != nil {
		log.Fatalln(err)
	}

	go server.Serve()
}
