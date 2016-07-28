package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

func getGFWlist() []string {
	var gfwlistFile = config.ConfigDir + "/gfwlist.txt"
	f, err := os.Stat(gfwlistFile)
	var needFetch, needUpdate = false, false
	if err != nil {
		if os.IsNotExist(err) {
			needFetch = true
		}
	} else if f.ModTime().Before(time.Now().AddDate(0, 0, -1)) {
		needUpdate = true
	}

	if needFetch {
		fmt.Println("首次下载gfwlist.txt")
		fetchGFWlist()
	}

	fl, err := os.Open(gfwlistFile)
	if err != nil {
		panic(err)
	}
	defer fl.Close()
	buf := bufio.NewReader(fl)
	var urlList []string
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
	if needUpdate {
		fmt.Println("重新下载gfwlist.txt")
		go fetchGFWlist()
	}

	return urlList

}

func fetchGFWlist() {
	var gfwUrl = "https://raw.githubusercontent.com/gfwlist/gfwlist/master/gfwlist.txt"
	//var gfwUrl = "http://127.0.0.1:8081/gfwlist.txt"
	resp, err := http.Get(gfwUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	data, err := base64.StdEncoding.DecodeString(string(body))

	urlList := strings.Split(string(data), "\n")

	newUrlList := parseGFWlist(urlList)

	f, err := os.Create(config.ConfigDir + "/gfwlist.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	for _, line := range newUrlList {
		f.WriteString(line + "\n")
	}
}

func parseGFWlist(urlList []string) []string {
	var newUrlList []string
	for _, line := range urlList {
		//过滤空行
		if len(line) == 0 {
			continue
		}
		//过滤注释和直连的
		match, err := regexp.MatchString("^[!,\\[,@].*", line)
		if err != nil || match {
			continue
		}
		//过滤关键字类型的 不含点的
		match, err = regexp.MatchString("[\\.]", line)
		if err != nil || !match {
			continue
		}
		//过滤链接
		match, err = regexp.MatchString("[\\/]", line)
		if err != nil || match {
			continue
		}

		//转换
		reg, err := regexp.Compile("^\\|{1,2}|^\\.")
		if err != nil {
			continue
		}

		newUrlList = append(newUrlList, reg.ReplaceAllString(line, ""))
	}
	return newUrlList
}
