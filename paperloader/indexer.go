package paperloader

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/dlclark/regexp2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/schollz/progressbar/v3"
)

var db *sql.DB

// EXCLUDE CHEMISTRY BECAUSE PDFTOTEXT CANT READ IT

var Exclude = map[string]struct{}{
	"chemistry": {},
}

func GenerateIndex() {
	IsLoading = true

	fmt.Println("PRESS ENTER TO REGENERATE INDEX. PRESS CTRL+C TO CANCEL")
	a := ""
	fmt.Scanln(&a)

	dbPath := path.Join("database", "index.db")
	fmt.Println("Rebuilding index...")

	os.RemoveAll("database")
	os.Mkdir("database", os.ModePerm)
	f, err := os.Create(dbPath)
	if err != nil {
		IsLoading = false
		panic(err)
	}
	f.Close()

	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		IsLoading = false
		panic(err)
	}
	defer db.Close()

	st, _ := db.Prepare("CREATE TABLE IF NOT EXISTS `tokens` (`term` VARCHAR NOT NULL,`file` VARCHAR NOT NULL, `subject` VARCHAR NOT NULL);")
	_, err = st.Exec()
	if err != nil {
		IsLoading = false
		panic(err)
	}
	st.Close()

	err = filepath.Walk(os.Getenv("PAPER_FOLDER"),
		func(fpath string, info fs.FileInfo, err error) error {
			if err != nil {
				IsLoading = false
				panic(err)
			}
			if info.IsDir() {
				return nil
			}
			if path.Ext(fpath) != ".pdf" {
				return nil
			}

			if _, ok := Exclude[strings.Split(fpath, string(os.PathSeparator))[3]]; ok {
				return nil
			}

			cmd := exec.Command("pdftotext", "-f", "2", "-q", fpath)

			err = cmd.Run()
			if err != nil {
				IsLoading = false
				panic(err)
			}
			ext := path.Ext(fpath)
			output := fpath[0:len(fpath)-len(ext)] + ".txt"

			data, err := os.ReadFile(output)
			if err != nil {
				IsLoading = false
				panic(err)
			}
			os.Remove(output)
			newData := filterText(string(data))

			indexAndWrite(newData, fpath, strings.Split(fpath, string(os.PathSeparator))[3])

			return nil
		})
	if err != nil {
		IsLoading = false
		panic(err)
	}
}

func filterText(dataIn string) []string {
	ePass1 := `(\*\S\S{8,11}\*)|(\(This page is for your |(first)|(second)|(third)|(fourth)|(fifth)| answer\))|(\x0c)|(\n[0-9]{1,2}\n)|(PMT)|(\.[ \.]{1,}[^A-Z][ \n])|(\(Total for Question [0-9] = [0-9]{1,2} marks{0,}\))|(Question [0-9]{1,2})|(Turn over)|(BLANK PAGE)|(TOTAL FOR PAPER = [0-9]{1,2} MARKS)|(\([0-9]{1,2}\))|(\([a-z]{1,3}\))|(DO NOT WRITE IN THIS AREA)|([0-9]{1,} )`
	ePass2 := `\n`
	ePass3 := `\.[^\d]`
	ePass4 := `,`
	ePass5 := `\s+`

	expr := regexp2.MustCompile(ePass1, 0)
	out, _ := expr.Replace(string(dataIn), "", -1, 1000000000)
	expr = regexp2.MustCompile(ePass2, 0)
	out, _ = expr.Replace(out, " ", -1, 1000000000)
	expr = regexp2.MustCompile(ePass3, 0)
	out, _ = expr.Replace(out, " ", -1, 1000000000)
	expr = regexp2.MustCompile(ePass4, 0)
	out, _ = expr.Replace(out, " ", -1, 1000000000)
	expr = regexp2.MustCompile(ePass5, 0)
	out, _ = expr.Replace(out, " ", -1, 1000000000)

	out = strings.ToLower(out)
	out = strings.TrimSpace(out)

	words := strings.Split(out, " ")
	fw := make([]string, 0, len(words))
	commonWords := map[string]struct{}{
		"the":  {},
		"of":   {},
		"and":  {},
		"a":    {},
		"to":   {},
		"in":   {},
		"is":   {},
		"you":  {},
		"that": {},
		"it":   {},
		"he":   {},
		"was":  {},
		"for":  {},
		"on":   {},
		"are":  {},
		"as":   {},
		"with": {},
		"his":  {},
		"they": {},
		"i":    {},
		"if":   {},
	}

	for _, word := range words {
		if _, ok := commonWords[word]; ok {
			continue
		}
		fw = append(fw, word)
	}

	fw = unique(fw)

	if len(fw) == 1 && fw[0] == "" {
		return []string{}
	}
	return fw
}

func indexAndWrite(tokens []string, fpath string, subject string) {
	fmt.Println(fpath)
	bar := progressbar.NewOptions(len(tokens), progressbar.OptionSetWidth(len(fpath)))
	for _, tok := range tokens {
		st, _ := db.Prepare("INSERT INTO tokens (term, file, subject) VALUES  (?, ?, ?);")
		_, err := st.Exec(tok, fpath, subject)
		if err != nil {
			panic(err)
		}
		bar.Add(1)
	}
	fmt.Println("")
}

func unique(s []string) []string {
	inResult := make(map[string]bool)
	var result []string
	for _, str := range s {
		if _, ok := inResult[str]; !ok {
			inResult[str] = true
			result = append(result, str)
		}
	}
	return result
}
