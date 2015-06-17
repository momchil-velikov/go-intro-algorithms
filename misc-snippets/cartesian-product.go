package main

import "fmt"

func gen(in [][]string, out chan []string) {
	n := len(in)
	cnt := make([]int, n)
	for {
		tmp := make([]string, n)
		for i := 0; i < n; i++ {
			tmp[i] = in[i][cnt[i]]
		}

		out <- tmp

		i := 0
		for i < n && cnt[i]+1 == len(in[i]) {
			cnt[i] = 0
			i++
		}
		if i >= n {
			break
		}
		cnt[i]++
	}
}

func permute(in [][]string) chan []string {
	out := make(chan []string)
	go func() {
		gen(in, out)
		close(out)
	}()
	return out
}

func main() {
	a := [][]string{{"A", "B", "C"}, {"D", "E"}, {"F"}}
	for b := range permute(a) {
		fmt.Println(b)
	}
}

