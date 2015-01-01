# httpdev

A simple CLI http server for dev purposes

## usage

```
usage: httpdev [options...] path_or_html_content?

  -h          Display the help message
  -address    Listening address
  -port       Listening port
  -status     HTTP status code if serving single file or arbitrary html content
  -prefix     Root path of the url

```

## examples

```
$ httpdev 
$ httpdev /tmp
$ httpdev -address localhost /tmp
$ httpdev -address localhost -port 8080 /tmp
$ httpdev -prefix /static/ /tmp
$ httpdev /tmp/500.html
$ httpdev -status 500 /tmp/500.html
$ httpdev '<h1>HELLO WORLD!</h1>'

```

## install
- from source
```
go get github.com/mathieuancelin/httpdev
```
