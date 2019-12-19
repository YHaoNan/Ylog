package utils

import (
	"os"
	"path/filepath"
)

//Get the absolute path where user executed ylog
func GetCurrentPath() (string,error){
	return os.Getwd()

}

//Get the absolute path where ylog installed
func GetInstallPath() (string,error){
	return filepath.Abs(filepath.Dir(os.Args[0]))
}