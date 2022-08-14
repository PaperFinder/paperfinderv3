package server

import (
	"encoding/json"
	"net/http"
	"os"
	"paperfinderv3/paperloader"
	"sort"

	"github.com/gin-gonic/gin"
)

type IndexTemplate struct {
	Subjects []string
}

func GenerateIndex(c *gin.Context) {
	i := IndexTemplate{
		Subjects: getSubjectList(),
	}

	c.HTML(http.StatusOK, "index_template.html", i)
}

func getSubjectList() []string {
	subjects := map[string]struct{}{}

	type ex map[string]map[string]struct {
		Template string `json:"template"`
		LastUnit int    `json:"lastunit"`
	}

	jsonfile, err := os.Open(os.Getenv("PAPER_CONFIG"))
	if err != nil {
		panic(err)
	}
	defer jsonfile.Close()

	papers := map[string]ex{}
	err = json.NewDecoder(jsonfile).Decode(&papers)
	if err != nil {
		panic(err)
	}

	for e := range papers {
		for b := range papers[e] {
			for s := range papers[e][b] {
				if _, ok := paperloader.Exclude[s]; !ok {
					subjects[s] = struct{}{}
				}
			}
		}
	}

	list := make([]string, 0, len(subjects))

	for k := range subjects {
		list = append(list, k)
	}

	sort.Strings(list)

	return list
}
