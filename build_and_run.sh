go build

if [[ $? -eq 0 ]]
then
    ./perceptron-go $@
fi

if [[ $? -eq 0 ]]
then
    ./plot2.sh
fi




