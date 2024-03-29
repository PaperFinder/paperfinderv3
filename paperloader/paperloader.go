package paperloader

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/schollz/progressbar/v3"
)

var IsLoading = false

type Exam map[string]map[string]struct {
	Template string `json:"template"`
	LastUnit int    `json:"lastunit"`
}

func LoadPapers() {
	IsLoading = true
	jsonfile, err := os.Open(os.Getenv("PAPER_CONFIG"))
	if err != nil {
		IsLoading = false
		panic(err)
	}
	defer jsonfile.Close()

	papers := map[string]Exam{}
	err = json.NewDecoder(jsonfile).Decode(&papers)
	if err != nil {
		IsLoading = false
		panic(err)
	}

	pfolder := os.Getenv("PAPER_FOLDER")
	fmt.Println("Redownloading papers...")

	pfoldertemp := pfolder + "_temp"
	os.RemoveAll(pfoldertemp)
	os.Mkdir(pfoldertemp, os.ModePerm)
	os.Chdir(pfoldertemp)

	for ename, exam := range papers {
		os.Mkdir(ename, os.ModePerm)
		os.Chdir(ename)
		for bname, board := range exam {
			os.Mkdir(bname, os.ModePerm)
			os.Chdir(bname)
			for sname, subject := range board {
				os.Mkdir(sname, os.ModePerm)
				os.Chdir(sname)

				fmt.Printf("\n%s %s\n", bname, sname)

				for unit := 1; unit <= subject.LastUnit; unit++ {
					os.Mkdir(fmt.Sprintf("u%d", unit), os.ModePerm)
					os.Chdir(fmt.Sprintf("u%d", unit))

					fmt.Println("\nUnit", unit)

					res, err := http.Get(fmt.Sprintf(subject.Template, unit))
					if err != nil {
						IsLoading = false
						panic(err)
					}

					doc, _ := goquery.NewDocumentFromReader(res.Body)
					sel := doc.Find(".files li")
					total := sel.Length()
					bar := progressbar.NewOptions(total, progressbar.OptionSetWidth(15))
					sel.Each(func(i int, s *goquery.Selection) {
						q, _ := s.Find("a").First().Attr("href")
						n := s.Text()

						filename := generateFilename(n)
						if filename == "" {
							total -= 1
							bar.ChangeMax(total)
							return
						}

						resp, _ := http.Get(q)
						if resp.StatusCode != 200 {
							fmt.Println("Error with", filename, ":", resp.Status)
						}

						f, _ := os.Create(filename)

						io.Copy(f, resp.Body)
						bar.Add(1)

						ext := path.Ext(filename)
						imgFile := filename[0:len(filename)-len(ext)] + ".jpg"

						cmd := exec.Command("convert", "-density", "300", filename+"[0]", "-background", "#ffffff", "-flatten", imgFile)

						err = cmd.Run()
						if err != nil {
							IsLoading = false
							panic(err)
						}

						resp.Body.Close()
						f.Close()
					})

					res.Body.Close()
					os.Chdir("..")
				}

				os.Chdir("..")
			}
			os.Chdir("..")
		}
		os.Chdir("..")
	}

	os.Chdir("..")
	os.RemoveAll(pfolder)
	os.Rename(pfoldertemp, pfolder)

	IsLoading = false
}

func generateFilename(name string) string {
	name = strings.TrimSpace(strings.Split(name, " - ")[0])

	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "(", "")
	name = strings.ReplaceAll(name, ")", "")
	name += ".pdf"

	if name == "grade_boundaries.pdf" {
		return ""
	}
	if strings.Contains(name, "combined") {
		return ""
	}
	if strings.Contains(name, "ms") {
		return ""
	}

	return name
}
