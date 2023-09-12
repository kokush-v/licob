package metrics

import (
	"fmt"
)

func Logger(m *Metrics) error {
	_, err := fmt.Printf(
		"Epoch: %d, Accuracy: %.16f Cost: %.14f\n",
		m.Epoch+1,
		m.Accuracy,
		m.Cost,
	)
	return err
}
