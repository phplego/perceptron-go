package perceptron_go

type Network struct {
	Name   string
	Layers []*Layer
}

func NewNetwork(name string) *Network {
	network := new(Network)
	network.Name = name
	return network
}

// CreateLayer - Create new network layer
func (this *Network) CreateLayer(name string, size int) {
	var previous *Layer = nil
	if len(this.Layers) > 0 {
		previous = this.Layers[len(this.Layers)-1]
	}
	layer := NewLayer(name, size, previous)
	this.Layers = append(this.Layers, layer)
}

func (this *Network) InputLayer() *Layer {
	if len(this.Layers) == 0 {
		return nil
	}
	return this.Layers[0]
}

func (this *Network) OutLayer() *Layer {
	if len(this.Layers) == 0 {
		return nil
	}
	return this.Layers[len(this.Layers)-1]
}

// Forward - Calculate results for each layer. From the input to the output
func (this *Network) Forward() {
	for i := 1; /*skip input layer*/ i < len(this.Layers); i++ {
		this.Layers[i].CalculateResult()
	}
}

// ErrorSum - Sum of all the layers error
func (this *Network) ErrorSum() FLOAT {
	var summary FLOAT = 0
	for i := 1; /* skip input layer */ i < len(this.Layers); i++ {
		summary += this.Layers[i].ErrorSum()
	}
	return summary
}

// Learn - Teach the network with target values.
// This method finds errors for each layer from
// the output to the input and then update weights.
func (this *Network) Learn(data []FLOAT) {
	// validations
	if this.OutLayer() == nil {
		PfRed("Error: outLayer is null.\n")
		return
	}
	if len(data) != this.OutLayer().Size {
		PfRed("Error: Data size (%d) doesn't match out layer size (%d).\n", len(data), this.OutLayer().Size)
		return
	}

	// set learn data to output layer: set errors as (target - result)
	for i := 0; i < this.OutLayer().Size; i++ {
		this.OutLayer().Perceptrons[i].Error = data[i] - this.OutLayer().Perceptrons[i].Result
	}

	// for each layer (in reverse order) find errors
	for i := len(this.Layers) - 2; /*skip out layer*/ i > 0; i-- {
		this.Layers[i].FindErrors()
	}

	// for each layer (in reverse order) update weights
	for i := len(this.Layers) - 1; i > 0; i-- {
		this.Layers[i].UpdateWeights()
	}
}

// PrintState - Print network state to the console
func (this *Network) PrintState() {
	for i := 0; i < len(this.Layers); i++ {
		layer := this.Layers[i]
		Pf("%-10s", layer.Name)
		for p := 0; p < layer.Size; p++ {
			if layer.Perceptrons[p].Result > 0.75 {
				PfGreen("%.2f ", layer.Perceptrons[p].Result)
			} else if layer.Perceptrons[p].Result > 0.5 {
				PfBlue("%.2f ", layer.Perceptrons[p].Result)
			} else if layer.Perceptrons[p].Result > 0.25 {
				PfYellow("%.2f ", layer.Perceptrons[p].Result)
			} else {
				PfRed("%.2f ", layer.Perceptrons[p].Result)
			}
		}
		Pf("\n")
	}
}

func (this *Network) SetInputValue(index int, value FLOAT) {
	if this.InputLayer() == nil {
		PfRed("Error: inputLayer is null.\n")
		return
	}

	if index < 0 || index >= this.InputLayer().Size {
		PfRed("Error: setInputValue: Out of bounds (%d).\n", index)
		return
	}

	this.InputLayer().Perceptrons[index].Result = value
}
