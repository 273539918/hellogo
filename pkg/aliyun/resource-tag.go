package aliyun

import vpc20160428 "github.com/alibabacloud-go/vpc-20160428/v2/client"

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
