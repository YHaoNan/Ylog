package utils

import "os"

func Stat(path string) error{
	_,err := os.Stat(path)
	return err
}

func RemoveIfExist(path string){
	if os.IsExist(Stat(path)){
		os.RemoveAll(path)
	}
}

func MkdirIfNotExist(path string){
	if os.IsNotExist(Stat(path)){
		os.Mkdir(path,0777)
	}
}
func CreateIfNotExist(path string){
	if os.IsNotExist(Stat(path)){
		os.Create(path)
	}
}