package aliyun

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	ecs20140526 "github.com/alibabacloud-go/ecs-20140526/v2/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/golang/glog"
	"golang/pkg/common"
	"time"
)

const (
	IMAGEID = "aliyun_2_1903_x64_20G_alibase_20200324.vhd"
	DRYRUN  = false
	RUNNING = "Running"
	STOP    = "Stopped"
)

var INSTANCETYPES = []string{"ecs.hfg6.large"}

func CreateECSClient(accessKeyId *string, accessKeySecret *string, region string) (_result *ecs20140526.Client, _err error) {
	config := &openapi.Config{
		// 您的 AccessKey ID
		AccessKeyId: accessKeyId,
		// 您的 AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	endpoint := "ecs." + region + ".aliyuncs.com"
	config.Endpoint = tea.String(endpoint)
	_result = &ecs20140526.Client{}
	_result, _err = ecs20140526.NewClient(config)
	return _result, _err
}

func GetEcsClient(deployInfo common.DeployInfo) (_result *ecs20140526.Client) {
	client, _err := CreateECSClient(tea.String(deployInfo.AK), tea.String(deployInfo.SK), deployInfo.REGION)
	if _err != nil {
		glog.Fatal(_err)
	}
	return client
}

func GetEcsByTag(deployInfo common.DeployInfo) (_result *ecs20140526.DescribeInstancesResponse, _err error) {
	client := GetEcsClient(deployInfo)

	tags := []*ecs20140526.DescribeInstancesRequestTag{}
	for k, v := range deployInfo.TAG {
		tag := ecs20140526.DescribeInstancesRequestTag{}
		tag.SetKey(k)
		tag.SetValue(v)
		tags = append(tags, &tag)
	}

	describeInstancesRequest := &ecs20140526.DescribeInstancesRequest{
		RegionId:     tea.String(deployInfo.REGION),
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

func GetECSDetailById(deployInfo common.DeployInfo, ecsIds []string) (_result *ecs20140526.DescribeInstancesResponse, _err error) {
	client := GetEcsClient(deployInfo)

	glog.Info(common.StringArrayToJsonArrStr(ecsIds))
	describeInstancesRequest := &ecs20140526.DescribeInstancesRequest{
		RegionId:    tea.String(deployInfo.REGION),
		InstanceIds: tea.String(common.StringArrayToJsonArrStr(ecsIds)),
		//InstanceIds: tea.String("[\"i-rj996a1oy510h4fwvstf\"]"),
	}

	resp, _err := client.DescribeInstancesWithOptions(describeInstancesRequest, &util.RuntimeOptions{})
	if _err != nil {
		return nil, _err
	}
	return resp, _err

}

func GetAvailableECSInstanceType() (intanceType string) {

	return "ecs.hfg6.large"

}

func CreateECSWithTag(deployInfo common.DeployInfo, vsw string, sg string) (_err error) {

	client := GetEcsClient(deployInfo)

	tags := []*ecs20140526.CreateInstanceRequestTag{}
	for k, v := range deployInfo.TAG {
		tag := ecs20140526.CreateInstanceRequestTag{}
		tag.SetKey(k)
		tag.SetValue(v)
		tags = append(tags, &tag)
	}
	createInstanceRequest := &ecs20140526.CreateInstanceRequest{
		InstanceName:    tea.String(deployInfo.REGION),
		Tag:             tags,
		VSwitchId:       tea.String(vsw),
		SecurityGroupId: tea.String(sg),
		ZoneId:          tea.String(deployInfo.ZONE),
		InstanceType:    tea.String(GetAvailableECSInstanceType()),
		ImageId:         tea.String(IMAGEID),
		RegionId:        tea.String(deployInfo.REGION),
		DryRun:          tea.Bool(DRYRUN),
	}
	resp, _create_err := client.CreateInstanceWithOptions(createInstanceRequest, &util.RuntimeOptions{})
	if _create_err != nil {
		glog.Fatal(_create_err)
	}

	glog.Info(*util.ToJSONString(tea.ToMap(resp)))
	return _create_err
}

func GetEcsStatus(deployInfo common.DeployInfo, instanceIds []*string) (result map[string]string, err error) {
	client := GetEcsClient(deployInfo)

	describeInstanceStatusRequest := &ecs20140526.DescribeInstanceStatusRequest{
		RegionId:   tea.String(deployInfo.REGION),
		InstanceId: instanceIds,
	}

	resp, _err := client.DescribeInstanceStatusWithOptions(describeInstanceStatusRequest, &util.RuntimeOptions{})
	if _err != nil {
		glog.Fatal(_err)
		return nil, _err
	}
	result = map[string]string{}
	for _, instance := range resp.Body.InstanceStatuses.InstanceStatus {
		result[*instance.InstanceId] = *instance.Status
	}

	return result, _err
}

func StartECSInstance(deployInfo common.DeployInfo, instanceId string) (success bool) {
	client := GetEcsClient(deployInfo)

	startInstanceRequest := &ecs20140526.StartInstanceRequest{
		DryRun:     tea.Bool(DRYRUN),
		InstanceId: tea.String(instanceId),
	}

	_, _start_err := client.StartInstanceWithOptions(startInstanceRequest, &util.RuntimeOptions{})
	if _start_err != nil {
		glog.Fatal(_start_err)
		return false
	}
	instanceStatus, _err := GetEcsStatus(deployInfo, []*string{&instanceId})
	if _err != nil {
		glog.Fatal(_err)
		return false
	}
	for {
		if instanceStatus[instanceId] == RUNNING {
			glog.Info("ECS已启动:", instanceId)
			return true
		} else {
			glog.Info("ECS还未启动:", instanceId)
			time.Sleep(time.Second * 10)
			instanceStatus, _ = GetEcsStatus(deployInfo, []*string{&instanceId})
		}
	}
	return false
}

func CreateCommand(deployInfo common.DeployInfo, name string, commandContent *string) (_err error) {
	client := GetEcsClient(deployInfo)
	createCommandRequest := &ecs20140526.CreateCommandRequest{
		RegionId:        tea.String(deployInfo.REGION),
		Name:            tea.String(name),
		Type:            tea.String("RunShellScript"),
		CommandContent:  commandContent,
		ContentEncoding: tea.String("PlainText"),
	}
	runtime := &util.RuntimeOptions{}
	resp, _err := client.CreateCommandWithOptions(createCommandRequest, runtime)
	if _err != nil {
		return _err
	}

	glog.Info(*util.ToJSONString(tea.ToMap(resp)))
	return _err
}

func GetCommandByName(deployInfo common.DeployInfo, commandName string) (_result *ecs20140526.DescribeCommandsResponse, _err error) {
	client := GetEcsClient(deployInfo)

	describeCommandsRequest := &ecs20140526.DescribeCommandsRequest{
		RegionId: tea.String(deployInfo.REGION),
		Name:     tea.String(commandName),
	}
	resp, _desc_err := client.DescribeCommandsWithOptions(describeCommandsRequest, &util.RuntimeOptions{})
	if _desc_err != nil {
		glog.Fatal(_desc_err)
	}
	return resp, _desc_err

}

func InvokeECSCommand(deployInfo common.DeployInfo, instanceId string, commanId string) (success bool) {
	glog.Infof("%s 执行 %s", instanceId, commanId)
	client := GetEcsClient(deployInfo)

	describeInvocationResultsRequest := &ecs20140526.DescribeInvocationResultsRequest{
		RegionId:           tea.String(deployInfo.REGION),
		InstanceId:         tea.String(instanceId),
		CommandId:          tea.String(commanId),
		InvokeRecordStatus: tea.String("Running"),
	}
	resp, _desc_error := client.DescribeInvocationResultsWithOptions(describeInvocationResultsRequest, &util.RuntimeOptions{})
	if _desc_error != nil {
		glog.Fatal(_desc_error)
	}

	if len(resp.Body.Invocation.InvocationResults.InvocationResult) <= 0 {
		glog.Info("无相同的命令在执行中，新建执行命令")
		invokeCommandRequest := &ecs20140526.InvokeCommandRequest{
			RegionId:   tea.String(deployInfo.REGION),
			InstanceId: []*string{&instanceId},
			CommandId:  tea.String(commanId),
		}
		client.InvokeCommandWithOptions(invokeCommandRequest, &util.RuntimeOptions{})
	}

	describeInvocationResultsRequest = &ecs20140526.DescribeInvocationResultsRequest{
		RegionId:   tea.String(deployInfo.REGION),
		InstanceId: tea.String(instanceId),
		CommandId:  tea.String(commanId),
	}
	resp, _desc_error = client.DescribeInvocationResultsWithOptions(describeInvocationResultsRequest, &util.RuntimeOptions{})
	if _desc_error != nil {
		glog.Fatal(_desc_error)
	}

	invokeRecordStatus := *new(string)
	for {
		invokeResult := resp.Body.Invocation.InvocationResults.InvocationResult[0]
		invokeRecordStatus = *invokeResult.InvokeRecordStatus
		if invokeRecordStatus == "Failed" || invokeRecordStatus == "Stopped" {
			glog.Fatalf("云助手命令执行失败 %+v", invokeResult)
			return false
		} else if invokeRecordStatus == "Running" {
			glog.Infof("云助手命令执行中")
			time.Sleep(time.Second * 10)
			resp, _ = client.DescribeInvocationResultsWithOptions(describeInvocationResultsRequest, &util.RuntimeOptions{})
		} else {
			glog.Infof("命令执行成功")
			break
		}
	}
	return true

}
