package common

import (
	"encoding/json"
	"github.com/golang/glog"
	"io"
	"net/http"
	"time"
)

type DeployInfo struct {
	VPC    string
	REGION string
	ZONE   string
	TAG    map[string]string
	AK     string
	SK     string
}

type NormalResponse struct {
	Message string
	Code    int
	Data    map[string]string
}

type ArrayResponse struct {
	Message string
	Code    int
	Data    []map[string]string
}

// 发起http请求，将请求的结果存入target 结构体
// data 是http的请求参数
func DoRequest(url string, method string, data io.Reader, header map[string]string, target interface{}, timeout time.Duration) error {

	glog.Info("request URL:", url)
	glog.Info("request data:", data)
	req, err := http.NewRequest(method, url, data)
	if err != nil {
		glog.Fatal(err)
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}
	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	glog.Info("response Status:", resp.Status)
	//glog.Info("response Headers:", resp.Header)
	//glog.Info("response Body:", resp.Body)
	return json.NewDecoder(resp.Body).Decode(target)
}
