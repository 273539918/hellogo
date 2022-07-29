package aliyun

import (
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	vpc20160428 "github.com/alibabacloud-go/vpc-20160428/v2/client"
	"github.com/golang/glog"
	"golang/pkg/common"
)

//Tag是否包含targetTAG的key,value
func TagIsInclude(Tag []*vpc20160428.DescribeVSwitchesResponseBodyVSwitchesVSwitchTagsTag, targetTAG map[string]string) (include bool) {

	flagOuter := true
	for k, v := range targetTAG {
		flagInner := false
		for _, tag := range Tag {
			key := *tag.Key
			value := *tag.Value
			if key == k && value == v {
				flagInner = true
				break
			}
		}
		if flagInner {
			continue
		} else {
			flagOuter = false
		}
	}
	return flagOuter
}

func TagVSW(deployInfo common.DeployInfo, vsw *string) (_err error) {
	client, _err := CreateVPClient(tea.String(deployInfo.AK), tea.String(deployInfo.SK))
	if _err != nil {
		return _err
	}
	tags := []*vpc20160428.TagResourcesRequestTag{}
	for k, v := range deployInfo.TAG {
		tag := vpc20160428.TagResourcesRequestTag{}
		tag.SetKey(k)
		tag.SetValue(v)
		tags = append(tags, &tag)
	}
	resourceId := []*string{vsw}

	tagResourcesRequest := &vpc20160428.TagResourcesRequest{
		RegionId:     tea.String(deployInfo.REGION),
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
