package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

func magicFile(width int, height int, speed int, stepcount int, seed int, positive string, negative string) string {
	return fmt.Sprint(height) + "\n" + fmt.Sprint(width) + "\n" + fmt.Sprint(speed) + "\n" + fmt.Sprint(stepcount) + "\n" + fmt.Sprint(seed) + "\n" + positive + "\n" + negative
}

func check(err error) {
	if err != nil {
		fmt.Println("An error occured:", err)
		os.Exit(1)
	}
}

func generate(w http.ResponseWriter, r *http.Request) {
	reqparams := r.URL.Query()

	width, err := strconv.Atoi(reqparams.Get("width"))
	check(err)
	height, err := strconv.Atoi(reqparams.Get("height"))
	check(err)
	steps, err := strconv.Atoi(reqparams.Get("numSteps"))
	check(err)
	promptP := reqparams.Get("positivePrompt")
	promptN := reqparams.Get("negativePrompt")

	SDGenerate(magicFile(width, height, 1, steps, 92587, promptP, promptN))
	http.Redirect(w, r, "/", http.StatusFound)
}

func SDGenerate(magicfile string) {
	err := os.WriteFile("../build/magic.txt", []byte(magicfile), 0644)
	check(err)

	cmd := exec.Command("../build/stable-diffusion-ncnn")
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	if _, err := os.Stat("../build/stable-diffusion-ncnn"); os.IsNotExist(err) {
		http.ServeFile(w, r, "../web/notpresent.html")
	} else {
		http.ServeFile(w, r, "../web/index.html")
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

	http.HandleFunc("/generate", generate)

	http.HandleFunc("/", getRoot)

	fmt.Println("Starting HTTP server on port 8080")
	err := http.ListenAndServe(":8080", nil)
	check(err)
}
