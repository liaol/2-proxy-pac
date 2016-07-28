#! /bin/bash

# 启动代理
shadowsocks-go -s ip -p 8388 -l 1081 -k password &
shadowsocks-go -s ip -p 8388 -l 1082 -k password &

# 生成pac
2-proxy-pac -b "SOCKS5 127.0.0.1:1082;SOCKS 127.0.0.1:1082;" -d "SOCKS5 127.0.0.1:1081;SOCKS 127.0.0.1:1081;" &

# 设置pac
networksetup -setautoproxyurl "Wi-Fi" "http://127.0.0.1:6789/pac"
