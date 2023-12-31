package model

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/therfoo/therfoo/metrics"
	"github.com/therfoo/therfoo/optimizers"
	"github.com/therfoo/therfoo/tensor"
)

type Model struct {
	accurate     func(yTrue, yEstimate *tensor.Vector) bool
	lossFunction func(yTrue, yEstimate *tensor.Vector) *tensor.Vector
	lossPrime    func(yTrue, yEstimate *tensor.Vector) *tensor.Vector
	generator    struct {
		testing    Generator
		training   Generator
		validating Generator
	}
	inputShape           tensor.Shape
	layers               []Layer
	layersCount          int
	learningRate         float64
	metricsConsumers     []metrics.Consumer
	neuronsCount         []int
	optimizer            optimizers.Optimizer
	epochs               int
	skipOutputDerivative bool
}

func (m *Model) activate(activation *tensor.Vector) *[]tensor.Vector {
	n := m.layersCount + 1
	activations := make([]tensor.Vector, n, n)

	activations[0] = *activation

	for l := 0; l < m.layersCount; l++ {
		activations[l+1] = *m.layers[l].Activate(&(activations[l]))
	}

	return &activations
}

func (m *Model) Add(neuronsCount int, layer Layer) {
	m.neuronsCount = append(m.neuronsCount, neuronsCount)
	m.layers = append(m.layers, layer)
	m.layersCount++
}

func (m *Model) Compile() error {
	neurons := make([][]int, m.layersCount, m.layersCount)

	for l := range m.layers {
		var weightsCount int
		if l == 0 {
			weightsCount = m.inputShape.Size()
		} else {
			weightsCount = m.neuronsCount[l-1]
		}
		m.layers[l].Init(m.neuronsCount[l], weightsCount)
		neurons[l] = []int{m.neuronsCount[l], weightsCount}
	}

	m.optimizer.Init(&neurons)

	return nil
}

func (m *Model) Fit() {

	for epoch := 0; epoch < m.epochs; epoch++ {
		rand.Seed(time.Now().Unix())
		for _, batch := range rand.Perm(m.generator.training.Len()) {
			xBatch, yBatch := m.generator.training.Get(batch)

			wg := sync.WaitGroup{}

			batchSize := len(*xBatch)

			wg.Add(batchSize)

			for i := range *xBatch {
				go func(s int) {
					x := (*xBatch)[s]
					m.minimize(&((*yBatch)[s]), m.activate(&x))
					wg.Done()
				}(i)
			}

			wg.Wait()

			m.optimize(m.optimizer.Optimizations())

		}

		m.validate(epoch)

	}
}

func (m *Model) Load(filename string) (err error) {
	f, err := os.Open(filename)
	var weights [][]byte
	gob.NewDecoder(f).Decode(&weights)
	for l := range m.layers {
		err = m.layers[l].Load(weights[l])
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Model) minimize(yTrue *tensor.Vector, activations *[]tensor.Vector) {

	yEstimate := (*activations)[m.layersCount]
	nablaLoss := m.lossPrime(yTrue, &yEstimate)

	if !m.skipOutputDerivative {
		m.layers[m.layersCount-1].Derive(&((*activations)[m.layersCount-1]), nablaLoss)
	}
	m.optimizer.Add(m.layersCount-1, &((*activations)[m.layersCount-1]), nablaLoss)

	ll := m.layersCount - 2
	for l := ll; l > -1; l-- {
		nablaLoss = m.layers[l+1].NextCost(nablaLoss)
		m.layers[l].Derive(&((*activations)[l+1]), nablaLoss)
		m.optimizer.Add(l, &((*activations)[l]), nablaLoss)
	}
}

func (m *Model) optimize(optimizations *[][][]float64) {
	for layer := range *optimizations {
		m.layers[layer].Adjust(&((*optimizations)[layer]))
	}
}

func (m *Model) Predict(xBatch *[]tensor.Vector) *[]tensor.Vector {
	predictions := make([]tensor.Vector, len(*xBatch), len(*xBatch))
	for i := range *xBatch {
		predictions[i] = (*m.activate(&(*xBatch)[i]))[m.layersCount]
	}
	return &predictions
}

func (m *Model) Save(filename string) (err error) {
	weights := make([][]byte, m.layersCount, m.layersCount)
	for l := range m.layers {
		weights[l], err = m.layers[l].Bytes()
		if err != nil {
			return err
		}
	}
	var b bytes.Buffer
	err = gob.NewEncoder(&b).Encode(weights)
	if err == nil {
		ioutil.WriteFile(filename, b.Bytes(), 0755)
	}
	return err
}

func (m *Model) validate(epoch int) {
	rand.Seed(time.Now().Unix())
	n := m.generator.validating.Len()
	mm := metrics.Metrics{Epoch: epoch}
	accurate, total := 0., 0.
	for _, batch := range rand.Perm(n) {
		xBatch, yBatch := m.generator.training.Get(batch)
		yBatchEstimate := m.Predict(xBatch)
		cost := 0.
		for i := range *yBatch {
			cost += m.lossFunction(&(*yBatch)[i], &(*yBatchEstimate)[i]).Sum()
			total++
			if m.accurate(&(*yBatch)[i], &(*yBatchEstimate)[i]) {
				accurate++
			}
		}
		mm.Cost += cost / float64(len(*yBatch))
	}
	mm.Accuracy, mm.Cost = accurate/total, mm.Cost/float64(n)
	for _, consume := range m.metricsConsumers {
		consume(&mm)
	}
}

func New(options ...Option) *Model {
	m := Model{}
	for i := range options {
		options[i](&m)
	}
	return &m
}
