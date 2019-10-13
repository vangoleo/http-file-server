package main

import (
    "fmt"
    "strings"
    "testing"
)

func Test_tmp(t *testing.T) {

    s := "/-/files/hello"
    fmt.Println(strings.Replace(s,"/-/files/","",1))
    s = "/-/files/hello/world/foo.txt"
    fmt.Println(strings.Replace(s,"/-/files/","",1))
}
