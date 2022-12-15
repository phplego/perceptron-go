package perceptron_go

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"unsafe"
)

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
	if !PRINT_ON {
		return
	}
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

// SaveWeights - Save network state (weights + biases) to the file
func (this *Network) SaveWeights(filename string) error {
	// calculate buffer size
	values_count := 0
	for l := 1; /*skip input layer*/ l < len(this.Layers); l++ {
		layer := this.Layers[l]
		for p := 0; p < len(layer.Perceptrons); p++ {
			perceptron := layer.Perceptrons[p]
			values_count += len(perceptron.Weights)
			values_count += 1 // and bias
		}
	}

	buf_size := values_count * int(unsafe.Sizeof(FLOAT(0)))

	//buf := make([]byte, buf_size)
	var buf bytes.Buffer

	for l := 1; /*skip input layer*/ l < len(this.Layers); l++ {
		layer := this.Layers[l]
		for p := 0; p < len(layer.Perceptrons); p++ {
			perceptron := layer.Perceptrons[p]

			for w := 0; w < len(perceptron.Weights); w++ { // save weights
				binary.Write(&buf, binary.LittleEndian, perceptron.Weights[w])
			}
			binary.Write(&buf, binary.LittleEndian, perceptron.Bias) // save bias
		}
	}

	if buf.Len() != buf_size {
		return errors.New("wrong size1")
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o664)
	defer file.Close()
	if err != nil {
		return err
	}
	file.Write(buf.Bytes())
	return nil
}

func (this *Network) LoadWeights(filename string) error {
	// calculate buffer size
	values_count := 0
	for l := 1; /*skip input layer*/ l < len(this.Layers); l++ {
		layer := this.Layers[l]
		for p := 0; p < len(layer.Perceptrons); p++ {
			perceptron := layer.Perceptrons[p]
			values_count += len(perceptron.Weights)
			values_count += 1 // and bias
		}
	}

	buf_size := values_count * int(unsafe.Sizeof(FLOAT(0)))

	file, err := os.OpenFile(filename, os.O_RDONLY, 0o664)
	if err != nil {
		return err
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	if fileInfo.Size() != int64(buf_size) {
		return fmt.Errorf("Size of network (%v) does not match file size (%v)", buf_size, fileInfo.Size())
	}

	for l := 1; /*skip input layer*/ l < len(this.Layers); l++ {
		layer := this.Layers[l]
		for p := 0; p < len(layer.Perceptrons); p++ {
			perceptron := layer.Perceptrons[p]

			for w := 0; w < len(perceptron.Weights); w++ { // save weights
				binary.Read(file, binary.LittleEndian, &perceptron.Weights[w])
			}
			binary.Read(file, binary.LittleEndian, &perceptron.Bias)
		}
	}

	return nil
}
