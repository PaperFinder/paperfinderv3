package paperloader

import (
	"database/sql"
	"fmt"
	"log"
	"path"
	"sort"
	"time"
)

func Search(query string, subject string) ([]string, error) {
	tokens := filterText(query)
	dbPath := path.Join("database", "index.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	resultsMap := make(map[string]int)

	for _, tok := range tokens {

		rows, err := db.Query("SELECT file FROM tokens WHERE term = ? AND subject = ?;", tok, subject)
		if err != nil || rows.Err() != nil {
			t := time.Now().Unix()
			erro := fmt.Errorf("E%x", t)
			log.Println(err)
			return nil, erro
		}

		for rows.Next() {
			fname := ""
			rows.Scan(&fname)
			resultsMap[fname]++
		}

		rows.Close()
	}

	papers := make([]string, 0, len(resultsMap))
	for k, v := range resultsMap {
		if float64(v) >= float64(len(tokens))*0.8 {
			papers = append(papers, k)
		}
	}

	sort.SliceStable(papers, func(i, j int) bool {
		return resultsMap[papers[i]] < resultsMap[papers[j]]
	})

	return papers, nil
}
