// +build !vfs
//go:generate go run assets_generate.go

package main

import "net/http"

var Assets http.FileSystem = http.Dir("assets")

