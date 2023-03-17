package main

import (
	"fmt"
	"net/http"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../web/index.html")
}

func getStyle(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../web/css")
}

func main() {
	http.HandleFunc("/", getRoot)
	http.HandleFunc("/style.css", getStyle)
	fmt.Println("Starting HTTP server on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Print("Error starting HTTP server: ", err)
	}
}
