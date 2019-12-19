package main
import (
	"Ylog/interpreter"
	_ "Ylog/plugins"
	"Ylog/scheme"
	"fmt"
)
func main() {
	config := scheme.GetYlogProjectConfig()
	fmt.Printf("%+v",config)
	fmt.Println(config.Label == "")

	interpreter.Run()
}
