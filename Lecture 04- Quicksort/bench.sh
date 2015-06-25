#! /bin/sh

go test --gccgoflags="-Ofast" --bench=$1 quicksort | awk '/^Benchmark/ { printf "%-40s\t% 15d\n", $1, $3; }'
