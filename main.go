package main

import (
    "github.com/alecthomas/kingpin"
    _ "github.com/alecthomas/kingpin"
    _ "github.com/gorilla/handlers"
    "github.com/gorilla/mux"
    _ "github.com/gorilla/mux"
    _ "github.com/shurcooL/vfsgen"
    "log"
    "net/http"
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

    //http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
    //    writer.Write([]byte("match /"))
    //})


    //r.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
    //   writer.Write([]byte("match /"))
    //})
    //r.HandleFunc("/foo", func(writer http.ResponseWriter, request *http.Request) {
    //   writer.Write([]byte("match /foo"))
    //})

    //r.HandleFunc("/bar", func(writer http.ResponseWriter, request *http.Request) {
    //   writer.Write([]byte("match /bar"))
    //})

    //r.HandleFunc("/-/user/{path:.*}", func(writer http.ResponseWriter, request *http.Request) {
    //    writer.Write([]byte("match /-/user/"))
    //})

    //r.HandleFunc("/{path:.*}", func(writer http.ResponseWriter, request *http.Request) {
    //    writer.Write([]byte("match /"))
    //})



    //http.Handle("/", r)



    if err := parseFlags(); err != nil {
       log.Fatal(err)
    }

    ss := NewHTTPStaticServer(gcfg.Root)

    var hdlr http.Handler = ss

    r := mux.NewRouter()
    r.Handle("/-/files/{path:.*}", hdlr)
    r.Handle("/{path:.*}",http.FileServer(Assets))
    http.ListenAndServe(":8000", r)

    ////http.Handle("/",hdlr)
    ////http.Handle("/",http.FileServer(Assets))
    ////http.Handle("/-/files",http.StripPrefix("/-/files/",hdlr))
    ////http.Handle("/-/assets", http.StripPrefix("/-/assets/", http.FileServer(Assets)))
    ////http.HandleFunc("/-/sysinfo", func(writer http.ResponseWriter, request *http.Request) {
    ////    writer.Header().Set("Content-Type", "application/json")
    ////    data, _ := json.Marshal(map[string]interface{}{
    ////        "version": "1.0.0",
    ////    })
    ////    writer.Write(data)
    ////})
    //
    ////http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
    ////    writer.Write([]byte("match /"))
    ////})
    ////
    ////http.HandleFunc("/-/files", func(writer http.ResponseWriter, request *http.Request) {
    ////    writer.Write([]byte("match /-/files"))
    ////})
    //
    //r := mux.NewRouter()
    //
    ////r.Handle("/", http.FileServer(Assets))
    //r.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
    //    writer.Write([]byte("match /"))
    //})
    //r.HandleFunc("/-/files", func(writer http.ResponseWriter, request *http.Request) {
    //    writer.Write([]byte("match /-/files"))
    //})
    //
    //http.Handle("/",r)
    //
    //if gcfg.Addr == "" {
    //    gcfg.Addr = fmt.Sprintf(":%d", gcfg.Port)
    //}
    //if !strings.Contains(gcfg.Addr, ":") {
    //    gcfg.Addr = ":" + gcfg.Addr
    //}
    //_, port, _ := net.SplitHostPort(gcfg.Addr)
    //log.Printf("listening on %s, local address http://%s:%s\n", strconv.Quote(gcfg.Addr), getLocalIP(), port)
    //
    //var err error
    //err = http.ListenAndServe(gcfg.Addr, nil)
    //log.Fatal(err)
}