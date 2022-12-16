package perceptron_go

import (
	"fmt"
	"os"
	. "perceptron-go/lib"
	"strconv"
	"testing"
)

func Test4(ot *testing.T) {
	PfGreen("Compare activation functions (chart)..\n")

	filename := "test_4.data"
	file, _ := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0o666)

	for i := 0; i < 100; i++ {
		value := -5 + float64(i)/10
		fmt.Fprintf(file, "%f ", value)

		for i := 0; i < len(ActivationBundles); i++ {
			//value = value + float64(i)*0.05 // overlap uglyfix
			fmt.Fprintf(file, "%f ", ActivationBundles[i].Activation(FLOAT(value)))
		}
		fmt.Fprintf(file, "\n")
	}
	file.Close()

	// build run gnuplot script
	plot_parts := ""
	for i := 0; i < len(ActivationBundles); i++ {
		if i > 0 {
			plot_parts += ",\n"
		}
		plot_parts += "'" + filename + "' using 1:" + strconv.Itoa(i+2) + " title '" + ActivationBundles[i].Name + "' with line linewidth 2"
	}
	cmd := "gnuplot -e \"set grid; set key left top; set yrange [-0.2:1.2];\n plot \n" + plot_parts + ";  pause mouse close;\""

	// save run gnuplot script
	ff, _ := os.OpenFile(filename+".sh", os.O_CREATE|os.O_WRONLY, 0o777)
	ff.Write([]byte(cmd))
	ff.Close()
}
