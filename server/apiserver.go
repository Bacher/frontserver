package main

import (
	"errors"
	"frontserver/future"
	"gorpc/rpc"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

var NoApiServers = errors.New("has no api servers")

var routes = make(map[uint64]*apiServer)
var apiServers = make(map[*rpc.Connection]*apiServer)

type apiServer struct {
	con          *rpc.Connection
	closing      bool
	currentUsers map[uint64]uint16
	rpm          uint32
	rpm1         uint32
	userMut      *sync.RWMutex
	callAfter    map[uint64]*future.Future
}

func initApiServers() {
	go func() {
		second := true

		for {
			var rps []uint32 = nil
			for _, api := range apiServers {
				rps = append(rps, atomic.LoadUint32(&api.rpm))
				atomic.SwapUint32(&api.rpm, api.rpm1)
				atomic.SwapUint32(&api.rpm1, 0)
			}

			second = !second

			//if second {
			//	log.Println("RPS:", rps)
			//}

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

	for _, s := range apiServers {
		if !s.closing && s.rpm < min {
			server = s
			min = s.rpm
		}
	}

	return server, nil
}

func addApiServer(con *rpc.Connection) {
	apiServer := &apiServer{
		con,
		false,
		make(map[uint64]uint16),
		0,
		0,
		&sync.RWMutex{},
		make(map[uint64]*future.Future),
	}

	apiServers[con] = apiServer
	log.Printf("Connected    | count: %d\n", len(apiServers))
}

func removeApiServer(con *rpc.Connection) {
	_, ok := apiServers[con]

	if ok {
		delete(apiServers, con)
		log.Printf("Disconnected | count: %d\n", len(apiServers))
	}
}

func onApiServerClosing(con *rpc.Connection) {
	apiServer, ok := apiServers[con]

	if ok {
		apiServer.closing = true
	}
}

func (s *apiServer) markRequest() {
	atomic.AddUint32(&s.rpm, 1)
	atomic.AddUint32(&s.rpm1, 1)
}

func (s *apiServer) request(userId uint64, apiName string, body []byte) ([]byte, error) {
	s.userMut.Lock()
	s.currentUsers[userId]++
	s.userMut.Unlock()

	res, err := s.con.Request(apiName, body)

	s.userMut.Lock()
	if s.currentUsers[userId] <= 1 {
		delete(s.currentUsers, userId)

		if s.closing {
			fut, ok := s.callAfter[userId]

			if ok {
				fut.Done()
			}

			delete(s.callAfter, userId)
		}

	} else {
		s.currentUsers[userId]--
	}

	s.userMut.Unlock()

	return res, err
}

func (s *apiServer) getCurrentCount(userId uint64) uint16 {
	s.userMut.RLock()
	count := s.currentUsers[userId]
	s.userMut.RUnlock()
	return count
}

func (s *apiServer) getUserDoneFuture(userId uint64) *future.Future {
	fut, ok := s.callAfter[userId]

	if ok {
		return fut
	} else {
		fut = future.New()
		s.callAfter[userId] = fut
		return fut
	}
}
