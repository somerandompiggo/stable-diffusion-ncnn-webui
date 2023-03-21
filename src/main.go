package main

import (
	"fmt"
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
	http.ServeFile(w, r, "../web/history.html")
}

func main() {
	if _, err := os.Stat("../web/images"); os.IsNotExist(err) {
		err := os.Mkdir("../web/images", 0777)
		check(err)
	}

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
