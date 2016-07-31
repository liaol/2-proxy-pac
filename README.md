# 2-proxy-pac
一个可以生成支持两个proxy的pac的小工具
两个proxy 一个访问国外网站(gfwlist)，另一个访问国内网站(默认)。

## 背景
公司网络可能会被监控，所以就有了全局代理的需求，但是公司相关网站需要直连。同时还有翻墙的需求。
1. 公司相关网站（内网）直连
2. gfwlist走国外代理
3. 其他的默认走国内代理

## 实现
抓取[gfwlist](https://raw.githubusercontent.com/gfwlist/gfwlist/master/gfwlist.txt)并和配置目录下的direct list,block list合并成pac
  
## 使用方法  
```go get github.com/liaol/2-proxy-pac```
配置目录在$HOME/.2-proxy-pac
将需要直连的域名放在direct文件里，需要翻墙的域名放在block文件里（格式请参考[cow](https://github.com/cyfdecyf/cow)）
修改start.sh

参数说明： ```2-proxy-pac --help```

如需在家里和公司之前切换pac，可参考switch_pac.sh

以上仅在OS X上测试过。

## 感谢
[cow](https://github.com/cyfdecyf/cow)
[shadowsocks-go](https://github.com/shadowsocks/shadowsocks-go)





