package plugins

import (
	"Ylog/utils"
	"bytes"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)
type post struct {
	title string
	urlRela string
	date string
}

var posts []*post = make([]*post,0)
type pageGenHandler struct {
	indexPageSrc *bytes.Buffer
}

func NewPageHandler() *pageGenHandler{
	var buf = bytes.NewBuffer(make([]byte,0))
	buf.WriteString(`---
title: Index
date: 2000-01-01 11:11:11
categories: DEF
---
`)
	h := &pageGenHandler{
		indexPageSrc:buf,
	}
	return h
}

func (c *pageGenHandler) Name()string{
	return "pagegen"
}
func (c *pageGenHandler) IsCommandLinePlugin() bool{
	return false
}
func (c *pageGenHandler) Exec(params map[string]interface{}){
	isPost := params["isPost"].(bool)
	if isPost {
		head := params["head"].(map[string]interface{})
		date := head["date"].(string)
		title := head["title"].(string)
		categorie := head["categories"].(string)
		file := params["dstF"].(string)
		file = file[strings.LastIndex(file,"/"):]
		posts = append(posts,&post{
			title:   title,
			urlRela: categorie+file,
			date:    date,
		})
	}

}

func (c *pageGenHandler) PostRenderOver() {
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].date > posts[j].date
	})
	c.indexPageSrc.WriteString(`<ul id="posts">`)
	for _,post := range posts{
		c.indexPageSrc.WriteString(fmt.Sprintf(`<li class="post"><a href="%s"><span class="post-title">%s</span></a><span class="post-date">%s</span></li>`,
			post.urlRela,post.title,post.date))
	}
	c.indexPageSrc.WriteString("</ul>")
	path ,err := utils.GetCurrentPath()
	if err != nil {
		log.Panicf("There is an error %v \n",err)
	}
	path += "/pages/index.md"
	f ,err:= os.OpenFile(path,os.O_CREATE|os.O_WRONLY,0777)
	defer f.Close()
	if err != nil {
		log.Panicf("There is an error %v \n",err)
	}
	f.WriteString(c.indexPageSrc.String())
}
