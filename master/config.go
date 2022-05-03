package master

import (
	"encoding/json"
	"io/ioutil"
)

var (
	G_conf *Config
)

//configure
type Config struct {
	ApiPort         int `json:"ApiPorta"`
	ApiReadTimeOut  int `json:"ApiReadTimeOut"`
	ApiWriteTimeOut int `json:"ApiWriteTimeOut"`
}

func InitConfig(fileName string) (err error) {
	// read file
	var (
		content []byte
		conf    Config
	)
	if content, err = ioutil.ReadFile(fileName); err != nil {

	}
	if err = json.Unmarshal(content, &conf); err != nil {

	}
	G_conf = &conf

	return nil
}
