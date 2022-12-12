go build

if [[ $? -eq 0 ]]
then
    ./perceptron-go $@
fi

if [[ $? -eq 0 ]]
then
    ./plot1.sh
fi




