package main

import (
	"errors"
	"gorpc/rpc"
	"log"
	"sync/atomic"
	"time"
)

var NoApiServers = errors.New("has no api servers")

var routes = make(map[uint64]*apiServer)
var apiServers = make(map[*apiServer]bool)

type apiServer struct {
	con  *rpc.Connection
	rpm  uint32
	rpm1 uint32
}

func initApiServers() {
	go func() {
		for {
			var rps []uint32 = nil
			for api := range apiServers {
				rps = append(rps, atomic.LoadUint32(&api.rpm))
				atomic.SwapUint32(&api.rpm, api.rpm1)
				atomic.SwapUint32(&api.rpm1, 0)
			}

			log.Println("RPS:", rps)

			time.Sleep(30 * time.Second)
		}
	}()
}

func getApiServer() (*apiServer, error) {
	var server *apiServer
	var min uint32 = 9999999

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
	apiServer := &apiServer{con, 0, 0}
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

func (s *apiServer) markRequest() {
	atomic.AddUint32(&s.rpm, 1)
	atomic.AddUint32(&s.rpm1, 1)
}
