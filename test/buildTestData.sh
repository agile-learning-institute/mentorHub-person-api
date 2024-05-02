sequence=${1:-1}
./test/test.sh $sequence 2
((sequence++))
./test/buildTestData.sh $sequence
