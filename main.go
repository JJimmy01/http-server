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


func healthHandler(w http.ResponseWriter, r *http.Request)  {
	resh := w.Header()
	for k := range r.Header {
		resh.Set(k, r.Header.Get(k))
	}
	resh.Set("Version", envVar("VERSION"))
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(r.RemoteAddr))
	if err != nil {
		return
	}
	glog.Infof("Req IP %s  HttpCode %d", r.RemoteAddr, http.StatusOK)
}


func main()  {
	http.HandleFunc("/healthz", healthHandler)
	glog.Fatal(http.ListenAndServe(":8080", nil))
}
