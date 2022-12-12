package perceptron_go

import (
	"math/rand"
	. "perceptron-go/lib"
	"testing"
)

func Test1(ot *testing.T) {
	Pf("Starting teaching one neuron...\n")

	inputs := []FLOAT{10, -2, 1.5}

	// teach one neuron
	rand.Seed(1)
	p := NewPerceptron("test", 3)

	p.Inputs[0] = &inputs[0]
	p.Inputs[1] = &inputs[1]
	p.Inputs[2] = &inputs[2]

	for i := 0; i < 3; i++ {
		PfBlue("INPUTS:  i%d: %0.2f  ", i, *p.Inputs[i])
	}

	p.CalculateAndUpdateResult()

	target := FLOAT(0.5)
	p.Error = target - p.Result
	Pf("initial result: %f\n", p.Result)
	Pf("target: %f\n", target)

	step := 0
	for p.Error*p.Error > 0.00001 {
		step += 1
		Pf("%.2d ", step)
		p.UpdateWeights()
		p.CalculateAndUpdateResult()
		p.Error = target - p.Result
		if step > 100 {
			break
		}
	}

	Pf("errSq: %f err: %f\n", p.Error*p.Error, p.Error)
	Pf("New Result: %f\n", p.Result)

}
