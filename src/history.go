package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"os"
)

type ImageStruct struct {
	Time      int64
	Image     string
	Width     int
	Height    int
	Speed     int
	Stepcount int
	Seed      int
	Positive  string
	Negative  string
}

var historyfile []ImageStruct

func loadHistoryFile() {
	f, err := os.ReadFile("./history.json")
	if err != nil {
		_, err := os.Create("./history.json")
		check(err)
	}
	err = json.Unmarshal(f, &historyfile)
	check(err)
}

func updateHistory() {
	f, err := os.Create("./history.json")
	check(err)

	marshalled, err := json.Marshal(historyfile)
	_, err = f.Write(marshalled)
	check(err)
}

func addEntry(s ImageStruct) {
	historyfile = append(historyfile, s)
	updateHistory()
}

func deleteEntry(i int) {
	removeFromSlice(historyfile, i)
	updateHistory()
}

func deleteAll() {
	historyfile = []ImageStruct{}
	updateHistory()
}

func getHistoryHTML() string {
	tmpl, err := template.ParseFiles("../web/historyimage.html")
	check(err)

	var buf bytes.Buffer
	var result string
	for i := 0; i < len(historyfile); i++ {
		err = tmpl.ExecuteTemplate(&buf, "layout", historyfile[i])
		result += buf.String()
	}
	return result
}
