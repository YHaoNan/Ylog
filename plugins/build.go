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
// Walk posts dir in current folder.
func walkPosts(dirToSave , dirToWalk string){
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
				currentSaveFile := currentSaveDir + "/" + htmlFileName
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
	for k,_ := range pluginsMap{
		plugin ,ok := GetPluginFromLib(k)
		if !ok {
			log.Panicf("Plugin [ %s ] is not found ! Try run `Ylog install %s to install it." , k,k)
		}
		plugin.Exec(params)
	}
}

func callBuildInPlugins(params map[string]interface{}) {
	headHandler := &HeadHandler{}
	headHandler.Exec(params)
	logHandler := &Logger{}
	logHandler.Exec(params)
}


func (c *Build) Exec(params map[string]interface{}){
	log.Println("Building")
	currentPath , err := utils.GetCurrentPath()
	if err != nil {
		log.Panicf("There is an error : %v\n",err)
	}else{
		postPath := currentPath + "/posts"
		_, err :=os.Stat(postPath)
		if os.IsNotExist(err){
			log.Panicf("Can't find post dir in %s. Run `mkdir posts`.",currentPath)
		}
		savePath := currentPath + "/out"
		walkPosts(savePath,postPath)
	}
}

func init(){
	log.Println("Init compiler(build)...")
	Registe(&Build{})
}