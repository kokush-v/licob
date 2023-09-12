package picker

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/cam-per/licob/utils"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/therfoo/therfoo/layers/dense"
	"github.com/therfoo/therfoo/model"
	"github.com/therfoo/therfoo/optimizers/sgd"
	"github.com/therfoo/therfoo/tensor"
)

const dict = "0123456789abcdef"

var nn *model.Model

func init() {
	nn = model.New(
		model.WithCategoricalAccuracy(),
		model.WithCrossEntropyLoss(),
		model.WithEpochs(1000),
		model.WithInputShape(tensor.Shape{24 * 52}),
		model.WithOptimizer(
			sgd.New(sgd.WithBatchSize(5), sgd.WithLearningRate(0.005)),
		),
		model.WithVerbosity(true),
	)
	nn.Add(24, dense.New(dense.WithReLU()))
	nn.Add(16, dense.New(dense.WithSigmoid()))

	nn.Compile()
	err := nn.Load("ocr.gob")
	if err != nil {
		log.Fatal(err)
	}
}

type CaptchaDecoder struct {
	src    *gg.Context
	canvas *gg.Context
	valid  bool
	areas  []*area
	codes  []string
}

type area struct {
	src    *gg.Context
	canvas *gg.Context
	rect   image.Rectangle
	pix    map[int]image.Point
	data   tensor.Vector
}

func newArea(canvas *gg.Context) *area {
	return &area{
		src: canvas,
		pix: make(map[int]image.Point),
		rect: image.Rectangle{
			Min: image.Pt(math.MaxInt, math.MaxInt),
			Max: image.Pt(-1, -1),
		},
	}
}

func NewCaptchaDecoder(r io.Reader) (*CaptchaDecoder, error) {
	captcha := &CaptchaDecoder{
		valid: false,
	}

	img, err := png.Decode(r)
	if err != nil {
		return nil, err
	}

	captcha.src = gg.NewContextForImage(img)
	return captcha, nil
}

func checked(c color.Color) bool {
	rgba := c.(color.RGBA)
	return rgba.R == 255 && rgba.G == 255 && rgba.B == 255
}

func (captcha *CaptchaDecoder) crop() bool {
	const WIDTH = 880
	captcha.canvas = gg.NewContext(WIDTH, captcha.src.Height())

	rect := image.Rectangle{
		Min: image.Pt(math.MaxInt, math.MaxInt),
		Max: image.Pt(-1, -1),
	}

	success := false
	for x := 0; x < WIDTH; x++ {
		for y := 0; y < captcha.src.Height(); y++ {

			if !checked(captcha.src.Image().(*image.RGBA).At(x, y)) {
				continue
			}
			success = true

			rect.Min.X = utils.Min(x, rect.Min.X)
			rect.Min.Y = utils.Min(y, rect.Min.Y)
			rect.Max.X = utils.Max(x, rect.Max.X)
			rect.Max.Y = utils.Max(y, rect.Max.Y)

			c := captcha.src.Image().(*image.RGBA).At(x, y).(color.RGBA)
			captcha.canvas.SetColor(c)
			captcha.canvas.SetPixel(x, y)
		}
	}

	if !success {
		return false
	}

	croped := gg.NewContext(rect.Dx(), rect.Dy())
	x := 0
	for X := rect.Min.X; X != rect.Max.X; X++ {
		y := 0
		for Y := rect.Min.Y; Y != rect.Max.Y; Y++ {
			c := captcha.canvas.Image().(*image.RGBA).At(X, Y).(color.RGBA)
			croped.SetColor(c)
			croped.SetPixel(x, y)
			y++
		}
		x++
	}
	captcha.canvas = croped
	return success
}

func (a *area) hash(x, y int) int {
	result := x
	result <<= 16
	result |= y
	return result
}

func (a *area) fill(x, y int) {
	if !image.Pt(x, y).In(a.src.Image().Bounds()) {
		return
	}

	if !checked(a.src.Image().(*image.RGBA).At(x, y)) {
		return
	}

	index := a.hash(x, y)
	if _, ok := a.pix[index]; ok {
		return
	}
	a.pix[index] = image.Pt(x, y)

	a.rect.Min.X = utils.Min(x, a.rect.Min.X)
	a.rect.Min.Y = utils.Min(y, a.rect.Min.Y)
	a.rect.Max.X = utils.Max(x, a.rect.Max.X)
	a.rect.Max.Y = utils.Max(y, a.rect.Max.Y)

	a.fill(x-1, y-1)
	a.fill(x, y-1)
	a.fill(x+1, y-1)
	a.fill(x-1, y+1)
	a.fill(x-1, y+1)
	a.fill(x, y+1)
	a.fill(x+1, y+1)
}

func (a *area) right(y int) int {
	result := -1
	for _, p := range a.pix {
		if p.Y != y {
			continue
		}
		result = utils.Max(result, p.X)
	}
	return result
}

func (a *area) draw(dx, dy int) {
	a.canvas = gg.NewContext(a.rect.Dx(), a.rect.Dy())
	a.canvas.SetColor(color.RGBA{255, 255, 255, 255})
	for _, p := range a.pix {
		point := p.Sub(a.rect.Min)
		a.canvas.SetPixel(point.X, point.Y)
	}

	img := imaging.Resize(a.canvas.Image(), dx, dy, imaging.NearestNeighbor)
	a.canvas = gg.NewContextForImage(img)
	for x := 0; x < dx; x++ {
		for y := 0; y < dy; y++ {
			var data float64 = 0.0
			if checked(a.canvas.Image().(*image.RGBA).At(x, y)) {
				data = 1.0
			}
			a.data = append(a.data, data)
		}
	}
}

func (a *area) save(path string) error {
	return a.canvas.SavePNG(path)
}

func (captcha *CaptchaDecoder) allocate() bool {
	x := 0
	y := captcha.canvas.Height() / 2

	for x < captcha.canvas.Width() {
		if !checked(captcha.canvas.Image().(*image.RGBA).At(x, y)) {
			x++
			continue
		}

		a := newArea(captcha.canvas)
		a.fill(x, y)
		x = a.right(y) + 1

		if x == 0 {
			break
		}

		a.draw(24, 52)
		captcha.areas = append(captcha.areas, a)
	}
	return len(captcha.areas) == 9
}

func (captcha *CaptchaDecoder) Decode() error {
	if !captcha.crop() {
		return errors.New("captcha empty")
	}

	if !captcha.allocate() {
		return errors.New("allocation areas failed")
	}
	captcha.predict()
	return nil
}

func (captcha *CaptchaDecoder) SaveAreas(dir string) error {
	for _, a := range captcha.areas {
		name := fmt.Sprint(time.Now().UnixMicro()) + ".png"
		path := filepath.Join(dir, name)
		err := a.save(path)
		if err != nil {
			return err
		}
	}
	return nil
}

func (captcha *CaptchaDecoder) Rects() []image.Rectangle {
	var result []image.Rectangle
	for _, a := range captcha.areas {
		result = append(result, a.rect)
	}
	return result
}

func (captcha *CaptchaDecoder) Size() image.Rectangle {
	return captcha.canvas.Image().Bounds()
}

func (captcha *CaptchaDecoder) Save(fname string) error {
	os.Mkdir("captcha", os.ModePerm)
	fname = filepath.Join(".", "captcha", fname+".png")
	return captcha.canvas.SavePNG(fname)
}

func (captcha *CaptchaDecoder) SaveDebug(fname string) {
	os.Mkdir("debug", os.ModePerm)
	fname = filepath.Join(".", "debug", fname)

	if captcha.src != nil {
		captcha.src.SavePNG(fname + ".png")
	}
	if captcha.canvas != nil {
		captcha.canvas.SavePNG(fname + ".crop.png")
	}
	for i, a := range captcha.areas {
		path := fmt.Sprintf("%s.[%d].png", fname, i)
		a.save(path)
	}
}

func (captcha *CaptchaDecoder) predict() {
	var data []tensor.Vector

	for i, it := range captcha.areas {
		if i == 4 {
			continue
		}
		data = append(data, it.data)
	}
	result := nn.Predict(&data)

	var s string
	for _, it := range *result {
		i, _ := it.Max()
		s += string(dict[i])
	}
	captcha.codes = []string{s[:4], s[4:]}
	captcha.valid = true
}

func (captcha *CaptchaDecoder) Codes() []string {
	return captcha.codes
}
