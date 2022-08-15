package server

import (
	"net/http"
	"os"
	"paperfinderv3/paperloader"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type SearchResultsTemplate struct {
	Subjects []string
	Subject  string
	Query    string
	Error    string
	Time     string
	Results  []SearchResult
}

type SearchResult struct {
	Specimen bool
	Exam     string
	Board    string
	Unit     string
	Month    string
	Year     string
	IAL      bool
	URL      string
}

func GenerateSearchResults(c *gin.Context) {
	subject, oks := c.GetQuery("s")
	query, okq := c.GetQuery("q")

	srT := SearchResultsTemplate{
		Subjects: getSubjectList(),
		Subject:  subject,
		Query:    query,
	}

	if !oks || !okq {
		c.JSON(http.StatusBadRequest, "query or subject is empty")
		return
	}

	start := time.Now()
	results, err := paperloader.Search(query, subject)
	if err != nil {
		srT.Error = err.Error()
		c.HTML(http.StatusOK, "search_template.html", srT)
	}
	elapsed := time.Since(start)
	srT.Time = elapsed.Round(time.Millisecond * 10).String()

	srT.Results = make([]SearchResult, 0, len(results))

	for _, res := range results {
		r := strings.Split(res, string(os.PathSeparator))
		date := strings.Split(r[5], "_")
		s := SearchResult{
			Exam:  r[1],
			Board: r[2],
			URL:   res,
		}
		if !strings.Contains(res, "specimen") {
			s.Unit = r[4]
			s.Month = date[0]
			s.Year = date[1]
			if len(date) == 4 {
				s.IAL = true
			} else {
				s.IAL = false
			}
		} else {
			s.Specimen = true
			if _, err := strconv.Atoi(date[1]); err != nil {
				s.IAL = true
			} else {
				if len(date) == 4 {
					s.IAL = true
				}
				s.Year = date[1]
			}
		}
		srT.Results = append(srT.Results, s)
	}

	c.HTML(http.StatusOK, "search_template.html", srT)
}
