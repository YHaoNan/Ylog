package main

import (
	"Ylog/scheme"
	"fmt"
)

func main(){
	config,error := scheme.GetConfigFromFile("../sample/config.yaml")
	fmt.Println(config,error)

}