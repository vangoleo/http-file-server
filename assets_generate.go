// +build ignore

package main

import (
	"github.com/shurcooL/vfsgen"
	_ "github.com/shurcooL/vfsgen"
	"log"
	"net/http"
)

func main() {
	//var fs http.FileSystem = http.Dir("assets")
	var fs http.FileSystem = http.Dir("vue-hello-world/dist")

	err := vfsgen.Generate(fs, vfsgen.Options{
		PackageName:     "main",
		BuildTags:       "vfs",
		VariableName:    "Assets",
	})

	if err != nil {
		log.Fatalln(err)
	}
}