go build

if [[ $? -eq 0 ]]
then
    ./perceptron-go $@
fi

