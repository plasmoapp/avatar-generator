package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

var eyes []string
var mouths []string

var light_colors []string
var dark_colors []string

func main() {
	/* Read and cache eyes */
	file, err := os.Open("./eyes")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		eyes = append(eyes, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	/* Read and cache mouths */
	file, err = os.Open("./mouths")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		mouths = append(mouths, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	/* Read and cache light colors */
	file1, err := os.Open("./light_colors")
	if err != nil {
		log.Fatal(err)
	}

	defer file1.Close()

	scanner = bufio.NewScanner(file1)
	for scanner.Scan() {
		light_colors = append(light_colors, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	/* Read and cache dark colors */
	file2, err := os.Open("./dark_colors")
	if err != nil {
		log.Fatal(err)
	}

	defer file2.Close()

	scanner = bufio.NewScanner(file2)
	for scanner.Scan() {
		dark_colors = append(dark_colors, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// router := mux.NewRouter()
	// router.HandleFunc("/", get).Methods("GET")
	// log.Fatal(http.ListenAndServe(":80", router))

	r := routing.New()

	r.Get("/", getFast)

	server := fasthttp.Server{
		Handler:      r.HandleRequest,
		Name:         "Plasmo",
		LogAllErrors: false,
	}

	//panic(server.ListenAndServe(":80"))
	panic(server.ListenAndServeTLS(":4097", "./certs/cert.crt", "./certs/s.key"))
}

func getFast(c *routing.Context) error {
	id := string(c.URI().QueryArgs().Peek("id"))
	theme := string(c.URI().QueryArgs().Peek("theme"))
	inv := string(c.URI().QueryArgs().Peek("inverted"))
	inverted := false

	if id == "" {
		fmt.Fprint(c, "ID can't be empty")
		return nil
	}

	if len(id) > 256 {
		fmt.Fprint(c, "ID max length is 256")
		return nil
	}

	if theme != "light" && theme != "dark" {
		theme = "light"
	}

	if inv != "" {
		var err error
		inverted, err = strconv.ParseBool(inv)
		if err != nil {
			inverted = false
		}
	}

	var hash int
	for i := 0; i < len(id); i++ {
		hash = int(id[i]) + ((hash << 8) - hash)
		hash = hash & hash
	}

	var generated [4]int // [Eyes, Mouth, MainColor, SecondColor]
	generated[0] = ((hash % len(eyes)) + len(eyes)) % len(eyes)
	generated[1] = ((hash % len(mouths)) + len(mouths)) % len(mouths)

	var th []string

	generated[2] = ((hash % (len(light_colors) - 1)) + (len(light_colors) - 1)) % (len(light_colors) - 1)
	generated[3] = len(light_colors) - 1

	if theme == "light" {
		th = light_colors
	} else {
		th = dark_colors
	}

	if inverted {
		t := generated[2]
		generated[2] = generated[3]
		generated[3] = t
	}

	var svg string
	svg += `<svg width="58" height="58" viewBox="0 0 58 58" fill="none" xmlns="http://www.w3.org/2000/svg">`
	svg += fmt.Sprintf(`<path d="M58 29C58 45.0163 45.0163 58 29 58C12.9837 58 0 45.0163 0 29C0 12.9837 12.9837 0 29 0C45.0163 0 58 12.9837 58 29Z" fill="%v"/>`, th[generated[2]])
	svg += `<g transform="translate(8,25)">`
	svg += strings.Replace(eyes[generated[0]], "%v", th[generated[3]], -1)
	svg += `</g>`
	svg += `<g transform="translate(5,33)">`
	svg += strings.Replace(mouths[generated[1]], "%v", th[generated[3]], -1)
	svg += `</g>`
	svg += `</svg>`

	c.Response.Header.Set("Content-Type", "image/svg+xml")
	c.Response.Header.Set("Content-Disposition", `inline; filename="avatar.svg"`)
	c.Response.Header.Set("Content-Length", strconv.Itoa(len([]byte(svg))))

	c.Response.BodyWriter().Write([]byte(svg))
	return nil
}

// func get(w http.ResponseWriter, r *http.Request) {
// 	id := r.URL.Query().Get("id")
// 	theme := r.URL.Query().Get("theme")
// 	inv := r.URL.Query().Get("inverted")
// 	inverted := false

// 	if id == "" {
// 		fmt.Fprint(w, "ID can't be empty")
// 		return
// 	}

// 	if theme != "light" && theme != "dark" {
// 		theme = "light"
// 	}

// 	if inv != "" {
// 		var err error
// 		inverted, err = strconv.ParseBool(inv)
// 		if err != nil {
// 			inverted = false
// 		}
// 	}

// 	var hash int
// 	for i := 0; i < len(id); i++ {
// 		hash = int(id[i]) + ((hash << 8) - hash)
// 		hash = hash & hash
// 	}

// 	var generated []int
// 	// [Eyes, Mouth, MainColor, SecondColor]
// 	generated = append(generated, (hash>>uint(0))&(len(eyes)-1))
// 	generated = append(generated, (hash>>uint(8))&(len(mouths)-1))

// 	var th []string

// 	if theme == "light" {
// 		th = light_colors
// 		generated = append(generated, (hash>>uint(16))&(len(light_colors)-2))
// 		generated = append(generated, len(light_colors)-1)
// 	} else {
// 		th = light_colors
// 		generated = append(generated, (hash>>uint(24))&(len(light_colors)-2))
// 		generated = append(generated, len(light_colors)-1)
// 	}

// 	if inverted {
// 		t := generated[2]
// 		generated[2] = generated[3]
// 		generated[3] = t
// 	}

// 	var svg string
// 	svg += `<svg width="58" height="58" viewBox="0 0 58 58" fill="none" xmlns="http://www.w3.org/2000/svg">`
// 	svg += fmt.Sprintf(`<path d="M58 29C58 45.0163 45.0163 58 29 58C12.9837 58 0 45.0163 0 29C0 12.9837 12.9837 0 29 0C45.0163 0 58 12.9837 58 29Z" fill="%v"/>`, th[generated[2]])
// 	svg += `<g transform="translate(8,25)">`
// 	svg += strings.Replace(eyes[generated[0]], "%v", th[generated[3]], -1)
// 	svg += `</g>`
// 	svg += `<g transform="translate(5,33)">`
// 	svg += strings.Replace(mouths[generated[1]], "%v", th[generated[3]], -1)
// 	svg += `</g>`
// 	svg += `</svg>`

// 	w.Header().Set("Content-Type", "image/svg+xml")
// 	w.Header().Set("Content-Disposition", `inline; filename="avatar.svg"`)
// 	w.Header().Set("Content-Length", strconv.Itoa(len([]byte(svg))))

// 	w.Write([]byte(svg))
// }
