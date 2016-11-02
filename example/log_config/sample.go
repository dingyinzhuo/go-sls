package main

import (
	"fmt"
	"os"
	sls "sls-sdk-go"
)

var project = &sls.LogProject{
	Name:            "loghub-test",
	Endpoint:        "cn-hangzhou.log.aliyuncs.com",
	AccessKeyID:     "xxx",
	AccessKeySecret: "xxx",
}
var logstore = "demo-store"

func main() {
	// log config sample
	testConf := "test-conf"
	testService := "demo-service"
	createConfig(testConf, logstore, testService)
	updateConfig(testConf)
	getConfig(testConf)

	exist, err := checkConfigExist(testConf)
	if err != nil {
		os.Exit(1)
	}
	if !exist {
		fmt.Println("config:" + testConf + " should be exist")
		os.Exit(1)
	}

	deleteConfig(testConf)

	exist, err = checkConfigExist(testConf)
	if err != nil {
		os.Exit(1)
	}
	if exist {
		fmt.Println("config:" + testConf + " should not be exist")
		os.Exit(1)
	}
	fmt.Println("log config sample end")
}

func checkConfigExist(confName string) (exist bool, err error) {
	exist, err = project.CheckConfigExist(confName)
	if err != nil {
		return false, err
	}
	return exist, nil
}

func deleteConfig(confName string) (err error) {
	err = project.DeleteConfig(confName)
	if err != nil {
		return err
	}
	return nil
}

func updateConfig(configName string) (err error) {
	config, _ := project.GetConfig(configName)
	config.InputDetail.FilePattern = "*.log"
	err = project.UpdateConfig(config)
	if err != nil {
		return err
	}
	return nil
}
func getConfig(configName string) (err error) {
	_, err = project.GetConfig(configName)
	if err != nil {
		return err
	}
	return nil
}
func createConfig(configName string, logstore string, serviceName string) (err error) {
	// 日志所在的父目录
	logPath := "/var/log/lambda/" + serviceName
	// 日志文件的pattern，如functionName.LOG
	filePattern := "*.LOG"
	// 日志时间格式
	timeFormat := "%Y/%m/%d %H:%M:%S"
	// 日志提取后所生成的Key
	key := make([]string, 1)
	// 用于过滤日志所用到的key，只有key的值满足对应filterRegex列中设定的正则表达式日志才是符合要求的
	filterKey := make([]string, 1)
	// 和每个filterKey对应的正正则表达式， filterRegex的长度和filterKey的长度必须相同
	filterRegex := make([]string, 1)
	// topicFormat
	// 1. 用于将日志文件路径的某部分作为topic
	// 2. none 表示topic为空
	// 3. default 表示将日志文件路径作为topic
	// 4. group_topic 表示将应用该配置的机器组topic属性作为topic
	// 以serviceName为topic的正则：/var/log/lambda/([^/]*)/.*
	// 日志路径: /var/log/lambda/my-service/fjaishgaidhfiajf2343/func1.LOG
	topicFormat := "/var/log/lambda/([^/]*)/.*" // topicFormat is right
	inputDetail := sls.InputDetail{
		LogType:       "common_reg_log",
		LogPath:       logPath,
		FilePattern:   filePattern,
		LocalStorage:  true,
		TimeFormat:    timeFormat,
		LogBeginRegex: "", // 日志首行特征
		Regex:         "", // 日志对提取正则表达式
		Keys:          key,
		FilterKeys:    filterKey,
		FilterRegex:   filterRegex,
		TopicFormat:   topicFormat,
	}
	outputDetail := sls.OutputDetail{
		Endpoint:     "",
		LogStoreName: logstore,
	}
	config := &sls.LogConfig{
		Name:         configName,
		InputType:    "file",       //现在只支持file
		OutputType:   "LogService", //现在只支持LogService
		InputDetail:  inputDetail,
		OutputDetail: outputDetail,
	}
	err = project.CreateConfig(config)
	if err != nil {
		return err
	}
	return nil
}
