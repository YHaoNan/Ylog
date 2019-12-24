package plugins

import (
	"Ylog/scheme"
	"github.com/PuerkitoBio/goquery"
)

type ThemeHandler struct {

}

func (c *ThemeHandler) Name()string{
	return "themehandler"
}
func (c *ThemeHandler) IsCommandLinePlugin() bool{
	return false
}
func (c *ThemeHandler) Exec(params map[string]interface{}){
	doc := params["document"].(*goquery.Document)
	//head := params["head"].(map[string]interface{})
	themeDir := scheme.GetYlogProjectConfig().Theme
	if params["isPost"].(bool) {
		doc.Find("head").AppendHtml("<link rel='stylesheet' href='../../themes/"+themeDir+"/public.css'/>")
		doc.Find("head").AppendHtml("<link rel='stylesheet' href='../../themes/"+themeDir+"/post.css'/>")
	}else{
		doc.Find("head").AppendHtml("<link rel='stylesheet' href='../themes/"+themeDir+"/public.css'/>")
		doc.Find("head").AppendHtml("<link rel='stylesheet' href='../themes/"+themeDir+"/post.css'/>")
	}
}
