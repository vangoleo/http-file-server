package main

import (
    "encoding/json"
    "github.com/gorilla/mux"
    _ "github.com/gorilla/mux"
    "html/template"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "strconv"
    "strings"
)

type IndexFileItem struct {
    Path  string
    Info  os.FileInfo
}

type HTTPStaticServer struct {
    Root       string
    indexes    []IndexFileItem
    m          *mux.Router
}

func NewHTTPStaticServer(root string) *HTTPStaticServer {
    if root == "" {
        root = "./"
    }
    root = filepath.ToSlash(root)
    if !strings.HasSuffix(root, "/") {
        root = root + "/"
    }
    log.Printf("root path: %s\n", root)
    m := mux.NewRouter()
    s := &HTTPStaticServer{
        Root: root,
        m:    m,
    }

    //m.HandleFunc("/{path:.*}",s.hIndex).Methods("GET", "HEAD")
    //m.Handle("/",http.FileServer(Assets))
    //m.Handle("/",http.StripPrefix("/", http.FileServer(Assets)))
    m.HandleFunc("/-/files/{path:.*}",s.dirOrFile)
    return s
}

func (s *HTTPStaticServer) dirOrFile(w http.ResponseWriter, r *http.Request) {
    path := strings.Replace(r.URL.Path,"/-/files/","",1)
    path = filepath.Join(s.Root, path)

    if isDir(path){
        s.dir(w,r)
    }else {
        s.file(w,r)
    }
    //w.Write([]byte(testdata))
}

func (s *HTTPStaticServer) dir(w http.ResponseWriter, r *http.Request){

    requestPath := strings.Replace(r.URL.Path,"/-/files/","",1)
    localPath := filepath.Join(s.Root, requestPath)

    //requestPath := mux.Vars(r)["path"]
    //localPath := filepath.Join(s.Root, requestPath)

    // path string -> info os.FileInfo
    fileInfoMap := make(map[string]os.FileInfo,0)

    infos, err := ioutil.ReadDir(localPath)
    if err != nil {
        http.Error(w,err.Error(), 500)
        return
    }
    for _, info := range infos {
        fileInfoMap[filepath.Join(requestPath,info.Name())] = info
    }

    // turn file list -> json
    lrs := make([]HTTPFileInfo2,0)
    for path, info := range fileInfoMap {
        lr := HTTPFileInfo2{
            Name:    info.Name(),
            Path:    path,
            ModTime: "2019-10-01 10:05:00",
        }
        if info.IsDir() {
            name := deepPath(localPath, info.Name())
            lr.Name = name
            lr.Path = filepath.Join(filepath.Dir(path), name)
            lr.Type = "dir"
            lr.Size = s.historyDirSize(lr.Path)
        } else {
            lr.Type = "file"
            lr.Size = info.Size()
        }
        lrs = append(lrs, lr)
    }

    data, _ := json.Marshal(lrs)
    w.Header().Set("Content-Type", "application/json")
    w.Write(data)
}

func (s *HTTPStaticServer) file(w http.ResponseWriter, r *http.Request){

    requestPath := strings.Replace(r.URL.Path,"/-/files/","",1)
    localPath := filepath.Join(s.Root, requestPath)

    w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(filepath.Base(requestPath)))
    http.ServeFile(w,r,localPath)

    //if r.FormValue("download") == "true" {
    //            w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(filepath.Base(path)))
    //        }
    //        http.ServeFile(w,r,relPath)

    // path string -> info os.FileInfo
    //fileInfoMap := make(map[string]os.FileInfo,0)

    //infos, err := ioutil.ReadDir(localPath)
    //if err != nil {
    //    http.Error(w,err.Error(), 500)
    //    return
    //}
    //for _, info := range infos {
    //    fileInfoMap[filepath.Join(requestPath,info.Name())] = info
    //}
    //
    //// turn file list -> json
    //lrs := make([]HTTPFileInfo2,0)
    //for path, info := range fileInfoMap {
    //    lr := HTTPFileInfo2{
    //        Name:    info.Name(),
    //        Path:    path,
    //        ModTime: "2019-10-01 10:05:00",
    //    }
    //    if info.IsDir() {
    //        name := deepPath(localPath, info.Name())
    //        lr.Name = name
    //        lr.Path = filepath.Join(filepath.Dir(path), name)
    //        lr.Type = "dir"
    //        lr.Size = s.historyDirSize(lr.Path)
    //    } else {
    //        lr.Type = "file"
    //        lr.Size = info.Size()
    //    }
    //    lrs = append(lrs, lr)
    //}
    //
    //data, _ := json.Marshal(lrs)
    //w.Header().Set("Content-Type", "application/json")
    //w.Write(data)
}

func (s *HTTPStaticServer) hIndex(w http.ResponseWriter, r *http.Request){
    path := mux.Vars(r)["path"]
    relPath := filepath.Join(s.Root, path)
    if r.FormValue("json") == "true" {
        s.hJSONList(w,r)
        return
    }

    if r.FormValue("raw") == "false" || isDir(relPath) {
        if r.Method == "HEAD" {
            return
        }
        renderHTML(w,"index.html",s)
    }else {
        if r.FormValue("download") == "true" {
            w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(filepath.Base(path)))
        }
        http.ServeFile(w,r,relPath)
    }
}

func (s *HTTPStaticServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    s.m.ServeHTTP(w, r)
}

func (s *HTTPStaticServer) hJSONList(w http.ResponseWriter, r *http.Request){
    requestPath := mux.Vars(r)["path"]
    localPath := filepath.Join(s.Root, requestPath)

    // path string -> info os.FileInfo
    fileInfoMap := make(map[string]os.FileInfo,0)

    infos, err := ioutil.ReadDir(localPath)
    if err != nil {
        http.Error(w,err.Error(), 500)
        return
    }
    for _, info := range infos {
        fileInfoMap[filepath.Join(requestPath,info.Name())] = info
    }

    // turn file list -> json
    lrs := make([]HTTPFileInfo,0)
    for path, info := range fileInfoMap {
        lr := HTTPFileInfo{
            Name:    info.Name(),
            Path:    path,
            ModTime: info.ModTime().UnixNano() / 1e6,
        }
        if info.IsDir() {
            name := deepPath(localPath, info.Name())
            lr.Name = name
            lr.Path = filepath.Join(filepath.Dir(path), name)
            lr.Type = "dir"
            lr.Size = s.historyDirSize(lr.Path)
        } else {
            lr.Type = "file"
            lr.Size = info.Size()
        }
        lrs = append(lrs, lr)
    }

    data, _ := json.Marshal(map[string]interface{}{
        "files": lrs,
    })
    w.Header().Set("Content-Type", "application/json")
    w.Write(data)
}

type HTTPFileInfo struct {
    Name      string  `json:"name"`
    Path      string  `json:"path"`
    Type      string  `json:"type"`
    Size      int64   `json:"size"`
    ModTime   int64   `json:"mtime"`
}

type HTTPFileInfo2 struct {
    Name      string `json:"name"`
    Path      string `json:"path"`
    Type      string `json:"type"`
    Size      int64  `json:"size"`
    ModTime   string `json:"modTime"`
    Actions   string `json:"actions"`
}

func deepPath(basedir, name string) string {
    isDir := true
    // loop max 5, incase of for loop not finished
    maxDepth := 5
    for depth := 0; depth <= maxDepth && isDir; depth += 1 {
        finfos, err := ioutil.ReadDir(filepath.Join(basedir, name))
        if err != nil || len(finfos) != 1 {
            break
        }
        if finfos[0].IsDir() {
            name = filepath.ToSlash(filepath.Join(name, finfos[0].Name()))
        } else {
            break
        }
    }
    return name
}

var dirSizeMap = make(map[string]int64)

var funcMap template.FuncMap

func (s *HTTPStaticServer) historyDirSize(dir string) int64 {
    var size int64
    if size, ok := dirSizeMap[dir]; ok {
        return size
    }
    for _, fitem := range s.indexes {
        if filepath.HasPrefix(fitem.Path, dir) {
            size += fitem.Info.Size()
        }
    }
    dirSizeMap[dir] = size
    return size
}

func isDir(path string) bool {
    info, err := os.Stat(path)
    return err == nil && info.Mode().IsDir()
}

func renderHTML(w http.ResponseWriter, name string, v interface{}){
    if _, ok := Assets.(http.Dir); ok {
        t := template.Must(template.New(name).Funcs(funcMap).Delims("[[","]]").Parse(assetsContent(name)))
        t.Execute(w, v)
    } else {
        executeTemplate(w,name,v)
    }
}

func assetsContent(name string) string {
    fd, err := Assets.Open(name)
    if err != nil {
        panic(err)
    }
    data, err := ioutil.ReadAll(fd)
    if err != nil {
        panic(err)
    }
    return string(data)
}

var _tmpls = make(map[string]*template.Template)
func executeTemplate(w http.ResponseWriter, name string, v interface{}){
    if t, ok := _tmpls[name]; ok {
        t.Execute(w,v)
        return
    }
    t := template.Must(template.New(name).Funcs(funcMap).Delims("[[", "]]").Parse(assetsContent(name)))
    _tmpls[name] = t
    t.Execute(w,v)
}

var testdata = `
[
    {
        "name": "hello.txt",
        "size": "10kb",
        "modTime": "2019-10-10 10:10:10",
        "actions": "Download, View"
    },
    {
        "name": "world.txt",
        "size": "15kb",
        "modTime": "2019-10-01 10:10:10",
        "actions": "Download, View"
    },
    {
        "name": "foo.pdf",
        "size": "20kb",
        "modTime": "2019-05-01 10:10:10",
        "actions": "Download, View"
    },
    {
        "name": "bar.pdf",
        "size": "28kb",
        "modTime": "2019-05-10 11:10:10",
        "actions": "Download, View"
    }
]
`