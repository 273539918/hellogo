package main

import (
	"flag"
	"github.com/golang/glog"
	"golang/pkg/abm"
	"golang/pkg/aliyun"
)

const (
	VPC    = "vpc-rj98d2apg2z97c94wl1o8"
	REGION = "us-west-1"
	ZONE   = "us-west-1b"
)

var TAG = map[string]string{
	"abm_console": "true",
}

var ACCOUNT = make(map[string]string, 10)

func init() {
	flag.Parse()
	flag.Set("alsologtostderr", "true")
	defer glog.Flush()
	ACCOUNT = abm.GetAccountInfo()
}

func GetSecurityIdsByTags() []string {
	glog.Info("查询安全组")
	resp, err := aliyun.GetSecurityByTags(ACCOUNT["accessKey"], ACCOUNT["accessSecret"], REGION, TAG, VPC)
	if err != nil {
		glog.Fatal(err)
	}
	//glog.Info(resp.Body.SecurityGroups.SecurityGroup)
	securityIds := []string{}
	for _, sg := range resp.Body.SecurityGroups.SecurityGroup {
		//glog.Info(sg)
		//glog.Info(*sg.SecurityGroupId)
		securityIds = append(securityIds, *sg.SecurityGroupId)
	}
	return securityIds
}

func CreateSecurityByTas() (_err error) {
	glog.Info("创建安全组")
	resp, err := aliyun.CreateSecurityByTags(ACCOUNT["accessKey"], ACCOUNT["accessSecret"], REGION, TAG, VPC)
	if err != nil {
		glog.Fatal(err)
	}
	glog.Info(resp.Body)
	return err
}

func GetVSWIdsByTags() []string {
	glog.Info("查询VSW")
	resp, _err := aliyun.GetVSW(ACCOUNT["accessKey"], ACCOUNT["accessSecret"], REGION, ZONE, VPC)
	if _err != nil {
		glog.Fatal(_err)
	}
	//glog.Info("VSW信息:", resp.Body.VSwitches.VSwitch)
	targetVSWId := []string{}
	//检查vsw里是tag 是否包含目标 TAG
	for _, vsw := range resp.Body.VSwitches.VSwitch {
		tagIsInclude := aliyun.TagIsInclude(vsw.Tags.Tag, TAG)
		if tagIsInclude {
			targetVSWId = append(targetVSWId, *vsw.VSwitchId)
		}
	}
	return targetVSWId
}

func CreateVswWithTag() (_err error) {
	glog.Info("创建VSW")
	err := aliyun.CreateVSW(ACCOUNT["accessKey"], ACCOUNT["accessSecret"], REGION, ZONE, TAG, VPC)
	if err != nil {
		glog.Info(err)
	}
	return err
}

func main() {
	glog.Info(ACCOUNT)
	securityIds := GetSecurityIdsByTags()
	glog.Info("安全组信息:", securityIds)
	if len(securityIds) < 1 {
		err := CreateSecurityByTas()
		if err != nil {
			glog.Fatal(err)
		}
		securityIds = GetSecurityIdsByTags()
		glog.Info("安全组信息:", securityIds)
	}
	securityId := securityIds[0]
	glog.Info("选择安全组:", securityId)
	vswIds := GetVSWIdsByTags()
	glog.Info("VSW信息:", vswIds)
	if len(vswIds) < 1 {
		err := CreateVswWithTag()
		if err != nil {
			glog.Fatal(err)
		}
		vswIds = GetVSWIdsByTags()
		glog.Info("VSW信息:", vswIds)
	}
	vswId := vswIds[0]
	glog.Info("选择VSW:", vswId)

}
