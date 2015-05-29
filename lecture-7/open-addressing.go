package main

import (
    "fmt"
    "math/rand"
)

const (
    present uint32 = 1
    deleted uint32 = 2
)

type hashEntry struct {
    key int64
    flags uint32
}

type IntSet struct {
    n   int          // number of elements in the set
    log uint         // logarithm of the size of hash table
    tab []hashEntry  // hash table
}

const (
    magic1 = 11400714819323198485
    magic2 =  5700357409661599242
)

func hash1(k int64, log uint) uint64 {
    return (uint64(k) * magic1) >> (64 - log)
}

func hash2(k int64, log uint) uint64 {
    h := (uint64(k) * magic2) >> (64 - log)
    return h | 1
}

func newIntSet() *IntSet {
    return &IntSet{n: 0, log: 2, tab: make([]hashEntry, 4)}
}

func (s *IntSet) probeForFind(k int64, h1, h2 uint64) *hashEntry {
    h := h1
    mask := uint64(len(s.tab)) - 1
    for {
        flags := s.tab[h].flags
        if (flags & deleted) == 0 && ((flags & present) == 0 || s.tab[h].key == k) {
            return &s.tab[h]
        }
        h = (h + h2) & mask
    }
}

func (s *IntSet) probeForInsert(k int64, h1, h2 uint64) *hashEntry {
    h := h1
    mask := uint64(len(s.tab)) - 1
    for {
        flags := s.tab[h].flags
        if (flags & deleted) != 0 || (flags & present) == 0 {
            return &s.tab[h]
        }
        h = (h + h2) & mask
    }
}

func (s *IntSet) rehash() {
    tt := IntSet{n: 0, log: s.log + 1, tab: make([]hashEntry, 2 * len(s.tab))}
    for i := range s.tab {
        if (s.tab[i].flags & present) != 0 {
            tt.Insert(s.tab[i].key)
        }
    }
    s.log = tt.log
    s.tab = tt.tab
}

func (s *IntSet) Find(k int64) bool {
    h1 := hash1(k, s.log)
    h2 := hash2(k, s.log)
    e := s.probeForFind(k, h1, h2)
    return (e.flags & present) != 0
}

func (s *IntSet) Insert(k int64) bool {
    h1 := hash1(k, s.log)
    h2 := hash2(k, s.log)
    e := s.probeForFind(k, h1, h2)
    if (e.flags & present) != 0 {
        return false
    }

    e = s.probeForInsert(k, h1, h2)
    e.key = k;
    e.flags &= ^deleted
    e.flags |= present
    s.n++
    if 4 * s.n > 3 * len(s.tab) {
        s.rehash()
    }
    return true
}

func (s *IntSet) Delete(k int64) bool {
    h1 := hash1(k, s.log)
    h2 := hash2(k, s.log)
    e := s.probeForFind(k, h1, h2)
    if (e.flags & present) == 0 {
        return false
    }

    e.flags &= ^present
    e.flags |= deleted
    s.n--
    return true
}

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

func main() {
    // Generate the values, randomly shuffled.
    const n = 100000000
    a := makeValues(n)

    // Insert all the things!
    s := newIntSet()
    for _, v := range a {
        s.Insert(v)
    }

    // Check all values are present.
    for _, v := range a {
        if !s.Find(v) {
            panic("key must be present")
        }
    }

    // Delete values at even positions.
    for i := range a {
        if (i & 1) == 0 {
            s.Delete(a[i])
        }
    }

    // Check values at even positions are absent and values at odd positions are present.
    for i,v := range a {
        if (i & 1) == 0 {
            if s.Find(v) {
                panic(fmt.Sprintf("key #%d must not be present", i))
            }
        } else {
            if !s.Find(v) {
                panic(fmt.Sprintf("key #%d must be present", i))
            }
        }
    }

    // Insert again the values at even positions.
    for i := range a {
        if (i & 1) == 0 {
            s.Insert(a[i])
        }
    }

    // Check all the values are present.
    for _, v := range a {
        if !s.Find(v) {
            panic("key must be present")
        }
    }
}
