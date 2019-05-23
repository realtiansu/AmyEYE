package util

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func GetTimeString() string {
	//const shortForm = "2006-01-01 15:04:05"
	//t := time.Now()
	//temp := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.Local)
	//str := temp.Format(shortForm)
	//return str
	return time.Now().String()[:19]
}

func GenHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknow host "
	}

	return hostname
}

func IsDir(f string) bool {
	fi, e := os.Stat(f)
	if e != nil {
		return false
	}
	return fi.IsDir()
}

func GetRunningDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}
