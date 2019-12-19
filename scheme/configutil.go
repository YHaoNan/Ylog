package scheme

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func GetConfigFromFile(path string) (*Config,error){
	config := &Config{}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data,config)
	if err != nil {
		return nil , err
	}
	return config,nil
}

func VaildateConfig(config *Config) bool{
	if config == nil || config.Label == ""{
		return false
	}
	return true
}