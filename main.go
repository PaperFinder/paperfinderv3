package main

import (
	"os"
	"paperfinderv3/paperloader"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	paperloader.LoadPapers()
	paperloader.GenerateIndex()

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.Run(":" + os.Getenv("HTTP_PORT"))
}
