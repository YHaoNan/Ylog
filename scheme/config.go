package scheme

import (
	"Ylog/utils"
	"gopkg.in/yaml.v2"
	"log"
)

var ylogProjectConfig *Config

type Config struct {
	Author string `yaml:"author"`
	Label string `yaml:"label"`
	Theme string `yaml:"theme"`
	Assets string `yaml:"assets"`
	CDN string `yaml:"cdn"`
	PluginsItem yaml.MapSlice `yaml:"plugins"`
	CmdPluginsItem yaml.MapSlice `yaml:"cmdPlugins"`
	PluginsMap map[string]map[string]interface{}
	CmdPluginsMap map[string]map[string]interface{}
}

func GetYlogProjectConfig() *Config{
	if ylogProjectConfig == nil {
		log.Println("Reading Ylog project config file...")
		path , err := utils.GetCurrentPath()
		if err != nil {
			log.Panicf("Can't get current path. Error info : %+v \n",err)
		}
		log.Println("Current Ylog project path is : " + path)
		ylogProjectConfig , err = GetConfigFromFile(path + "/config.yaml")
		if err != nil{
			log.Panicf("Error to read Ylog config file. Error : %v",err)
		}
		if !VaildateConfig(ylogProjectConfig) {
			log.Panicf("Config file missing fields [ label ]\n")
		}
		if err != nil {
			log.Panicf("`config.yaml is not found! %+v",err)
		}
	}
	ylogProjectConfig.PluginsMap = mapsliceToPluginMap(ylogProjectConfig.PluginsItem)
	ylogProjectConfig.CmdPluginsMap = mapsliceToPluginMap(ylogProjectConfig.CmdPluginsItem)

	return ylogProjectConfig
}

func mapsliceToPluginMap(mapslice yaml.MapSlice) map[string]map[string]interface{}{
	plugin := make(map[string]map[string]interface{})
	for _,item := range mapslice{
		pluginParams := make(map[string]interface{})
		plugin[item.Key.(string)] = pluginParams
		for _,subItem := range item.Value.(yaml.MapSlice){
			pluginParams[subItem.Key.(string)] = subItem.Value
		}
	}
	return plugin
}


