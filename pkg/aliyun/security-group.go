package aliyun

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	ecs20140526 "github.com/alibabacloud-go/ecs-20140526/v2/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/golang/glog"
	"golang/pkg/common"
)

func CreateECSClient(accessKeyId *string, accessKeySecret *string) (_result *ecs20140526.Client, _err error) {
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

func GetSecurityByTags(deployInfo common.DeployInfo) (result *ecs20140526.DescribeSecurityGroupsResponse, _err error) {

	client, _err := CreateECSClient(tea.String(deployInfo.AK), tea.String(deployInfo.SK))
	if _err != nil {
		return nil, _err
	}
	tags := []*ecs20140526.DescribeSecurityGroupsRequestTag{}
	for k, v := range deployInfo.TAG {
		tag := ecs20140526.DescribeSecurityGroupsRequestTag{}
		tag.SetKey(k)
		tag.SetValue(v)
		tags = append(tags, &tag)
	}

	describeSecurityGroupsRequest := ecs20140526.DescribeSecurityGroupsRequest{
		RegionId: tea.String(deployInfo.REGION),
		Tag:      tags,
		VpcId:    tea.String(deployInfo.VPC),
	}
	runtime := &util.RuntimeOptions{}
	resp, _err := client.DescribeSecurityGroupsWithOptions(&describeSecurityGroupsRequest, runtime)
	if _err != nil {
		return nil, _err
	}

	respStr := util.ToJSONString(tea.ToMap(resp))
	glog.Infof("%+v \n", *respStr)
	return resp, _err

}

func CreateSecurityByTags(deployInfo common.DeployInfo) (_result *ecs20140526.CreateSecurityGroupResponse, _err error) {
	client, _err := CreateECSClient(tea.String(deployInfo.AK), tea.String(deployInfo.SK))
	if _err != nil {
		return nil, _err
	}

	tags := []*ecs20140526.CreateSecurityGroupRequestTag{}
	for k, v := range deployInfo.TAG {
		tag := ecs20140526.CreateSecurityGroupRequestTag{}
		tag.SetKey(k)
		tag.SetValue(v)
		tags = append(tags, &tag)
	}

	createSecurityGroupRequest := &ecs20140526.CreateSecurityGroupRequest{
		RegionId: tea.String(deployInfo.REGION),
		Tag:      tags,
		VpcId:    tea.String(deployInfo.VPC),
	}
	runtime := &util.RuntimeOptions{}
	resp, _err := client.CreateSecurityGroupWithOptions(createSecurityGroupRequest, runtime)
	if _err != nil {
		return nil, _err
	}
	glog.Info(*util.ToJSONString(tea.ToMap(resp)))
	return resp, _err
}