package main

import (
	"github.com/golang/glog"
	"os"
	"strings"
)

func main() {

	environ := os.Environ()
	//fmt.Printf("environ is %v", environ)
	for _, part := range environ {
		env := strings.SplitN(part, "=", 2)
		if len(env) != 2 {
			glog.Warningf("env formatter error %v", env)
			continue
		} else {
			key := env[0]
			value := env[1]
			glog.Infof("%s:%s \n", key, value)
		}

	}

}
