package main

import (
	"fmt"
	"html"
	"net/http"
	"strconv"
)

const APP_NAME = "2-proxy-pac"

func init() {
	//获取pac
	initConfig()
}

func main() {
	http.HandleFunc("/pac", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-ns-proxy-autoconfig;")
		w.Write([]byte(html.UnescapeString(string(GenPac())))) //这里要将html标记转换回来 不然"会变成&#34;
	})
	var port = strconv.Itoa(config.Port)
	fmt.Printf("pac url is http:127.0.0.1:%s/pac\n", port)
	http.ListenAndServe(":"+port, nil)
}
