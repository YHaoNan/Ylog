package plugins

import (
	"log"
)
var pluginMap map[string]Plugin

func InitPluginLib(){
	log.Println("Init PluginLib...")
	pluginMap = make(map[string]Plugin)
	log.Println("PluginLib init over...")
}
func Registe(plugin Plugin) {
	pluginMap[plugin.Name()] = plugin
	log.Printf("Add a plugin to PluginLib : [ %s ] \n" , plugin.Name())
}
func GetPluginFromLib(name string) (Plugin,bool) {
	plugin, ok := pluginMap[name]
	return plugin,ok
}