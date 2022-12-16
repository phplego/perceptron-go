package perceptron_go

import "math"

type ActivationBundle struct {
	Name       string
	Activation func(FLOAT) FLOAT
	Derivative func(FLOAT) FLOAT
}

var ActivationBundles = []ActivationBundle{
	ActivationBundle{ // Sigmoid / Logistic Function
		Name: "Sigmoid",
		Activation: func(x FLOAT) FLOAT {
			return FLOAT(1.0 / (1.0 + math.Exp(float64(-x))))
		},
		Derivative: func(y FLOAT) FLOAT { return y * (1.0 - y) },
	},
	ActivationBundle{ // Tanh(x)
		Name: "Tanh(x)",
		Activation: func(x FLOAT) FLOAT {
			return FLOAT(1.0 / (1.0 + math.Exp(float64(-2*x))))
		},
		Derivative: func(y FLOAT) FLOAT { return y * (1.0 - y*y) },
	},
	ActivationBundle{ // Leaky ReLU
		Name: "L-ReLU",
		Activation: func(x FLOAT) FLOAT {
			if x < 0.0 {
				return 0.01
			}
			return x
		},
		Derivative: func(y FLOAT) FLOAT {
			if y < 0.0 {
				return 0.01
			}
			return 1.0
		},
	},
	ActivationBundle{ // Leaky Capped ReLU
		Name: "LC-ReLU",
		Activation: func(x FLOAT) FLOAT {
			if x >= 0.0 {
				if x > 1.0 {
					return 1 + 0.01*(x-1)
				} else {
					return x
				}
			} else {
				return 0.01 * x
			}
		},
		Derivative: func(y FLOAT) FLOAT {
			if y < 0.0 || y > 1.0 {
				return 0.01
			} else {
				return 1.0
			}
		},
	},

	ActivationBundle{ // Experimental
		Name: "Experimental",
		Activation: func(x FLOAT) FLOAT {
			if x >= 0.0 {
				if x > 1.0 {
					return 1
				} else {
					return x
				}
			} else {
				return 0.01 * x
			}
		},
		Derivative: func(y FLOAT) FLOAT {
			if y < 0.0 || y > 1.0 {
				return 0.01
			} else {
				return 1.0
			}
		},
	},
}

var CurrentActivationBundleIndex = 0

func GetCurrentActivationBundle() ActivationBundle {
	return ActivationBundles[CurrentActivationBundleIndex]
}
