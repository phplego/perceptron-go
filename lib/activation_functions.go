package perceptron_go

import "math"

type ActivationBundle struct {
	Name       string
	activation func(FLOAT) FLOAT
	derivative func(FLOAT) FLOAT
}

var ActivationBundles = []ActivationBundle{
	ActivationBundle{ // Sigmoid / Logistic Function
		Name:       "Sigmoid",
		activation: func(x FLOAT) FLOAT { return FLOAT(1.0 / (1.0 + math.Exp(float64(-x)))) },
		derivative: func(y FLOAT) FLOAT { return y * (1.0 - y) },
	},
}

var CurrentActivationBundleIndex = 0

func GetCurrentActivationBundle() ActivationBundle {
	return ActivationBundles[CurrentActivationBundleIndex]
}
