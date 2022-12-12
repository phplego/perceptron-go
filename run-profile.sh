go build
./perceptron-go 1 1 10000 --cpuprofile=profile.data
go tool pprof --top profile.data