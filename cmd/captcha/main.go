package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"time"

	"github.com/cam-per/licob/picker"
)

var (
	sizes = make(map[int]image.Rectangle)
)

func runes(path string) {
	//fmt.Println(path)

	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	captcha, err := picker.NewCaptchaDecoder(f)
	if err != nil {
		log.Fatal(err)
	}

	start := time.Now()
	err = captcha.Decode()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(time.Since(start))

	for _, rect := range captcha.Rects() {
		hash := rect.Dx()
		hash <<= 16
		hash |= rect.Dy()

		if _, ok := sizes[hash]; ok {
			continue
		}
		sizes[hash] = rect
	}

	captcha.SaveAreas("runes")
}

func main() {
	// dir := filepath.Join(".", "images")
	// files, err := os.ReadDir(dir)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// for _, info := range files {
	// 	if info.IsDir() {
	// 		return
	// 	}
	// 	runes(filepath.Join(dir, info.Name()))
	// }

	f, err := os.Open("1147761101791051857.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	captcha, err := picker.NewCaptchaDecoder(f)
	if err != nil {
		log.Fatal(err)
	}

	err = captcha.Decode()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(captcha.Codes())
}
