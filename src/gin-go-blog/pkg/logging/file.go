package logging

import (
	"fmt"
	setting "gin-blog/pkg/settting"
	"log"
	"os"

	"time"
)

// var (
// 	LogSavePath = "runtime/logs/"
// 	LogSaveName = "log"
// 	LogFileExt  = "log"
// 	TimeFormat  = "20060102"
// )

func getLogFilePath() string {
	return fmt.Sprintf("%s%s", setting.AppSetting.RuntimeRootPath, setting.AppSetting.LogSavePath)
}

func getLogFileName() string {
	return fmt.Sprintf("%s%s.%s",
		setting.AppSetting.LogSaveName,
		time.Now().Format(setting.AppSetting.TimeFormat),
		setting.AppSetting.LogFileExt)
}

// func getLogFileFullPath() string {
// 	prefixPath := getLogFilePath()
// 	suffixPath := getLogFileName()

// 	return fmt.Sprintf("%s%s", prefixPath, suffixPath)
// }

func openLogFile(filepath string) *os.File {
	_, err := os.Stat(filepath)
	switch {
	case os.IsNotExist(err):
		mkDir()
	case os.IsPermission(err):
		log.Fatalf("Permission :%v", err)
	}

	handle, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Fail to OpenFile :%v", err)
	}

	return handle
}

func mkDir() {
	dir, _ := os.Getwd()
	err := os.MkdirAll(dir+"/"+getLogFilePath(), os.ModePerm)
	if err != nil {
		panic(err)
	}
}
