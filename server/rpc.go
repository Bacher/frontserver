package main

import (
	"gorpc/rpc"
	"log"
	"os"
)

func startRpc() {
	addr := os.Getenv("RPC_ADDR")

	if addr == "" {
		addr = "localhost:9999"
	}

	server := rpc.NewServer(addr, func(con *rpc.Connection, apiName string, body []byte) ([]byte, error) {
		if apiName == "disconnect" {
			onApiServerClosing(con)
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
