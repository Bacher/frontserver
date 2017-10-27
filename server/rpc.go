package main

import (
	"gorpc/rpc"
	"log"
)

func startRpc() {
	server := rpc.NewServer(func(con *rpc.Connection, apiName string, body []byte) ([]byte, error) {
		if apiName == "disconnect" {
			removeApiServer(con)
			return nil, nil
		}

		log.Printf("Unknown apiMethod for parse %s\n", apiName)
		return nil, rpc.ApiNotFound
	})

	server.SetHandlers(addApiServer, removeApiServer)

	err := server.Listen()

	if err != nil {
		log.Fatalln(err)
	}

	go server.Serve()
}
