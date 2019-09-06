package main

import (
    "encoding/json"
    "fmt"
    "github.com/alecthomas/kingpin"
    _ "github.com/alecthomas/kingpin"
    _ "github.com/gorilla/handlers"
    _ "github.com/shurcooL/vfsgen"
    "log"
    "net"
    "net/http"
    "strconv"
    "strings"
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

    ss := NewHTTPStaticServer(gcfg.Root)

    var hdlr http.Handler = ss

    http.Handle("/",hdlr)
    http.Handle("/-/assets", http.StripPrefix("/-/assets/", http.FileServer(Assets)))
    http.HandleFunc("/-/sysinfo", func(writer http.ResponseWriter, request *http.Request) {
        writer.Header().Set("Content-Type", "application/json")
        data, _ := json.Marshal(map[string]interface{}{
            "version": "1.0.0",
        })
        writer.Write(data)
    })

    if gcfg.Addr == "" {
        gcfg.Addr = fmt.Sprintf(":%d", gcfg.Port)
    }
    if !strings.Contains(gcfg.Addr, ":") {
        gcfg.Addr = ":" + gcfg.Addr
    }
    _, port, _ := net.SplitHostPort(gcfg.Addr)
    log.Printf("listening on %s, local address http://%s:%s\n", strconv.Quote(gcfg.Addr), getLocalIP(), port)

    var err error
    err = http.ListenAndServe(gcfg.Addr, nil)
    log.Fatal(err)
}