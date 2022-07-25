package main

import (
	"encoding/json"
	"flag"
	"github.com/golang/glog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
)

func main() {
	flag.Parse()
	defer glog.Flush()

	mux := http.NewServeMux()
	mux.HandleFunc("/header", headerHandler)
	mux.HandleFunc("/env/version", envVersionHandler)
	mux.HandleFunc("/healthz", healthzHandler)
	glog.Info("Start Server")
	err := http.ListenAndServe(":80", mux)
	if err != nil {
		glog.Fatal(err)
	}

}

func MapToJsonEncode(m map[string]interface{}) []byte {
	jsonResp, err := json.Marshal(m)
	if err != nil {
		glog.Fatal(err)
	}
	return jsonResp
}

func envVersionHandler(w http.ResponseWriter, req *http.Request) {
	glog.Info("envVersion handler")
	resMap := make(map[string]interface{})
	for _, part := range os.Environ() {
		env := strings.SplitN(part, "=", 2)
		if len(env) != 2 {
			glog.Warningf("env formatter error %v", env)
		} else {
			resMap[env[0]] = env[1]
		}
	}
	w.Write(MapToJsonEncode(resMap))
	RecordResq(w, req)
}

func headerHandler(w http.ResponseWriter, req *http.Request) {
	glog.Info("header handler")
	w.Header().Set("Content-Type", "application/json")
	respMap := make(map[string]interface{})
	for k, v := range req.Header {
		respMap[k] = strings.Join(v, ",")
	}
	respMap["a"] = "A"
	w.Write(MapToJsonEncode(respMap))
	RecordResq(w, req)

}

func healthzHandler(w http.ResponseWriter, req *http.Request) {
	glog.Info("healthz handler")
	w.WriteHeader(http.StatusOK)
	RecordResq(w, req)
}

//server端记录访问日志：客户端ip，http返回码
func RecordResq(w http.ResponseWriter, req *http.Request) {
	wRecord := httptest.NewRecorder().Result()
	clientIP := req.Header.Get("X-FORWARDED-FOR")
	if clientIP == "" {
		clientIP = req.RemoteAddr
	}
	glog.Infof("client ip: %s code: %d", clientIP, wRecord.StatusCode)
}
