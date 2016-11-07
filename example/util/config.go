package util

import sls "go-sls"

// Project define Project for test
var Project = &sls.LogProject{
	Name:            "loghub-test",
	Endpoint:        "cn-hangzhou.log.aliyuncs.com",
	AccessKeyID:     "xxx",
	AccessKeySecret: "xxx",
}
