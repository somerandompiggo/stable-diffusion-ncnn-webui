package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"
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

func genFileName(stepcount int, seed int, height int, width int) string {
	return "result_" + fmt.Sprint(stepcount) + "_" + fmt.Sprint(seed) + "_" + fmt.Sprint(height) + "x" + fmt.Sprint(width) + ".png"
}

func setShownImageTo(image string) {
	if _, err := os.Stat("../build/stable-diffusion-ncnn"); os.IsNotExist(err) == false {
		os.Remove("../web/images/image.png")
	}

	err := os.Link("../build/"+image, "../web/images/image.png")
	check(err)
}

func generate(w http.ResponseWriter, r *http.Request) {
	randseed := rand.New(rand.NewSource(time.Now().UnixNano()))

	reqparams := r.URL.Query()

	width, err := strconv.Atoi(reqparams.Get("width"))
	check(err)
	height, err := strconv.Atoi(reqparams.Get("height"))
	check(err)
	stepcount, err := strconv.Atoi(reqparams.Get("numSteps"))
	check(err)
	promptP := reqparams.Get("positivePrompt")
	promptN := reqparams.Get("negativePrompt")
	seed := randseed.Intn(100000)

	go SDGenerate(width, height, 1, stepcount, seed, promptP, promptN)
	http.Redirect(w, r, "/", http.StatusFound)
}

func SDGenerate(width int, height int, speed int, stepcount int, seed int, positive string, negative string) {
	err := os.WriteFile("../build/magic.txt", []byte(magicFile(width, height, speed, stepcount, seed, positive, negative)), 0644)
	check(err)

	err = os.Chdir("../build")
	check(err)

	cmd := exec.Command("../build/stable-diffusion-ncnn")
	err = cmd.Run()
	check(err)

	err = os.Chdir("../src")
	check(err)

	setShownImageTo(genFileName(stepcount, seed, height, width))
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	if _, err := os.Stat("../build/stable-diffusion-ncnn"); os.IsNotExist(err) {
		http.ServeFile(w, r, "../web/notpresent.html")
	} else {
		http.ServeFile(w, r, "../web/index.html")
	}
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
		http.ServeFile(w, r, "../web/images/image.png")
	})

	http.HandleFunc("/generate", generate)

	http.HandleFunc("/", getRoot)

	fmt.Println("Starting HTTP server on port 8080")
	err := http.ListenAndServe(":8080", nil)
	check(err)
}
