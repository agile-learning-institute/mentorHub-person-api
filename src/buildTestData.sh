sequence=${1:-1}
./test.sh $sequence
((sequence++))
./buildTestData.sh $sequence
