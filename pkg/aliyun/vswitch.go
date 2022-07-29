package aliyun

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	vpc20160428 "github.com/alibabacloud-go/vpc-20160428/v2/client"
	"github.com/golang/glog"
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

func CreateVSW(AK string, SK string, region string, zone string, _tag map[string]string, vpc string) (_err error) {
	client, _err := CreateVPClient(tea.String(AK), tea.String(SK))
	if _err != nil {
		return _err
	}

	createVSwitchRequest := &vpc20160428.CreateVSwitchRequest{
		CidrBlock:   tea.String(VSWCidr),
		ZoneId:      tea.String(zone),
		VpcId:       tea.String(vpc),
		VSwitchName: tea.String(VSWName),
	}
	runtime := &util.RuntimeOptions{}

	resp, _create_error := client.CreateVSwitchWithOptions(createVSwitchRequest, runtime)
	if _create_error != nil {
		return _create_error
	}
	glog.Info(*util.ToJSONString(tea.ToMap(resp)))

	_tag_error := TagVSW(AK, SK, region, _tag, resp.Body.VSwitchId)
	if _tag_error != nil {
		glog.Fatal(_tag_error)
	}
	return _create_error
}

func TagVSW(AK string, SK string, region string, _tag map[string]string, vsw *string) (_err error) {
	client, _err := CreateVPClient(tea.String(AK), tea.String(SK))
	if _err != nil {
		return _err
	}

	tags := []*vpc20160428.TagResourcesRequestTag{}
	for k, v := range _tag {
		tag := vpc20160428.TagResourcesRequestTag{}
		tag.SetKey(k)
		tag.SetValue(v)
		tags = append(tags, &tag)
	}
	resourceId := []*string{vsw}

	tagResourcesRequest := &vpc20160428.TagResourcesRequest{
		RegionId:     tea.String(region),
		Tag:          tags,
		ResourceType: tea.String("VSWITCH"),
		ResourceId:   resourceId,
	}
	runtime := &util.RuntimeOptions{}
	_, _tag_error := client.TagResourcesWithOptions(tagResourcesRequest, runtime)
	if _tag_error != nil {
		glog.Fatal(_tag_error)
	}
	return _tag_error
}

func GetVSW(AK string, SK string, region string, zone string, vpc string) (_result *vpc20160428.DescribeVSwitchesResponse, _err error) {
	client, _err := CreateVPClient(tea.String(AK), tea.String(SK))
	if _err != nil {
		glog.Fatal(_err)
		return nil, _err
	}

	describeVSwitchesRequest := &vpc20160428.DescribeVSwitchesRequest{
		RegionId:    tea.String(region),
		ZoneId:      tea.String(zone),
		VpcId:       tea.String(vpc),
		VSwitchName: tea.String(VSWName),
		PageNumber:  tea.Int32(1),
		PageSize:    tea.Int32(1),
	}
	runtime := &util.RuntimeOptions{}
	resp, _err := client.DescribeVSwitchesWithOptions(describeVSwitchesRequest, runtime)
	if _err != nil {
		glog.Fatal(_err)
	}

	return resp, _err
}
