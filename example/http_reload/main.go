package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/lyp256/pnet"
)

var tplSet = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>http_reload_demo</title>
</head>
<body>
<div>this server start at {{.StartAt}}</div>
<div>This is the {{.Count}}th visit！</div>
<div>
	<a href="./reload">reload this server</a>
	<a href=".">refresh</a>
	<a href="./stop">stop</a>
</div>
</body>
</html>
`

func main() {
	stopped := make(chan struct{})
	startAt := time.Now()
	count := 0

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
		count++
		err := tpl.Execute(writer, map[string]interface{}{
			"StartAt": startAt,
			"Count":   count,
		})
		if err != nil {
			logrus.Error(err)
		}
	})

	mux.HandleFunc("/reload", func(writer http.ResponseWriter, request *http.Request) {
		cmd, err := pnet.CurrentCommand()
		if err != nil {
			fmt.Fprintf(writer, "get CurrentCommand failed:%s", err)
			return
		}
		err = cmd.Start()
		if err != nil {
			fmt.Fprintf(writer, "commond start failed:%s", err)
			return
		}
		go ser.Shutdown(context.Background())
		fmt.Fprintf(writer, "success")
		return
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
