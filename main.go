package main

import (
	"paperfinderv3/paperloader"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	paperloader.LoadPapers()
}
