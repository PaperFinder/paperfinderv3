package server

import (
	"encoding/base64"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func MangleFilename(in string) string {
	key := os.Getenv("MANGLE_KEY")

	data := make([]byte, 0, len(in))

	for i := 0; i < len(in); i++ {
		b := in[i] ^ key[i%len(key)]
		data = append(data, b)
	}

	result := base64.URLEncoding.EncodeToString(data)

	return result
}

func UnmangleFilename(in string) (string, error) {
	key := os.Getenv("MANGLE_KEY")
	decoded, err := base64.URLEncoding.DecodeString(in)
	if err != nil {
		return "", err
	}

	data := make([]byte, 0)

	for i := 0; i < len(decoded); i++ {
		b := decoded[i] ^ key[i%len(key)]
		data = append(data, b)
	}
	return string(data), nil
}

func GetPaper(c *gin.Context) {
	mangled := c.Param("filename")

	data, err := UnmangleFilename(mangled)
	if err != nil {
		NotFound(c)
		return
	}

	fdata, err := os.ReadFile(data)
	if err != nil {
		NotFound(c)
		return
	}

	c.Data(http.StatusOK, "application/pdf", fdata)
}

func GetImage(c *gin.Context) {
	mangled := c.Param("filename")

	data, err := UnmangleFilename(mangled)
	if err != nil {
		NotFound(c)
		return
	}

	fdata, err := os.ReadFile(data)
	if err != nil {
		NotFound(c)
		return
	}

	c.Data(http.StatusOK, "image/jpg", fdata)
}

func NotFound(c *gin.Context) {
	c.Data(http.StatusNotFound, "text/html", []byte{})
}
