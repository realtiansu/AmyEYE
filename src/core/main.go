package main

import (
	"core/config/jsonloader"
	"core/dataManager"
	"core/kernel"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var (
	configFile = flag.String("config", "", "Config file for AmyEYE.")
	version    = flag.Bool("version", false, "Show current version of AmyEYE.")
	manual     = flag.Bool("manual", false, "Show come detail, can't use when running as openfalcon plugins")
)

func fileExists(file string) bool {
	info, err := os.Stat(file)
	return err == nil && !info.IsDir()
}

func getConfigFilePath() string {
	if len(*configFile) > 0 {
		return *configFile
	} else if workingDir, err := os.Getwd(); err == nil {
		configFile := filepath.Join(workingDir, "agent/plugins/AmyEYE/config/defaults.json")
		if fileExists(configFile) {
			return  configFile
		}
	}

	return ""
}

func loadConfig() (*jsonloader.Config, error) {
	configFile := getConfigFilePath()
	config, err := jsonloader.LoadConfig(configFile)

	return config, err
}

func run(runningTime int) error {
	//读取配置
	config, err := loadConfig()
	if err != nil {
		return err
	}

	//初始化核心
	core, err := kernel.New(&config.Points, &config.Data, *manual)
	if err != nil {
		return err
	}

	//初始化数据管理模块
	dataman, err := dataManager.New(core.Channel, &config.Data)
	if err != nil {
		return err
	}

	//核心启动
	err = core.Start()
	if err != nil {
		return err
	}

	//data manager 启动
	dataman.Start(*manual)

	//http服务
	//httpserver.Start(dataman.Cmd)
	//for {
	//	if !core.Running || !dataman.Running {
	//		break
	//	}
	//	time.Sleep(time.Second * 3)
	//}


	time.Sleep( time.Duration(runningTime) * time.Second )
	defer dataman.Stop(runningTime)
	return nil
}

func main() {
	flag.Parse()
	if *version {
		fmt.Println("1.6")
		return
	}

	err := run( 115 )
	if err != nil {
		os.Exit(23)
	}

}
