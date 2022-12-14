package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	. "perceptron-go/lib"
	"runtime/pprof"
	"strconv"
	"time"
)

type Sample struct {
	data           [9]FLOAT
	expectedResult [2]FLOAT
}

const _0 = 0.1
const _1 = 0.9

var learn_data = []Sample{
	Sample{[9]FLOAT{
		_0, _0, _0,
		_0, _0, _0,
		_0, _0, _0}, [2]FLOAT{_0 /*vert*/, _0 /*horiz*/}},
	Sample{[9]FLOAT{
		_0, _0, _0,
		_0, _1, _0,
		_0, _0, _0}, [2]FLOAT{_0 /*vert*/, _0 /*horiz*/}},
	Sample{[9]FLOAT{
		_0, _0, _0,
		_1, _1, _1,
		_0, _0, _0}, [2]FLOAT{_0 /*vert*/, _1 /*horiz*/}},
	Sample{[9]FLOAT{
		_1, _1, _1,
		_0, _0, _0,
		_0, _0, _0}, [2]FLOAT{_0 /*vert*/, _1 /*horiz*/}},
	Sample{[9]FLOAT{
		_0, _0, _0,
		_0, _0, _0,
		_1, _1, _1}, [2]FLOAT{_0 /*vert*/, _1 /*horiz*/}},
	Sample{[9]FLOAT{
		_1, _0, _0,
		_1, _0, _0,
		_1, _0, _0}, [2]FLOAT{_1 /*vert*/, _0 /*horiz*/}},
	Sample{[9]FLOAT{
		_0, _1, _0,
		_0, _1, _0,
		_0, _1, _0}, [2]FLOAT{_1 /*vert*/, _0 /*horiz*/}},
	Sample{[9]FLOAT{
		_0, _0, _1,
		_0, _0, _1,
		_0, _0, _1}, [2]FLOAT{_1 /*vert*/, _0 /*horiz*/}},
	Sample{[9]FLOAT{
		_1, _1, _1,
		_0, _0, _1, // corner 1
		_0, _0, _1}, [2]FLOAT{_1 /*vert*/, _1 /*horiz*/}},
	Sample{[9]FLOAT{
		_1, _0, _0,
		_1, _0, _0, // corner 2
		_1, _1, _1}, [2]FLOAT{_1 /*vert*/, _1 /*horiz*/}},
	Sample{[9]FLOAT{
		_1, _1, _1,
		_1, _0, _0, // corner 3
		_1, _0, _0}, [2]FLOAT{_1 /*vert*/, _1 /*horiz*/}},
	Sample{[9]FLOAT{
		_0, _0, _1,
		_0, _0, _1, // corner 4
		_1, _1, _1}, [2]FLOAT{_1 /*vert*/, _1 /*horiz*/}},
	Sample{[9]FLOAT{
		_1, _1, _1,
		_0, _1, _0, // T
		_0, _1, _0}, [2]FLOAT{_1 /*vert*/, _1 /*horiz*/}},
	Sample{[9]FLOAT{
		_0, _1, _0,
		_0, _1, _0, // T upside
		_1, _1, _1}, [2]FLOAT{_1 /*vert*/, _1 /*horiz*/}},
	Sample{[9]FLOAT{
		_0, _0, _1,
		_1, _1, _1, // T right
		_0, _0, _1}, [2]FLOAT{_1 /*vert*/, _1 /*horiz*/}},
	Sample{[9]FLOAT{
		_1, _0, _0,
		_1, _1, _1, // T left
		_1, _0, _0}, [2]FLOAT{_1 /*vert*/, _1 /*horiz*/}},
	Sample{[9]FLOAT{
		_0, _1, _0,
		_1, _1, _1, // cross
		_0, _1, _0}, [2]FLOAT{_1 /*vert*/, _1 /*horiz*/}},
	Sample{[9]FLOAT{
		_0, _1, _0,
		_0, _1, _0, // short v-line
		_0, _0, _0}, [2]FLOAT{_1 /*vert*/, _0 /*horiz*/}},
	Sample{[9]FLOAT{
		_0, _0, _0,
		_1, _1, _0, // short h-line
		_0, _0, _0}, [2]FLOAT{_0 /*vert*/, _1 /*horiz*/}},
	Sample{[9]FLOAT{
		_0, _1, _1,
		_0, _0, _0, // short h-line
		_0, _0, _0}, [2]FLOAT{_0 /*vert*/, _1 /*horiz*/}},
}

func trinary(condition bool, str1, str2 string) string {
	if condition {
		return str1
	}
	return str2
}

func print_results(vert_value FLOAT, horz_value FLOAT) {
	Pf("  Results:  %.2f %.2f  ", vert_value, horz_value)
	PfBold(trinary(vert_value > 0.5, "??? vertical   ", ""))
	PfBold(trinary(horz_value > 0.5, "??? horizontal ", ""))
	Pf("\n")
}

var gCPUprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {

	flag.Parse()
	if *gCPUprofile != "" {
		f, err := os.Create(*gCPUprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	Pf("Usage: %s [activation-function] [seed] [epoches] \n", os.Args[0])
	start_time := time.Now()
	seed := time.Now().Unix()
	total_epoches := 3000
	train_count := 0

	if len(os.Args) > 1 {
		Pf("Using activation function: %s\n", os.Args[1])
		index, _ := strconv.Atoi(os.Args[1])
		CurrentActivationBundleIndex = index
	}
	if len(os.Args) > 2 {
		Pf("Using seed: %s\n", os.Args[2])
		intSeed, _ := strconv.Atoi(os.Args[2])
		seed = int64(intSeed)
	}
	if len(os.Args) > 3 {
		total_epoches, _ = strconv.Atoi(os.Args[3])
		Pf("Using total_epoches: %d\n", total_epoches)
	}

	// initialize random generator with seed
	rand.Seed(seed)

	net := NewNetwork("net")
	net.CreateLayer("input", 9) // 9 pixels input
	net.CreateLayer("hidd1", 5)
	net.CreateLayer("hidd2", 3)
	net.CreateLayer("out", 2) // two neurons at the output

	if _, err := os.Stat("weights.save"); err == nil {
		err := net.LoadWeights("weights.save")
		if err != nil {
			panic(err)
		}
	}

	file_errors_by_sample_, _ := os.OpenFile("plot1.data", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o664)
	file_errors_by_sample := bufio.NewWriter(file_errors_by_sample_)
	file_errors_summary_, _ := os.OpenFile("plot2.data", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o664)
	file_errors_summary := bufio.NewWriter(file_errors_summary_)
	defer file_errors_by_sample_.Close()
	defer file_errors_summary_.Close()

	// learn cycle
	for epoch := 0; epoch < total_epoches; epoch++ {
		epoch_out_err_max := 0.0

		for sample := 0; sample < len(learn_data); sample++ {
			//PRINT_ON = epoch > total_epoches - 10;    // print only last 10 epoches
			PRINT_ON = (epoch+1)%1000 == 0 // print every 1000-th epoch

			if PRINT_ON {
				PfGreen("\nepoch %d  sample #%d\n", epoch, sample)
				for k := 0; k < 9; k++ { // print sample square
					if (k)%3 == 0 {
						Pf("\n")
					}
					if learn_data[sample].data[k] > 0.5 {
						Pf("%s ", "???")
					} else {
						Pf("%s ", "???")
					}
				}
			}

			for k := 0; k < 9; k++ { // fill input layer with sample
				net.SetInputValue(k, learn_data[sample].data[k])
			}

			// provide signal through the network
			net.Forward()

			print_results(net.OutLayer().Perceptrons[0].Result, net.OutLayer().Perceptrons[1].Result)
			net.PrintState()

			// lear the sample
			net.Learn(learn_data[sample].expectedResult[:])
			train_count++

			if PRINT_ON {
				Pf("error sum: "+C_RED+"%+.3f  "+C_RST+" outerr:"+C_YELLOW+" %.3f"+C_RST+"\n", net.ErrorSum(), net.OutLayer().ErrorSum())
			}

			// save charts
			fmt.Fprintf(file_errors_summary, "%f %f\n", net.ErrorSum(), net.OutLayer().ErrorSum())
			fmt.Fprintf(file_errors_by_sample, "%f ", net.OutLayer().ErrorAbsSum())
			epoch_out_err_max = math.Max(float64(epoch_out_err_max), float64(net.OutLayer().ErrorAbsSum()))
		}
		fmt.Fprintf(file_errors_by_sample, "\n")

		if PRINT_ON {
			Pf("epoch_out_err_max: "+C_BG_BLUE+" %.3f "+C_RST+" seed: %d\n", epoch_out_err_max, seed)
		}
	}

	err := net.SaveWeights("weights.save")
	if err != nil {
		panic(err)
	}

	PRINT_ON = true

	speed := float64(train_count) / time.Now().Sub(start_time).Seconds() / 1000
	speed = math.Round(speed)
	Pf("Samples count: %d\n", len(learn_data))
	Pf("Used Activation function: "+C_BG_RED+" %s "+C_RST, GetCurrentActivationBundle().Name)
	Pf("  LR: "+C_BG_YELL+" %g "+C_RST, G_learning_rate)
	Pf("  speed: "+C_BG_MAG+" %g "+C_RST+" train/ms\n", speed)
	Pf("train count: %d\n", train_count)
	Pf("exec time: %g sec\n", time.Now().Sub(start_time).Seconds())
}
