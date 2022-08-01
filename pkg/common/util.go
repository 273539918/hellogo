package common

import (
	"encoding/json"
	"github.com/golang/glog"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"
)

type DeployInfo struct {
	VPC     string
	REGION  string
	ZONE    string
	TAG     map[string]string
	AK      string
	SK      string
	BASTION string //堡垒机
	CLUSTER string
	COMMAND string
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

func GetFileContent(fileName string) (str string) {
	filePath := "../../cmd/shell/" + fileName
	absPath, _ := filepath.Abs(filePath)
	//glog.Info(absPath)
	glog.Info("loading file :", absPath)
	content, _err := ioutil.ReadFile(absPath)
	if _err != nil {
		glog.Fatal(_err)
		return ""
	}
	str = string(content)
	return str
}

// input ["a","b"]
// output "["a","b"]"
func StringArrayToJsonArrStr(arr []string) (str string) {
	//fmt.Println(arr)
	r := "["
	for _, str := range arr {
		r = r + "\"" + str + "\"" + ","
	}
	r = r[:len(r)-1]
	r = r + "]"
	return r
}
