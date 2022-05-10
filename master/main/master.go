package main

import (
	"crontab/master"
	"flag"
	"fmt"
	"runtime"
	"time"
)

var (
	confFile string //path
)

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

// command
func initArgs() {
	flag.StringVar(&confFile, "config", "./master/main/master.json", "指定master.json")
	flag.Parse()
}

func main() {
	var (
		err error
	)

	//init args
	initArgs()

	initEnv()

	//init thread
	//load configure
	if err = master.InitConfig(confFile); err != nil {
		goto ERR

	}

	// 启动管理器
	if err = master.InitJobMgr(); err != nil {
		goto ERR
	}

	// 	start API http service
	if err = master.InitServer(); err != nil {
		goto ERR
	}

	for {
		time.Sleep(1 * time.Second)
	}

	return
ERR:
	fmt.Println(err)
}
