package plugins

import (
	//"github.com/PuerkitoBio/goquery"
	"fmt"
	"github.com/PuerkitoBio/goquery"
)

type HeadHandler struct {

}

func (c *HeadHandler) Name()string{
	return "headhandler"
}
func (c *HeadHandler) IsCommandLinePlugin() bool{
	return false
}
func (c *HeadHandler) Exec(params map[string]interface{}){
	doc := params["document"].(*goquery.Document)
	head := params["head"].(map[string]interface{})
	headEle := doc.Find("head")
	headEle.AppendHtml(fmt.Sprintf("<title>%s</title>",head["title"].(string)))
	headEle.AppendHtml(`<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />`)
	headEle.AppendHtml(`<meta name="viewport" content="width=device-width, initial-scale=1.0">`)
}
