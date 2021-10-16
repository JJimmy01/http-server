package main

import (
	"github.com/gogf/gf/os/glog"
	"net/http"
	"os"
)

func envVar(name string) string {
	val := os.Getenv(name)
	return val
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	resh := w.Header()
	for k := range r.Header {
		resh.Set(k, r.Header.Get(k))
	}
	resh.Set("Version", envVar("VERSION"))
	resh.Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	resh.Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	resh.Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(r.RemoteAddr))
	if err != nil {
		return
	}
	glog.Infof("Req IP %s  HttpCode %d", r.RemoteAddr, http.StatusOK)
}

func main() {
	http.HandleFunc("/healthz", healthHandler)
	glog.Fatal(http.ListenAndServe(":8082", nil))
}
