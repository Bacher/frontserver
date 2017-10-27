package main

import (
	"errors"
	"fmt"
	"frontserver/proto"
	"github.com/golang/protobuf/proto"
	"gorpc/rpc"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"
)

const JsonContentType = "application/json"

var TimeoutError = errors.New("timeout")
var getTokenRx = regexp.MustCompile(`^\[\s*"([A-Za-z0-9]{1,16})"`)

func startHttp(address string) {
	srv := &http.Server{Addr: address}

	http.HandleFunc("/", handler)

	fmt.Printf("Server is started at: %s\n", address)

	err := srv.ListenAndServe()

	if err != nil {
		log.Fatalln(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fail(w, 404, "Not found")
		return
	}

	err, body := readBodyWithTimeout(r, 5*time.Second)

	if err == TimeoutError {
		fail(w, 408, "Request Timeout")
		return
	} else if err != nil {
		fail(w, 400, "Bad Request")
		return
	}

	match := getTokenRx.FindSubmatch(body[:32])

	if len(match) == 0 {
		fail(w, 401, "Unauthorized")
		return
	}

	token := string(match[1])

	var userId uint64
	var found bool
	if userId, found = tokensMap[token]; !found {
		userId, err = getUserId(token)

		if err != nil {
			fail(w, 401, "Unauthorized")
			return
		}
	}

	var apiServer *apiServer
	if apiServer, found = routes[userId]; !found {
		apiServer, err = getApiServer()

		if err != nil {
			log.Println(err)
			fail(w, 500, "Internal Server Error")
			return
		}

		routes[userId] = apiServer
	}

	if apiServer == nil {
		log.Println("Api server not accessable")
		fail(w, 500, "Internal Server Error")
		return
	}

	apiServer.markRequest()

	apiCallStruct := &pb.ApiCall{userId, "getInitialData", body}

	bytes, _ := proto.Marshal(apiCallStruct)

	res, err := apiServer.con.Request("apiCall", bytes)

	if err != nil {
		log.Println("Api failed:", err)

		if err == rpc.ResponseTimeout {
			fail(w, 504, "Gateway Timeout")
		} else {
			fail(w, 500, "Internal Server Error")
		}
		return
	}

	w.Header().Add("Content-Type", JsonContentType)
	w.WriteHeader(200)
	w.Write(res)
}

func readBodyWithTimeout(r *http.Request, timeout time.Duration) (error, []byte) {
	receiveBodyTimeout := time.NewTimer(timeout)

	bodyChan := make(chan []byte)
	errChan := make(chan error)

	go readBody(r.Body, bodyChan, errChan)

	select {
	case err := <-errChan:
		receiveBodyTimeout.Stop()
		return err, nil
	case body := <-bodyChan:
		receiveBodyTimeout.Stop()
		return nil, body
	case <-receiveBodyTimeout.C:
		return TimeoutError, nil
	}
}

func readBody(stream io.ReadCloser, resCh chan []byte, errCh chan error) {
	body, err := ioutil.ReadAll(stream)

	if err != nil {
		errCh <- err
	} else {
		resCh <- body
	}
}

func fail(w http.ResponseWriter, code int, message string) {
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(code)
	w.Write([]byte(message))
}
