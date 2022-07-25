package abm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"golang/pkg/common"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

var CONFIG = map[string]string{}

func init() {
	absPath, _ := filepath.Abs("../../pkg/abm/secret.json")
	fmt.Println("loading file :", absPath)
	content, err := ioutil.ReadFile(absPath)
	if err != nil {
		glog.Fatal("Error when opening file: ", err)
	}
	err = json.Unmarshal(content, &CONFIG)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

}

func GETABMToken() string {
	data := url.Values{}
	data.Set("username", CONFIG["username"])
	data.Set("password", "44079ec1a1eba6099746bc374b62339e")
	data.Set("client_id", "tesla_blink")
	data.Set("client_secret", "eda473e0-c77f-44d1-938e-5ee364740f6d")
	data.Set("grant_type", "password")
	encodedData := data.Encode()
	header := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}
	bodyJson := &common.NormalResponse{}
	//bodyJson := &NormalResponse{}
	common.DoRequest(CONFIG["abm_gw"]+"v2/common/authProxy/oauth/token", http.MethodPost,
		strings.NewReader(encodedData), header, bodyJson, time.Second*10)
	//glog.Infof("get response data:%v\n", bodyJson.Data)
	token, exists := bodyJson.Data["token"]
	if exists == false {
		glog.Error("get token fail")
	}
	return token
}

func ABMRequest(url string, data string, target interface{}) {
	token := GETABMToken()
	//glog.Info("abm token:", token)

	header := map[string]string{
		"Authorization": token,
		"Content-Type":  "application/json",
	}
	common.DoRequest(url, http.MethodPost, bytes.NewBuffer([]byte(data)), header, target, time.Second*10)
}

func SelectCMDB(product string, data string, target interface{}) {
	url := CONFIG["abm_endpoint"] + "gateway/v2/foundation/teslacmdb/entity/product/selectAll?product=" + product
	ABMRequest(url, data, target)
}

func GetAccountInfo() map[string]string {
	jsonBody := `{
	"fields": [],
	"filters": {},
	"sort": {},
	"table": "DELIVERY_PLATFORM",
	"withBaseCols": true
	}`
	target := &common.ArrayResponse{}
	SelectCMDB("blink", jsonBody, target)
	for _, m := range target.Data {
		value, exists := m["account"]
		if exists && value == CONFIG["account"] {
			return m
		}
		continue
	}
	glog.Fatal("cannot found target account :", CONFIG["account"])
	return nil

}
