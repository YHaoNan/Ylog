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
	doc.Find("head").AppendHtml(fmt.Sprintf("<title>%s</title>",head["title"].(string)))
}
