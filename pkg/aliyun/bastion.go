package aliyun

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	console "github.com/alibabacloud-go/tea-console/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	yundun_bastionhost20191209 "github.com/alibabacloud-go/yundun-bastionhost-20191209/client"
	"golang/pkg/common"
)

func CreateBastionClient(accessKeyId *string, accessKeySecret *string) (_result *yundun_bastionhost20191209.Client, _err error) {
	config := &openapi.Config{
		// 您的 AccessKey ID
		AccessKeyId: accessKeyId,
		// 您的 AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("bastionhost.cn-hongkong.aliyuncs.com")
	_result = &yundun_bastionhost20191209.Client{}
	_result, _err = yundun_bastionhost20191209.NewClient(config)
	return _result, _err
}

func BastionAddEcs(deployInfo common.DeployInfo, ecs string, eip string) (_err error) {
	client, _err := CreateBastionClient(tea.String(deployInfo.AK), tea.String(deployInfo.SK))
	if _err != nil {
		return _err
	}

	createHostRequest := &yundun_bastionhost20191209.CreateHostRequest{
		InstanceId:         tea.String("bastionhost-cn-2r42h8ac825"),
		ActiveAddressType:  tea.String("Private"),
		HostName:           tea.String(deployInfo.CLUSTER),
		Source:             tea.String("ECS"),
		OSType:             tea.String("Linux"),
		SourceInstanceId:   tea.String(ecs),
		InstanceRegionId:   tea.String("us-west-1"),
		HostPrivateAddress: tea.String(eip),
		Comment:            tea.String(deployInfo.COMMAND),
	}
	runtime := &util.RuntimeOptions{}
	resp, _err := client.CreateHostWithOptions(createHostRequest, runtime)
	if _err != nil {
		return _err
	}

	console.Log(util.ToJSONString(tea.ToMap(resp)))
	return _err
}
