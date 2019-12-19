package utils

import (
	"bufio"
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"gopkg.in/russross/blackfriday.v2"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"strings"
)



func getHead(md []byte,path string) (map[string]interface{},int){
	inHead := false
	buff := bufio.NewReader(bytes.NewReader(md))
	yamlBuf := make([]byte,0)

	n := 0
	for l := 0 ; true ; l++{
		line , err := buff.ReadString('\n')
		line = strings.Trim(line,"\n")
		if err == io.EOF{
			break
		}else if err != nil{
			log.Panicf("There is an error : %v \n",err)
		}else{
			n += len(line) + 1
			if l==0 && len(line) >=3 && line[0:3] == "---" {
				inHead = true
				continue
			}
			if l>0 && inHead && len(line) >= 3 && line[0:3] == "---" {
				inHead = false
				break
			}
			if inHead {
				yamlBuf = append(yamlBuf,[]byte(line+"\n")...)
			}
		}
	}
	if inHead || len(yamlBuf) == 0{
		log.Panicf("Markdown head is not vaild, file : %s",path)
	}

	items := yaml.MapSlice{}
	yaml.Unmarshal(yamlBuf,&items)
	retMap := make(map[string]interface{})
	for _,item := range items{
		retMap[item.Key.(string)] = item.Value
	}
	return retMap,n
}

func ConvertMarkdownToHTML(md []byte,path string) (map[string]interface{},*goquery.Document){
	head ,i := getHead(md,path)
	description := blackfriday.Run(md[i:])
	doc, error := goquery.NewDocumentFromReader(bytes.NewReader(description))
	if error != nil{
		log.Panicf("There is an error : %v\n",error)
	}
	return head,doc
}

func ConvertMarkdownToHTMLFromFile(path string) (map[string]interface{},*goquery.Document){
	datas, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panicf("There is an error : %v \n",err)
	}
	return ConvertMarkdownToHTML(datas,path)
}