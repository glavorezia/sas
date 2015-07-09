package webengine

import (
	"gopkg.in/yaml.v2"
	"fmt"
	"io/ioutil"
)

var config Config

type Config struct {
	Http struct {
		Host string `yaml: host`
		Port int `yaml: port`
		Secure bool `yaml: secure`
	}
	Database []struct {
		Name string `yaml: name`
		Driver string `yaml: driver`
		Host string
		User string
		Pass string
		Database string
		Port int
	}
	Cache []struct {
		Name string
		Driver string
		Host string
		Port int
	}
	Session struct {
		Driver string
		Ttl string
	}
	Domains []struct {
		Name string
		Domain string
		Cookie string
		Alias []string
		Root string
	}
}


func init(){
	config = Config{}
	/*data, err := ioutil.ReadFile("conf/webengine.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		fmt.Printf("Error Loading config %v", err)
	}*/
	LoadConfig( "conf/webengine.yaml", &config )
}

func LoadConfig(file string, target interface{}) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal([]byte(data), target)
	if err != nil {
		fmt.Printf("Error Loading config %v: %v", file, err)
	}
}
