package main

import (
	"fmt"
	"net/http"
	"os"
)

func magicFile(width int, height int, speed int, stepcount int, seed int, positive string, negative string) string {
	return fmt.Sprint(width) + "\n" + fmt.Sprint(height) + "\n" + fmt.Sprint(speed) + "\n" + fmt.Sprint(stepcount) + "\n" + fmt.Sprint(seed) + "\n" + positive + "\n" + negative
}

func check(err error) {
	if err != nil {
		fmt.Println("An error occured:", err)
		os.Exit(1)
	}
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	if _, err := os.Stat("../stable-diffusion-ncnn"); os.IsNotExist(err) {
		http.ServeFile(w, r, "../web/notpresent.html")
	} else {
		http.ServeFile(w, r, "../web/index.html")
		err := os.WriteFile("./magic.txt", []byte(magicFile(256, 256, 1, 15, 92587, "cat", "dog")), 0644)
		check(err)
	}
}
func main() {
	http.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../web/css/main.css")
	})
	http.HandleFunc("/font/inter.ttf", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../web/font/inter.ttf")
	})

	http.HandleFunc("/result.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../web/images/demoimg.png")
	})

	http.HandleFunc("/", getRoot)

	fmt.Println("Starting HTTP server on port 8080")
	err := http.ListenAndServe(":8080", nil)
	check(err)
}
