package main

import (
    "github.com/alecthomas/kingpin"
    _ "github.com/alecthomas/kingpin"
    _ "github.com/gorilla/handlers"
    _ "github.com/shurcooL/vfsgen"
    "log"
)
type Configure struct {
    Addr      string     `yaml:"addr"`
    Port      int        `yam:"port"`
    Root      string     `yaml:"root"`
}
var (
    gcfg = Configure{}
)

func parseFlags() error {
    // init default value
    gcfg.Root = "./"
    gcfg.Port = 8000
    gcfg.Addr = ""

    // get config from command line flags
    kingpin.HelpFlag.Short('h')
    kingpin.Flag("root","root directory").Short('r').StringVar(&gcfg.Root)
    kingpin.Flag("port","listen port").Short('p').IntVar(&gcfg.Port)
    kingpin.Flag("addr","listen addr").Short('a').StringVar(&gcfg.Addr)

    kingpin.Parse()

    return nil
}

func main() {
    if err := parseFlags(); err != nil {
        log.Fatal(err)
    }



}









}


