package plugins

import (
	"Ylog/scheme"
	"Ylog/utils"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Build struct {

}
func (c *Build) Name()string{
	return "build"
}
func (c *Build) IsCommandLinePlugin() bool{
	return true
}

func getFileNameWithoutExtName(name string) string{
	return name[:strings.LastIndex(name,".")]
}

func getFileNameAndExtName(name string) (string,string){
	lastI := strings.LastIndex(name,".")
	if lastI == -1{
		return "",""
	}
	return name[:lastI],name[lastI + 1:]
}
func walkPages(dirToSave,dirToWalk string){
	walkMarkdown(dirToSave,dirToWalk,false)
}

func walkPosts(dirToSave , dirToWalk string){
	walkMarkdown(dirToSave,dirToWalk,true)
}

func walkMarkdown(dirToSave,dirToWalk string,isPost bool){
	utils.RemoveIfExist(dirToSave)
	os.Mkdir(dirToSave,0777)
	var currentSaveDir string

	filepath.Walk(dirToWalk, func(path string, info os.FileInfo, err error) error {
		if path == dirToWalk {
			return nil
		}

		if info.IsDir() {
			currentSaveDir = dirToSave + "/" + info.Name()
			utils.MkdirIfNotExist(currentSaveDir)
		}else {
			htmlFileName,extName := getFileNameAndExtName(info.Name())
			htmlFileName = htmlFileName + ".html"
			if extName == "md"{
				// Get html source
				head,doc := utils.ConvertMarkdownToHTMLFromFile(path)
				currentSaveFile := ""
				if isPost {
					currentSaveFile = currentSaveDir + "/" + htmlFileName
				}else{
					currentSaveFile = dirToSave + "/" + htmlFileName
				}
				utils.CreateIfNotExist(currentSaveFile)
				w,err := os.OpenFile(currentSaveFile,os.O_WRONLY,0666)
				if err != nil {
					log.Panicf("There is an error : %v ." , err)
				}

				params := make(map[string]interface{})
				params["head"] = head
				params["srcF"] = path
				params["dstF"] = currentSaveFile
				params["document"] = doc
				params["isPost"] = isPost

				callEachPlugins(params)

				html ,_:= doc.Html()
				io.WriteString(w,html)
			}
		}
		return nil
	})

}

func callEachPlugins(params map[string]interface{}) {
	callBuildInPlugins(params)
	pluginsMap := scheme.GetYlogProjectConfig().PluginsMap
	for k,v := range pluginsMap{
		plugin ,ok := GetPluginFromLib(k)

		if !ok {
			log.Panicf("Plugin [ %s ] is not found ! Try run `Ylog install %s to install it." , k,k)
		}
		if !plugin.IsCommandLinePlugin() {
			for sk,sv := range v{
				params[sk] = sv
			}
			plugin.Exec(params)
		}
	}
}

var pageHandler = NewPageHandler()
var buildInPlugins = []Plugin{
	// Handle head of markdown file
	&HeadHandler{},
	// Add theme to html file
	&ThemeHandler{},
	// Print log to screen
	&Logger{},
	// To add some elements to post html file
	&PostInit{},
	// Create default page(index,categoires...)
	pageHandler,
}

func callBuildInPlugins(params map[string]interface{}) {
	for _ , handler := range buildInPlugins{
		handler.Exec(params)
	}
}


func (c *Build) Exec(params map[string]interface{}){
	log.Println("Building")
	currentPath , err := utils.GetCurrentPath()
	if err != nil {
		log.Panicf("There is an error : %v\n",err)
	}else{
		postPath := currentPath + "/posts"
		pagePath := currentPath + "/pages"
		_, err :=os.Stat(postPath)
		if os.IsNotExist(err){
			log.Panicf("Can't find post dir in %s. Run `mkdir posts`.",currentPath)
		}
		_,err = os.Stat(pagePath)
		if os.IsNotExist(err) {
			log.Panicf("Can't find page dir in %s. Run `mkdir pages`.",currentPath)
		}
		savePath := currentPath + "/out"
		walkPosts(savePath,postPath)
		pageHandler.PostRenderOver()
		walkPages(savePath,pagePath)
	}
}

func init(){
	log.Println("Init compiler(build)...")
	Registe(&Build{})
}