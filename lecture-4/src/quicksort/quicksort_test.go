package quicksort

import (
    "math/rand"
    "sort"
    "testing"
)

const array_size = 1000000
const key_range = 500000000

var (
    workArray []int = nil
    randomArray []int = nil
    sortedRangesArray10 []int = nil
    sortedRangesArray100 []int = nil
    sortedRangesArray1000 []int = nil
    mostlySortedArray10 []int = nil
    mostlySortedArray100 []int = nil
    mostlySortedArray1000 []int = nil
    _01Array[]int = nil
)

func makeRandomArray(k int) []int {
    a := make([]int, array_size)
    for i := 0; i < array_size; i++ {
        a[i] = rand.Intn(k)
    }
    return a
}

func makeSortedRangesArray(k, p int) []int {
    a := makeRandomArray(k)
    m := array_size / p
    i := 0
    for i + m < array_size {
        QuicksortH(a[i:i+m])
        i += m
    }
    QuicksortH(a[i:])
    return a
}

func randomShuffle(a []int) {
    for i := len(a); i > 1; i-- {
        j := rand.Intn(i)
        a[i - 1], a[j] = a[j], a[i - 1]
    }
}

func makeMostlySortedArray(k, p int) []int {
    a := makeRandomArray(k)
    QuicksortH(a)
    m := array_size / p
    i := 0
    for i + m < array_size {
        randomShuffle(a[i:i+m])
        i += m
    }
    randomShuffle(a[i:])
    return a

}

func initRandomArray() {
    if randomArray == nil {
        randomArray = makeRandomArray(key_range)
    }
    if workArray == nil {
        workArray = make([]int, array_size)
    }
}

func initRandom01Array() {
    if _01Array == nil {
        _01Array = makeRandomArray(1)
    }
    if workArray == nil {
        workArray = make([]int, array_size)
    }
}

func initSortedRangesArray10() {
    if sortedRangesArray10 == nil {
        sortedRangesArray10 =  makeSortedRangesArray(key_range, 10)
    }
    if workArray == nil {
        workArray = make([]int, array_size)
    }
}

func initSortedRangesArray100() {
    if sortedRangesArray100 == nil {
        sortedRangesArray100 =  makeSortedRangesArray(key_range, 100)
    }
    if workArray == nil {
        workArray = make([]int, array_size)
    }
}

func initSortedRangesArray1000() {
    if sortedRangesArray1000 == nil {
        sortedRangesArray1000 =  makeSortedRangesArray(key_range, 1000)
    }
    if workArray == nil {
        workArray = make([]int, array_size)
    }
}

func initMostlySortedArray10() {
    if mostlySortedArray10 == nil {
        mostlySortedArray10 =  makeMostlySortedArray(key_range, 10)
    }
    if workArray == nil {
        workArray = make([]int, array_size)
    }
}

func initMostlySortedArray100() {
    if mostlySortedArray100 == nil {
        mostlySortedArray100 =  makeMostlySortedArray(key_range, 100)
    }
    if workArray == nil {
        workArray = make([]int, array_size)
    }
}

func initMostlySortedArray1000() {
    if mostlySortedArray1000 == nil {
        mostlySortedArray1000 =  makeMostlySortedArray(key_range, 1000)
    }
    if workArray == nil {
        workArray = make([]int, array_size)
    }
}

func checkSorted(a []int) bool {
    n := len(a)
    for i := 1; i < n; i++ {
        if a[i] < a[i-1] {
            return false
        }
    }
    return true
}

func testInternal(a []int, sort func([]int), t *testing.T) {
    copy(workArray, a)
    sort(workArray)
    if !checkSorted(workArray) {
        t.Error("Array not sorted")
    }
}

func Test_QS(t *testing.T) {
    initRandomArray()
    testInternal(randomArray, Quicksort, t)
}

func Test_QSH(t *testing.T) {
    initRandomArray()
    testInternal(randomArray, QuicksortH, t)
}

func Test_QSM(t *testing.T) {
    initRandomArray()
    testInternal(randomArray, QuicksortM, t)
}

func Test_QSR(t *testing.T) {
    initRandomArray()
    testInternal(randomArray, QuicksortR, t)
}

func Test_QSF(t *testing.T) {
    initRandomArray()
    testInternal(randomArray, QuicksortF, t)
}

func benchInternal(a []int, sort func([]int), b *testing.B) {
    for i := 0; i < b.N; i++ {
        copy(workArray, a)
        sort(workArray)
    }
}

// Random --------------------------------------------------------------------------------
func Benchmark_QS_Random(b *testing.B) {
    initRandomArray()
    benchInternal(randomArray, Quicksort, b)
}

func Benchmark_QSH_Random(b *testing.B) {
    initRandomArray()
    benchInternal(randomArray, QuicksortH, b)
}

func Benchmark_QSM_Random(b *testing.B) {
    initRandomArray()
    benchInternal(randomArray, QuicksortM, b)
}

func Benchmark_QSR_Random(b *testing.B) {
    initRandomArray()
    benchInternal(randomArray, QuicksortR, b)
}

func Benchmark_QSF_Random(b *testing.B) {
    initRandomArray()
    benchInternal(randomArray, QuicksortF, b)
}

func Benchmark_GO_Random(b *testing.B) {
    initRandomArray()
    for i := 0; i < b.N; i++ {
        copy(workArray, randomArray)
        sort.Sort(sort.IntSlice(workArray))
    }
}

// Random 01 -----------------------------------------------------------------------------
// func Benchmark_QS_Random01(b *testing.B) {
//     initRandom01Array()
//     benchInternal(_01Array, Quicksort, b)
// }

func Benchmark_QSH_Random01(b *testing.B) {
    initRandom01Array()
    benchInternal(_01Array, QuicksortH, b)
}

func Benchmark_QSM_Random01(b *testing.B) {
    initRandom01Array()
    benchInternal(_01Array, QuicksortM, b)
}

func Benchmark_QSR_Random01(b *testing.B) {
    initRandom01Array()
    benchInternal(_01Array, QuicksortR, b)
}

func Benchmark_QSF_Random01(b *testing.B) {
    initRandom01Array()
    benchInternal(_01Array, QuicksortF, b)
}

func Benchmark_GO_Random01(b *testing.B) {
    initRandom01Array()
    for i := 0; i < b.N; i++ {
        copy(workArray, _01Array)
        sort.Sort(sort.IntSlice(workArray))
    }
}

// 10 ------------------------------------------------------------------------------------
func Benchmark_QS_SortedRanges10(b *testing.B) {
    initSortedRangesArray10()
    benchInternal(sortedRangesArray10, Quicksort, b)
}

func Benchmark_QSH_SortedRanges10(b *testing.B) {
    initSortedRangesArray10()
    benchInternal(sortedRangesArray10, QuicksortH, b)
}

func Benchmark_QSM_SortedRanges10(b *testing.B) {
    initSortedRangesArray10()
    benchInternal(sortedRangesArray10, QuicksortM, b)
}

func Benchmark_QSR_SortedRanges10(b *testing.B) {
    initSortedRangesArray10()
    benchInternal(sortedRangesArray10, QuicksortR, b)
}

func Benchmark_QSF_SortedRanges10(b *testing.B) {
    initSortedRangesArray10()
    benchInternal(sortedRangesArray10, QuicksortF, b)
}

func Benchmark_GO_SortedRanges10(b *testing.B) {
    initSortedRangesArray10()
    for i := 0; i < b.N; i++ {
        copy(workArray, sortedRangesArray10)
        sort.Sort(sort.IntSlice(workArray))
    }
}

// 100 -----------------------------------------------------------------------------------
func Benchmark_QS_SortedRanges100(b *testing.B) {
    initSortedRangesArray100()
    benchInternal(sortedRangesArray100, Quicksort, b)
}

func Benchmark_QSH_SortedRanges100(b *testing.B) {
    initSortedRangesArray100()
    benchInternal(sortedRangesArray100, QuicksortH, b)
}

func Benchmark_QSM_SortedRanges100(b *testing.B) {
    initSortedRangesArray100()
    benchInternal(sortedRangesArray100, QuicksortM, b)
}

func Benchmark_QSR_SortedRanges100(b *testing.B) {
    initSortedRangesArray100()
    benchInternal(sortedRangesArray100, QuicksortR, b)
}

func Benchmark_QSF_SortedRanges100(b *testing.B) {
    initSortedRangesArray100()
    benchInternal(sortedRangesArray100, QuicksortF, b)
}

func Benchmark_GO_SortedRanges100(b *testing.B) {
    initSortedRangesArray100()
    for i := 0; i < b.N; i++ {
        copy(workArray, sortedRangesArray100)
        sort.Sort(sort.IntSlice(workArray))
    }
}

// 1000 ----------------------------------------------------------------------------------
func Benchmark_QS_SortedRanges1000(b *testing.B) {
    initSortedRangesArray1000()
    benchInternal(sortedRangesArray1000, Quicksort, b)
}

func Benchmark_QSH_SortedRanges1000(b *testing.B) {
    initSortedRangesArray1000()
    benchInternal(sortedRangesArray1000, QuicksortH, b)
}

func Benchmark_QSM_SortedRanges1000(b *testing.B) {
    initSortedRangesArray1000()
    benchInternal(sortedRangesArray1000, QuicksortM, b)
}

func Benchmark_QSR_SortedRanges1000(b *testing.B) {
    initSortedRangesArray1000()
    benchInternal(sortedRangesArray1000, QuicksortR, b)
}

func Benchmark_QSF_SortedRanges1000(b *testing.B) {
    initSortedRangesArray1000()
    benchInternal(sortedRangesArray1000, QuicksortF, b)
}

func Benchmark_GO_SortedRanges1000(b *testing.B) {
    initSortedRangesArray1000()
    for i := 0; i < b.N; i++ {
        copy(workArray, sortedRangesArray1000)
        sort.Sort(sort.IntSlice(workArray))
    }
}

// 10 ------------------------------------------------------------------------------------
func Benchmark_QS_MostlySorted10(b *testing.B) {
    initMostlySortedArray10()
    benchInternal(mostlySortedArray10, Quicksort, b)
}

func Benchmark_QSH_MostlySorted10(b *testing.B) {
    initMostlySortedArray10()
    benchInternal(mostlySortedArray10, QuicksortH, b)
}

func Benchmark_QSM_MostlySorted10(b *testing.B) {
    initMostlySortedArray10()
    benchInternal(mostlySortedArray10, QuicksortM, b)
}

func Benchmark_QSR_MostlySorted10(b *testing.B) {
    initMostlySortedArray10()
    benchInternal(mostlySortedArray10, QuicksortR, b)
}

func Benchmark_QSF_MostlySorted10(b *testing.B) {
    initMostlySortedArray10()
    benchInternal(mostlySortedArray10, QuicksortF, b)
}

func Benchmark_GO_MostlySorted10(b *testing.B) {
    initMostlySortedArray10()
    for i := 0; i < b.N; i++ {
        copy(workArray, mostlySortedArray10)
        sort.Sort(sort.IntSlice(workArray))
    }
}

// 100 -----------------------------------------------------------------------------------
func Benchmark_QS_MostlySorted100(b *testing.B) {
    initMostlySortedArray100()
    benchInternal(mostlySortedArray100, Quicksort, b)
}

func Benchmark_QSH_MostlySorted100(b *testing.B) {
    initMostlySortedArray100()
    benchInternal(mostlySortedArray100, QuicksortH, b)
}

func Benchmark_QSM_MostlySorted100(b *testing.B) {
    initMostlySortedArray100()
    benchInternal(mostlySortedArray100, QuicksortM, b)
}

func Benchmark_QSR_MostlySorted100(b *testing.B) {
    initMostlySortedArray100()
    benchInternal(mostlySortedArray100, QuicksortR, b)
}

func Benchmark_QSF_MostlySorted100(b *testing.B) {
    initMostlySortedArray100()
    benchInternal(mostlySortedArray100, QuicksortF, b)
}

func Benchmark_GO_MostlySorted100(b *testing.B) {
    initMostlySortedArray100()
    for i := 0; i < b.N; i++ {
        copy(workArray, mostlySortedArray100)
        sort.Sort(sort.IntSlice(workArray))
    }
}

// 1000 ----------------------------------------------------------------------------------
func Benchmark_QS_MostlySorted1000(b *testing.B) {
    initMostlySortedArray1000()
    benchInternal(mostlySortedArray1000, Quicksort, b)
}

func Benchmark_QSH_MostlySorted1000(b *testing.B) {
    initMostlySortedArray1000()
    benchInternal(mostlySortedArray1000, QuicksortH, b)
}

func Benchmark_QSM_MostlySorted1000(b *testing.B) {
    initMostlySortedArray1000()
    benchInternal(mostlySortedArray1000, QuicksortM, b)
}

func Benchmark_QSR_MostlySorted1000(b *testing.B) {
    initMostlySortedArray1000()
    benchInternal(mostlySortedArray1000, QuicksortR, b)
}

func Benchmark_QSF_MostlySorted1000(b *testing.B) {
    initMostlySortedArray1000()
    benchInternal(mostlySortedArray1000, QuicksortF, b)
}

func Benchmark_GO_MostlySorted1000(b *testing.B) {
    initMostlySortedArray1000()
    for i := 0; i < b.N; i++ {
        copy(workArray, mostlySortedArray1000)
        sort.Sort(sort.IntSlice(workArray))
    }
}
