package aliyun

import (
	ecs20140526 "github.com/alibabacloud-go/ecs-20140526/v2/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/golang/glog"
	"golang/pkg/common"
)

func GetEcsByTag(deployInfo common.DeployInfo) (_result *ecs20140526.DescribeInstancesResponse, _err error) {
	client, _err := CreateECSClient(tea.String(deployInfo.AK), tea.String(deployInfo.SK))
	if _err != nil {
		return nil, _err
	}

	tags := []*ecs20140526.DescribeInstancesRequestTag{}
	for k, v := range deployInfo.TAG {
		tag := ecs20140526.DescribeInstancesRequestTag{}
		tag.SetKey(k)
		tag.SetValue(v)
		tags = append(tags, &tag)
	}

	describeInstancesRequest := &ecs20140526.DescribeInstancesRequest{
		RegionId:     tea.String("us-west-1"),
		Tag:          tags,
		VpcId:        tea.String(deployInfo.VPC),
		ZoneId:       tea.String(deployInfo.ZONE),
		InstanceName: tea.String(deployInfo.REGION), //使用region名作为ECS名字
	}
	runtime := &util.RuntimeOptions{}
	resp, _err := client.DescribeInstancesWithOptions(describeInstancesRequest, runtime)
	if _err != nil {
		return nil, _err
	}

	glog.Info(*util.ToJSONString(tea.ToMap(resp)))
	return resp, nil
}
