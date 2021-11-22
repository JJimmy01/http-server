package main

import (
	"context"
	"encoding/json"
	"github.com/gogf/gf/os/glog"
	"http-server/message"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type healthzHandler struct{}

func (healthzHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	respHeader := w.Header()
	respHeader.Set("Version", envVar("VERSION"))
	respHeader.Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	respHeader.Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	respHeader.Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(message.RestReply{Code: "200", Msg: "Ok", Data: nil})
	if err != nil {
		return
	}
	glog.Infof("Req IP %s  HttpCode %d", req.RemoteAddr, http.StatusOK)
}

func envVar(name string) string {
	val := os.Getenv(name)
	return val
}

func main() {
	// Conns close chan
	idleConnsClosed := make(chan struct{})
	graceQuitChan := make(chan os.Signal)
	signal.Notify(graceQuitChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	mux := http.NewServeMux()
	mux.Handle("/healthz", healthzHandler{})
	server := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}

	go func(graceQuitChan chan os.Signal, closeNotifyChan chan struct{}) {
		<-graceQuitChan
		if err := server.Shutdown(context.Background()); err != nil {
			glog.Infof("HTTP Server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}(graceQuitChan, idleConnsClosed)

	glog.Info("HTTP Server start serving at 8080")
	glog.Fatal(server.ListenAndServe())
	<-idleConnsClosed
}
