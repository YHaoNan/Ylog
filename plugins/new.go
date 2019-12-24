package plugins

import (
	"Ylog/utils"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type NewHandler struct {

}

func (c *NewHandler) Name()string{
	return "new"
}
func (c *NewHandler) IsCommandLinePlugin() bool{
	return true
}
func (c *NewHandler) Exec(params map[string]interface{}){
	cmdParams := params["cmdParams"].([]string)
	if len(cmdParams) != 2 {
		log.Println("Usage : Ylog new <type> <name>")
	}else{
		_type , name := cmdParams[0], cmdParams[1]
		dir , err:= utils.GetCurrentPath()
		if  err != nil{
			log.Panicf("There is an error %v \n",err)
		}
		if _type == "post"{
			dir += "/posts"
		}else if _type == "page" {
			dir += "/pages"
		}else{
			log.Panicf("Unknown type [ %s ] . try `Ylog new post %s` or `Ylog new page %s`",_type,name,name)
		}

		splLast := strings.LastIndex(name,"/")
		categories := ""
		title := ""
		if splLast != -1{
			dir += "/" + name[:splLast]
			categories = "categories: "+name[:splLast]
			utils.MkdirIfNotExist(dir)
			dir += "/"+name[splLast + 1:] + ".md"
			title = name[splLast + 1:]
		}else{
			dir += "/"+name+".md"
			title = name
		}

		f, err := os.Create(dir)
		defer f.Close()
		if err!=nil {
			log.Printf("Can't craete file , because : %v\n",err)
		}else{
			content := `---
title: %s
date: %s
%s
---`
			timeStr := time.Now().Format("2006-01-02 15:04:05")
			content = fmt.Sprintf(content,title,timeStr,categories)
			_,err = f.Write([]byte(content))

		}
	}
}

func init(){
	handler := &NewHandler{}
	Registe(handler)
}