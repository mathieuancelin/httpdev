package main

import (
  "log"
  "net/http"
  "flag"
  "fmt"
  "os"
  "path/filepath"
  "io/ioutil"
)

var (
  address = flag.String("address", "0.0.0.0", "Listening address")
  port    = flag.String("port", "8000", "Listening port")
  status  = flag.Int("status", 200, "HTTP status code")
  prefix  = flag.String("root", "/", "Root path of the url")
  help    = flag.Bool("h", false, "Display the help message")
)

type bytesHandler []byte

type FakeResponseWriter struct {
  underlying http.ResponseWriter
  statusCode int
}

func (w *FakeResponseWriter) Header() http.Header {
  return w.underlying.Header()
}

func (w *FakeResponseWriter) Write(bytes []byte) (int, error) {
  return w.underlying.Write(bytes)
}

func (w *FakeResponseWriter) WriteHeader(code int) {
  w.statusCode = code
  w.underlying.WriteHeader(code) 
}

func (h bytesHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
  w.Header().Set("Content-Type", "text/html; charset=utf-8")
  w.WriteHeader(*status)
  w.Write(h)
}

type LoggableHandler struct {
  underlying http.Handler
}

func (h LoggableHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
  fakeW := &FakeResponseWriter{w, 200,}
  h.underlying.ServeHTTP(fakeW, req)
  log.Println("(", req.Method, ":", fakeW.statusCode, ")\t=>\t", req.RequestURI)
}

func main() {
  flag.Parse()
  if (*help) {
    fmt.Println("\nhttpdev is an HTTP server for devs\n")
    fmt.Println("usage: httpdev [options...] path?\n")
    flag.PrintDefaults()
    os.Exit(0)
  }
  var root = "/"
  if (len(flag.Args()) < 1) {
    dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
      log.Fatal(err)
    }
    root = dir
  } else {
    dir := flag.Arg(0) 
    root = dir
  }

  if fi, err := os.Stat(root); err == nil {
    switch mode := fi.Mode(); {
      case mode.IsDir():
        log.Println("Listening at", *address, ":", *port, "and serving directory", root)
        fs := LoggableHandler{http.FileServer(http.Dir(root))}
        http.Handle(*prefix, http.StripPrefix(*prefix, fs))
        log.Fatal(http.ListenAndServe(*address + ":" + *port, nil))
      case mode.IsRegular():
        if content, err := ioutil.ReadFile(root); err != nil {
          log.Fatal("Error reading file: ", err)
        } else {
          log.Println("Listening at", *address, ":", *port, "and serving file", root)
          log.Fatal(http.ListenAndServe(*address + ":" + *port, LoggableHandler{bytesHandler(content)}))
        }
      }
  } else {
    log.Println("Listening at", *address, ":", *port, "and serving some content")
    log.Fatal(http.ListenAndServe(*address + ":" + *port, LoggableHandler{bytesHandler(root)}))
  }
  //fs := http.FileServer(http.Dir(root))
  //http.Handle(*prefix, http.StripPrefix(*prefix, fs))
  //log.Println("Listening at ", *address, ":", *port, " and serving ", root)
  //log.Fatal(http.ListenAndServe(*address + ":" + *port, nil))
  //log.Fatal(http.ListenAndServe(":" + *port, http.FileServer(http.Dir(root))))
}
