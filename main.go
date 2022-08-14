package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"paperfinderv3/paperloader"
	"paperfinderv3/server"
	"path"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Cannot find .env file")
		os.Exit(1)
	}

	_, err = strconv.Atoi(os.Getenv("SEARCH_THRESHOLD"))
	if err != nil {
		fmt.Println("Invalid search threshold")
		os.Exit(1)
	}
}

func main() {

	papers := flag.Bool("papers", false, "If set will redownload papers")
	index := flag.Bool("index", false, "If set will regenerate index")

	flag.Parse()

	if *papers {
		paperloader.LoadPapers()
	}

	if *index {
		paperloader.GenerateIndex()
	}

	gin.SetMode(gin.ReleaseMode)
	fmt.Println("Starting server...")

	r := gin.New()
	r.SetFuncMap(template.FuncMap{
		"title": cases.Title(language.English).String,
		"len": func(a []server.SearchResult) int {
			return len(a)
		},
		"mangle": server.MangleFilename,
		"img": func(in string) string {
			ext := path.Ext(in)
			return in[0:len(in)-len(ext)] + ".jpg"
		},
	})
	r.LoadHTMLGlob("web/*_template.html")

	r.GET("/", server.GenerateIndex)
	r.GET("/search", server.GenerateSearchResults)
	r.GET("/file/:filename", server.GetPaper)
	r.GET("/image/:filename", server.GetImage)

	r.NoRoute(server.NotFound)

	r.StaticFile("/css/index.css", "web/css/index.css")
	r.StaticFile("/css/search.css", "web/css/search.css")

	r.StaticFS("/assets", http.Dir("web/assets"))

	r.Run("0.0.0.0:" + os.Getenv("HTTP_PORT"))
}
