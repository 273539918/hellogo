package aliyun

import (
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	vpc20160428 "github.com/alibabacloud-go/vpc-20160428/v2/client"
	"github.com/golang/glog"
	"golang/pkg/common"
)

type NatInfo struct {
	SnatTableId string
	SnatIp      string
}

func CreateSnatEntry(deployInfo common.DeployInfo, vsw string) (_err error) {
	glog.Info("新建SNAT")
	client, _err := CreateVPClient(tea.String(deployInfo.AK), tea.String(deployInfo.SK))
	if _err != nil {
		return _err
	}

	createSnatEntryRequest := &vpc20160428.CreateSnatEntryRequest{
		RegionId:        tea.String(deployInfo.REGION),
		SnatTableId:     tea.String("stb-rj9n6fxhe0lswkpvv7hla"),
		SnatIp:          tea.String("47.88.63.30"),
		SourceVSwitchId: tea.String(vsw),
	}
	resp, _create_err := client.CreateSnatEntryWithOptions(createSnatEntryRequest, &util.RuntimeOptions{})
	if _create_err != nil {
		glog.Fatal(_create_err)
	}
	glog.Info(*util.ToJSONString(tea.ToMap(resp)))
	return _create_err
}

func GetSnatEntry(deployInfo common.DeployInfo, vsw string) (_result *vpc20160428.DescribeSnatTableEntriesResponse, _err error) {
	glog.Info("查询SNAT")
	client, _err := CreateVPClient(tea.String(deployInfo.AK), tea.String(deployInfo.SK))
	if _err != nil {
		return nil, _err
	}

	describeSnatTableEntriesRequest := &vpc20160428.DescribeSnatTableEntriesRequest{
		RegionId:        tea.String(deployInfo.REGION),
		SnatTableId:     tea.String("stb-rj9n6fxhe0lswkpvv7hla"),
		SourceVSwitchId: tea.String(vsw),
	}
	resp, _desc_err := client.DescribeSnatTableEntriesWithOptions(describeSnatTableEntriesRequest, &util.RuntimeOptions{})
	if _desc_err != nil {
		glog.Fatal(_desc_err)
		return nil, _desc_err
	}
	return resp, _desc_err

}
