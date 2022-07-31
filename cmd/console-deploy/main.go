package main

import (
	"flag"
	"github.com/golang/glog"
	"golang/pkg/abm"
	"golang/pkg/aliyun"
	"golang/pkg/common"
)

var DEPLOYINFO = common.DeployInfo{
	VPC:    "vpc-rj98d2apg2z97c94wl1o8",
	REGION: "us-west-1",
	ZONE:   "us-west-1b",
	TAG: map[string]string{
		"abm_console": "true",
	},
}

func init() {
	flag.Parse()
	flag.Set("alsologtostderr", "true")
	defer glog.Flush()
	ACCOUNT := abm.GetAccountInfo()
	DEPLOYINFO.AK = ACCOUNT["accessKey"]
	DEPLOYINFO.SK = ACCOUNT["accessSecret"]

}

func GetSecurityIdsByTags() []string {
	glog.Info("查询安全组")
	resp, err := aliyun.GetSecurityByTags(DEPLOYINFO)
	if err != nil {
		glog.Fatal(err)
	}
	//glog.Info(resp.Body.SecurityGroups.SecurityGroup)
	securityIds := []string{}
	for _, sg := range resp.Body.SecurityGroups.SecurityGroup {
		//glog.Info(sg)
		//glog.Info(*sg.SecurityGroupId)
		securityIds = append(securityIds, *sg.SecurityGroupId)
	}
	return securityIds
}

func CreateSecurityByTas() (_err error) {
	glog.Info("创建安全组")
	resp, err := aliyun.CreateSecurityByTags(DEPLOYINFO)
	if err != nil {
		glog.Fatal(err)
	}
	glog.Info(resp.Body)
	return err
}

func GetVSWIdsByTags() []string {
	glog.Info("查询VSW")
	resp, _err := aliyun.GetVSW(DEPLOYINFO)
	if _err != nil {
		glog.Fatal(_err)
	}
	//glog.Info("VSW信息:", resp.Body.VSwitches.VSwitch)
	targetVSWId := []string{}
	//检查vsw里是tag 是否包含目标 TAG
	for _, vsw := range resp.Body.VSwitches.VSwitch {
		tagIsInclude := aliyun.TagIsInclude(vsw.Tags.Tag, DEPLOYINFO.TAG)
		if tagIsInclude {
			targetVSWId = append(targetVSWId, *vsw.VSwitchId)
		}
	}
	return targetVSWId
}

func CreateVswWithTag() (_err error) {
	glog.Info("创建VSW")
	err := aliyun.CreateVSW(DEPLOYINFO)
	if err != nil {
		glog.Info(err)
	}
	return err
}

func GetEcsIdsByTag() (ecs []string) {
	glog.Info("查询ECS")
	resp, err := aliyun.GetEcsByTag(DEPLOYINFO)
	if err != nil {
		glog.Fatal(err)
	}
	ecsIds := []string{}
	for _, ecs := range resp.Body.Instances.Instance {
		ecsIds = append(ecsIds, *ecs.InstanceId)
	}
	return ecsIds
}

func CreateECSWithTag(vsw string, sg string) {
	glog.Info("创建ECS")
	err := aliyun.CreateECSWithTag(DEPLOYINFO, vsw, sg)
	if err != nil {
		glog.Fatal(err)
	}

}

//根据命令名获取命令Id，如果不存在会创建，脚本内容会在 golang/cmd/shell 按照名称查找
func CreateCommandIdByNameIfNotExist(commandName string) (commandId string) {
	resp, _err := aliyun.GetCommandByName(DEPLOYINFO, commandName)
	if _err != nil {
		glog.Fatal(_err)
	}
	if *resp.Body.TotalCount <= 0 {
		commandContent := common.GetFileContent(commandName)
		_err := aliyun.CreateCommand(DEPLOYINFO, commandName, &commandContent)
		if _err != nil {
			glog.Fatal(_err)
		}
		resp, _err = aliyun.GetCommandByName(DEPLOYINFO, commandName)
	}
	for _, command := range resp.Body.Commands.Command {
		commandId = *command.CommandId
	}
	return commandId
}

func CreateSnatIfNotExist(vsw string) {

	natInfo := aliyun.GetNatInfo(DEPLOYINFO)

	resp, _err := aliyun.GetSnatEntry(DEPLOYINFO, vsw, natInfo)
	if _err != nil {
		glog.Fatal(_err)
	}
	if *resp.Body.TotalCount < 1 {
		aliyun.CreateSnatEntry(DEPLOYINFO, vsw, natInfo)
	}
}

func GetEcsVSW(ecsId string) (vsw string) {
	glog.Info("获取ECS的VSW")
	ecsDetail, _err := aliyun.GetECSDetailById(DEPLOYINFO, []string{ecsId})
	if _err != nil {
		glog.Fatal(_err)
	}
	if len(ecsDetail.Body.Instances.Instance) < 1 {
		glog.Fatal("找不到ECS信息")
	}

	return *ecsDetail.Body.Instances.Instance[0].VpcAttributes.VSwitchId
}

func main() {
	glog.Info(DEPLOYINFO)
	//获取ECS
	ecsIds := GetEcsIdsByTag()
	glog.Info("ECS信息:", ecsIds)
	if len(ecsIds) < 1 {
		//获取安全组
		securityIds := GetSecurityIdsByTags()
		glog.Info("安全组信息:", securityIds)
		if len(securityIds) < 1 {
			err := CreateSecurityByTas()
			if err != nil {
				glog.Fatal(err)
			}
			securityIds = GetSecurityIdsByTags()
			glog.Info("安全组信息:", securityIds)
		}
		securityId := securityIds[0]
		glog.Info("选择安全组:", securityId)
		//获取VSW
		vswIds := GetVSWIdsByTags()
		glog.Info("VSW信息:", vswIds)
		if len(vswIds) < 1 {
			err := CreateVswWithTag()
			if err != nil {
				glog.Fatal(err)
			}
			vswIds = GetVSWIdsByTags()
			glog.Info("VSW信息:", vswIds)
		}
		vswId := vswIds[0]
		glog.Info("选择VSW:", vswId)
		//创建ECS
		CreateECSWithTag(vswId, securityId)
		ecsIds = GetEcsIdsByTag()
		glog.Info("ECS信息:", ecsIds)
	}
	//选择ECS
	ecsId := ecsIds[0]
	glog.Info("选择ECS:", ecsId)
	glog.Info("启动ECS:", ecsId)
	aliyun.StartECSInstance(DEPLOYINFO, ecsId)
	//创建SNAT
	CreateSnatIfNotExist(GetEcsVSW(ecsId))
	//执行云助手命令
	commandId := CreateCommandIdByNameIfNotExist("console-init.sh")
	glog.Info("云助手命令Id:", commandId)
	aliyun.InvokeECSCommand(DEPLOYINFO, ecsId, commandId)

}
