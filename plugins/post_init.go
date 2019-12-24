package plugins

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
)

type PostInit struct {

}

func (c *PostInit) Name()string{
	return "postinit"
}
func (c *PostInit) IsCommandLinePlugin() bool{
	return false
}
func (c *PostInit) Exec(params map[string]interface{}){
	if params["isPost"].(bool) {
		head := params["head"].(map[string]interface{})
		title ,ok := head["title"]
		titleHTML := ""
		if ok {
			titleHTML = `<div id="title">
	<span>%s</span>
<div>`
			titleHTML = fmt.Sprintf(titleHTML,title)
		}

		categorieHTML := getCategorieHTML(head)
		doc := params["document"].(*goquery.Document)

		doc.Find("body").PrependHtml(categorieHTML)
		doc.Find("body").PrependHtml(titleHTML)
	}
}

func genTag(tag,id string,text interface{}) string{
	return fmt.Sprintf(`<%s id="%s">%v</%s>`,tag,id,text,tag)
}

func getCategorieHTML(head map[string]interface{}) string{
	categorie ,ok:= head["categories"]
	if  ok {
		categorieHTMLTemplate := `<div id="categorie-bar">%s%s</div>`
		categorieTag := genTag("span","categorie",categorie)
		time ,ok:= head["date"]
		if ok {
			return fmt.Sprintf(categorieHTMLTemplate,genTag("span","time",time.(string)),categorieTag)
		}else{
			return fmt.Sprintf(categorieHTMLTemplate,categorieTag,"")
		}
	}
	return ""

}
