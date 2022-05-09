package master

import (
	"net"
	"net/http"
	"strconv"
	"time"
)

var (
	//single obj
	G_apiServer *ApiServer
)

//task Http interface
type ApiServer struct {
	httpServer *http.Server
}

//save task interface
func handleJobSave(resp http.ResponseWriter, req *http.Request) {
	var (
		err error
	)

	// parse the post form
	if err = req.ParseForm(); err != nil {
		goto ERR
	}

ERR:
}

//init server
func InitServer() (err error) {
	var (
		mux        *http.ServeMux
		listener   net.Listener
		httpServer *http.Server
	)

	//configure route
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)

	//start tcp listen
	if listener, err = net.Listen("tcp", ":"+strconv.Itoa(G_conf.ApiPort)); err != nil {
		return
	}

	//create http server
	httpServer = &http.Server{
		ReadTimeout:  time.Duration(G_conf.ApiReadTimeOut) * time.Millisecond,
		WriteTimeout: time.Duration(G_conf.ApiWriteTimeOut) * time.Millisecond,
		Handler:      mux,
	}

	//
	G_apiServer = &ApiServer{
		httpServer: httpServer,
	}

	//start service
	go httpServer.Serve(listener)
	return nil
}
