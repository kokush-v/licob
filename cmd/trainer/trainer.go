package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/therfoo/therfoo/tensor"
)

const dict = "0123456789abcdef"

type data struct {
	in  tensor.Vector
	out tensor.Vector
}

type trainer struct {
	items []data
}

func (t *trainer) Get(index int) (*[]tensor.Vector, *[]tensor.Vector) {
	data := t.items[index]
	return &[]tensor.Vector{data.in}, &[]tensor.Vector{data.out}
}

func (t *trainer) Len() int {
	return len(t.items)
}

func read(path string) ([]float64, error) {
	var result []float64
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		return nil, err
	}

	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			rgba := img.(*image.NRGBA).At(x, y).(color.NRGBA)
			var data float64 = 0.0
			if rgba.R == 255 && rgba.G == 255 && rgba.B == 255 {
				data = 1.0
			}
			result = append(result, data)
		}
	}
	fmt.Println(len(result))
	return result, nil
}

func (t *trainer) fetch(dir string) error {
	return filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		d, err := read(path)
		if err != nil {
			return err
		}

		response := make([]float64, 16)
		dir := strings.Split(path, string(os.PathSeparator))
		index := strings.Index(dict, dir[1])
		if index != -1 {
			response[index] = 1.0
		}
		t.items = append(t.items, data{
			in:  d,
			out: response,
		})
		return nil
	})
}
