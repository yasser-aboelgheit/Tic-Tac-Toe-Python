#!/bin/sh
set -o pipefail

usage() {

	echo "
 Run test cases for go project.
 -j         generate JUnit report.
 -t         expecting test tags you want to run.
 -r         create HTML report.
 -o         expecting the output file name.
 -n         expecting the number of parallel tests, defaults to 1
 -c       testCase name

 -h         print this help message.
 "
	exit
}

# flags
report=""
junit=""
test_flags=""
testcase=""
tests=${TESTS:="unit"}
output="${tests}_coverage.out"
parallel=1

while getopts j:r:o:t:f:n:h:c: flag; do
	case "$flag" in
	j) junit="true" ;;
	t)
		tests=${OPTARG}
		output="${tests}_coverage.out"
		;;
	r) report=${OPTARG} ;;
	o) output="$OPTARG" ;;
	f) test_flags=${OPTARG} ;;
	n) parallel=${OPTARG} ;;
	c) testcase=${OPTARG} ;;
	h) usage ;;
	*) usage ;;
	esac
done

echo "coverage=$output"

go env -w GONOSUMDB="gitlab.com/startupbuilder"

command="go test $test_flags -coverpkg=./... --tags=$tests -coverprofile=$output -p $parallel -v ./..."

if [ "$testcase" ]; then
	command="$command -run $testcase"
fi

echo "command: $command"
sh -c "$command" 2>&1 | tee "$output".txt
exitcode=$?

if [ "$report" ]; then
	echo "Reporting HTML $report"
	go tool cover -html="$output"
fi

if [ "$junit" ]; then
	echo "Junit reporting "
	go-junit-report -set-exit-code -subtest-mode=exclude-parents >report.xml <"$output".txt
fi

# sleep 3
go tool cover -func "$output" | grep "total:"

exit "${exitcode}"
