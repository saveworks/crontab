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
	ApiPort         int      `json:"apiPort"`
	ApiReadTimeOut  int      `json:"apiReadTimeOut"`
	ApiWriteTimeOut int      `json:"apiWriteTimeOut"`
	EtcdEndpoints   []string `json:"etcdEndpoints"`
	EtcdDialTimeout int      `json:"etcdDialTimeout"`
	WebRoot         string   `json:"webRoot"`
}

func InitConfig(fileName string) (err error) {
	// read file
	var (
		content []byte
		conf    Config
	)

	if content, err = ioutil.ReadFile(fileName); err != nil {
		return
	}
	if err = json.Unmarshal(content, &conf); err != nil {
		return
	}
	G_conf = &conf

	return nil
}
