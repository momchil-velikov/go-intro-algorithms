package main

import (
    "fmt"
    "rbt"
)

func main() {
    a := []rbt.KeyType{26, 17, 14, 10, 7, 3, 12, 16, 15, 21, 19, 20, 23, 41, 30, 28, 38, 35, 39, 47 }
    t := rbt.New()
    for _, v := range a {
        t.Insert(v)
    }
    t.Delete(30)
    fmt.Println(t.WriteDot())
}
