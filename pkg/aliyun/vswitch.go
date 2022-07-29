package aliyun

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	vpc20160428 "github.com/alibabacloud-go/vpc-20160428/v2/client"
	"github.com/golang/glog"
	"golang/pkg/common"
)

const (
	VSWCidr = "172.20.0.0/16" //如果冲突了可以修改
	VSWName = "compute-node-vsw"
)

func CreateVPClient(accessKeyId *string, accessKeySecret *string) (_result *vpc20160428.Client, _err error) {
	config := &openapi.Config{
		// 您的 AccessKey ID
		AccessKeyId: accessKeyId,
		// 您的 AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("vpc.aliyuncs.com")
	_result = &vpc20160428.Client{}
	_result, _err = vpc20160428.NewClient(config)
	return _result, _err
}

func CreateVSW(deployInfo common.DeployInfo) (_err error) {
	client, _err := CreateVPClient(tea.String(deployInfo.AK), tea.String(deployInfo.SK))
	if _err != nil {
		return _err
	}

	createVSwitchRequest := &vpc20160428.CreateVSwitchRequest{
		CidrBlock:   tea.String(VSWCidr),
		ZoneId:      tea.String(deployInfo.ZONE),
		VpcId:       tea.String(deployInfo.VPC),
		VSwitchName: tea.String(VSWName),
	}
	runtime := &util.RuntimeOptions{}

	resp, _create_error := client.CreateVSwitchWithOptions(createVSwitchRequest, runtime)
	if _create_error != nil {
		return _create_error
	}
	glog.Info(*util.ToJSONString(tea.ToMap(resp)))

	_tag_error := TagVSW(deployInfo, resp.Body.VSwitchId)
	if _tag_error != nil {
		glog.Fatal(_tag_error)
	}
	return _create_error
}

func GetVSW(deployInfo common.DeployInfo) (_result *vpc20160428.DescribeVSwitchesResponse, _err error) {
	client, _err := CreateVPClient(tea.String(deployInfo.AK), tea.String(deployInfo.SK))
	if _err != nil {
		glog.Fatal(_err)
		return nil, _err
	}

	describeVSwitchesRequest := &vpc20160428.DescribeVSwitchesRequest{
		RegionId:    tea.String(deployInfo.REGION),
		ZoneId:      tea.String(deployInfo.ZONE),
		VpcId:       tea.String(deployInfo.VPC),
		VSwitchName: tea.String(VSWName),
		PageNumber:  tea.Int32(1),
		PageSize:    tea.Int32(1),
	}
	runtime := &util.RuntimeOptions{}
	resp, _err := client.DescribeVSwitchesWithOptions(describeVSwitchesRequest, runtime)
	if _err != nil {
		glog.Fatal(_err)
	}
	glog.Info(*util.ToJSONString(tea.ToMap(resp)))
	return resp, _err
}
