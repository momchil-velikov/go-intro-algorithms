package openaddr

import (
    "math/rand"
    "testing"
)

func randomShuffle(a []int64) []int64 {
    for i := len(a); i > 1; i-- {
        j := rand.Intn(i)
        a[i - 1], a[j] = a[j], a[i - 1]
    }
    return a
}

func makeValues(n uint) []int64 {
    a := make([]int64, n)
    for i := uint(0); i < n; i++ {
        a[i] = int64(i)
    }
    return randomShuffle(a)
}

const _N = 1000000
var values []int64 = makeValues(_N)

func TestHashTable(t *testing.T) {
    // Insert all the things!
    s := newIntSet()
    for _, v := range values {
        s.Insert(v)
    }

    // Check all values are present.
    for i, v := range values {
        if !s.Find(v) {
            t.Errorf("key #%d must be present", i)
        }
    }

    // Delete values at even positions.
    for i := range values {
        if (i & 1) == 0 {
            s.Delete(values[i])
        }
    }

    // Check values at even positions are absent and values at odd positions are present.
    for i, v := range values {
        if (i & 1) == 0 {
            if s.Find(v) {
                t.Errorf("key #%d must not be present", i)
            }
        } else if !s.Find(v) {
            t.Errorf("key #%d must be present", i)
        }
    }

    // Delete values at even positions, again.
    for i := range values {
        if (i & 1) == 0 {
            s.Delete(values[i])
        }
    }

    // Check values at even positions are absent and values at odd positions are present,
    // again
    for i, v := range values {
        if (i & 1) == 0 {
            if s.Find(v) {
                t.Errorf("key #%d must not be present", i)
            }
        } else if !s.Find(v) {
            t.Errorf("key #%d must be present", i)
        }
    }

    // Insert again values at even positions.
    for i := range values {
        if (i & 1) == 0 {
            s.Insert(values[i])
        }
    }

    // Check all values are present.
    for i, v := range values {
        if !s.Find(v) {
            t.Errorf("key #%d must be present", i)
        }
    }

    // Insert again values at even positions, again
    for i := range values {
        if (i & 1) == 0 {
            s.Insert(values[i])
        }
    }

    // Check all values are present, again
    for i, v := range values {
        if !s.Find(v) {
            t.Errorf("key #%d must be present", i)
        }
    }
}

func benchmarkHashTable(b *testing.B) {
    // Insert all the things!
    s := newIntSet()
    for _, v := range values {
        s.Insert(v)
    }

    for _, v := range values {
        s.Find(v)
    }

    // Delete values at even positions.
    for i := range values {
        if (i & 1) == 0 {
            s.Delete(values[i])
        }
    }

    for _, v := range values {
        s.Find(v)
    }

    // Delete values at even positions, again.
    for i := range values {
        if (i & 1) == 0 {
            s.Delete(values[i])
        }
    }

    for _, v := range values {
        s.Find(v)
    }

    // Insert again values at even positions.
    for i := range values {
        if (i & 1) == 0 {
            s.Insert(values[i])
        }
    }

    // Check all values are present.
    for _, v := range values {
        s.Find(v)
    }

    // Insert again values at even positions, again
    for i := range values {
        if (i & 1) == 0 {
            s.Insert(values[i])
        }
    }

    // Check all values are present, again
    for _, v := range values {
        s.Find(v)
    }
}

func BenchmarkHashTable(b *testing.B) {
    for i := 0; i < b.N; i++ {
        benchmarkHashTable(b)
    }
}

func benchmarkMap(b *testing.B) {
    // Insert all the things!
    s := make(map[int64]bool)
    for _, v := range values {
        s[v] = true
    }

    for _, v := range values {
        _, _ = s[v]
    }

    // Delete values at even positions.
    for i := range values {
        if (i & 1) == 0 {
            delete(s, values[i])
        }
    }

    for _, v := range values {
        _, _ = s[v]
    }

    // Delete values at even positions, again.
    for i := range values {
        if (i & 1) == 0 {
            delete(s, values[i])
        }
    }

    for _, v := range values {
        _, _ = s[v]
    }

    // Insert again values at even positions.
    for i := range values {
        if (i & 1) == 0 {
            s[values[i]] = true
        }
    }

    for _, v := range values {
        _, _ = s[v]
    }

    // Insert again values at even positions, again
    for i := range values {
        if (i & 1) == 0 {
            s[values[i]] = true
        }
    }

    for _, v := range values {
        _, _ = s[v]
    }
}

func BenchmarkMap(b *testing.B) {
    for i := 0; i < b.N; i++ {
        benchmarkMap(b)
    }
}

