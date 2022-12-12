package perceptron_go

import "math"

type Layer struct {
	Name        string
	Size        int
	Perceptrons []*Perceptron
	NextLayer   *Layer
	PrevLayer   *Layer
}

func NewLayer(name string, size int, previous *Layer) *Layer {
	layer := new(Layer)
	layer.Name = name
	layer.Size = size
	layer.Perceptrons = make([]*Perceptron, size)
	layer.PrevLayer = previous
	if previous != nil {
		// save next layer
		previous.NextLayer = layer

		// create perceptrons
		for i := 0; i < size; i++ {
			layer.Perceptrons[i] = NewPerceptron(name, previous.Size)

			// bind results of previous layer to all input of all perceptrons of this layer
			for k := 0; k < previous.Size; k++ {
				layer.Perceptrons[i].Inputs[k] = &previous.Perceptrons[k].Result
				//Pf("bind item %d input %d %f\n", i, k, *layer.Perceptrons[i].Inputs[k])
			}
		}
	} else {
		// input layer
		for i := 0; i < size; i++ {
			layer.Perceptrons[i] = NewPerceptron(name, 0 /*has no inputs*/)
		}
	}
	return layer
}

func (this *Layer) isOutput() bool {
	return this.NextLayer == nil
}

func (this *Layer) isInput() bool {
	return this.PrevLayer == nil
}

func (this *Layer) CalculateResult() {
	for i := 0; i < this.Size; i++ {
		this.Perceptrons[i].CalculateAndUpdateResult()
	}
}

func (this *Layer) FindErrors() {
	if this.NextLayer == nil {
		// last layer (out layer)
		// errors for last layer must be set manually (err = target - result)
		return
	}

	// for each this layer perceptron
	for i := 0; i < this.Size; i++ {
		var calculated_error FLOAT = 0

		// for each next layer perceptron
		for n := 0; n < this.NextLayer.Size; n++ {
			// find my weight
			var myweight FLOAT = this.NextLayer.Perceptrons[n].Weights[i /*!!!*/]

			// sum next layer (errors*weight)
			calculated_error += myweight * this.NextLayer.Perceptrons[n].Error
			//pf_green("findErrors: [%s] i=%d n=%d calculated=%f myweight=%f %s_Layer->perceptrons[%d]->error=%f\n", this->Name, i, n, calculated_error, myweight, this->nextLayer->Name, n, this->nextLayer->perceptrons[n]->error);
		}

		// set new error value
		this.Perceptrons[i].Error = calculated_error
	}
}

func (this *Layer) UpdateWeights() {
	for i := 0; i < this.Size; i++ {
		this.Perceptrons[i].UpdateWeights()
	}
}

func (this *Layer) ErrorSum() FLOAT {
	var sum FLOAT = 0
	for i := 0; i < this.Size; i++ {
		sum += this.Perceptrons[i].Error
	}
	return sum
}

func (this *Layer) ErrorAbsSum() FLOAT {
	var sum FLOAT = 0
	for i := 0; i < this.Size; i++ {
		sum += FLOAT(math.Abs(float64(this.Perceptrons[i].Error)))
	}
	return sum
}

func (this *Layer) ErrorSqrtSum() FLOAT {
	var sum FLOAT = 0
	for i := 0; i < this.Size; i++ {
		sum += FLOAT(math.Sqrt(float64(this.Perceptrons[i].Error)))
	}
	return sum
}
