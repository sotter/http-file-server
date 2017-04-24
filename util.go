package main

import (
	"net"
	"regexp"
)

func GetIp(interfaceName string) string{
	_, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	//通过正则将IP拿出来
	ip_regexp, _ := regexp.Compile(`\d+\.\d+\.\d+\.\d+`)
	ifaces, _ := net.InterfaceByName(interfaceName)
	addrs, _ := ifaces.Addrs()
	for _, addr := range addrs {
		ip := ip_regexp.FindString(addr.String())
		if len(ip) > 0 {
			return ip
		}
	}

	return ""
}