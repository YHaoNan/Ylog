package plugins

import (
	"log"
)

type Logger struct {
}
func (c *Logger) Name()string{
	return "logger"
}
func (c *Logger) IsCommandLinePlugin() bool{
	return false
}
func (c *Logger) Exec(params map[string]interface{}){
	log.Printf("Building file %s...\n",params["srcF"].(string))
}

func init(){
	log.Println("Init help ...")
	Registe(&Help{})
}
