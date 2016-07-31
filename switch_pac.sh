#! /bin/bash
# 如果在不在公司则切换pac
function SetPac() {
    pac=`/usr/sbin/networksetup -getautoproxyurl Wi-Fi | grep $1 | wc -l`
    if [ $pac -eq 0 ]
    then
        /usr/sbin/networksetup -setautoproxyurl "Wi-Fi" $1
    fi
}

pac_company=""
pac_home=""

ip=`/sbin/ifconfig | grep 172 | wc -l`

if [ $ip -eq 0 ]
then
    SetPac $pac_home
else 
    SetPac $pac_company
fi
