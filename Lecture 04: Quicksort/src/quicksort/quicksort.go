package quicksort

import "math/rand"

func insertionSort(a []int) {
    n := len(a)
    if n <= 1 {
        return;
    }
    for i := 0; i < n - 1; i++ {
        j := i + 1
        key := a[j]
        for ;j > 0 && key < a[j-1]; j-- {
            a[j] = a[j - 1]
        }
        a[j] = key
    }
}

func min(a, b int) int {
    if  a < b {
        return a
    } else {
        return b
    }
}

func partition(a[] int) int {
    n := len(a)
    a[0], a[n/2] = a[n/2], a[0]
    i := 0
    x := a[0]
    for j := 1; j < n; j++ {
        if a[j] <= x {
            i++
            a[i], a[j] = a[j], a[i]
        }
    }
    a[0], a[i] = a[i], a[0]
    return i
}

func hpartitionInternal(k int, a[] int) (int, int) {
    i, j := 0, len(a) - 1
    for {
        for k < a[j] {
            j--
        }
        for a[i] < k {
            i++
        }
        if i < j {
            a[i], a[j] = a[j], a[i]
            i++
            j--
        } else {
            return i, j
        }
    }
}

func hpartition(a []int) (int, int) {
    return hpartitionInternal(a[len(a)/2], a)
}

func median(x, y, z int) int {
    if x > y {
        x, y = y, x
    }
    if y > z {
        y, z = z, y
    }
    if x > y {
        return x
    } else {
        return y;
    }
}

func mpartition(a []int) (int, int) {
    n := len(a)
    m := n / 2
    return hpartitionInternal(median(a[0], a[m], a[n-1]), a)
}

func rpartition(r *rand.Rand, a[] int) (int, int) {
    k := r.Intn(len(a))
    return hpartitionInternal(a[k], a)
}

func fpartitionInternal(k int, v []int) (int, int) {
	a, b := 0, 0
	c, d := len(v)-1, len(v)-1
	// v[0, a) == k
	// v[a, b) < k
	// v(d, n) == k
	// v(c, d] > k
	// b <= c
	for {
		for b <= c && v[b] <= k {
			if k == v[b] {
				v[a], v[b] = v[b], v[a]
				a++
			}
			b++
		}
		for b <= c && k <= v[c] {
			if k == v[c] {
				v[d], v[c] = v[c], v[d]
				d--
			}
			c--
		}
		if b > c {
			break
		}
		v[b], v[c] = v[c], v[b]
		b++
		c--
	}
	s := min(a, b-a)
    lo, hi := 0, c
	for i := 0; i < s; i++ {
		v[i], v[hi - i] = v[hi - i], v[i]
	}
	s = min(d-c, len(v)-d-1)
    lo, hi = b, len(v) - 1
	for i := 0; i < s; i++ {
		v[lo + i], v[hi - i] = v[hi - i], v[lo + i]
	}
	return b - a, len(v) - d + c
}

func fpartition(a []int) (int, int) {
    y := a[len(a) / 2]
    // n := len(a)
    // m := n / 2
    // y := median(a[0], a[m], a[n-1])
    return fpartitionInternal(y, a)
}

const sizeThreshold = 32
func Quicksort(a []int) {
    for len(a) > sizeThreshold {
        m := partition(a)
        if m > len(a) / 2 {
            Quicksort(a[m+1:])
            a = a[:m]
        } else {
            Quicksort(a[:m])
            a = a[m+1:]
        }
    }
    insertionSort(a)
}

func QuicksortH(a []int) {
    for len(a) > sizeThreshold {
        i, j := hpartition(a)
        if i > len(a) - j - 1 {
            QuicksortH(a[j+1:])
            a = a[:i]
        } else {
            QuicksortH(a[:i])
            a = a[j+1:]
        }
    }
    insertionSort(a)
}

func QuicksortM(a []int) {
    for len(a) > sizeThreshold {
        i, j := mpartition(a)
        if i > len(a) - j - 1 {
            QuicksortM(a[j+1:])
            a = a[:i]
        } else {
            QuicksortM(a[:i])
            a = a[j+1:]
        }
    }
    insertionSort(a)
}

func quicksortR(r *rand.Rand, a []int) {
    for len(a) > sizeThreshold {
        i, j := rpartition(r, a)
        if i > len(a) - j - 1 {
            quicksortR(r, a[j+1:])
            a = a[:i]
        } else {
            quicksortR(r, a[:i])
            a = a[j+1:]
        }
    }
    insertionSort(a)
}

func QuicksortR(a []int) {
    r := rand.New(rand.NewSource(1))
    quicksortR(r, a)
}

func QuicksortF(a []int) {
    for len(a) > sizeThreshold {
        i, j := fpartition(a)
        if i > len(a) - j {
            QuicksortF(a[j:])
            a = a[:i]
        } else {
            QuicksortF(a[:i])
            a = a[j:]
        }
    }
    insertionSort(a)
}
