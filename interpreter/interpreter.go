package interpreter

import (
	"Ylog/plugins"
	"Ylog/scheme"
	"Ylog/utils"
	"log"
	"os"
)


func Run(){
	log.Println("Init commandline interpreter...")
	pluginToExec,err := utils.CmdGetSafe(0)
	if err != nil{
		pluginToExec = "help"
	}

	plugin,ok := plugins.GetPluginFromLib(pluginToExec)
	if ok && plugin.IsCommandLinePlugin() {
		pluginConfig , ok := scheme.GetYlogProjectConfig().CmdPluginsMap[pluginToExec]
		if !ok{
			pluginConfig = make(map[string]interface{})
		}
		if utils.CmdLen() > 2 {
			pluginConfig["cmdParams"] = os.Args[2:]
		}else{
			pluginConfig["cmdParams"] = []string{}
		}
		plugin.Exec(pluginConfig)
	}else{
		log.Panicf("There is no plugin to handle command : [ %s ] \n",pluginToExec)
	}
}
