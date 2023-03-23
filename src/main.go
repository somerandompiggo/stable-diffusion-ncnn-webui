package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	if _, err := os.Stat("../build/stable-diffusion-ncnn"); os.IsNotExist(err) {
		http.ServeFile(w, r, "../web/notpresent.html")
	} else {
		http.ServeFile(w, r, "../web/index.html")
	}
}

func getHistoryPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("../web/history.html")
	check(err)

	if len(historyfile) == 0 {
		err = tmpl.Execute(w, "You haven't generated any images yet!")
	} else {
		err = tmpl.Execute(w, template.HTML(getHistoryHTML()))
	}
	check(err)
}

func main() {
	loadHistoryFile()
	if _, err := os.Stat("../web/images"); os.IsNotExist(err) {
		err := os.Mkdir("../web/images", 0777)
		check(err)
	}

	sdfs := http.FileServer(http.Dir("../build"))
	http.Handle("/sd/", http.StripPrefix("/sd/", sdfs))

	http.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../web/css/main.css")
	})
	http.HandleFunc("/font/inter.ttf", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../web/font/inter.ttf")
	})

	http.HandleFunc("/result.png", func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat("../web/images/image.png"); os.IsNotExist(err) {
			http.ServeFile(w, r, "../web/images/fallback.png")
		} else {
			http.ServeFile(w, r, "../web/images/image.png")
		}
	})

	http.HandleFunc("/generate", generate)
	http.HandleFunc("/history", getHistoryPage)

	http.HandleFunc("/", getRoot)

	fmt.Println("Starting HTTP server on port 8080")
	err := http.ListenAndServe(":8080", nil)
	check(err)
}
