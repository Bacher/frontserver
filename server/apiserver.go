package main

import (
	"errors"
	"gorpc/rpc"
)

var NoApiServers = errors.New("has no api servers")

var routes = make(map[uint64]*apiServer)
var apiServers = make(map[*apiServer]bool)

type apiServer struct {
	con *rpc.Connection
	rpm int
}

func getApiServer() (*apiServer, error) {
	var server *apiServer
	min := 9999999

	if len(apiServers) == 0 {
		return nil, NoApiServers
	}

	for s := range apiServers {
		if s.rpm < min {
			server = s
			min = s.rpm
		}
	}

	return server, nil
}

func addApiServer(con *rpc.Connection) {
	apiServer := &apiServer{con, 0}
	apiServers[apiServer] = true
}

func removeApiServer(con *rpc.Connection) {
	for s := range apiServers {
		if s.con == con {
			delete(apiServers, s)
			break
		}
	}
}
