package main

import (
	"flag"
	"github.com/golang/glog"
)

// 1、运行：mkdir -p log && go run glog.go -log_dir=log -alsologtostderr
// 打印日志将会同时打印在 log/ 目录和标准错误输出中（-alsologtostderr）

// 2、运行：mkdir -p log && go run glog.go -v=4 -log_dir=log -alsologtostderr
// 日志级别小于等于4的日志会被输出
func main() {
	flag.Parse()
	defer glog.Flush()

	glog.Info("This is info message ")
	glog.Infof("This is info message: %v", 123456)
	glog.InfoDepth(1, "This is info message", 12345)

	glog.Warning("This is warning message")
	glog.Warningf("This is warning message: %v", 123456)
	glog.WarningDepth(1, "This is warning message", 123456)

	glog.Error("This is error message")
	glog.Errorf("This is error message:%v", 123456)
	glog.ErrorDepth(1, "This is error message", 123456)

	//调用fatal后进程会调用os.exit退出
	//glog.Fatal("This is fatal message")
	//glog.Fatalf("This is fatal message", 123456)
	//glog.FatalDepth(1, "This is fatal message", 123456)
	
	glog.V(3).Info("LEVEL 3 message") // 使用日志级别 3
	glog.V(4).Info("LEVEL 4 message") // 使用日志级别 4
	glog.V(5).Info("LEVEL 5 message") // 使用日志级别 5
	glog.V(8).Info("LEVEL 8 message") // 使用日志级别 8

}
