package utils

import (
	"errors"
	"os"
)

func CmdGetSafe(index int) (string,error){
	if CmdLen() < index + 2 {
		return "",errors.New("")
	}
	return os.Args[index + 1],nil
}


func CmdLen() int{
	return len(os.Args)
}