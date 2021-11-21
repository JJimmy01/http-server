package main

import (
	"context"
	"encoding/json"
	"github.com/gogf/gf/os/glog"
	pb "jimmy.com/http-server/message"
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
	for k := range req.Header {
		respHeader.Set(k, req.Header.Get(k))
	}
	w.WriteHeader(http.StatusOK)
	res, _ := json.Marshal(pb.RestReply{Code: "200", Msg: "Ok", Data: nil})
	_, err := w.Write(res)
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
			glog.Info("HTTP Server Shutdown: %v", err)
		}

	}(graceQuitChan, idleConnsClosed)
	glog.Fatal(server.ListenAndServe())
}
