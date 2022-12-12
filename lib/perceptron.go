package perceptron_go

import (
	"math/rand"
)

const G_learning_rate = 0.05

type FLOAT float32

type Perceptron struct {
	Name       string   // Perceptron shiny Name.
	Inputs     []*FLOAT // Array of pointers(!) Because they will refer to previous perceptron 'result'
	InputCount int
	Weights    []FLOAT
	Bias       FLOAT
	Error      FLOAT
	Result     FLOAT
	//ActivationFunc func(FLOAT) FLOAT //todo: initialize
}

func NewPerceptron(name string, inputCount int) *Perceptron {
	p := new(Perceptron)
	p.Name = name
	p.InputCount = inputCount
	p.Inputs = make([]*FLOAT, inputCount)
	p.Weights = make([]FLOAT, inputCount)
	p.Result = 0

	for i := 0; i < inputCount; i++ {
		p.Weights[i] = FLOAT(rand.Float32()) - 0.5
	}
	p.Bias = FLOAT(rand.Float32()) - 0.5

	return p
}

func (this *Perceptron) CalculateAndUpdateResult() {
	if this.InputCount == 0 {
		PfYellow("[Name=%s] Warning: unable to calculate result. No inputs. Skipping.\n", this.Name)
		panic(111)
		return
	}
	this.Result = this._calculateResult()
}

func (this *Perceptron) _calculateResult() FLOAT {
	if this.InputCount == 0 {
		PfRed("Error: unable to calculate result. No inputs. Returning zero.\n")
		return 0
	}
	sum := FLOAT(0)

	for i := 0; i < this.InputCount; i++ {
		value := *this.Inputs[i] * this.Weights[i]
		sum += value
	}
	sum += this.Bias

	return this.ActivationFunc(sum)
}

func (this *Perceptron) UpdateWeights() {
	Pf(C_BLUE+C_BOLD+"%-7s"+C_RST, this.Name)
	Pf(C_MAGENTA+" update_weights: error: "+C_RST+C_RED+"%+05.2f"+C_RST+" weights: ", this.Error)

	// weight correction formula
	correct_weight := func(rate, old_weight, err, result, input FLOAT) FLOAT {
		return old_weight + rate*err*this.DerivativeFunc(result)*input
	}

	// update weights
	for i := 0; i < this.InputCount; i++ {
		old_weight := this.Weights[i]
		this.Weights[i] = correct_weight(G_learning_rate, this.Weights[i], this.Error, this.Result, *this.Inputs[i])
		Pf("w%+5.3f", this.Weights[i])
		Pf(C_CYAN+"Δ%+.0fm "+C_RST, (this.Weights[i]-old_weight)*1000)
	}

	// correct bias same way as other weights have been corrected
	old_bias := this.Bias
	this.Bias = correct_weight(G_learning_rate, this.Bias, this.Error, this.Result, 1)
	PfGray("b%+.2f", this.Bias)
	PfGray(C_CYAN2+"Δ%+.0fm "+C_RST, (this.Bias-old_bias)*1000)

	new_result := this._calculateResult()
	PfBlue("NR: %+.2fΔ%+0.fm \n", new_result, (this.Result-new_result)*1000)
}

func (this *Perceptron) ActivationFunc(x FLOAT) FLOAT {
	return GetCurrentActivationBundle().activation(x)
	//return FLOAT(1.0 / (1.0 + math.Exp(float64(-x))))
}

func (this *Perceptron) DerivativeFunc(y FLOAT) FLOAT {
	return GetCurrentActivationBundle().derivative(y)
	//return y * (1.0 - y)
}
