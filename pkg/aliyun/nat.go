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

func GetNatInfo(deployInfo common.DeployInfo) (natInfo NatInfo) {

	natInfo = NatInfo{}
	client, _err := CreateVPClient(tea.String(deployInfo.AK), tea.String(deployInfo.SK))
	if _err != nil {
		glog.Fatal(_err)
	}

	describeNatGatewaysRequest := &vpc20160428.DescribeNatGatewaysRequest{
		RegionId: tea.String(deployInfo.REGION),
		VpcId:    tea.String(deployInfo.VPC),
		ZoneId:   tea.String(deployInfo.ZONE),
	}
	resp, _err := client.DescribeNatGatewaysWithOptions(describeNatGatewaysRequest, &util.RuntimeOptions{})
	if _err != nil {
		glog.Fatal(_err)
	}
	if len(resp.Body.NatGateways.NatGateway) <= 0 {
		glog.Fatalf("不存在nat: %+v", deployInfo)
	}
	natGateway := resp.Body.NatGateways.NatGateway[0]
	if len(natGateway.SnatTableIds.SnatTableId) <= 0 {
		glog.Fatal("不存在SnatTable")
	}
	if len(natGateway.IpLists.IpList) <= 0 {
		glog.Fatal("找不到公网ip")
	}

	natInfo.SnatTableId = *natGateway.SnatTableIds.SnatTableId[0]
	natInfo.SnatIp = *natGateway.IpLists.IpList[0].IpAddress

	return natInfo

}

func CreateSnatEntry(deployInfo common.DeployInfo, vsw string, natInfo NatInfo) (_err error) {
	glog.Info("新建SNAT")
	client, _err := CreateVPClient(tea.String(deployInfo.AK), tea.String(deployInfo.SK))
	if _err != nil {
		return _err
	}

	createSnatEntryRequest := &vpc20160428.CreateSnatEntryRequest{
		RegionId:        tea.String(deployInfo.REGION),
		SnatTableId:     tea.String(natInfo.SnatTableId),
		SnatIp:          tea.String(natInfo.SnatIp),
		SourceVSwitchId: tea.String(vsw),
	}
	resp, _create_err := client.CreateSnatEntryWithOptions(createSnatEntryRequest, &util.RuntimeOptions{})
	if _create_err != nil {
		glog.Fatal(_create_err)
	}
	glog.Info(*util.ToJSONString(tea.ToMap(resp)))
	return _create_err
}

func GetSnatEntry(deployInfo common.DeployInfo, vsw string, natInfo NatInfo) (_result *vpc20160428.DescribeSnatTableEntriesResponse, _err error) {
	glog.Info("查询SNAT")
	client, _err := CreateVPClient(tea.String(deployInfo.AK), tea.String(deployInfo.SK))
	if _err != nil {
		return nil, _err
	}

	describeSnatTableEntriesRequest := &vpc20160428.DescribeSnatTableEntriesRequest{
		RegionId:        tea.String(deployInfo.REGION),
		SnatTableId:     tea.String(natInfo.SnatTableId),
		SourceVSwitchId: tea.String(vsw),
	}
	resp, _desc_err := client.DescribeSnatTableEntriesWithOptions(describeSnatTableEntriesRequest, &util.RuntimeOptions{})
	if _desc_err != nil {
		glog.Fatal(_desc_err)
		return nil, _desc_err
	}
	return resp, _desc_err

}
