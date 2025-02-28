package main

import (
	"fmt"
	"html/template"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
)

var (
	tmpl     *template.Template
	gradient = map[string][]string{
		"default": {" ", ".", ":", "!", "/", "r", "(", "l", "Z", "4", "H", "9", "W", "8", "$", "@"},
		"reverse": {"@", "$", "8", "W", "9", "H", "4", "Z", "l", "(", "r", "/", "!", ":", ".", " "},
	}
)

type PageData struct {
	ASCII string
	Error string
	Form  struct {
		Width    int
		Gradient string
	}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/convert", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		r.ParseMultipartForm(10 << 20)
		width, _ := strconv.Atoi(r.FormValue("width"))
		style := r.FormValue("style")

		if width < 20 || width > 500 {
			http.Error(w, "Invalid width (20-500 allowed)", http.StatusBadRequest)
			return
		}

		file, _, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "Error reading file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			http.Error(w, "Error decoding image", http.StatusInternalServerError)
			return
		}

		resized := resizeImage(img, width)
		chars := gradient[style]
		if chars == nil {
			chars = gradient["default"]
		}
		asciiArt := imageToASCII(resized, chars)

		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, asciiArt)
	})

	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func resizeImage(img image.Image, newWidth int) image.Image {
	bounds := img.Bounds()
	originalWidth := bounds.Dx()
	originalHeight := bounds.Dy()

	aspectRatio := float64(originalHeight) / float64(originalWidth)
	newHeight := int(math.Round(float64(newWidth) * aspectRatio * 0.5))

	resized := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	xScale := float64(originalWidth) / float64(newWidth)
	yScale := float64(originalHeight) / float64(newHeight)

	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			srcX := int(math.Round(float64(x) * xScale))
			srcY := int(math.Round(float64(y) * yScale))
			if srcX >= originalWidth {
				srcX = originalWidth - 1
			}
			if srcY >= originalHeight {
				srcY = originalHeight - 1
			}
			resized.Set(x, y, img.At(srcX, srcY))
		}
	}
	return resized
}

func imageToASCII(img image.Image, gradient []string) string {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	maxIndex := float64(len(gradient) - 1)
	var sb strings.Builder

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			luminance := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
			pixelValue := luminance / 256

			index := int((pixelValue * maxIndex) / 255)
			if index >= len(gradient) {
				index = len(gradient) - 1
			}
			sb.WriteString(gradient[index])
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
