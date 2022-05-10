package master

import (
	"crontab/master/common"
	"encoding/json"
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
		err     error
		postJob string
		job     common.Job
		oldJob  *common.Job
		bytes   []byte
	)

	// parse the post form
	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	// job
	postJob = req.PostForm.Get("job")

	//mash
	if err = json.Unmarshal([]byte(postJob), &job); err != nil {
		goto ERR
	}

	//save
	if oldJob, err = G_jobMgr.SaveJob(&job); err != nil {
		goto ERR
	}

	//resp
	if bytes, err = common.BuildResponse(0, "succc", oldJob); err == nil {
		resp.Write(bytes)
	}
	return
ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err != nil {
		resp.Write(bytes)
	}
}
func handleJobDelete(resp http.ResponseWriter, req *http.Request) {
	var (
		err    error
		name   string
		oldjob *common.Job
		bytes  []byte
	)

	// parse the post form
	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	//delete
	name = req.PostForm.Get("name")

	//delete job
	if oldjob, err = G_jobMgr.DeleteJob(name); err != nil {
		goto ERR
	}

	// succ resp
	if bytes, err = common.BuildResponse(0, "succ", oldjob); err == nil {
		resp.Write(bytes)
	}
	return
ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

func handleJobList(resp http.ResponseWriter, req *http.Request) {
	var (
		jobList []*common.Job
		err     error
		bytes   []byte
	)
	if jobList, err = G_jobMgr.ListJobs(); err != nil {
		goto ERR
	}
	//resp
	if bytes, err = common.BuildResponse(0, "succ", jobList); err == nil {
		resp.Write(bytes)
	}
	return
ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

func handleJobKill(resp http.ResponseWriter, req *http.Request) {

	var (
		err   error
		name  string
		bytes []byte
	)

	// parse the post form
	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	name = req.PostForm.Get("name")

	if err = G_jobMgr.KillJobs(name); err != nil {
		goto ERR
	}

	if bytes, err = common.BuildResponse(0, "succ", nil); err == nil {
		resp.Write(bytes)
	}
	return
ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

//init server
func InitServer() (err error) {
	var (
		mux        *http.ServeMux
		listener   net.Listener
		httpServer *http.Server
		staticDir  http.Dir

		staticHandler http.Handler
	)

	//configure route
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)
	mux.HandleFunc("/job/delete", handleJobDelete)
	mux.HandleFunc("/job/list", handleJobList)
	mux.HandleFunc("/job/kill", handleJobKill)

	staticDir = http.Dir(G_conf.WebRoot)
	staticHandler = http.FileServer(staticDir)
	mux.Handle("/", http.StripPrefix("/", staticHandler))

	//start tcp listen
	if listener, err = net.Listen("tcp", ":"+strconv.Itoa(G_conf.ApiPort)); err != nil {
		return
	}

	//fmt.Printf("端口:%s", strconv.Itoa(G_conf.ApiPort))

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
