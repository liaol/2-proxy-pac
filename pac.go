package main

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
)

func GenPac() []byte {
	var pacRawTmpl = `
var direct = 'DIRECT';
var proxyCN = "{{.ProxyDefault}}";
var proxyUS = "{{.ProxyBlock}}";

var directList = [
"",
"{{.DirectList}}"
];

var directAcc = {};
for (var i = 0; i < directList.length; i += 1) {
	directAcc[directList[i]] = true;
}

var blockList = [
	"{{.BlockList}}"
];

var blockAcc = {};
for (var i = 0; i < blockList.length; i += 1) {
	blockAcc[blockList[i]] = true;
}

var topLevel = {
	"com": true,
	"edu": true,
	"gov": true,
	"net": true,
	"org": true,
	"ac": true,
	"co": true
};

// hostIsIP determines whether a host address is an IP address and whether
// it is private. Currenly only handles IPv4 addresses.
function hostIsIP(host) {
	var part = host.split('.');
	if (part.length != 4) {
		return [false, false];
	}
	var n;
	for (var i = 3; i >= 0; i--) {
		if (part[i].length === 0 || part[i].length > 3) {
			return [false, false];
		}
		n = Number(part[i]);
		if (isNaN(n) || n < 0 || n > 255) {
			return [false, false];
		}
	}
	if (part[0] == '127' || part[0] == '10' || (part[0] == '192' && part[1] == '168')) {
		return [true, true];
	}
	if (part[0] == '172') {
		n = Number(part[1]);
		if (16 <= n && n <= 31) {
			return [true, true];
		}
	}
	return [true, false];
}

function host2Domain(host) {
	var arr, isIP, isPrivate;
	arr = hostIsIP(host);
	isIP = arr[0];
	isPrivate = arr[1];
	if (isPrivate) {
		return "";
	}
	if (isIP) {
		return host;
	}

	var lastDot = host.lastIndexOf('.');
	if (lastDot === -1) {
		return ""; // simple host name has no domain
	}
	// Find the second last dot
	dot2ndLast = host.lastIndexOf(".", lastDot-1);
	if (dot2ndLast === -1)
		return host;
	
	//use the top domain
	return host.substring(dot2ndLast+1);

	var part = host.substring(dot2ndLast+1, lastDot);
	if (topLevel[part]) {
		var dot3rdLast = host.lastIndexOf(".", dot2ndLast-1);
		if (dot3rdLast === -1) {
			return host;
		}
		return host.substring(dot3rdLast+1);
	}
	return host.substring(dot2ndLast+1);
}

function FindProxyForURL(url, host) {
	if (url.substring(0,4) == "ftp:")
		return direct;
	if (host.indexOf(".local", host.length - 6) !== -1) {
		return direct;
	}
	var domain = host2Domain(host);
	if (host.length == domain.length) {
		if (directAcc[host]) {
			return direct;
		}
		if (blockAcc[host]) {
			return proxyUS;
		}
	} else {
		if (directAcc[host] || directAcc[domain]) {
			return direct;
		}

		if (blockAcc[host] || blockAcc[domain]) {
			return proxyUS;
		}
	}
	return proxyCN;
}
`
	pacTmpl, err := template.New("pac").Parse(pacRawTmpl)
	if err != nil {
		panic(err)
	}

	//获取gfwlist
	gfwlist := getGFWlist()

	//获取blocklist并合并
	blockList := getBlockList()
	blockList = append(gfwlist, blockList...)

	//获取directlist
	directList := getDirectList()

	buf := new(bytes.Buffer)
	data := struct {
		ProxyBlock   string
		ProxyDefault string
		DirectList   string
		BlockList    string
	}{
		config.ProxyBlock,
		config.ProxyDefault,
		strings.Join(directList, "\",\n\""),
		strings.Join(blockList, "\",\n\""),
	}

	if err := pacTmpl.Execute(buf, data); err != nil {
		fmt.Println(err)
		panic("Error generating pac file")
	}
	return buf.Bytes()
}
