package perceptron_go

import (
	"math/rand"
	. "perceptron-go/lib"
	"testing"
)

func Test8(ot *testing.T) {
	Pf("Random sequence for seed 0...\n")

	rand.Seed(0)
	for i := 0; i < 10; i++ {
		Pf("val %d\n", rand.Uint32())
	}

}
