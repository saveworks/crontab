package main

import (
	"crontab/master"
	"flag"
	"fmt"
	"runtime"
)

var (
	confFile string //path
)

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

// command
func initArgs() {
	flag.StringVar(&confFile, "config", "./master.json", "指定master.json")
	flag.Parse()
}

func main() {
	var (
		err error
	)

	//init args
	initArgs()

	//init thread
	initEnv()

	//load configure
	if err = master.InitConfig(confFile); err != nil {
		goto ERR

	}

	// 	start API http service
	if err = master.InitServer(); err != nil {
		goto ERR
	}

	return
ERR:
	fmt.Println(err)
}
