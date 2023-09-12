package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/cam-per/licob/picker"
	"github.com/therfoo/therfoo/layers/dense"
	"github.com/therfoo/therfoo/model"
	"github.com/therfoo/therfoo/optimizers/sgd"
	"github.com/therfoo/therfoo/tensor"
)

var (
	generator = new(trainer)

	nn = model.New(
		model.WithCategoricalAccuracy(),
		model.WithCrossEntropyLoss(),
		model.WithEpochs(1000),
		model.WithInputShape(tensor.Shape{24 * 52}),
		model.WithOptimizer(
			sgd.New(sgd.WithBatchSize(5), sgd.WithLearningRate(0.005)),
		),
		model.WithTrainingGenerator(generator),
		model.WithValidatingGenerator(generator),
		model.WithVerbosity(true),
	)
)

func init() {
	nn.Add(24, dense.New(dense.WithReLU()))
	nn.Add(16, dense.New(dense.WithSigmoid()))

	nn.Compile()
}

func main() {
	//fit()

	f, err := os.Open("1148764307262283816.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	captcha, err := picker.NewCaptchaDecoder(f)
	if err != nil {
		log.Fatal(err)
	}
	captcha.Decode()
}

func fit() {
	generator.fetch("data")
	nn.Fit()
	nn.Save(filepath.Join(".", "ocr.gob"))

	for _, it := range generator.items {
		predict := *nn.Predict(&[]tensor.Vector{it.in})
		for _, p := range predict[0] {
			fmt.Printf("%.4f ", p)
		}
		fmt.Println(" ", it.out)
	}
}

func predict(dir string) {
	var data []tensor.Vector
	filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
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

		data = append(data, d)
		return nil
	})

	nn.Load("ocr.gob")
	var s string
	for _, it := range data {
		predict := *nn.Predict(&[]tensor.Vector{it})
		for _, p := range predict[0] {
			fmt.Printf("%.4f ", p)
		}
		fmt.Println()

		i, _ := predict[0].Max()
		s += string(dict[i])
	}

	fmt.Println(s)
}
