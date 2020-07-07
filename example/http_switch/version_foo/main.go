package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os/exec"

	"github.com/sirupsen/logrus"

	"github.com/lyp256/pnet"
)

var tplSet = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>http_switch_demo</title>
</head>
<body>
<div>this foo server</div>
<div><a href="./switch_bar">switch to bar server</a></div>
<div><a href=".">refresh</a></div>
<div><a href="./stop">stop</a></div>
</body>
</html>
`

func main() {
	stopped := make(chan struct{})
	mux := http.NewServeMux()
	ser := http.Server{
		Handler: mux,
	}

	l, err := pnet.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	tpl, err := template.New("tpl").Parse(tplSet)
	if err != nil {
		panic(err)
	}
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		err := tpl.Execute(writer, nil)
		if err != nil {
			logrus.Error(err)
		}
	})

	mux.HandleFunc("/switch_foo", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "bar is currently running")
	})

	mux.HandleFunc("/switch_bar", func(writer http.ResponseWriter, request *http.Request) {
		execPath, err := exec.LookPath("./bar")
		if err != nil {
			fmt.Fprintf(writer, "get bar failed:%s", err)
			return
		}
		cmd := exec.Command(execPath)
		err = cmd.Start()
		if err != nil {
			fmt.Fprintf(writer, "bar start failed:%s", err)
			return
		}
		go ser.Shutdown(context.Background())
		fmt.Fprintf(writer, "bar start success")
	})

	mux.HandleFunc("/stop", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "bye")
		go func() {
			ser.Shutdown(context.Background())
			close(stopped)
		}()
	})

	ser.Serve(l)
	logrus.Println("server shutdown")
	<-stopped
}
