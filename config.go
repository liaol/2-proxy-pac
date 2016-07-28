package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
)

type Config struct {
	ProxyBlock   string
	ProxyDefault string
	ConfigDir    string
	Port         int
}

var config Config

func initConfig() {
	flag.StringVar(&config.ProxyBlock, "b", "", "国外代理地址")
	flag.StringVar(&config.ProxyDefault, "d", "", "国内代理地址")
	flag.StringVar(&config.ConfigDir, "c", "", "配置目录,默认为$HOME/.2-proxy-pac")
	flag.IntVar(&config.Port, "p", 6789, "port")
	flag.Parse()

	//默认路径为$HOME/2-proxy-pac
	if config.ConfigDir == "" {
		config.ConfigDir = getDefaultConfigDir()
	}
	//准备配置信息
	prepareConfigDir()

}
func prepareConfigDir() {
	dir, err := os.Stat(config.ConfigDir)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(config.ConfigDir, 0755)
		}
	} else {
		if !dir.IsDir() {
			fmt.Sprintln("%s不是目录", config.ConfigDir)
			os.Exit(1)
		}
		//todo 判断可写
	}
}

func getBlockList() []string {
	var urlList []string
	var blockListFile = config.ConfigDir + "/block"

	fl, err := os.Open(blockListFile)
	defer fl.Close()

	if err == nil {
		buf := bufio.NewReader(fl)
		for {
			line, err := buf.ReadString('\n')
			if err != nil {
				break
			}
			line = strings.TrimSpace(line)
			if len(line) > 0 {
				urlList = append(urlList, line)
			}
		}
	}
	return urlList
}

func getDirectList() []string {
	var urlList []string
	var blockListFile = config.ConfigDir + "/direct"

	fl, err := os.Open(blockListFile)
	defer fl.Close()

	if err == nil {
		buf := bufio.NewReader(fl)
		for {
			line, err := buf.ReadString('\n')
			if err != nil {
				break
			}
			line = strings.TrimSpace(line)
			if len(line) > 0 {
				urlList = append(urlList, line)
			}
		}
	}
	return urlList
}

func getDefaultConfigDir() string {
	return path.Join(path.Join(os.Getenv("HOME"), "."+APP_NAME))
}
