package main

import (
	"flag"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	ecs20140526 "github.com/alibabacloud-go/ecs-20140526/v2/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/golang/glog"
	"golang/pkg/abm"
	"os"
)

const (
	VPC    = "vpc-rj98d2apg2z97c94wl1o8"
	REGION = "us-west-1"
)

var ACCOUNT = make(map[string]string, 10)

func init() {
	flag.Parse()
	flag.Set("alsologtostderr", "true")
	defer glog.Flush()
	ACCOUNT = abm.GetAccountInfo()
}

func CreateClient(accessKeyId *string, accessKeySecret *string) (_result *ecs20140526.Client, _err error) {
	config := &openapi.Config{
		// 您的 AccessKey ID
		AccessKeyId: accessKeyId,
		// 您的 AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("ecs.us-west-1.aliyuncs.com")
	_result = &ecs20140526.Client{}
	_result, _err = ecs20140526.NewClient(config)
	return _result, _err
}

func _main(args []*string) (respMap map[string]interface{}, _err error) {
	client, _err := CreateClient(tea.String(ACCOUNT["accessKey"]), tea.String(ACCOUNT["accessSecret"]))
	if _err != nil {
		return nil, _err
	}

	tag := ecs20140526.DescribeSecurityGroupsRequestTag{}
	tag.SetKey("abm_console")
	tag.SetValue("true")

	tags := []*ecs20140526.DescribeSecurityGroupsRequestTag{}
	tags = append(tags, &tag)

	describeSecurityGroupsRequest := ecs20140526.DescribeSecurityGroupsRequest{
		RegionId: tea.String("us-west-1"),
		Tag:      tags,
	}
	runtime := &util.RuntimeOptions{}
	resp, _err := client.DescribeSecurityGroupsWithOptions(&describeSecurityGroupsRequest, runtime)
	if _err != nil {
		return nil, _err
	}

	respMap = tea.ToMap(resp)
	respStr := util.ToJSONString(tea.ToMap(resp))
	glog.Infof("%+v \n", *respStr)
	return respMap, _err
}

func main() {
	glog.Info(ACCOUNT)
	respMap, err := _main(tea.StringSlice(os.Args[1:]))
	if err != nil {
		panic(err)
	}
	body := respMap["body"]
	securityGroups := tea.ToMap(body)["SecurityGroups"]
	//glog.Info(securityGroups)
	securityGroup := tea.ToMap(securityGroups)["SecurityGroup"].([]interface{})
	//glog.Info(securityGroup)
	for _, v := range securityGroup {
		m := tea.ToMap(v)
		glog.Info(m["SecurityGroupId"])
	}
}
