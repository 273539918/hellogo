package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", jsonHandler)
	fmt.Println("Start Server")
	err := http.ListenAndServe(":80", mux)
	if err != nil {
		fmt.Println(err)
	}
}

func jsonHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("default handler")
	//fmt.Printf("%v \n", req.Header)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	respMap := make(map[string][]string)
	for k, v := range req.Header {
		respMap[k] = v
	}
	respMap["a"] = []string{"A"}
	jsonResp, err := json.Marshal(respMap)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(jsonResp)

	//w.Write(json.Marshal(resp))
}
