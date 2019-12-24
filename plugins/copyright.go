package plugins

import (
	"Ylog/scheme"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
)

type Copyright struct {

}

func (c *Copyright) Name()string{
	return "copyright"
}
func (c *Copyright) IsCommandLinePlugin() bool{
	return false
}
func (c *Copyright) Exec(params map[string]interface{}){
	if params["isPost"].(bool) {
		doc := params["document"].(*goquery.Document)
		template , ok:= params["template"]
		if !ok {
			log.Panicf("Plugin [ copyright ] needs a argument named [ template ]\n")
		}
		templateStr := template.(string)


		var support string
		if params["showYlogSupport"].(bool) {
			support="Powerd by Ylog, Theme by "+scheme.GetYlogProjectConfig().Theme
		}else{
			support = ""
		}
		doc.Find("body").AppendHtml(fmt.Sprintf(`<div id="copyright">%s<br>%s</div>`,templateStr,support))
	}
}

func init(){
	plugin := &Copyright{}
	Registe(plugin)
}