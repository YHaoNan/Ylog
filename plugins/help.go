package plugins

import (
	"fmt"
	"log"
)

type Help struct {

}

func (c *Help) Name()string{
	return "help"
}
func (c *Help) IsCommandLinePlugin() bool{
	return true
}
func (c *Help) Exec(params map[string]interface{}){
	fmt.Println(`Usage:
	help -- get help
	build [path] -- build static blog
	serve [port] -- serve static blog in local`)
}

func init(){
	log.Println("Init help ...")
	Registe(&Help{})
}