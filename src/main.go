package main

import (
	"fmt"
	"net/http"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../web/index.html")
}

func main() {
	http.HandleFunc("/", getRoot)
	http.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../web/css/main.css")
	})
	http.HandleFunc("/font/inter.ttf", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../web/font/inter.ttf")
	})

	fs := http.FileServer(http.Dir("../web/images/"))
	http.Handle("/images", fs)

	fmt.Println("Starting HTTP server on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Print("Error starting HTTP server: ", err)
	}
}
