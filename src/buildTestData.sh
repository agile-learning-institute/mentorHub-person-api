sequence=${1:-1}
./test.sh $sequence 2
((sequence++))
./buildTestData.sh $sequence
